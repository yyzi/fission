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
	// Functions and packages
	//

	// Checksum of package contents when the contents are stored
	// outside the Package struct. Type is the checksum algorithm;
	// "sha256" is the only currently supported one. Sum is hex
	// encoded.
	Checksum struct {
		Type string `json:"type"`
		Sum  string `json:"sum"`
	}

	PackageType string

	// Package contains or references a collection of source or
	// binary files.
	Package struct {
		// Type specifies how the package is stored: literal, URL, etc.
		Type PackageType `json:"type"`

		// Literal contents of the package. Can be used for
		// encoding packages below TODO (256KB?) size.
		Literal []byte `json:"literal"`

		// URL references a package.
		URL string `json:"url"`

		// Checksum ensures the integrity of packages
		// refereced by URL. Ignored for literals.
		Checksum Checksum `json:"checksum"`

		// EntryPoint optionally specifies an entry point in
		// the package. Each environment defines a default
		// entry point, but that can be overridden here.
		EntryPoint string `json:"entrypoint"`
	}

	// FunctionSpec describes the contents of the function.
	FunctionSpec struct {
		// EnvironmentName is the name of the environment that this function is associated
		// with. An Environment with this name should exist, otherwise the function cannot
		// be invoked.
		EnvironmentName string `json:"environmentName"`

		// Source is an source package for this function; it's used for the build step if
		// the environment defines a build container.
		Source Package `json:"source"`

		// Deployment is a deployable package for this function. This is the package that's
		// loaded into the environment's runtime container.
		Deployment Package `json:"deployment"`
	}

	Function struct {
		Metadata
		Spec FunctionSpec
	}

	FunctionReference struct {
		// Type indicates whether this function reference is by name or selector. For now,
		// the only supported reference type is by name.
		Type FunctionReferenceType `json:"type"`

		// Name of the function.
		Name string `json:"name"`
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
		Runtime Runtime `json:"runtime"`

		// Optional
		Builder Builder `json:"builder"`
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
	// PackageTypeLiteral means the package contents are specified in the Literal field of
	// resource itself.
	PackageTypeLiteral PackageType = "literal"

	// PackageTypeUrl means the package contents are at the specified URL.
	PackageTypeUrl PackageType = "url"
)

const (
	// FunctionReferenceTypeName means that the function reference is simply by name.
	FunctionReferenceTypeName = "name"
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
