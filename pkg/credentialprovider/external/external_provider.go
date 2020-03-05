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
	"fmt"
	"os"

	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/credentialprovider"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/pkg/apis/registrycredential/v1alpha1"
)

var (
	scheme *runtime.Scheme
	codecs serializer.CodecFactory
)

func init() {
	scheme = runtime.NewScheme()
	v1alpha1.AddToScheme(scheme)
	codecs = serializer.NewCodecFactory(scheme)
}

func ConvertStringToRegistryCredentialPluginRequest(pluginRequest string) (*v1alpha1.RegistryCredentialPluginRequest, error) {

	obj, gvk, err := codecs.UniversalDecoder(schema.GroupVersion{Group: "registrycredentials.k8s.io", Version: "v1alpha1"}).Decode([]byte(pluginRequest), nil, nil)
	if err != nil {
		return nil, err
	}

	if gvk.Kind != "RegistryCredentialPluginRequest" {
		return nil, fmt.Errorf("failed to decode %q (missing Kind)", gvk.Kind)
	}
	config, err := scheme.ConvertToVersion(obj, v1alpha1.SchemeGroupVersion)
	if err != nil {
		return nil, err
	}
	if registryCredentialPluginRequest, ok := config.(*v1alpha1.RegistryCredentialPluginRequest); ok {
		return registryCredentialPluginRequest, nil
	}
	return nil, fmt.Errorf("unable to convert %T to *RegistryCredentialPluginRequest", config)
}

func ConvertDockerConfigToRegistryCredentialPluginResponse(cfg credentialprovider.DockerConfig) (*v1alpha1.RegistryCredentialPluginResponse, error) {
	registryCredentialResponse := &v1alpha1.RegistryCredentialPluginResponse{}

	if cfg == nil {
		klog.Fatalf("DockerConfig value is  nil")
		return nil, nil
	}
	// expect only one entry
	if len(cfg) > 1 {
		klog.Fatalf("There should ne only one entry in DockerConfig")
		os.Exit(1)
	}

	for _, v := range cfg {
		registryCredentialResponse.Username = &v.Username
		klog.Info("registryCredentialResponse.Username is ", *registryCredentialResponse.Username)
		registryCredentialResponse.Password = &v.Password
	}
	return registryCredentialResponse, nil
}

func newRegistryCredentialJSONEncoder(targetVersion schema.GroupVersion) (runtime.Encoder, error) {
	mediaType := "application/json"
	info, ok := runtime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, fmt.Errorf("unsupported media type %q", mediaType)
	}
	return codecs.EncoderForVersion(info.Serializer, targetVersion), nil
}

func ProvideCreds(pluginRequest string, provider credentialprovider.DockerConfigProvider) (string, error) {
	klog.V(4).Infof("RegistryCredentialPluginRequest from stdIn ", pluginRequest)
	registryCredentailPluginRequest, err := ConvertStringToRegistryCredentialPluginRequest(pluginRequest)

	if err != nil {
		klog.Fatalf("Failed to read the request obj ", err)
		return "Failed to read request the obj", err
	}
	klog.V(4).Infof("Credentials requested for ECR Image ", registryCredentailPluginRequest.Image)
	img := registryCredentailPluginRequest.Image
	cfg := provider.Provide(img)
	klog.V(4).Infof("encoding to RegistryPluginResponse type")
	registryCredentialPluginResponse, err := ConvertDockerConfigToRegistryCredentialPluginResponse(cfg)
	// add additional fields
	encoder, err := newRegistryCredentialJSONEncoder(v1alpha1.SchemeGroupVersion)
	if err != nil {
		klog.Fatalf("Failed to set up Json encoder: %v", err)
		return "Failed to set up Json encoder", err
	}
	data, err := runtime.Encode(encoder, registryCredentialPluginResponse)
	if err != nil {
		klog.Fatalf("error %v", data)
		return "could not encode", err
	}
	return string(data), nil
}
