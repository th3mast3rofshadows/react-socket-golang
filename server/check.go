package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func checkPods(clientset *kubernetes.Clientset, app string) (int32, int32, []interface{}) {
	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := apiv1.NamespaceDefault
	deploymentData := getDeployment(clientset, namespace, app)
	podList := getPodList(clientset, namespace, app)

	log.Printf("%d pods in podlist", len(podList.Items))
	return deploymentData.Status.Replicas, deploymentData.Status.ReadyReplicas, getPodData(podList.Items)
}

func getPodList(clientset *kubernetes.Clientset, namespace string, app string) *v1.PodList {
	log.Printf("Getting podlist for app: %s", app)
	options := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s-k8s", app),
	}
	podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), options)
	if errors.IsNotFound(err) {
		log.Printf("No pods found in namespace %s", namespace)
	} else if err != nil {
		panic(err.Error())
	}
	return podList
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, app string) *appsv1.Deployment {
	deployment := fmt.Sprintf("%s-deployment", app)
	deploymentData, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Deployment %s in namespace %s not found\n", deployment, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting deployment %s in namespace %s: %v\n",
			deployment, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("Deployment %s has %d ready pods\n", deployment, deploymentData.Status.ReadyReplicas)
	}
	return deploymentData
}

func getPodData(data []v1.Pod) []interface{} {
	var mapped []interface{}

	for _, pod := range data {
		if pod.DeletionTimestamp == nil {
			mapped = append(mapped, map[string]interface{}{
				"name":    pod.Name,
				"running": pod.Status.Phase == "Running",
			})
		}
	}

	return mapped
}
