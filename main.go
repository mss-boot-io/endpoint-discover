/*
Copyright 2016 The Kubernetes Authors.
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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//clusterUrl, err := url.Parse(os.Getenv("cluster_url"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//config := &rest.Config{
	//	Host:    clusterUrl.Host,
	//	APIPath: clusterUrl.Path,
	//	TLSClientConfig: rest.TLSClientConfig{
	//		Insecure: true,
	//	},
	//	BearerToken: os.Getenv("token"),
	//}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	cm := &corev1.ConfigMap{}
	cm.Name = os.Getenv("configmap_name")
	cm.Namespace = os.Getenv("namespace")
	if cm.Namespace == "" {
		cm.Namespace = "default"
	}

	serviceList, err := clientset.CoreV1().Services(cm.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d services in the cluster\n", len(serviceList.Items))

	endpoints := make(map[string][]Endpoint)
	protocols := strings.Split(os.Getenv("protocols"), ",")
	for i := range protocols {
		protocols[i] = strings.Trim(protocols[i], " ")
		endpoints[protocols[i]] = make([]Endpoint, 0)
	}
	for i := range serviceList.Items {
		for j := range serviceList.Items[i].Spec.Ports {
			for n := range protocols {
				if strings.Index(serviceList.Items[i].Spec.Ports[j].Name, protocols[n]) > -1 {
					var port int
					switch serviceList.Items[i].Spec.Ports[j].TargetPort.String() {
					case "http":
						port = 80
					case "https":
						port = 443
					default:
						port = serviceList.Items[i].Spec.Ports[j].TargetPort.IntValue()
					}
					endpoints[protocols[n]] = append(endpoints[protocols[n]], Endpoint{
						Name:     serviceList.Items[i].Name,
						Endpoint: fmt.Sprintf("%s.%s:%d", serviceList.Items[i].Name, cm.Namespace, port),
					})
				}
			}

		}
	}
	out, err := yaml.Marshal(endpoints)
	if err != nil {
		panic(err)
	}
	cm.Data = map[string]string{
		os.Getenv("config_name"): string(out),
	}

	_, err = clientset.CoreV1().ConfigMaps(cm.Namespace).Get(context.TODO(), cm.Name, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		cm, err = clientset.CoreV1().ConfigMaps(cm.Namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	} else {
		cm, err = clientset.CoreV1().ConfigMaps(cm.Namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	}
	//查询是否存在
	if err != nil {
		panic(err)
	}
}

type Endpoint struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}
