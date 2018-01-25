package web

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

		templateName := strings.TrimPrefix(strings.TrimPrefix(file, directory), "/")
		templateText += "{{ define \"" + templateName + "\" }}" + string(contents) + "{{ end }}\n"
	}

	return &Renderer{
		templateText: templateText,
	}, nil
}
