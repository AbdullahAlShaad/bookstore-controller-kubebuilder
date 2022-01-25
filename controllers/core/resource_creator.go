package core

import (
	customcorev1 "github.com/Shaad7/bookstore-controller-kubebuilder/apis/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(bookstore *customcorev1.Bookstore) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "bookstore-api-server",
		"controller": bookstore.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookstore.Spec.Name + "-deployment",
			Namespace: bookstore.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bookstore, customcorev1.GroupVersion.WithKind("Bookstore")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: bookstore.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  bookstore.Spec.ImageName,
							Image: "shaad7/bookstore-api-server:latest",
						},
					},
				},
			},
		},
	}
}

func newService(bookstore *customcorev1.Bookstore) *corev1.Service {
	labels := map[string]string{
		"app":        "bookstore-api-server",
		"controller": bookstore.Name,
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "k8s.io/api/core/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: bookstore.Spec.Name + "-service",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bookstore, customcorev1.GroupVersion.WithKind("Bookstore")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labels,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port:       bookstore.Spec.ContainerPort,
					TargetPort: intstr.FromInt(int(bookstore.Spec.ContainerPort)),
					NodePort:   30000,
				},
			},
		},
	}

}
