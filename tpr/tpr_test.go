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

	// assert that we get a watch event for the new function
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

func environmentTests(tprClient *rest.RESTClient) {
	// sample environment object
	environment := &Environment{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Environment",
			APIVersion: "fission.io/v1",
		},
		Metadata: api.ObjectMeta{
			Name: "hello",
		},
		Spec: EnvironmentSpec{
			Runtime: Runtime{
				Image: "xxx",
			},
			Builder: Builder{
				Image:   "yyy",
				Command: "zzz",
			},
			FilenameExtensions: []string{"js"},
		},
	}

	// Test environment CRUD
	ei := MakeEnvironmentInterface(tprClient, api.NamespaceDefault)

	// cleanup from old crashed tests, ignore errors
	ei.Delete(environment.Metadata.Name, nil)

	// create
	e, err := ei.Create(environment)
	panicIf(err)
	if e.Metadata.Name != environment.Metadata.Name {
		log.Panicf("Bad result from create: %v", e)
	}

	// read
	e, err = ei.Get(environment.Metadata.Name)
	panicIf(err)
	if len(e.Spec.Runtime.Image) != len(environment.Spec.Runtime.Image) {
		log.Panicf("Bad result from Get: %#v", e)
	}

	// update
	environment.Spec.Runtime.Image = "www"
	e, err = ei.Update(environment)
	panicIf(err)

	// list
	el, err := ei.List(api.ListOptions{})
	panicIf(err)
	if len(el.Items) != 1 {
		log.Panicf("wrong count from list: %v", el)
	}
	if el.Items[0].Spec.Runtime.Image != environment.Spec.Runtime.Image {
		log.Panicf("bad object from list: %v", el.Items[0])
	}

	// delete
	err = ei.Delete(e.Metadata.Name, nil)
	panicIf(err)

	// start a watch
	wi, err := ei.Watch(api.ListOptions{})
	panicIf(err)

	start := time.Now()
	e, err = ei.Create(environment)
	panicIf(err)
	defer ei.Delete(e.Metadata.Name, nil)

	// assert that we get a watch event for the new environment
	recvd := false
	select {
	case <-time.NewTimer(1 * time.Second).C:
		if !recvd {
			log.Panicf("Didn't get watch event")
		}
	case ev := <-wi.ResultChan():
		obj, ok := ev.Object.(*Environment)
		if !ok {
			log.Panicf("Can't cast to Environment")
		}
		if obj.Spec.Runtime.Image != environment.Spec.Runtime.Image {
			log.Panicf("Bad object from watch: %#v", obj)
		}
		log.Printf("watch event took %v", time.Now().Sub(start))
		recvd = true
	}

}

func httpTriggerTests(tprClient *rest.RESTClient) {
	// sample httpTrigger object
	httpTrigger := &Httptrigger{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Httptrigger",
			APIVersion: "fission.io/v1",
		},
		Metadata: api.ObjectMeta{
			Name: "hello",
		},
		Spec: HttptriggerSpec{
			RelativeURL: "/hi",
			Method:      "GET",
			FunctionReference: FunctionReference{
				Selector: map[string]string{
					"name": "hello",
					"uid":  "42",
				},
			},
		},
	}

	// Test httpTrigger CRUD
	ei := MakeHttptriggerInterface(tprClient, api.NamespaceDefault)

	// cleanup from old crashed tests, ignore errors
	ei.Delete(httpTrigger.Metadata.Name, nil)

	// create
	e, err := ei.Create(httpTrigger)
	panicIf(err)
	if e.Metadata.Name != httpTrigger.Metadata.Name {
		log.Panicf("Bad result from create: %v", e)
	}

	// read
	e, err = ei.Get(httpTrigger.Metadata.Name)
	panicIf(err)
	if len(e.Spec.Method) != len(httpTrigger.Spec.Method) {
		log.Panicf("Bad result from Get: %#v", e)
	}

	// update
	httpTrigger.Spec.Method = "POST"
	e, err = ei.Update(httpTrigger)
	panicIf(err)

	// list
	el, err := ei.List(api.ListOptions{})
	panicIf(err)
	if len(el.Items) != 1 {
		log.Panicf("wrong count from list: %v", el)
	}
	if el.Items[0].Spec.Method != httpTrigger.Spec.Method {
		log.Panicf("bad object from list: %v", el.Items[0])
	}

	// delete
	err = ei.Delete(e.Metadata.Name, nil)
	panicIf(err)

	// start a watch
	wi, err := ei.Watch(api.ListOptions{})
	panicIf(err)

	start := time.Now()
	e, err = ei.Create(httpTrigger)
	panicIf(err)
	defer ei.Delete(e.Metadata.Name, nil)

	// assert that we get a watch event for the new httpTrigger
	recvd := false
	select {
	case <-time.NewTimer(1 * time.Second).C:
		if !recvd {
			log.Panicf("Didn't get watch event")
		}
	case ev := <-wi.ResultChan():
		obj, ok := ev.Object.(*Httptrigger)
		if !ok {
			log.Panicf("Can't cast to Httptrigger")
		}
		if obj.Spec.Method != httpTrigger.Spec.Method {
			log.Panicf("Bad object from watch: %#v", obj)
		}
		log.Printf("watch event took %v", time.Now().Sub(start))
		recvd = true
	}

}

func kubernetesWatchTriggerTests(tprClient *rest.RESTClient) {
	// sample kubernetesWatchTrigger object
	kubernetesWatchTrigger := &Kuberneteswatchtrigger{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Kuberneteswatchtrigger",
			APIVersion: "fission.io/v1",
		},
		Metadata: api.ObjectMeta{
			Name: "hello",
		},
		Spec: KuberneteswatchtriggerSpec{
			Namespace: "foo",
			Type:      "pod",
			LabelSelector: map[string]string{
				"x": "y",
			},
			FunctionReference: FunctionReference{
				Selector: map[string]string{
					"name": "foo",
				},
			},
		},
	}

	// Test kubernetesWatchTrigger CRUD
	ei := MakeKuberneteswatchtriggerInterface(tprClient, api.NamespaceDefault)

	// cleanup from old crashed tests, ignore errors
	ei.Delete(kubernetesWatchTrigger.Metadata.Name, nil)

	// create
	e, err := ei.Create(kubernetesWatchTrigger)
	panicIf(err)
	if e.Metadata.Name != kubernetesWatchTrigger.Metadata.Name {
		log.Panicf("Bad result from create: %v", e)
	}

	// read
	e, err = ei.Get(kubernetesWatchTrigger.Metadata.Name)
	panicIf(err)
	if e.Spec.Type != kubernetesWatchTrigger.Spec.Type {
		log.Panicf("Bad result from Get: %#v", e)
	}

	// update
	kubernetesWatchTrigger.Spec.Type = "service"
	e, err = ei.Update(kubernetesWatchTrigger)
	panicIf(err)

	// list
	el, err := ei.List(api.ListOptions{})
	panicIf(err)
	if len(el.Items) != 1 {
		log.Panicf("wrong count from list: %v", el)
	}
	if el.Items[0].Spec.Type != kubernetesWatchTrigger.Spec.Type {
		log.Panicf("bad object from list: %v", el.Items[0])
	}

	// delete
	err = ei.Delete(e.Metadata.Name, nil)
	panicIf(err)

	// start a watch
	wi, err := ei.Watch(api.ListOptions{})
	panicIf(err)

	start := time.Now()
	e, err = ei.Create(kubernetesWatchTrigger)
	panicIf(err)
	defer ei.Delete(e.Metadata.Name, nil)

	// assert that we get a watch event for the new kubernetesWatchTrigger
	recvd := false
	select {
	case <-time.NewTimer(1 * time.Second).C:
		if !recvd {
			log.Panicf("Didn't get watch event")
		}
	case ev := <-wi.ResultChan():
		obj, ok := ev.Object.(*Kuberneteswatchtrigger)
		if !ok {
			log.Panicf("Can't cast to Kuberneteswatchtrigger")
		}
		if obj.Spec.Type != kubernetesWatchTrigger.Spec.Type {
			log.Panicf("Bad object from watch: %#v", obj)
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
	environmentTests(tprClient)
	httpTriggerTests(tprClient)
	kubernetesWatchTriggerTests(tprClient)
}
