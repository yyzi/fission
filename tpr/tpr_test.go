/*
Copyright 2016 The Fission Authors.

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

package tpr

import (
	"log"
	"testing"
	"time"

	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

func panicIf(err error) {
	if err != nil {
		log.Panicf("err: %v", err)
	}
}

func functionTests(tprClient *rest.RESTClient) {
	// sample function object
	function := &Function{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Function",
			APIVersion: "fission.io/v1",
		},
		Metadata: api.ObjectMeta{
			Name: "hello",
		},
		Spec: FunctionSpec{
			Source: Package{},
			Deployment: Package{
				Literal: []byte("hi"),
			},
			EnvironmentUid: "xxx",
		},
	}

	// Test function CRUD
	fi := MakeFunctionInterface(tprClient, api.NamespaceDefault)

	// cleanup from old crashed tests, ignore errors
	fi.Delete(function.Metadata.Name, nil)

	// create
	f, err := fi.Create(function)
	panicIf(err)
	if f.Metadata.Name != function.Metadata.Name {
		log.Panicf("Bad result from create: %v", f)
	}

	// read
	f, err = fi.Get(function.Metadata.Name)
	panicIf(err)
	if len(f.Spec.Deployment.Literal) != len(function.Spec.Deployment.Literal) {
		log.Panicf("Bad result from Get: %v", f)
	}

	// update
	function.Spec.EnvironmentUid = "yyy"
	f, err = fi.Update(function)
	panicIf(err)

	// list
	fl, err := fi.List(api.ListOptions{})
	panicIf(err)
	if len(fl.Items) != 1 {
		log.Panicf("wrong count from list: %v", fl)
	}
	if fl.Items[0].Spec.EnvironmentUid != function.Spec.EnvironmentUid {
		log.Panicf("bad object from list: %v", fl.Items[0])
	}

	// delete
	err = fi.Delete(f.Metadata.Name, nil)
	panicIf(err)

	// start a watch
	wi, err := fi.Watch(api.ListOptions{})
	panicIf(err)

	start := time.Now()
	f, err = fi.Create(function)
	panicIf(err)
	defer fi.Delete(f.Metadata.Name, nil)

	recvd := false
	select {
	case <-time.NewTimer(1 * time.Second).C:
		if !recvd {
			log.Panicf("Didn't get watch event")
		}
	case ev := <-wi.ResultChan():
		wf, ok := ev.Object.(*Function)
		if !ok {
			log.Panicf("Can't cast to Function")
		}
		if wf.Spec.EnvironmentUid != function.Spec.EnvironmentUid {
			log.Panicf("Bad object from watch: %#v", wf)
		}
		log.Printf("watch event took %v", time.Now().Sub(start))
		recvd = true
	}

}

func TestTpr(t *testing.T) {

	// TODO skip test if we're not configured with a cluster

	// Create the client config. Use kubeconfig at ./config
	config, err := clientcmd.BuildConfigFromFlags("", "config")
	panicIf(err)

	// kube client
	clientset, err := kubernetes.NewForConfig(config)
	panicIf(err)

	// init our types
	err = EnsureFissionTPRs(clientset)
	panicIf(err)

	// The "right" way is to poll for the tpr api endpoint to stop
	// 404'ing, or something like that.
	time.Sleep(10 * time.Second)

	// rest client with knowledge about our tpr types
	tprClient, err := getTprClient(config)
	panicIf(err)

	functionTests(tprClient)
}
