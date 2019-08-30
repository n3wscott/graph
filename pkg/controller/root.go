package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/n3wscott/graph/pkg/graph"
	"github.com/n3wscott/graph/pkg/knative"
)

var once sync.Once
var t *template.Template

func getQueryParam(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return ""
	}
	return keys[0]
}

// TODO: support just fetching the graph image

var defaultFormat = "svg"    // or png
var defaultFocus = "trigger" // or png

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		var err error
		t, err = template.ParseFiles(
			c.root+"/templates/index.html",
			c.root+"/templates/main.html",
		)
		if err != nil {
			log.Printf("Failed to parse template: %v\n", err)
		}
	})

	fmt.Println("handling", r.URL)

	format := getQueryParam(r, "format")
	if format == "" {
		format = defaultFormat
	}

	focus := getQueryParam(r, "focus")
	if focus == "" {
		focus = defaultFocus
	}

	var dotGraph string
	var yv []knative.YamlView

	switch focus {
	case "sub", "subs", "subscription", "subscriptions":
		dotGraph = graph.ForSubscriptions(c.client, c.namespace)
	case "broker", "trigger", "triggers":
		fallthrough
	default:
		dotGraph, yv = graph.ForTriggers(c.client, c.namespace)
	}

	file, err := dotToImage(format, []byte(dotGraph))
	if err != nil {
		log.Printf("dotToImage error %s\n\n%s\n", err, dotGraph)
		return
	}
	img, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("read file error %s", err)
	}

	defer os.Remove(file) // clean up

	var data map[string]interface{}
	log.Printf("Image is %s", format)
	if format == "svg" {
		data = map[string]interface{}{
			"svg":    true,
			"Image":  template.HTML(string(img)),
			"Format": format,
			"yv":     yv,
		}
	} else {
		data = map[string]interface{}{
			"Image":  base64.StdEncoding.EncodeToString(img),
			"Format": fmt.Sprintf("image/%s;base64", format),
		}
	}
	data["Dot"] = template.HTML(dotGraph)

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("template execute error %s", err)
	}
}

var dot string

func dotToImage(format string, b []byte) (string, error) {
	if dot == "" {
		var err error
		dot, err = exec.LookPath("dot")
		if err != nil {
			log.Fatalln("unable to find program 'dot', please install it or check your PATH")
		}
	}

	var img = filepath.Join(os.TempDir(), fmt.Sprintf("graph.%s", format))

	cmd := exec.Command(dot, fmt.Sprintf("-T%s", format), "-o", img)
	cmd.Stdin = bytes.NewBuffer(b)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return img, nil
}
