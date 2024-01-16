package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

var functions = template.FuncMap{}

// RenderTemplate render templates using html/templates
func RenderTemplate(w http.ResponseWriter, path string) {
	tc, err := CreateTemplateCache(w)
	if err != nil {
		log.Fatal(err)
	}
	t, ok := tc[path]
	if !ok {
		fmt.Println("Parsing template error: ", err)
		return
	}
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, nil)

	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error write tempalate to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache(w http.ResponseWriter) (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
