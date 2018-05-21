package shortcuts  

import (
	"html/template"
	"net/http"
	"fmt"
)

func Render(w http.ResponseWriter, file string, basefile string, data interface{}) {
	if basefile == "" {
		tmpl := template.Must(template.ParseFiles(file))
		err := tmpl.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		tmpl := template.Must(template.ParseFiles(file, basefile))
		err := tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			fmt.Println(err)
		}
	}
} 

func HttpError(w http.ResponseWriter, code int, message string) {
	type exception struct {
		Message string
	}
	switch code {
	case 404:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error 404: %s\n", message)
		return
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error 500: %s\n", message)
		return
	case 405:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Error 405: %s\n", message)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error 400: %s\n", message)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Error 401: %s\n", message)
	}
}