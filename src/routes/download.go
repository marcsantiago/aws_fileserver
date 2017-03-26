package routes

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func findFile(fileName string) (string, error) {
	files, err := ioutil.ReadDir(staticPath)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if f.Name() == fileName {
			return filepath.Join(staticPath, fileName), nil
		}
	}
	return "", err
}

// Downloader ...
func Downloader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fileName := r.FormValue("file_name")
		if fileName != "" {
			path, err := findFile(fileName)
			if err == nil {
				if path != "" {
					var buf bytes.Buffer
					w.Header().Set("Content-Type", "text/csv")
					w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
					file, err := os.Open(path)
					defer file.Close()
					if err == nil {
						io.Copy(&buf, file)
						w.Write(buf.Bytes())
						buf.Reset()
					}
					return
				}
			}
		}
		return
	}
}
