package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

func killRandomPod(clientset *kubernetes.Clientset, app string) {
	rand.Seed(time.Now().UnixNano())
	log.Println("Deleting pod...")
	deletePolicy := metav1.DeletePropagationForeground
	options := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s-k8s", app),
	}
	podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), options)

	randomPod := rand.IntnRange(0, len(podList.Items)-1)
	podName := podList.Items[randomPod].Name

	err = clientset.CoreV1().Pods("default").Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if errors.IsNotFound(err) {
		log.Printf("Pod %s in namespace %s not found\n", podName, "default")
	} else if err != nil {
		log.Println(err)
	} else {
		log.Println("Deleted pod.")
	}
}

func killPod(clientset *kubernetes.Clientset, name string) {
	log.Printf("Deleting pod %s", name)
	deletePolicy := metav1.DeletePropagationForeground
	err := clientset.CoreV1().Pods("default").Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if errors.IsNotFound(err) {
		log.Printf("Pod %s in namespace %s not found\n", name, "default")
	} else if err != nil {
		log.Println(err)
	} else {
		log.Println("Deleted pod.")
	}
}
