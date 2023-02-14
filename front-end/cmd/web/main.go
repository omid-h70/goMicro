package main

import (
	"errors"
	"fmt"
	"helpers"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

const defaultPortNumber = "80"

// Filename is the __filename equivalent
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Dirname is the __dirname equivalent
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}

func main() {
	cwd, _ := Dirname();

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "index.page.gohtml")
	})

	destMapPath := filepath.Join(cwd, "/assets/css/")
	http.Handle("/css/", http.StripPrefix("/css/",http.FileServer(http.Dir(destMapPath))))

	serverAddr := fmt.Sprintf(":%s", helpers.GetEnvVar("FRONT_END_PORT", defaultPortNumber))
	fmt.Printf("Starting front end service on %s \n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, page string) {
	var templateSlice []string
	partials := []string{
		"./templates/base.layout.gohtml",
		"./templates/header.partial.gohtml",
		"./templates/footer.partial.gohtml",
	}

	cwd, _ := Dirname();

	/*Whatever Pages Come - use these Templates*/
	templateSlice = append(templateSlice, filepath.Join(cwd, fmt.Sprintf("./templates/%s", page)))
	/* /Whatever Pages Come - use these Templates */

	for _, partial := range partials {
		relativePath := filepath.Join(cwd, partial)
		templateSlice = append(templateSlice, relativePath)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
