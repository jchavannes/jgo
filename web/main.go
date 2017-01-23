package web

import (
	"net/http"
	"path/filepath"
	"os"
	"strings"
	"io/ioutil"
	"log"
)

func GetFilenameFromRequest(r *http.Request) string {
	return r.RequestURI[1:len(r.RequestURI)]
}

func GetRenderer(directory string) (*Renderer, error) {
	fileList := []string{}
	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	check(err)

	templateText := ""
	for _, file := range fileList {
		if !strings.HasSuffix(file, ".html") {
			continue
		}

		contents, err := ioutil.ReadFile(file)
		check(err)

		filename := file[len(directory) + 1:]
		templateText += "{{ define \"" + filename + "\" }}" + string(contents) + "{{ end }}\n"
	}

	return &Renderer{
		templateText: templateText,
	}, nil
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
