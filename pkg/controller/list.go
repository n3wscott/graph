package controller

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

var listOnce sync.Once
var listT *template.Template

func (c *Controller) ListHandler(w http.ResponseWriter, r *http.Request) {
	listOnce.Do(func() {
		var err error
		listT, err = template.ParseFiles(
			c.root+"/templates/list/index.html",
			c.root+"/templates/list/main.html",
		)
		if err != nil {
			log.Printf("Failed to parse template: %v\n", err)
		}
	})

	var data map[string]interface{}

	err := listT.Execute(w, data)
	if err != nil {
		log.Printf("template execute error %s", err)
	}
}
