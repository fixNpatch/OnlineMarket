package Router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("index called")
	tmpl := template.Must(template.ParseFiles("./front/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func User(w http.ResponseWriter, req *http.Request) {
	fmt.Println("user called")
}

func ManageStatic(w http.ResponseWriter, req *http.Request) {
	fileBytes, err := ioutil.ReadFile(StaticPath + req.URL.Path[7:])
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
	}
}