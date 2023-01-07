package main

import (
	"flag"
	f "github.com/ambelovsky/gosf"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func setupK8s() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	createGoDeployment(clientset)
	createPhpDeployment(clientset)

	return clientset
}

func cleanupK8s(clientset *kubernetes.Clientset) {
	deleteDeployment(clientset, "go-deployment")
	deleteDeployment(clientset, "php-deployment")
}

func main() {
	c := make(chan os.Signal)
	clientset := setupK8s()
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanupK8s(clientset)
		f.Shutdown()
		os.Exit(1)
	}()

	server := setupSocketServer()
	setupRoutes(clientset, server)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)

	log.Println("Serving at localhost:3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))

}

func int32Ptr(i int32) *int32 { return &i }
