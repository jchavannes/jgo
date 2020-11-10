package db_util

import (
	"errors"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"reflect"
	"sort"
	"strings"

	"github.com/jchavannes/gorm"
)

// FROM: https://github.com/bombsimon/gorm-bulk

type execFunc func(scope *gorm.Scope, columnNames, groups []string)

// InsertFunc is the default insert func. It will pass a gorm.Scope pointer
// which holds all the vars in scope.SQLVars. The value set to scope.SQL
// will be used as SQL and the variables in scope.SQLVars will be used as
// values.
//
//  INSERT INTO `tbl`
//    (col1, col2)
//  VALUES
//    (?, ?), (?, ?)
func insertFunc(scope *gorm.Scope, columnNames, groups []string) {
	defaultWithFormat(scope, columnNames, groups, "INSERT INTO %s (%s) VALUES %s")
}

// InsertIgnoreFunc will run INSERT IGNORE with all the records and values set
// on the passed scope pointer.
//
//  INSERT IGNORE INTO `tbl`
//    (col1, col2)
//  VALUES
//    (?, ?), (?, ?)
func insertIgnoreFunc(scope *gorm.Scope, columnNames, groups []string) {
	defaultWithFormat(scope, columnNames, groups, "INSERT IGNORE INTO %s (%s) VALUES %s")
}

// InsertOnDuplicateKeyUpdateFunc will perform a bulk insert but on duplicate key
// perform an update.
//
//  INSERT INTO `tbl`
//    (col1, col2)
//  VALUES
//    (?, ?), (?, ?)
//  ON DUPLICATE KEY UPDATE
//    col1 = VALUES(col1),
//    col2 = VALUES(col2)
func insertOnDuplicateKeyUpdateFunc(scope *gorm.Scope, columnNames, groups []string) {
	var duplicateUpdates []string

	for i := range columnNames {
		// Don't update created at on duplicate.
		if columnNames[i] == "`created_at`" {
			continue
		}

		duplicateUpdates = append(
			duplicateUpdates,
			fmt.Sprintf("%s = VALUES(%s)", columnNames[i], columnNames[i]),
		)
	}

	// This is not SQL string formatting, prepare statements is in use.
	// nolint: gosec
	scope.Raw(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
		scope.QuotedTableName(),
		strings.Join(columnNames, ", "),
		strings.Join(groups, ", "),
		strings.Join(duplicateUpdates, ", "),
	))
}

func defaultWithFormat(scope *gorm.Scope, columnNames, groups []string, format string) {
	var (
		extraOptions string
		sqlFormat    = fmt.Sprintf("%s%%s", format)
	)

	if insertOption, ok := scope.Get("gorm:insert_option"); ok {
		// Add the extra insert option
		extraOptions = fmt.Sprintf(" %s", insertOption)
	}

	scope.Raw(fmt.Sprintf(
		sqlFormat,
		scope.QuotedTableName(),
		strings.Join(columnNames, ", "),
		strings.Join(groups, ", "),
		extraOptions,
	))
}

// BulkInsert will call BulkExec with the default InsertFunc.
func bulkInsert(db *gorm.DB, objects []interface{}) error {
	errs := bulkExecChunk(db, objects, insertFunc, LargeLimit)
	if len(errs) == 0 {
		return nil
	}
	return jerr.Combine(errs...)
}

// BulkInsertIgnore will call BulkExec with the default InsertFunc.
func bulkInsertIgnore(db *gorm.DB, objects []interface{}) error {
	return bulkExec(db, objects, insertIgnoreFunc)
}

// BulkInsertOnDuplicateKeyUpdate will call BulkExec with the default InsertFunc.
func bulkInsertOnDuplicateKeyUpdate(db *gorm.DB, objects []interface{}) error {
	return bulkExec(db, objects, insertOnDuplicateKeyUpdateFunc)
}

// BulkExecChunk will split the objects passed into the passed chunk size. A
// slice of errors will be returned (if any).
func bulkExecChunk(db *gorm.DB, objects []interface{}, execFunc execFunc, chunkSize int) []error {
	var allErrors []error

	for {
		var chunkObjects []interface{}

		if len(objects) <= chunkSize {
			chunkObjects = objects
			objects = []interface{}{}
		} else {
			chunkObjects = objects[:chunkSize]
			objects = objects[chunkSize:]
		}

		if err := bulkExec(db, chunkObjects, execFunc); err != nil {
			allErrors = append(allErrors, err)
		}

		// Nothing more to do
		if len(objects) < 1 {
			break
		}
	}

	if len(allErrors) > 0 {
		return allErrors
	}

	return nil
}

// BulkExec will convert a slice of interface to bulk SQL statement. The final
// SQL will be determined by the ExecFunc passed.
func bulkExec(db *gorm.DB, objects []interface{}, execFunc execFunc) error {
	scope, err := scopeFromObjects(db, objects, execFunc)
	if err != nil {
		return err
	}

	// No scope and no error means nothing to do
	if scope == nil {
		return nil
	}

	return db.Exec(scope.SQL, scope.SQLVars...).Error
}

func scopeFromObjects(db *gorm.DB, objects []interface{}, execFunc execFunc) (*gorm.Scope, error) {
	// No objects passed, nothing to do.
	if len(objects) < 1 {
		return nil, nil
	}

	var (
		columnNames       []string
		quotedColumnNames []string
		placeholders      []string
		groups            []string
		scope             = db.NewScope(objects[0])
		bulkNow           = gorm.NowFunc()
	)

	// Get a map of the first element to calculate field names and number of
	// placeholders.
	firstObjectFields, err := objectToMap(objects[0])
	if err != nil {
		return nil, err
	}

	for k := range firstObjectFields {
		// Add raw column names to use for iteration over each row later to get
		// the correct order of columns.
		columnNames = append(columnNames, k)

		// Add as many placeholders (question marks) as there are columns.
		placeholders = append(placeholders, "?")

		// Sort the column names to ensure the right order.
		sort.Strings(columnNames)
	}

	// We must setup quotedColumnNames after sorting columnNames since sorting
	// of quoted fields might differ from sorting without. This way we know that
	// columnNames is the master of the order and will be used both when setting
	// field and values order.
	for i := range columnNames {
		quotedColumnNames = append(quotedColumnNames, scope.Quote(columnNames[i]))
	}

	for _, r := range objects {
		objectScope := db.NewScope(r)

		row, err := objectToMap(r)
		if err != nil {
			return nil, err
		}

		for _, key := range columnNames {
			field := row[key]
			value := field.Field.Interface()

			switch field.Struct.Name {
			// Column CreatedAt and UpdatedAt with zero value will be set to same time
			case "CreatedAt", "UpdatedAt":
				if field.IsBlank {
					value = bulkNow
				}
			}

			objectScope.AddToVars(value)
		}

		groups = append(
			groups,
			fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")),
		)

		// Add object vars to the outer scope vars
		scope.SQLVars = append(scope.SQLVars, objectScope.SQLVars...)
	}

	execFunc(scope, quotedColumnNames, groups)

	return scope, nil
}

// ObjectToMap takes any object of type <T> and returns a map with the gorm
// field DB name as key and the value as value. Special fields and actions
//  * Foreign keys - Will be left out
//  * Relationship fields - Will be left out
//  * Fields marked to be ignored - Will be left out
//  * Fields named ID with auto increment - Will be left out
//  * Fields named ID set as primary key with blank value - Will be left out
//  * Blank fields with default value - Will be set to the default value
func objectToMap(object interface{}) (map[string]*gorm.Field, error) {
	var (
		attributes = map[string]*gorm.Field{}
	)

	// De-reference pointers (and it's values)
	rv := reflect.ValueOf(object)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		object = rv.Interface()
	}

	if rv.Kind() != reflect.Struct {
		return nil, errors.New("value must be kind of Struct")
	}

	for _, field := range (&gorm.Scope{Value: object}).Fields() {
		// Exclude relational record because it's not directly contained in database columns
		_, hasForeignKey := field.TagSettings["FOREIGNKEY"]
		if hasForeignKey {
			continue
		}

		if field.StructField.Relationship != nil {
			continue
		}

		if field.IsIgnored {
			continue
		}

		// Let the DBM set the default values since these might be meta values
		// such as 'CURRENT_TIMESTAMP'. Has default will be set to true also for
		// 'AUTO_INCREMENT' fields which is not primary keys so we must check
		// that we've ACTUALLY configured a default value and uses the tag
		// before we skip it.
		if field.StructField.HasDefaultValue && field.IsBlank {
			if _, ok := field.TagSettings["DEFAULT"]; ok {
				continue
			}
		}

		// Skip blank primary key fields named ID. They're probably coming from
		// `gorm.Model` which doesn't have the AUTO_INCREMENT tag.
		if field.DBName == "id" && field.IsPrimaryKey && field.IsBlank {
			continue
		}

		// Check if auto increment is set (but not set to false). If so skip the
		// field and let the DBM auto increment the value.
		if value, ok := field.TagSettings["AUTO_INCREMENT"]; ok {
			if !strings.EqualFold(value, "false") {
				continue
			}
		}

		attributes[field.DBName] = field
	}

	return attributes, nil
}
