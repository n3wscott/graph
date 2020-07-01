package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/n3wscott/graph/pkg/knative"
)

type Task struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

var taskGVR = schema.GroupVersionResource{
	Group:    "n3wscott.com",
	Version:  "v1alpha1",
	Resource: "tasks",
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
	addrs := knative.New(c.client).Addressable(c.namespace, taskGVR)

	list := make([]Task, 0)
	for _, a := range addrs {
		if a.Status.Address == nil {
			continue
		}
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
	resp.Header().Set("content-type", "application/json")
	_, _ = resp.Write(b)
}

func (c *Controller) DeleteTasksHandler(resp http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")
	id := parts[len(parts)-1]

	addrs := knative.New(c.client).Addressable(c.namespace, taskGVR)
	for _, a := range addrs {
		if string(a.UID) == id {
			fmt.Println("will delete ", a)

			if c.CE != nil {
				event := cloudevents.NewEvent()
				event.SetType("com.n3wscott.target")
				event.SetSource("n3wscott/graph")
				event.SetExtension("target", fmt.Sprintf("%s/%s", a.Namespace, a.Name))
				if result := c.CE.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
					fmt.Printf("failed to send: %s\n", result)
				} else {
					fmt.Printf("sent: %s, %s\n", event.String(), result)
				}
			}
		}
	}

	resp.Header().Set("content-type", "application/json")
	_, _ = resp.Write([]byte("{}"))
}
