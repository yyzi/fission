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
	"encoding/json"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/meta"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
)

//
// To add a Fission TPR type:
//   0. Create a "spec" type, for everything in the type except metadata
//   1. Create the type with metadata + the spec
//   2. Create a list type (for example see FunctionList and Function, below)
//   3. Add methods at the bottom of this file for satisfying Object and List interfaces
//   4. Add the type to configureClient in tpr.go
//   5. Add the type to EnsureFissionTPRs in tpr.go
//   6. Add tests to tpr_test.go
//

type (
	Package struct {
		Literal         []byte
		URL             string
		PackageStoreRef string

		entryPoint string // optional
	}
	FunctionSpec struct {
		Source         Package
		Deployment     Package
		EnvironmentUid string
	}
	Function struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta `json:"metadata"`
		Spec                 FunctionSpec   `json:"spec"`
	}
	FunctionList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Function `json:"items"`
	}

	// environment
	Runtime struct {
		Image string
	}
	Builder struct {
		Image   string
		Command string
	}
	EnvironmentSpec struct {
		Runtime
		Builder
		FilenameExtensions []string
		LanguageName       []string
	}
	Environment struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta  `json:"metadata"`
		Spec                 EnvironmentSpec `json:"spec"`
	}
	EnvironmentList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Environment `json:"items"`
	}

	// HTTP Triggers
	FunctionReference struct {
		// Selector selects a function by labels.  Functions have auto-assigned labels by Fission.
		Selector map[string]string
	}
	HTTPTriggerSpec struct {
		Host        string
		RelativeURL string
		Method      string
		FunctionReference
	}
	HTTPTrigger struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta  `json:"metadata"`
		Spec                 HTTPTriggerSpec `json:"spec"`
	}
	HTTPTriggerList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []HTTPTrigger `json:"items"`
	}

	// Kubernetes Watches as function triggers
	KubernetesWatchTriggerSpec struct {
		Namespace     string
		Type          string
		LabelSelector map[string]string
		FunctionReference
	}
	KubernetesWatchTrigger struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta             `json:"metadata"`
		Spec                 KubernetesWatchTriggerSpec `json:"spec"`
	}
	KubernetesWatchTriggerList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []KubernetesWatchTrigger `json:"items"`
	}
)

// Each TPR type needs:
//   GetObjectKind (to satisfy the Object interface)
//
// In addition, each singular TPR type needs:
//   GetObjectMeta (to satisfy the ObjectMetaAccessor interface)
//
// And each list TPR type needs:
//   GetListMeta (to satisfy the ListMetaAccessor interface)

func (f Function) GetObjectKind() unversioned.ObjectKind {
	return &f.TypeMeta
}
func (e *Environment) GetObjectKind() unversioned.ObjectKind {
	return &e.TypeMeta
}
func (ht *HTTPTrigger) GetObjectKind() unversioned.ObjectKind {
	return &ht.TypeMeta
}
func (w *KubernetesWatchTrigger) GetObjectKind() unversioned.ObjectKind {
	return &w.TypeMeta
}

func (f Function) GetObjectMeta() meta.Object {
	return &f.Metadata
}
func (e *Environment) GetObjectMeta() meta.Object {
	return &e.Metadata
}
func (ht *HTTPTrigger) GetObjectMeta() meta.Object {
	return &ht.Metadata
}
func (w *KubernetesWatchTrigger) GetObjectMeta() meta.Object {
	return &w.Metadata
}

func (fl FunctionList) GetObjectKind() unversioned.ObjectKind {
	return &fl.TypeMeta
}
func (el *EnvironmentList) GetObjectKind() unversioned.ObjectKind {
	return &el.TypeMeta
}
func (hl *HTTPTriggerList) GetObjectKind() unversioned.ObjectKind {
	return &hl.TypeMeta
}
func (wl *KubernetesWatchTriggerList) GetObjectKind() unversioned.ObjectKind {
	return &wl.TypeMeta
}

func (fl FunctionList) GetListMeta() unversioned.List {
	return &fl.Metadata
}
func (el *EnvironmentList) GetListMeta() unversioned.List {
	return &el.Metadata
}
func (hl *HTTPTriggerList) GetListMeta() unversioned.List {
	return &hl.Metadata
}
func (wl *KubernetesWatchTriggerList) GetListMeta() unversioned.List {
	return &wl.Metadata
}

// XXX do we need the TPR unmarshalling workaround?

type FunctionListCopy FunctionList
type FunctionCopy Function

func (e *Function) UnmarshalJSON(data []byte) error {
	tmp := FunctionCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Function(tmp)
	*e = tmp2
	return nil
}

func (el *FunctionList) UnmarshalJSON(data []byte) error {
	tmp := FunctionListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := FunctionList(tmp)
	*el = tmp2
	return nil
}
