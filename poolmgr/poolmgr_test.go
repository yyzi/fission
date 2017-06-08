package poolmgr

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/labels"

	"github.com/fission/fission"
	"github.com/fission/fission/poolmgr/client"
	"github.com/fission/fission/tpr"
	"io/ioutil"
)

type MockController struct {
	port int
}

func MakeMockController(port int) *MockController {
	funcBody := `
module.exports = async function(context) {
    return {
        status: 200,
        body: "Hello, world!\n"
    };
}
`
	mc := &MockController{port: port}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ignore request, always respond with the same
		// response. The mock controller is only used to get
		// the function code.
		log.Printf("Serving %v", r.URL)
		_, err := w.Write([]byte(funcBody))
		if err != nil {
			log.Panicf("mock controller failed to write response: %v", err)
		}
	})
	go http.ListenAndServe(fmt.Sprintf(":%v", mc.port), nil)
	return mc
}

func (mc *MockController) Url() string {
	return fmt.Sprintf("http://localhost:%v", mc.port)
}

// return the number of pods in the given namespace matching the given labels
func countPods(kubeClient *kubernetes.Clientset, ns string, labelz map[string]string) int {
	pods, err := kubeClient.Pods(ns).List(api.ListOptions{
		LabelSelector: labels.Set(labelz).AsSelector(),
	})
	if err != nil {
		log.Panicf("Failed to list pods: %v", err)
	}
	return len(pods.Items)
}

func createTestNamespace(kubeClient *kubernetes.Clientset, ns string) {
	_, err := kubeClient.Namespaces().Create(&v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: ns,
		},
	})
	if err != nil {
		log.Panicf("failed to create ns %v: %v", ns, err)
	}
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("HTTP Get failed: URL %v: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("HTTP Get failed to read body: URL %v: %v", url, err)
	}
	return string(body)
}

func TestPoolmgr(t *testing.T) {
	// run in a random namespace so we can have concurrent tests
	// on a given cluster
	rand.Seed(time.Now().UTC().UnixNano())
	fissionNs := fmt.Sprintf("test-%v", rand.Intn(999))
	functionNs := fmt.Sprintf("test-function-%v", rand.Intn(999))

	// skip test if no cluster available for testing
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		t.Skip("Skipping test, no kubernetes cluster")
		return
	}

	// connect to k8s
	// and get TPR client
	fissionClient, kubeClient, err := tpr.MakeFissionClient(fissionNs)
	if err != nil {
		log.Panicf("failed to connect: %v", err)
	}

	// create the test's namespaces
	createTestNamespace(kubeClient, fissionNs)
	defer kubeClient.Namespaces().Delete(fissionNs, nil)

	createTestNamespace(kubeClient, functionNs)
	defer kubeClient.Namespaces().Delete(functionNs, nil)

	// make sure TPR types exist on cluster
	err = tpr.EnsureFissionTPRs(kubeClient)
	if err != nil {
		log.Panicf("failed to ensure tprs: %v", err)
	}
	//	tpr.WaitForTPRs(kubeClient)

	// create a mock controller for fetcher
	mockController := MakeMockController(9000)
	mcResp := httpGet(mockController.Url())
	log.Printf("mock controller response: %v", mcResp)

	functionUrl := fmt.Sprintf("%v/v1/functions/%v?uid=%v&raw=1",
		mockController.Url(), "hi", "42")
	mcResp = httpGet(functionUrl)
	log.Printf("mock controller response: %v", mcResp)

	// create an env on the cluster
	env, err := fissionClient.Environments.Create(&tpr.Environment{
		Metadata: api.ObjectMeta{
			Name:      "nodejs",
			Namespace: fissionNs,
		},
		Spec: fission.EnvironmentSpec{
			Version: 1,
			Runtime: fission.Runtime{
				Image: "fission/node-env",
			},
			Builder: fission.Builder{},
		},
	})
	if err != nil {
		log.Panicf("failed to create env: %v", err)
	}

	// create poolmgr
	port := 9999
	err = StartPoolmgr(mockController.Url(), fissionNs, functionNs, port)
	if err != nil {
		log.Panicf("failed to start poolmgr: %v", err)
	}

	// connect poolmgr client
	poolmgrClient := client.MakeClient(fmt.Sprintf("http://localhost:%v", port))

	// Wait for pool to be created (we don't actually need to do
	// this, since the API should do the right thing in any case).
	// waitForPool(functionNs, "nodejs")

	// create a function
	f := &tpr.Function{
		Metadata: api.ObjectMeta{
			Name:      "hello",
			Namespace: fissionNs,
		},
		Spec: fission.FunctionSpec{
			Source:          fission.Package{},
			Deployment:      fission.Package{},
			EnvironmentName: env.Metadata.Name,
		},
	}
	_, err = fissionClient.Functions.Create(f)
	if err != nil {
		log.Panicf("failed to create function: %v", err)
	}

	// the main test: get a service for a given function
	svc, err := poolmgrClient.GetServiceForFunction(&f.Metadata)
	if err != nil {
		log.Panicf("failed to get func svc: %v", err)
	}
	log.Printf("svc for function created at: %v", svc)

	// ensure that a pod with the label functionName=f.Metadata.Name exists
	podCount := countPods(kubeClient, functionNs, map[string]string{"funcName": f.Metadata.Name})
	if podCount != 1 {
		log.Panicf("expected 1 function pod, found %v", podCount)
	}

	// call the service to ensure it works

	// wait for a bit

	// tap service to simulate calling it again

	// make sure the same pod is still there

	// wait for idleTimeout to ensure the pod is removed

	// remove env

	// wait for pool to be destroyed

	// that's it
}
