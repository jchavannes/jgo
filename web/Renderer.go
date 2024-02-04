package web

import (
	"errors"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jfmt"
	"github.com/jchavannes/jgo/jutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Renderer struct {
	funcMap      template.FuncMap
	templateText string
}

var defaultFuncMap = template.FuncMap{
	"loop": func(n uint) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i
		}
		return s
	},
	"contains": func(s, substr string) bool {
		return strings.Contains(s, substr)
	},
	"substr": func(s string, ints ...int) string {
		var start, end int
		if len(ints) >= 2 {
			start = ints[0]
			end = ints[1]
		} else {
			end = ints[0]
		}
		if start > end || end > len(s) {
			return s
		}
		return s[start:end]
	},
	// Allows passing in multiple variables into a template.
	// Pass in pairs of keys and values.
	// Example: {{ template "index.html" dict "MyVar" .SomeVar "MySecondVar" .SomeOtherVar }}
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	},
	"inSlice": func(needle string, haystack ...string) bool {
		for _, hay := range haystack {
			if needle == hay {
				return true
			}
		}
		return false
	},
	"formatFloat": func(f float32, decimals ...int) string {
		return jfmt.AddCommasFloat(float64(f), decimals...)
	},
	"formatBigFloat": func(f float64, decimals ...int) string {
		return jfmt.AddCommasFloat(f, decimals...)
	},
	"formatBigUInt": func(i uint64) string {
		return jfmt.AddCommasUint(i)
	},
	"formatBigInt": func(i int64) string {
		return jfmt.AddCommas(i)
	},
	"formatUInt": func(i uint) string {
		return jfmt.AddCommasUint(uint64(i))
	},
	"formatInt": func(i int) string {
		return jfmt.AddCommas(int64(i))
	},
	"formatInt32": func(i int32) string {
		return jfmt.AddCommas(int64(i))
	},
	"getUnique": func(n int) string {
		return jutil.GetUnique(n)
	},
	"mod": func(i, j int) bool {
		return i%j == 0
	},
	"add": func(a, b int) int {
		return a + b
	},
	"minus": func(a, b int) int {
		return a - b
	},
	"before": func(a string) bool {
		t := getTime(a)
		return time.Now().Before(t)
	},
	"after": func(a string) bool {
		t := getTime(a)
		return time.Now().After(t)
	},
	"isNil": func(a interface{}) bool {
		return jutil.IsNil(a)
	},
}

func getTime(a string) time.Time {
	var fmt string
	if len(a) == 16 {
		fmt = "2006-01-02 15:04"
	} else if len(a) == 10 {
		fmt = "2006-01-02"
	} else {
		return time.Time{}
	}
	t, _ := time.ParseInLocation(fmt, a, time.Local)
	return t
}

func (r *Renderer) getTemplate() (*template.Template, error) {
	funcMap := make(template.FuncMap)
	for k, v := range defaultFuncMap {
		funcMap[k] = v
	}
	for k, v := range r.funcMap {
		funcMap[k] = v
	}
	t, err := template.New("_base").Funcs(funcMap).Parse(r.templateText)
	if err != nil {
		return nil, jerr.Get("error parsing template", err)
	}
	return t, nil
}

func (r *Renderer) SetFuncMap(funcMap map[string]interface{}) {
	r.funcMap = template.FuncMap(funcMap)
}

const UnableToFindTemplateErrorMsg = "unable to find template"
const WriteTimeoutErrorMsgPart = "write: broken pipe"

func (r *Renderer) Render(names []string, writer io.Writer, data interface{}) error {
	t, err := r.getTemplate()
	if err != nil {
		return jerr.Get("error getting template", err)
	}
	for _, name := range names {
		if t.Lookup(name) != nil {
			return t.ExecuteTemplate(writer, name, data)
		}
	}
	return errors.New(UnableToFindTemplateErrorMsg)
}

func GetRenderer(directory string) (*Renderer, error) {
	var fileList []string
	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	templateText := ""
	for _, file := range fileList {
		if !strings.HasSuffix(file, ".html") {
			continue
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}

		file = strings.Replace(file, "\\", "/", -1)
		templateName := strings.TrimPrefix(strings.TrimPrefix(file, directory), "/")
		templateText += "{{ define \"" + templateName + "\" }}" + string(contents) + "{{ end }}\n"
	}

	return &Renderer{
		templateText: templateText,
	}, nil
}
