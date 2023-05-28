package ingress

import (
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
	hyperv1 "github.com/openshift/hypershift/api/v1beta1"
	"github.com/openshift/hypershift/control-plane-operator/hostedclusterconfigoperator/controllers/resources/manifests"
)

func TestReconcileDefaultIngressController(t *testing.T) {
	fakeIngressDomain := "example.com"
	fakeInputReplicas := int32(3)
	testsCases := []struct {
		name                      string
		inputIngressController    *operatorv1.IngressController
		inputIngressDomain        string
		inputPlatformType         hyperv1.PlatformType
		inputReplicas             int32
		inputIsIBMCloudUPI        bool
		inputIsPrivate            bool
		expectedIngressController *operatorv1.IngressController
	}{
		{
			name:                   "IBM Cloud UPI uses Nodeport publishing strategy",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.IBMCloudPlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     true,
			inputIsPrivate:         false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type: operatorv1.NodePortServiceStrategyType,
						NodePort: &operatorv1.NodePortStrategy{
							Protocol: operatorv1.TCPProtocol,
						},
					},
					NodePlacement: &operatorv1.NodePlacement{
						Tolerations: []corev1.Toleration{
							{
								Key:   "dedicated",
								Value: "edge",
							},
						},
					},
				},
			},
		},
		{
			name:                   "IBM Cloud Non-UPI uses LoadBalancer publishing strategy",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.IBMCloudPlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type: operatorv1.LoadBalancerServiceStrategyType,
						LoadBalancer: &operatorv1.LoadBalancerStrategy{
							Scope: operatorv1.ExternalLoadBalancer,
						},
					},
					NodePlacement: &operatorv1.NodePlacement{
						Tolerations: []corev1.Toleration{
							{
								Key:   "dedicated",
								Value: "edge",
							},
						},
					},
				},
			},
		},
		{
			name:                   "Kubevirt uses NodePort publishing strategy",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.KubevirtPlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type: operatorv1.NodePortServiceStrategyType,
					},
					DefaultCertificate: &corev1.LocalObjectReference{
						Name: manifests.IngressDefaultIngressControllerCert().Name,
					},
				},
			},
		},
		{
			name:                   "None Platform uses HostNetwork publishing strategy",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.NonePlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type: operatorv1.HostNetworkStrategyType,
					},
					DefaultCertificate: &corev1.LocalObjectReference{
						Name: manifests.IngressDefaultIngressControllerCert().Name,
					},
				},
			},
		},
		{
			name:                   "AWS uses Loadbalancer publishing strategy",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.AWSPlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type: operatorv1.LoadBalancerServiceStrategyType,
					},
					DefaultCertificate: &corev1.LocalObjectReference{
						Name: manifests.IngressDefaultIngressControllerCert().Name,
					},
				},
			},
		},
		{
			name:                   "Private Publishing Strategy on IBM Cloud",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputPlatformType:      hyperv1.IBMCloudPlatform,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         true,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type:    operatorv1.PrivateStrategyType,
						Private: &operatorv1.PrivateStrategy{},
					},
					NodePlacement: &operatorv1.NodePlacement{
						Tolerations: []corev1.Toleration{
							{
								Key:   "dedicated",
								Value: "edge",
							},
						},
					},
				},
			},
		},
		{
			name:                   "Private Publishing Strategy on other Platforms",
			inputIngressController: manifests.IngressDefaultIngressController(),
			inputIngressDomain:     fakeIngressDomain,
			inputReplicas:          fakeInputReplicas,
			inputIsIBMCloudUPI:     false,
			inputIsPrivate:         true,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: manifests.IngressDefaultIngressController().ObjectMeta,
				Spec: operatorv1.IngressControllerSpec{
					Domain:   fakeIngressDomain,
					Replicas: &fakeInputReplicas,
					EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
						Type:    operatorv1.PrivateStrategyType,
						Private: &operatorv1.PrivateStrategy{},
					},
					DefaultCertificate: &corev1.LocalObjectReference{
						Name: manifests.IngressDefaultIngressControllerCert().Name,
					},
				},
			},
		},
		{
			name: "Existing ingress controller",
			inputIngressController: func() *operatorv1.IngressController {
				ic := manifests.IngressDefaultIngressController()
				ic.ResourceVersion = "1"
				return ic
			}(),
			inputIngressDomain: fakeIngressDomain,
			inputReplicas:      fakeInputReplicas,
			inputIsIBMCloudUPI: false,
			inputIsPrivate:     false,
			expectedIngressController: &operatorv1.IngressController{
				ObjectMeta: func() metav1.ObjectMeta {
					m := manifests.IngressDefaultIngressController().ObjectMeta
					m.ResourceVersion = "1"
					return m
				}(),
				Spec: operatorv1.IngressControllerSpec{},
			},
		},
	}
	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			err := ReconcileDefaultIngressController(tc.inputIngressController, tc.inputIngressDomain, tc.inputPlatformType, tc.inputReplicas, tc.inputIsIBMCloudUPI, tc.inputIsPrivate)
			g.Expect(err).To(BeNil())
			g.Expect(tc.inputIngressController).To(BeEquivalentTo(tc.expectedIngressController))
		})
	}
}

func TestReconcileDefaultIngressPassthroughService(t *testing.T) {
	const (
		infraId        = "12345678"
		nodePort       = 8080
		secureNodePort = 6443
	)

	hcp := &hyperv1.HostedControlPlane{
		Spec: hyperv1.HostedControlPlaneSpec{
			InfraID: infraId,
		},
	}

	for _, tc := range []struct {
		name          string
		defSvc        func() *corev1.Service
		shouldSucceed bool
	}{
		{
			name: "valid use case",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port:     443,
						NodePort: secureNodePort,
					},
					{
						Port:     80,
						NodePort: nodePort,
					},
				}
				return defSvs
			},
			shouldSucceed: true,
		},
		{
			name: "valid use case with additional ports",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port:     22,
						NodePort: 2222,
					},
					{
						Port:     443,
						NodePort: secureNodePort,
					},
					{
						Port:     80,
						NodePort: nodePort,
					},
					{
						Port:     9999,
						NodePort: 9999,
					},
				}
				return defSvs
			},
			shouldSucceed: true,
		},
		{
			name: "error: missing secure node port",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port: 443,
					},
					{
						Port:     80,
						NodePort: nodePort,
					},
				}
				return defSvs
			},
			shouldSucceed: false,
		},
		{
			name: "error: missing node port",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port:     443,
						NodePort: secureNodePort,
					},
					{
						Port: 80,
					},
				}
				return defSvs
			},
			shouldSucceed: false,
		},
		{
			name: "error: missing secure port",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port:     80,
						NodePort: nodePort,
					},
				}
				return defSvs
			},
			shouldSucceed: false,
		},
		{
			name: "error: missing port",
			defSvc: func() *corev1.Service {
				defSvs := manifests.IngressDefaultIngressNodePortService()
				defSvs.Spec.Ports = []corev1.ServicePort{
					{
						Port:     443,
						NodePort: secureNodePort,
					},
				}
				return defSvs
			},
			shouldSucceed: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			resSvc := &corev1.Service{}
			err := ReconcileDefaultIngressPassthroughService(resSvc, tc.defSvc(), hcp)

			if tc.shouldSucceed {
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(resSvc.Spec.Selector).To(HaveKeyWithValue("kubevirt.io", "virt-launcher"))
				g.Expect(resSvc.Spec.Selector).To(HaveKeyWithValue(hyperv1.InfraIDLabel, infraId))
				g.Expect(resSvc.Spec.Type).To(Equal(corev1.ServiceTypeClusterIP))
				g.Expect(resSvc.Labels).To(HaveKeyWithValue(hyperv1.InfraIDLabel, infraId))
				g.Expect(resSvc.Spec.Ports).To(HaveLen(2))
				g.Expect(resSvc.Spec.Ports).To(ContainElements(
					corev1.ServicePort{
						Port:       80,
						Name:       "http-80",
						Protocol:   corev1.ProtocolTCP,
						TargetPort: intstr.FromInt(nodePort),
					},
					corev1.ServicePort{
						Port:       443,
						Name:       "https-443",
						Protocol:   corev1.ProtocolTCP,
						TargetPort: intstr.FromInt(secureNodePort),
					}),
				)
			} else {
				g.Expect(err).To(HaveOccurred())
			}

		})
	}
}
