package main

import (
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/dynamic"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/n3wscott/graph/pkg/config"
	"github.com/n3wscott/graph/pkg/controller"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type envConfig struct {
	FilePath  string `envconfig:"FILE_PATH" default:"/var/run/ko/" required:"true"`
	Namespace string `envconfig:"NAMESPACE" default:"default" required:"true"`
	Port      int    `envconfig:"PORT" default:"8080" required:"true"`
}

var (
	cluster    string
	kubeconfig string
)

func init() {
	flag.StringVar(&cluster, "cluster", "",
		"Provide the cluster to test against. Defaults to the current cluster in kubeconfig.")

	var defaultKubeconfig string
	if usr, err := user.Current(); err == nil {
		defaultKubeconfig = path.Join(usr.HomeDir, ".kube/config")
	}

	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig,
		"Provide the path to the `kubeconfig` file.")
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

	cfg, err := config.BuildClientConfig(kubeconfig, cluster)
	if err != nil {
		log.Fatalf("Error building kubeconfig", err)
	}
	client := dynamic.NewForConfigOrDie(cfg)

	c := controller.New(env.FilePath, env.Namespace, client)

	c.Mux().Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(env.FilePath+"static"))))

	c.Mux().HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.RootHandler(w, r)
	})

	log.Println("Listening on", env.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), c.Mux()))
}
