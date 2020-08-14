/*


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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LoadServiceSpec defines the desired state of LoadService
type LoadServiceSpec struct {
	Namespace      string          `json:"namespace"`
	DeploymentName string          `json:"deployment_name"`
	Delay          metav1.Duration `json:"delay"`
}

// LoadServiceStatus defines the observed state of LoadService
type LoadServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LoadService is the Schema for the loadservices API
type LoadService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadServiceSpec   `json:"spec,omitempty"`
	Status LoadServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadServiceList contains a list of LoadService
type LoadServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadService{}, &LoadServiceList{})
}
