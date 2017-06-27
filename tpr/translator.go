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
	//"github.com/satori/go.uuid"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/types"

	"github.com/fission/fission"
)

/*

Throughout this file, instances of TPR types have the "t_" prefix.
Vars with the f_ prefix are fission types.  We don't follow that
convention anywhere else, just in this file.

*/

func OwnerReferenceFromTPR(t_o *api.OwnerReference) *fission.OwnerReference {
	o := &fission.OwnerReference{
		APIVersion: t_o.APIVersion,
		Kind:       t_o.Kind,
		Name:       t_o.Name,
		UID:        string(t_o.UID),
	}
	b := *t_o.Controller
	o.Controller = &b
	return o
}

func OwnerReferenceToTPR(o *fission.OwnerReference) *api.OwnerReference {
	t_o := &api.OwnerReference{
		APIVersion: o.APIVersion,
		Kind:       o.Kind,
		Name:       o.Name,
		UID:        types.UID(o.UID),
	}

	b := *t_o.Controller
	t_o.Controller = &b
	return t_o
}

// MetadataFromTPR translates a Kubernetes API metadata object to a
// Fission API metadata object.
func MetadataFromTPR(t_metadata *api.ObjectMeta) *fission.Metadata {
	m := fission.Metadata{
		Name:              t_metadata.Name,
		GenerateName:      t_metadata.GenerateName,
		Namespace:         t_metadata.Namespace,
		SelfLink:          t_metadata.SelfLink,
		UID:               string(t_metadata.UID),
		ResourceVersion:   t_metadata.ResourceVersion,
		Generation:        t_metadata.Generation,
		CreationTimestamp: fission.Time(t_metadata.CreationTimestamp),
		Labels:            t_metadata.Labels,
		Annotations:       t_metadata.Annotations,
		Finalizers:        t_metadata.Finalizers,
		ClusterName:       t_metadata.ClusterName,
	}

	t := fission.Time(*t_metadata.DeletionTimestamp)
	m.DeletionTimestamp = &t

	s := *t_metadata.DeletionGracePeriodSeconds
	m.DeletionGracePeriodSeconds = &s

	os := make([]fission.OwnerReference, 0)
	for _, o := range t_metadata.OwnerReferences {
		os = append(os, *OwnerReferenceFromTPR(&o))
	}
	m.OwnerReferences = os

	return &m
}

// MetadataToTPR translates a Fission API metadata object to a
// Kubernetes API metadata object.
func MetadataToTPR(m *fission.Metadata) *api.ObjectMeta {
	t_m := api.ObjectMeta{
		Name:              m.Name,
		GenerateName:      m.GenerateName,
		Namespace:         m.Namespace,
		SelfLink:          m.SelfLink,
		UID:               types.UID(m.UID),
		ResourceVersion:   m.ResourceVersion,
		Generation:        m.Generation,
		CreationTimestamp: unversioned.Time(m.CreationTimestamp),
		Labels:            m.Labels,
		Annotations:       m.Annotations,
		Finalizers:        m.Finalizers,
		ClusterName:       m.ClusterName,
	}

	t := unversioned.Time(*m.DeletionTimestamp)
	t_m.DeletionTimestamp = &t

	s := *m.DeletionGracePeriodSeconds
	t_m.DeletionGracePeriodSeconds = &s

	os := make([]api.OwnerReference, 0)
	for _, o := range m.OwnerReferences {
		os = append(os, *OwnerReferenceToTPR(&o))
	}
	t_m.OwnerReferences = os

	return &t_m
}

func FunctionFromTPR(t_function *Function) *fission.Function {
	return &fission.Function{
		Metadata: *MetadataFromTPR(&t_function.Metadata),
		Spec:     t_function.Spec,
	}
}

func FunctionToTPR(f_function *fission.Function) *Function {
	return &Function{
		Metadata: *MetadataToTPR(&f_function.Metadata),
		Spec:     f_function.Spec,
	}
}

func FunctionListFromTPR(t_functions *FunctionList) []fission.Function {
	fs := make([]fission.Function, len(t_functions.Items))
	for _, t := range t_functions.Items {
		fs = append(fs, *FunctionFromTPR(&t))
	}
	return fs
}

func EnvironmentFromTPR(t_environment *Environment) *fission.Environment {
	return &fission.Environment{
		Metadata: *MetadataFromTPR(&t_environment.Metadata),
		Spec:     t_environment.Spec,
	}
}

func EnvironmentToTPR(f_environment *fission.Environment) *Environment {
	return &Environment{
		Metadata: *MetadataToTPR(&f_environment.Metadata),
		Spec:     f_environment.Spec,
	}
}

func EnvironmentListFromTPR(t_environments *EnvironmentList) []fission.Environment {
	fs := make([]fission.Environment, len(t_environments.Items))
	for _, t := range t_environments.Items {
		fs = append(fs, *EnvironmentFromTPR(&t))
	}
	return fs
}

func HttpTriggerFromTPR(t_httptrigger *Httptrigger) *fission.HTTPTrigger {
	return &fission.HTTPTrigger{
		Metadata: *MetadataFromTPR(&t_httptrigger.Metadata),
		Spec:     t_httptrigger.Spec,
	}
}

func HttpTriggerToTPR(f_httptrigger *fission.HTTPTrigger) *Httptrigger {
	return &Httptrigger{
		Metadata: *MetadataToTPR(&f_httptrigger.Metadata),
		Spec:     f_httptrigger.Spec,
	}
}

func HttpTriggerListFromTPR(t_httptriggers *HttptriggerList) []fission.HTTPTrigger {
	fs := make([]fission.HTTPTrigger, len(t_httptriggers.Items))
	for _, t := range t_httptriggers.Items {
		fs = append(fs, *HttpTriggerFromTPR(&t))
	}
	return fs
}

func KubernetesWatchTriggerFromTPR(t_kuberneteswatchtrigger *Kuberneteswatchtrigger) *fission.KubernetesWatchTrigger {
	return &fission.KubernetesWatchTrigger{
		Metadata: *MetadataFromTPR(&t_kuberneteswatchtrigger.Metadata),
		Spec:     t_kuberneteswatchtrigger.Spec,
	}
}

func KubernetesWatchTriggerToTPR(f_kuberneteswatchtrigger *fission.KubernetesWatchTrigger) *Kuberneteswatchtrigger {
	return &Kuberneteswatchtrigger{
		Metadata: *MetadataToTPR(&f_kuberneteswatchtrigger.Metadata),
		Spec:     f_kuberneteswatchtrigger.Spec,
	}
}

func KubernetesWatchTriggerListFromTPR(t_kuberneteswatchtriggers *KuberneteswatchtriggerList) []fission.KubernetesWatchTrigger {
	fs := make([]fission.KubernetesWatchTrigger, len(t_kuberneteswatchtriggers.Items))
	for _, t := range t_kuberneteswatchtriggers.Items {
		fs = append(fs, *KubernetesWatchTriggerFromTPR(&t))
	}
	return fs
}
