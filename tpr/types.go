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
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/meta"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
)

//
// To add a Fission TPR type:
//   1. Create a "spec" type, for everything in the type except metadata
//   2. Create the type with metadata + the spec
//   3. Create a list type (for example see FunctionList and Function, below)
//   4. Add methods at the bottom of this file for satisfying Object and List interfaces
//   5. Add the type to configureClient in tpr.go
//   6. Add the type to EnsureFissionTPRs in tpr.go
//   7. Add tests to tpr_test.go
//

type (
	// Function

	// Package contains or references a collection of source or
	// binary files.
	Package struct {
		// Literal can be used for encoding packages below a certain size.
		Literal []byte `json:"literal"`

		// URL can be used to reference
		URL             string `json:"url"`
		PackageStoreRef string `json:"packagestoreref"`

		// Optional
		entryPoint string `json:"entrypoint"`
	}
	FunctionSpec struct {
		Source         Package `json:"source"`
		Deployment     Package `json:"deployment"`
		EnvironmentUid string  `json:"environmentuid"`
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

	// Environment
	Runtime struct {
		Image string `json:"image"`
	}
	Builder struct {
		Image   string `json:"image"`
		Command string `json:"command"`
	}
	EnvironmentSpec struct {
		// Environment version
		Version string `json:"version"`

		// Runtime container image etc.; required
		Runtime `json:"runtime"`

		// Optional
		Builder `json:"builder"`

		// FilenameExtensions can be used by CLI/UI tooling
		// (e.g. auto-detect env by filename, syntax
		// highlighting etc.)  It isn't enforced by fission
		// itself in any way.
		FilenameExtensions []string `json:"filenameextensions"`
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

	// HTTP Triggers.  (Something in the TPR reflection stuff wants
	// it to be spelled "Httptrigger" not "HTTPTrigger" or even
	// "HttpTrigger".  Bleh.)
	FunctionReference struct {
		// Selector selects a function by labels.  Functions
		// have auto-assigned labels in addition to user
		// labels.
		Selector map[string]string `json:"selector"`
	}
	HttptriggerSpec struct {
		Host        string `json:"host"`
		RelativeURL string `json:"relativeurl"`
		Method      string `json:"method"`
		FunctionReference
	}
	Httptrigger struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta  `json:"metadata"`
		Spec                 HttptriggerSpec `json:"spec"`
	}
	HttptriggerList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Httptrigger `json:"items"`
	}

	// Kubernetes Watches as function triggers
	KuberneteswatchtriggerSpec struct {
		Namespace         string            `json:"namespace"`
		Type              string            `json:"type"`
		LabelSelector     map[string]string `json:"labelselector"`
		FunctionReference `json:"functionref"`
	}
	Kuberneteswatchtrigger struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             api.ObjectMeta             `json:"metadata"`
		Spec                 KuberneteswatchtriggerSpec `json:"spec"`
	}
	KuberneteswatchtriggerList struct {
		unversioned.TypeMeta `json:",inline"`
		Metadata             unversioned.ListMeta `json:"metadata"`

		Items []Kuberneteswatchtrigger `json:"items"`
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

func (f *Function) GetObjectKind() unversioned.ObjectKind {
	return &f.TypeMeta
}
func (e *Environment) GetObjectKind() unversioned.ObjectKind {
	return &e.TypeMeta
}
func (ht *Httptrigger) GetObjectKind() unversioned.ObjectKind {
	return &ht.TypeMeta
}
func (w *Kuberneteswatchtrigger) GetObjectKind() unversioned.ObjectKind {
	return &w.TypeMeta
}

func (f *Function) GetObjectMeta() meta.Object {
	return &f.Metadata
}
func (e *Environment) GetObjectMeta() meta.Object {
	return &e.Metadata
}
func (ht *Httptrigger) GetObjectMeta() meta.Object {
	return &ht.Metadata
}
func (w *Kuberneteswatchtrigger) GetObjectMeta() meta.Object {
	return &w.Metadata
}

func (fl *FunctionList) GetObjectKind() unversioned.ObjectKind {
	return &fl.TypeMeta
}
func (el *EnvironmentList) GetObjectKind() unversioned.ObjectKind {
	return &el.TypeMeta
}
func (hl *HttptriggerList) GetObjectKind() unversioned.ObjectKind {
	return &hl.TypeMeta
}
func (wl *KuberneteswatchtriggerList) GetObjectKind() unversioned.ObjectKind {
	return &wl.TypeMeta
}

func (fl *FunctionList) GetListMeta() unversioned.List {
	return &fl.Metadata
}
func (el *EnvironmentList) GetListMeta() unversioned.List {
	return &el.Metadata
}
func (hl *HttptriggerList) GetListMeta() unversioned.List {
	return &hl.Metadata
}
func (wl *KuberneteswatchtriggerList) GetListMeta() unversioned.List {
	return &wl.Metadata
}

// In the client-go TPR example, UnmarshalJSON is defined here for the
// singular and list types.  That's supposed to be a workaround for
// some ugorji bug.  But we don't seem to need it, and all our tests
// pass without it, so we don't define any UnmarshalJSON methods.
