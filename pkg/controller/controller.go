package controller

import (
	"k8s.io/client-go/dynamic"
	"net/http"
	"sync"
)

type Controller struct {
	root string
	mux  *http.ServeMux
	once sync.Once

	namespace string
	client    dynamic.Interface
}

func New(root, namespace string, client dynamic.Interface) *Controller {
	return &Controller{root: root, namespace: namespace, client: client}
}

func (c *Controller) Mux() *http.ServeMux {
	c.once.Do(func() {
		m := http.NewServeMux()
		m.HandleFunc("/ui", c.RootHandler)
		c.mux = m
	})

	return c.mux
}
