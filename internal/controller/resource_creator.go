package controller

import (
	controllerv1 "github.com/Neaj-Morshad-101/crd-kubebuilder-controller/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(kluster *controllerv1.Kluster, deploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: kluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{

				*metav1.NewControllerRef(kluster, controllerv1.GroupVersion.WithKind("Kluster")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: kluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "my-app" + "-" + deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "my-app" + "-" + deploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-app",
							Image: kluster.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: kluster.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func newService(kluster *controllerv1.Kluster, name string, dep_name string) *corev1.Service {
	labels := map[string]string{
		"app": "my-app" + "-" + dep_name,
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: kluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(kluster, controllerv1.GroupVersion.WithKind("Kluster")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       kluster.Spec.Container.Port,
					TargetPort: intstr.FromInt(int(kluster.Spec.Container.Port)),
				},
			},
		},
	}
}
