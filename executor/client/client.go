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

package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fission/fission"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	executorUrl string
	tappedByUrl map[string]bool
	requestChan chan string
	k8sClient *kubernetes.Clientset
}

func MakeClient(executorUrl string, k8sClient *kubernetes.Clientset) *Client {
	c := &Client{
		executorUrl: strings.TrimSuffix(executorUrl, "/"),
		k8sClient: k8sClient,
		tappedByUrl: make(map[string]bool),
		requestChan: make(chan string),
	}
	go c.service()
	return c
}



func (c *Client) debugPrintServiceObject(metadata *metav1.ObjectMeta) {
	log.Printf("debugPrintServiceObject from GetServiceForFunction for function: %s exeecutorUrl : %s", metadata.Name, c.executorUrl)
	executorService := strings.TrimPrefix(c.executorUrl, "http://")
	service := strings.Split(executorService, ".")
	if len(service) == 0 {
		log.Printf("Unable to split executor service into name and namespace")
		return
	}
	if c.k8sClient == nil {
		log.Printf("Unit test, k8sclient set nil")
		return
	}

	log.Printf("ns: %s, name: %s", service[1], service[0])
	svcObj, err := c.k8sClient.CoreV1().Services(service[1]).Get(service[0], metav1.GetOptions{})
	if err == nil {
		log.Printf("service address : %s", svcObj.Spec.ClusterIP)
		log.Printf("also dumping service object: %v", svcObj)
		return
	}

	log.Printf("error getting service object name: %s and namespace: %s, err : %v", service[0], service[1], err)
}



func (c *Client) GetServiceForFunction(metadata *metav1.ObjectMeta) (string, error) {
	executorUrl := c.executorUrl + "/v2/getServiceForFunction"

	c.debugPrintServiceObject(metadata)

	body, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", executorUrl, bytes.NewReader(body))
	if err != nil {
		log.Printf("error making a http request object in executor client")
	} else {
		log.Printf("http request object: %v", req)
	}


	resp, err := http.Post(executorUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fission.MakeErrorFromHTTP(resp)
	}

	svcName, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(svcName), nil
}

func (c *Client) service() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case serviceUrl := <-c.requestChan:
			c.tappedByUrl[serviceUrl] = true
		case <-ticker.C:
			urls := c.tappedByUrl
			c.tappedByUrl = make(map[string]bool)
			if len(urls) > 0 {
				go func() {
					for u := range c.tappedByUrl {
						c._tapService(u)
					}
					log.Printf("Tapped %v services in batch", len(urls))
				}()
				log.Printf("Tapped %v services in batch", len(urls))
			}
		}
	}
}

func (c *Client) TapService(serviceUrl *url.URL) {
	c.requestChan <- serviceUrl.String()
}

func (c *Client) _tapService(serviceUrlStr string) error {
	executorUrl := c.executorUrl + "/v2/tapService"

	resp, err := http.Post(executorUrl, "application/octet-stream", bytes.NewReader([]byte(serviceUrlStr)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fission.MakeErrorFromHTTP(resp)
	}
	return nil
}
