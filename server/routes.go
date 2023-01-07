package main

import (
	socketio "github.com/googollee/go-socket.io"
	"k8s.io/client-go/kubernetes"
	"log"
)

func setupRoutes(clientset *kubernetes.Clientset, server *socketio.Server) {
	goRoutes(clientset, server)
	phpRoutes(clientset, server)
}

func goRoutes(clientset *kubernetes.Clientset, server *socketio.Server) {
	app := "go"
	namespace := "/go"
	server.OnEvent(namespace, "get_pods", func(s socketio.Conn, msg string) {
		log.Println("getting go pods:", msg)
		totalReplicas, availableReplicas, podList := checkPods(clientset, app)
		s.Emit("list", map[string]interface{}{
			"total":     totalReplicas,
			"available": availableReplicas,
			"podList":   podList,
		})
	})

	server.OnEvent(namespace, "kill_random_pod", func(s socketio.Conn, msg string) {
		log.Println("killing random go pod")
		killRandomPod(clientset, app)
	})

	server.OnEvent(namespace, "kill_pod", func(s socketio.Conn, msg string) {
		log.Printf("killing pod: %s", msg)
		killPod(clientset, msg)
	})
}

func phpRoutes(clientset *kubernetes.Clientset, server *socketio.Server) {
	app := "php"
	namespace := "/php"
	server.OnEvent(namespace, "get_pods", func(s socketio.Conn, msg string) {
		log.Println("getting php pods:", msg)
		totalReplicas, availableReplicas, podList := checkPods(clientset, app)
		s.Emit("list", map[string]interface{}{
			"total":     totalReplicas,
			"available": availableReplicas,
			"podList":   podList,
		})
	})

	server.OnEvent(namespace, "kill_random_pod", func(s socketio.Conn, msg string) {
		log.Println("killing random php pod")
		killRandomPod(clientset, app)
	})

	server.OnEvent(namespace, "kill_pod", func(s socketio.Conn, msg string) {
		log.Printf("killing pod: %s", msg)
		killPod(clientset, msg)
	})
}
