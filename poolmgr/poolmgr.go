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

package poolmgr

import (
	"fmt"
	"log"
	"strings"

	"github.com/dchest/uniuri"
	"k8s.io/client-go/1.5/pkg/api"

	"github.com/fission/fission/tpr"
)

// Given metadata, create a key that uniquely identifies the contents
// of the object. Since resourceVersion changes on every update and
// names are unique, name+resourceVersion identifies the
// content. (ResourceVersion may also update on status updates, so
// this will result in some unnecessary cache misses. That should be
// ok.)
func cacheKey(metadata *api.ObjectMeta) string {
	return fmt.Sprintf("%v_%v", metadata.UID, metadata.ResourceVersion)
}

// Start the poolmgr service.
func StartPoolmgr(controllerUrl string, fissionNamespace string, functionNamespace string, port int) error {
	controllerUrl = strings.TrimSuffix(controllerUrl, "/")

	fissionClient, kubernetesClient, err := tpr.MakeFissionClient()
	if err != nil {
		log.Printf("Failed to get kubernetes client: %v", err)
		return err
	}

	instanceId := uniuri.NewLen(8)
	cleanupOldPoolmgrResources(kubernetesClient, functionNamespace, instanceId)

	fsCache := MakeFunctionServiceCache()
	gpm := MakeGenericPoolManager(
		controllerUrl, fissionClient, kubernetesClient, fissionNamespace,
		functionNamespace, fsCache, instanceId)

	api := MakePoolmgr(gpm, fissionClient, fissionNamespace, fsCache)
	go api.Serve(port)

	return nil
}
