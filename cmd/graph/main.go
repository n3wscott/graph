package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/graph/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/clients/dynamicclient"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type envConfig struct {
	FilePath  string `envconfig:"FILE_PATH" default:"/var/run/ko/" required:"true"`
	Namespace string `envconfig:"NAMESPACE" default:"default" required:"true"`
	Port      int    `envconfig:"PORT" default:"8080" required:"true"`
	Target    string `envconfig:"TARGET"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	if !strings.HasSuffix(env.FilePath, "/") {
		env.FilePath = env.FilePath + "/"
	}

	ctx, _ := injection.EnableInjectionOrDie(nil, nil)
	client := dynamicclient.Get(ctx)

	c := controller.New(env.FilePath, env.Namespace, client)

	if env.Target != "" {
		t, err := cloudevents.NewHTTP(
			cloudevents.WithTarget(env.Target),
		)
		if err != nil {
			panic(err)
		}
		ce, err := cloudevents.NewClient(t,
			cloudevents.WithTimeNow(),
			cloudevents.WithUUIDs(),
		)
		if err != nil {
			panic(err)
		}
		c.CE = ce
		log.Println("Will target", env.Target)
	}

	c.Mux().Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(env.FilePath+"static"))))

	c.Mux().HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.RootHandler(w, r)
	})

	log.Println("Listening on", env.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), c.Mux()))
}
