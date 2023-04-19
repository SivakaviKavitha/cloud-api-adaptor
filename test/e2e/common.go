package e2e

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// newPod returns a new Pod object.
func newPod(namespace string, name string, containerName string, runtimeclass string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.PodSpec{
			Containers:       []corev1.Container{{Name: containerName, Image: "nginx"}},
			DNSPolicy:        "ClusterFirst",
			RestartPolicy:    "Never",
			RuntimeClassName: &runtimeclass,
		},
	}
}
func newDeployment(namespace string, deploymentname string, runtimeclass string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentname, Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentname,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "publicpod",
					Namespace: namespace,
					Labels: map[string]string{
						"app": deploymentname,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "busybox",
							Ports: []corev1.ContainerPort{
								{
									HostPort:      80,
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

}
func newService(namespace string, deploymentname string, servicename string, ipaddress string) *corev1.Service {

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: servicename, Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 80, Protocol: "TCP", TargetPort: intstr.IntOrString{IntVal: 80}}},
			Selector: map[string]string{"app": deploymentname},
			// ExternalIPs: []string{ipaddress},
		},
	}
}

// CloudAssert defines assertions to perform on the cloud provider.
type CloudAssert interface {
	HasPodVM(t *testing.T, id string) // Assert there is a PodVM with `id`.
}
