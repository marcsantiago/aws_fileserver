package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	staticPath    string
	indexTemplate *template.Template
)

func init() {
	path, _ := os.Getwd()
	staticPath = filepath.Join(path, "/static")
	indexTemplate = template.Must(template.New("").ParseFiles(filepath.Join(path, "/src/templates/index.tmpl")))
}

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		files, _ := ioutil.ReadDir(staticPath)
		var fileNames []string
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}
		data := map[string]interface{}{"Files": fileNames}
		err := indexTemplate.ExecuteTemplate(w, "index", data)
		if err != nil {
			panic(err)
		}
	}
}
