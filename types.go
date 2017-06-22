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
	// Metadata for any resource. This is compatible with
	// Kubernetes API's v1.ObjectMeta. Only a subset (XXX) of
	// these fields are used in the Fission API; the rest are
	// included for type compatibility.
	//
	// See Kubernetes API docs for documentation on the meanings
	// of these fields.
	Metadata struct {
		Name                       string            `json:"name,omitempty"`
		GenerateName               string            `json:"generateName,omitempty"`
		Namespace                  string            `json:"namespace,omitempty"`
		SelfLink                   string            `json:"selfLink,omitempty"`
		UID                        types.UID         `json:"uid,omitempty"`
		ResourceVersion            string            `json:"resourceVersion,omitempty"`
		Generation                 int64             `json:"generation,omitempty"`
		CreationTimestamp          unversioned.Time  `json:"creationTimestamp,omitempty"`
		DeletionTimestamp          *unversioned.Time `json:"deletionTimestamp,omitempty"`
		DeletionGracePeriodSeconds *int64            `json:"deletionGracePeriodSeconds,omitempty"`
		Labels                     map[string]string `json:"labels,omitempty"`
		Annotations                map[string]string `json:"annotations,omitempty"`
		OwnerReferences            []OwnerReference  `json:"ownerReferences,omitempty"`
		Finalizers                 []string          `json:"finalizers,omitempty"`
		ClusterName                string            `json:"clusterName,omitempty"`
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

	// PackageType is either literal or URL, indicating whether
	// the package is specified in the Package struct or
	// externally.
	PackageType string

	// Package contains or references a collection of source or
	// binary files.
	Package struct {
		// Type defines how the package is specified: literal or URL.
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
		// the only supported reference type is by name.  Future reference types:
		//   * Function by label or annotation
		//   * Branch or tag of a versioned function
		//   * A "rolling upgrade" from one version of a function to another
		Type FunctionReferenceType `json:"type"`

		// Name of the function.
		Name string `json:"name"`
	}

	//
	// Environments
	//

	Runtime struct {
		// Image for containing the language runtime.
		Image string `json:"image"`

		// LoadEndpointPort defines the port on which the
		// server listens for function load
		// requests. Optional; default 8888.
		LoadEndpointPort int32 `json:"loadendpointport"`

		// LoadEndpointPath defines the relative URL on which
		// the server listens for function load
		// requests. Optional; default "/specialize".
		LoadEndpointPath string `json:"loadendpointpath"`

		// FunctionEndpointPort defines the port on which the
		// server listens for function requests. Optional;
		// default 8888.
		FunctionEndpointPort int32 `json:"functionendpointport"`
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

		// Optional, but strongly encouraged. Used to populate
		// links from UI, CLI, etc.
		DocumentationURL string `json:"documentationurl"`
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

	//
	// Fission-Environment interface. The following types are not
	// exposed in the Fission API, but rather used by Fission to
	// talk to environments.
	//
	FunctionLoadRequest struct {
		// FilePath is an absolute filesystem path to the
		// function. What exactly is stored here is
		// env-specific. Optional.
		FilePath string `json:"filepath"`

		// Entrypoint has an environment-specific meaning;
		// usually, it defines a function within a module
		// containing multiple functions. Optional; default is
		// environment-specific.
		EntryPoint string `json:"entrypoint"`

		// URL to expose this function at. Optional; defaults
		// to "/".
		URL string `json:"url"`
	}
)

const (
	// PackageTypeLiteral means the package contents are specified in the Literal field of
	// resource itself.
	PackageTypeLiteral PackageType = "literal"

	// PackageTypeUrl means the package contents are at the specified URL.
	PackageTypeUrl PackageType = "url"
)

const (
	// FunctionReferenceFunctionName means that the function
	// reference is simply by function name.
	FunctionReferenceTypeFunctionName = "name"
)

const (
	ErrorInternal = iota

	ErrorNotAuthorized
	ErrorNotFound
	ErrorNameExists
	ErrorInvalidArgument
	ErrorNoSpace
	ErrorNotImplmented
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
	"Checksum verification failed",
}
