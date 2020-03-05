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
// The aws-credentials-provider binary is responsible for providing
// ecr credentials for pulling down the docker image

package app

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog"
	credentials "k8s.io/kubernetes/pkg/credentialprovider/aws"
	"k8s.io/kubernetes/pkg/credentialprovider/external"
)

const (
	// aws credentials provider binary  name
	awsecrcredentialsprovider = "ecr-creds"
	getcredentials            = "get-credentials"
)

func NewEcrCredentialsProviderCommand() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}

var rootCmd = &cobra.Command{
	Use:  awsecrcredentialsprovider,
	Long: "AWS ECR Credentials provider for fetching ecr creds for kubelet",
	Run: func(cmd *cobra.Command, args []string) {
		verflag.PrintAndExitIfRequested()
	},
}

var tokenCmd = &cobra.Command{
	Use:  getcredentials,
	Long: "AWS ECR Credentials provider for fetching ecr creds for kubelet",
	Run: func(cmd *cobra.Command, args []string) {
		verflag.PrintAndExitIfRequested()
		ecrProvider := credentials.NewEcrProvider()
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			klog.Fatalf("error while reading from stdin %v", err)
			os.Exit(1)
		}
		klog.V(4).Infof("Input from stdin: %s", string(buf))
		data, err := external.ProvideCreds(string(buf), ecrProvider)
		if err != nil {
			klog.Fatalf("could not get credentials %v", err)
		}
		os.Stdout.Write([]byte(data))
	},
}
