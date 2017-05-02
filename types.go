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
		Metadata
		Spec FunctionSpec
	}

	FunctionReference struct {
		// Selector selects a function by labels.  Functions
		// have auto-assigned labels in addition to user
		// labels.
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
		// Environment version
		Version string `json:"version"`

		// Runtime container image etc.; required
		Runtime Runtime `json:"runtime"`

		// Optional
		Builder Builder `json:"builder"`

		// FilenameExtensions can be used by CLI/UI tooling
		// (e.g. auto-detect env by filename, syntax
		// highlighting etc.)  It isn't enforced by fission
		// itself in any way.
		FilenameExtensions []string `json:"filenameextensions"`
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

	// MessageQueueTrigger invokes a function on events in a
	// message queue.
	MessageQueueTrigger struct {
		Metadata         `json:"metadata"`
		Function         Metadata `json:"function"`
		MessageQueueType string   `json:"messageQueueType"`
		Topic            string   `json:"topic"`
		ResponseTopic    string   `json:"respTopic,omitempty"`
	}

	// TimeTrigger invokes the specific function at a time or
	// times specified by a cron string.
	TimeTrigger struct {
		Metadata `json:"metadata"`

		Cron string `json:"cron"`

		Function Metadata `json:"function"`
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
