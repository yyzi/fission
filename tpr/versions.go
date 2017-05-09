/*
Copyright 2017 The Fission Authors.

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
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/rest"
)

type (
	versionClient struct {
		client    *rest.RESTClient
		namespace string
	}
)

func MakeVersionClient(tprClient *rest.RESTClient, namespace string) *versionClient {
	return &versionClient{
		client:    tprClient,
		namespace: namespace,
	}
}

func (vc *versionClient) objName(resourceType, objName string) string {
	return resourceType + "__" + objName
}

func (vc *versionClient) createDefault(resourceType, resourceName, version string) error {
	obj := &Resourceversion{
		Metadata: api.ObjectMeta{
			Name: vc.objName(resourceType, resourceName),
		},
		ResourceType: resourceType,
		ResourceName: resourceName,
		Version:      version,
	}
	var result Resourceversion
	err := vc.client.Post().
		Resource("resourceversions").
		Namespace(vc.namespace).
		Body(obj).
		Do().Into(&result)
	return err
}

func (vc *versionClient) updateDefault(resourceType, resourceName, version string) error {
	obj := &Resourceversion{
		Metadata: api.ObjectMeta{
			Name: vc.objName(resourceType, resourceName),
		},
		ResourceType: resourceType,
		ResourceName: resourceName,
		Version:      version,
	}
	var result Resourceversion
	err := vc.client.Put().
		Resource("resourceversions").
		Namespace(vc.namespace).
		Body(obj).
		Do().Into(&result)
	return err
}

func (vc *versionClient) getDefault(resourceType, resourceName string) (string, error) {
	name := vc.objName(resourceType, resourceName)
	var result Resourceversion
	err := vc.client.Get().
		Resource("resourceversions").
		Namespace(vc.namespace).
		Name(name).
		Do().Into(&result)
	if err != nil {
		return "", err
	}
	return result.Version, nil
}

func (vc *versionClient) deleteDefault(resourceType, resourceName string) error {
	return vc.client.Delete().
		Namespace(vc.namespace).
		Resource("resourceversions").
		Name(vc.objName(resourceType, resourceName)).
		Body(&api.DeleteOptions{}).
		Do().
		Error()
}
