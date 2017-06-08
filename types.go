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

package fission

type (
	// Metadata is used as the general identifier for all kinds of
	// resources managed by the controller.
	//
	// Name is mutable. Uid is immutable. Version must be updated
	// on any change to the object.  (Uid, Version) uniquely
	// identifies the contents of the object.
	//
	// Labels are entirely optional; they're for higher level
	// tools to organize functions or attach information to them.
	Metadata struct {
		Name        string            `json:"name"`
		Namespace   string            `json:"namespace,omitempty"`
		Uid         string            `json:"uid,omitempty"`
		Version     string            `json:"version,omitempty"`
		Labels      map[string]string `json:"labels,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}

	//
	// Functions
	//

	// Used to checksum the contents of a package when it is
	// stored outside the Package struct.  Sum is hex encoded.
	Checksum struct {
		Type string `json:"checksumType"`
		Sum  string `json:"checksum"`
	}

	// Package contains or references a collection of source or
	// binary files.
	Package struct {
		// Literal contents of the package.  Can be used for
		// encoding packages below TODO (256KB?) size.
		Literal []byte `json:"literal"`

		// Reference to a package, with a checksum. ChecksumType is
		URL      string   `json:"url"`
		Checksum Checksum `json:"checksum"`

		// Optional entry point in the package. Each
		// environment defines a default function entry point
		// name, but its name can be overridden here.
		entryPoint string `json:"entrypoint"`
	}

	FunctionSpec struct {
		Source          Package `json:"source"`
		Deployment      Package `json:"deployment"`
		EnvironmentName string  `json:"environmentName"`
	}

	Function struct {
		Metadata
		Spec FunctionSpec
	}

	FunctionReference struct {
		// Selector selects a function by labels.  Functions
		// have auto-assigned name, uid, and version labels in
		// addition to user labels.
		Selector map[string]string `json:"selector"`
	}

	//
	// Environments
	//

	Runtime struct {
		Image string `json:"image"`
	}
	Builder struct {
		Image   string `json:"image"`
		Command string `json:"command"`
	}
	EnvironmentSpec struct {
		// Environment API version
		Version int `json:"version"`

		// Runtime container image etc.; required
		Runtime `json:"runtime"`

		// Optional
		Builder `json:"builder"`
	}
	Environment struct {
		Metadata
		Spec EnvironmentSpec
	}

	//
	// Triggers
	//

	HTTPTriggerSpec struct {
		Host              string `json:"host"`
		RelativeURL       string `json:"relativeurl"`
		Method            string `json:"method"`
		FunctionReference `json:"functionref"`
	}
	HTTPTrigger struct {
		Metadata
		Spec HTTPTriggerSpec
	}

	KubernetesWatchTriggerSpec struct {
		Namespace         string            `json:"namespace"`
		Type              string            `json:"type"`
		LabelSelector     map[string]string `json:"labelselector"`
		FunctionReference `json:"functionref"`
	}
	KubernetesWatchTrigger struct {
		Metadata
		Spec KubernetesWatchTriggerSpec
	}

	// TimeTrigger invokes the specific function at a time or
	// times specified by a cron string.
	TimeTriggerSpec struct {
		Cron              string `json:"cron"`
		FunctionReference `json:"functionref"`
	}
	TimeTrigger struct {
		Metadata
		Spec TimeTriggerSpec
	}

	// Errors returned by the Fission API.
	Error struct {
		Code    errorCode `json:"code"`
		Message string    `json:"message"`
	}

	errorCode int
)

const (
	ErrorInternal = iota

	ErrorNotAuthorized
	ErrorNotFound
	ErrorNameExists
	ErrorInvalidArgument
	ErrorNoSpace
	ErrorNotImplmented
	ErrorTypeTranslation
	ErrorChecksumFail
)

// must match order and len of the above const
var errorDescriptions = []string{
	"Internal error",
	"Not authorized",
	"Resource not found",
	"Resource exists",
	"Invalid argument",
	"No space",
	"Not implemented",
}
