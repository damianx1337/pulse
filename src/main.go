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
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load Kubernetes configuration from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
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
							Command: []string{"sh", "-c", "echo Start; sleep 10"}, // Simulate startup
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

	var imagePullTime, startupTime time.Duration

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
				if status.State.Waiting == nil && status.State.Running != nil {
					imagePullTime = time.Since(startTime)
					fmt.Printf("Image pull time: %v\n", imagePullTime)
				}
				if status.State.Terminated != nil {
					startupTime = time.Since(startTime)
					fmt.Printf("Application startup time: %v\n", startupTime)
					return
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

