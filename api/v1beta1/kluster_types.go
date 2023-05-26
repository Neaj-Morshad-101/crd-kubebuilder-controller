/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

////KlusterSpec defines the desired state of Kluster
//type KlusterSpec struct {
//
//	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
//	// Important: Run "make" to regenerate code after modifying this file
//
//	// Foo is an example field of Kluster. Edit kluster_types.go to remove/update
//	Foo string `json:"foo,omitempty"`
//}

type ContainerSpec struct {
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

type ServiceSpec struct {
	//+optional
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType string `json:"serviceType"`
	//+optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
	ServicePort     int32 `json:"servicePort"`
}

// KlusterSpec defines the desired state of KlusterCRD
type KlusterSpec struct {
	Replicas  *int32        `json:"replicas"`
	Container ContainerSpec `json:"container,container"`
	Service   ServiceSpec   `json:"service,omitempty"`
}

//// KlusterStatus defines the observed state of Kluster
//type KlusterStatus struct {
//	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
//	// Important: Run "make" to regenerate code after modifying this file
//}

type KlusterStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Kluster is the Schema for the klusters API
type Kluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KlusterSpec   `json:"spec,omitempty"`
	Status KlusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KlusterList contains a list of Kluster
type KlusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kluster{}, &KlusterList{})
}
