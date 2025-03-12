package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Create in-cluster Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create in-cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	namespace := "default"
	image := "nginx:latest" // Change this to your desired image
	jobName := "image-pull-test"

	// Define job spec
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test-container",
							Image: image,
							Command: []string{"sleep", "10"}, // Short-running container
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	// Record start time
	startTime := time.Now()

	// Create job
	_, err = clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}

	// Wait for pod to start running
	for {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "job-name=" + jobName,
		})
		if err != nil {
			log.Fatalf("Failed to list pods: %v", err)
		}

		for _, pod := range pods.Items {
			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Running != nil || status.State.Terminated != nil {
					duration := time.Since(startTime)
					fmt.Printf("Image pull time: %v\n", duration)
					return
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

