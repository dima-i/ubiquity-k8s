/**
 * Copyright 2017 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"

	"fmt"
	"os/user"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	master = flag.String("master", "", "Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.")
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Failed to get the active user: %v", err))
	}
	homedir := usr.HomeDir

	kubeconfig := flag.String("kubeconfig", fmt.Sprintf("%s/.kube/config", homedir), "Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.")

	var config *rest.Config
	if *kubeconfig != "" {

		config, err = clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to create config: %v", err))
	}
	clientset, err := kubernetes.NewForConfig(config)

	volumes, err := clientset.Core().PersistentVolumes().List(v1.ListOptions{})
	if err != nil {
		fmt.Printf("error %#v", err)
	}
	for _, volume := range volumes.Items {

		if contains(volume.Labels, "static") {
			fmt.Printf("%#v\n", volume.Name)
			clientset.Core().PersistentVolumes().Delete(volume.Name, nil)
		}
	}

}

func contains(m map[string]string, s string) bool {
	for _, value := range m {
		if value == s {
			return true
		}
	}
	return false
}
