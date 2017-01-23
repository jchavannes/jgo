package web

import (
	"io"
	"text/template"
	"errors"
)

type Renderer struct {
	funcMap      template.FuncMap
	templateText string
}

func (r *Renderer) getTemplate() *template.Template {
	return template.Must(template.New("_base").Funcs(r.funcMap).Parse(r.templateText))
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
	return errors.New("Unable to find template.")
}
