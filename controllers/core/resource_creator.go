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

func getServicePort(bookstore *customcorev1.Bookstore) (corev1.ServiceType, corev1.ServicePort) {

	var servicePort corev1.ServicePort
	var serviceType corev1.ServiceType
	var containerPort int32 = 8081

	if bookstore.Spec.ServiceType == customcorev1.NodePort {
		serviceType = corev1.ServiceTypeNodePort
		servicePort = corev1.ServicePort{
			Port:       containerPort,
			TargetPort: intstr.FromInt(int(containerPort)),
			NodePort:   bookstore.Spec.Port,
		}
	} else if bookstore.Spec.ServiceType == customcorev1.ClusterIP {
		serviceType = corev1.ServiceTypeClusterIP
		servicePort = corev1.ServicePort{
			Protocol:   corev1.ProtocolTCP,
			Port:       containerPort,
			TargetPort: intstr.FromInt(int(containerPort)),
		}
	} else {
		serviceType = corev1.ServiceTypeLoadBalancer
		servicePort = corev1.ServicePort{
			Protocol:   corev1.ProtocolTCP,
			Port:       containerPort,
			TargetPort: intstr.FromInt(int(containerPort)),
		}
	}
	return serviceType, servicePort
}

func newService(bookstore *customcorev1.Bookstore) *corev1.Service {

	labels := map[string]string{
		"app":        "bookstore-api-server",
		"controller": bookstore.Name,
	}

	serviceType, servicePort := getServicePort(bookstore)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookstore.Spec.Name + "-service",
			Namespace: bookstore.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bookstore, customcorev1.GroupVersion.WithKind("Bookstore")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Selector: labels,
			Ports: []corev1.ServicePort{
				servicePort,
			},
		},
	}

}
