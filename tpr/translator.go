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

	"github.com/fission/fission"
)

/*

Throughout this file, instances of fission types have the "f_" prefix
in their name, and instances of TPR types have the "t_" prefix.

We don't follow that convention anywhere else, just in this file.

*/

// MetadataFromTPR translates a Kubernetes API metadata object to a
// Fission API metadata object.
func MetadataFromTPR(t_metadata *api.ObjectMeta) (*fission.Metadata, error) {
	labels := t_metadata.Labels

	name, ok := labels["fission-name"]
	if !ok {
		return nil, fission.MakeError(fission.ErrorTypeTranslation, "Missing annotation: fission.io/name")
	}

	uid, ok := labels["fission-uid"]
	if !ok {
		return nil, fission.MakeError(fission.ErrorTypeTranslation, "Missing annotation: fission.io/uid")
	}

	version, ok := labels["fission-version"]
	if !ok {
		return nil, fission.MakeError(fission.ErrorTypeTranslation, "Missing annotation: fission.io/version")
	}

	m := &fission.Metadata{
		Name:        name,
		Uid:         uid,
		Version:     version,
		Labels:      labels,
		Annotations: t_metadata.Annotations,
	}
	return m, nil
}

// MetadataToTPR translates a fission metadata object to a Kubernetes
// metadata object.  If the t_metadata parameter is provided, its
// labels and annotations are merged with the f_metadata parameter's
// labels and annotations (f_metadata overrides).
//
// We store fission identifying metadata in labels rather than
// annotations, because we need to query by name, uid, and version.
func MetadataToTPR(f_metadata *fission.Metadata, t_metadata *api.ObjectMeta) *api.ObjectMeta {
	var t_m api.ObjectMeta

	// K8s requires unique names, but we may have multiple
	// versions of the same object. So we use the f_metadata UID
	// as the TPR's name.
	t_m.Name = f_metadata.Uid

	// Annotations: preserve existing annotations; merge fission
	// metadata annotations into the tpr object.
	var annotations map[string]string
	if t_metadata != nil {
		for k, v := range t_metadata.Annotations {
			annotations[k] = v
		}
	}
	for k, v := range f_metadata.Annotations {
		annotations[k] = v
	}
	t_m.Annotations = annotations

	// Preserve existing labels, store our metadata as labels, and
	// merge fission metadata labels
	var labels map[string]string
	if t_metadata != nil {
		for k, v := range t_metadata.Labels {
			labels[k] = v
		}
	}
	labels = map[string]string{
		"fission-name":    f_metadata.Name,
		"fission-uid":     f_metadata.Uid,
		"fission-version": f_metadata.Version,
	}
	for k, v := range f_metadata.Labels {
		labels[k] = v
	}
	t_m.Labels = labels

	// Namespace
	t_m.Namespace = f_metadata.Namespace

	return t_metadata
}

func FunctionFromTPR(t_function *Function) (*fission.Function, error) {
	metadata, err := MetadataFromTPR(&t_function.Metadata)
	if err != nil {
		return nil, err
	}
	f_function := &fission.Function{
		Metadata: *metadata,
		Spec:     t_function.Spec,
	}
	return f_function, nil
}

func FunctionToTPR(f_function *fission.Function, t_metadata *api.ObjectMeta) *Function {
	t_function := &Function{
		Metadata: *MetadataToTPR(&f_function.Metadata, t_metadata),
		Spec:     f_function.Spec,
	}
	return t_function
}

func FunctionListFromTPR(t_functions *FunctionList) (*[]fission.Function, error) {
	f_functions := make([]fission.Function, len(t_functions.Items))
	for _, t_f := range t_functions.Items {
		f_f, err := FunctionFromTPR(&t_f)
		if err != nil {
			return nil, err
		}
		f_functions = append(f_functions, *f_f)
	}
	return &f_functions, nil
}

func EnvironmentFromTPR(t_environment *Environment) (*fission.Environment, error) {
	metadata, err := MetadataFromTPR(&t_environment.Metadata)
	if err != nil {
		return nil, err
	}
	f_environment := &fission.Environment{
		Metadata: *metadata,
		Spec:     t_environment.Spec,
	}
	return f_environment, nil
}

func EnvironmentToTPR(f_environment *fission.Environment, t_metadata *api.ObjectMeta) *Environment {
	t_environment := &Environment{
		Metadata: *MetadataToTPR(&f_environment.Metadata, t_metadata),
		Spec:     f_environment.Spec,
	}
	return t_environment
}

func EnvironmentListFromTPR(t_environments *EnvironmentList) (*[]fission.Environment, error) {
	f_environments := make([]fission.Environment, len(t_environments.Items))
	for _, t_e := range t_environments.Items {
		f_e, err := EnvironmentFromTPR(&t_e)
		if err != nil {
			return nil, err
		}
		f_environments = append(f_environments, *f_e)
	}
	return &f_environments, nil
}

func HttpTriggerFromTPR(t_httptrigger *Httptrigger) (*fission.HTTPTrigger, error) {
	metadata, err := MetadataFromTPR(&t_httptrigger.Metadata)
	if err != nil {
		return nil, err
	}
	f_httptrigger := &fission.HTTPTrigger{
		Metadata: *metadata,
		Spec:     t_httptrigger.Spec,
	}
	return f_httptrigger, nil
}

func HttpTriggerToTPR(f_httptrigger *fission.HTTPTrigger, t_metadata *api.ObjectMeta) *Httptrigger {
	t_httptrigger := &Httptrigger{
		Metadata: *MetadataToTPR(&f_httptrigger.Metadata, t_metadata),
		Spec:     f_httptrigger.Spec,
	}
	return t_httptrigger
}

func HttpTriggerListFromTPR(t_httpTriggers *HttptriggerList) (*[]fission.HTTPTrigger, error) {
	f_httpTriggers := make([]fission.HTTPTrigger, len(t_httpTriggers.Items))
	for _, t_ht := range t_httpTriggers.Items {
		f_ht, err := HttpTriggerFromTPR(&t_ht)
		if err != nil {
			return nil, err
		}
		f_httpTriggers = append(f_httpTriggers, *f_ht)
	}
	return &f_httpTriggers, nil
}

func KubernetesWatchTriggerFromTPR(t_kuberneteswatchtrigger *Kuberneteswatchtrigger) (*fission.KubernetesWatchTrigger, error) {
	metadata, err := MetadataFromTPR(&t_kuberneteswatchtrigger.Metadata)
	if err != nil {
		return nil, err
	}
	f_kuberneteswatchtrigger := &fission.KubernetesWatchTrigger{
		Metadata: *metadata,
		Spec:     t_kuberneteswatchtrigger.Spec,
	}
	return f_kuberneteswatchtrigger, nil
}

func KubernetesWatchTriggerToTPR(
	f_kuberneteswatchtrigger *fission.KubernetesWatchTrigger,
	t_metadata *api.ObjectMeta) *Kuberneteswatchtrigger {

	t_kuberneteswatchtrigger := &Kuberneteswatchtrigger{
		Metadata: *MetadataToTPR(&f_kuberneteswatchtrigger.Metadata, t_metadata),
		Spec:     f_kuberneteswatchtrigger.Spec,
	}
	return t_kuberneteswatchtrigger
}

func KubernetesWatchTriggerListFromTPR(t_kubernetesWatchTriggers *KuberneteswatchtriggerList) (*[]fission.KubernetesWatchTrigger, error) {
	f_kubernetesWatchTriggers := make([]fission.KubernetesWatchTrigger, len(t_kubernetesWatchTriggers.Items))
	for _, t_wt := range t_kubernetesWatchTriggers.Items {
		f_wt, err := KubernetesWatchTriggerFromTPR(&t_wt)
		if err != nil {
			return nil, err
		}
		f_kubernetesWatchTriggers = append(f_kubernetesWatchTriggers, *f_wt)
	}
	return &f_kubernetesWatchTriggers, nil
}
