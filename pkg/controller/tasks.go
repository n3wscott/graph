package controller

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"strings"

	"github.com/n3wscott/graph/pkg/knative"
)

type Task struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (c *Controller) TasksHandler(resp http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		c.GetTasksHandler(resp, req)
		return
	case http.MethodPost:
	case http.MethodDelete:
		c.DeleteTasksHandler(resp, req)
		return
	}

	resp.WriteHeader(404)
}

func (c *Controller) GetTasksHandler(resp http.ResponseWriter, req *http.Request) {

	kn := knative.New(c.client)
	gvr := schema.GroupVersionResource{
		Group:    "n3wscott.com",
		Version:  "v1alpha1",
		Resource: "tasks",
	}

	addrs := kn.Addressable(c.namespace, gvr)

	list := make([]Task, 0)
	for _, a := range addrs {
		list = append(list, Task{
			ID:        string(a.UID),
			Name:      a.Name,
			Namespace: a.Namespace,
		})
	}

	b, err := json.Marshal(list)
	if err != nil {
		resp.WriteHeader(500)
		return
	}

	_, _ = resp.Write(b)
}

func (c *Controller) DeleteTasksHandler(resp http.ResponseWriter, req *http.Request) {

	parts := strings.Split(req.URL.Path, "/")
	id := parts[len(parts)-1]

	fmt.Println("here", id)
}
