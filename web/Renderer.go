package web

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Renderer struct {
	funcMap      template.FuncMap
	templateText string
}

var defaultFuncMap = template.FuncMap{
	"loop": func(n uint) []struct{} {
		return make([]struct{}, n)
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
	"formatFloat": func(f float32) string {
		str := strconv.FormatFloat(float64(f), 'f', 2, 32)
		re := regexp.MustCompile("(\\d+)(\\d{3})")
		for i := 0; i < (len(str)-1)/3; i++ {
			str = re.ReplaceAllString(str, "$1,$2")
		}
		return str
	},
	"formatBigFloat": func(f float64) string {
		str := strconv.FormatFloat(f, 'f', 8, 64)
		return str
	},
	"formatBigInt": func(f int64) string {
		str := strconv.Itoa(int(f))
		re := regexp.MustCompile("(\\d+)(\\d{3})")
		for i := 0; i < (len(str)-1)/3; i++ {
			str = re.ReplaceAllString(str, "$1,$2")
		}
		return str
	},
	"formatUInt": func(f uint) string {
		str := strconv.Itoa(int(f))
		re := regexp.MustCompile("(\\d+)(\\d{3})")
		for i := 0; i < (len(str)-1)/3; i++ {
			str = re.ReplaceAllString(str, "$1,$2")
		}
		return str
	},
	"formatInt": func(f int) string {
		str := strconv.Itoa(f)
		re := regexp.MustCompile("(\\d+)(\\d{3})")
		for i := 0; i < (len(str)-1)/3; i++ {
			str = re.ReplaceAllString(str, "$1,$2")
		}
		return str
	},
	"getUnique": func(n int) string {
		const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		return string(b)
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
}

func (r *Renderer) getTemplate() *template.Template {
	funcMap := make(template.FuncMap)
	for k, v := range defaultFuncMap {
		funcMap[k] = v
	}
	for k, v := range r.funcMap {
		funcMap[k] = v
	}
	return template.Must(template.New("_base").Funcs(funcMap).Parse(r.templateText))
}

func (r *Renderer) SetFuncMap(funcMap map[string]interface{}) {
	r.funcMap = template.FuncMap(funcMap)
}

func (r *Renderer) Render(names []string, writer io.Writer, data interface{}) error {
	t := r.getTemplate()
	for _, name := range names {
		if t.Lookup(name) != nil {
			return t.ExecuteTemplate(writer, name, data)
		}
	}
	return errors.New("unable to find template")
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
