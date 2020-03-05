/*
Copyright 2014 The Kubernetes Authors.

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

package external

import (
	"testing"

	"k8s.io/kubernetes/pkg/credentialprovider/external"
)

const image = "602401143452.dkr.ecr.us-west-2.amazonaws.com/eks/kube-scheduler"

func TestConvertStringToRegistryCredentialPluginRequest(t *testing.T) {

	pluginRequest := `{
			"kind" : "RegistryCredentialPluginRequest" , 
			"apiVersion": "registrycredential.k8s.io/v1alpha1", 
			"image" : "602401143452.dkr.ecr.us-west-2.amazonaws.com/eks/kube-scheduler"
	}`
	registryCredentialPluginRequest, err := external.ConvertStringToRegistryCredentialPluginRequest(pluginRequest)

	if err != nil {
		t.Errorf("unexpected error while decoding input  %v", pluginRequest)
	}
	expectedImage := image
	if registryCredentialPluginRequest.Image != expectedImage {
		t.Errorf("expected output %v doesn't match with actual out  %v",
			expectedImage, registryCredentialPluginRequest.Image)
	}

}
