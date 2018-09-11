// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pingcap/tidb-operator/pkg/apis/pingcap.com/v1alpha1"
)

// TiDBControlInterface is the interface that knows how to manage tidb peers
type TiDBControlInterface interface {
	// GetHealth returns tidb's health info
	GetHealth(tc *v1alpha1.TidbCluster) map[string]bool
}

// defaultTiDBControl is default implementation of TiDBControlInterface.
type defaultTiDBControl struct {
	httpClient *http.Client
}

// NewDefaultTiDBControl returns a defaultTiDBControl instance
func NewDefaultTiDBControl() TiDBControlInterface {
	httpClient := &http.Client{Timeout: timeout}
	return &defaultTiDBControl{httpClient: httpClient}
}

func (tdc *defaultTiDBControl) GetHealth(tc *v1alpha1.TidbCluster) map[string]bool {
	tcName := tc.GetName()
	ns := tc.GetNamespace()

	result := map[string]bool{}
	for i := 0; i < int(tc.Spec.TiDB.Replicas); i++ {
		hostName := fmt.Sprintf("%s-%d", TiDBMemberName(tcName), i)
		url := fmt.Sprintf("http://%s.%s-tidb-peer.%s:10080/status", hostName, tcName, ns)
		_, err := tdc.getBodyOK(url)
		if err != nil {
			result[hostName] = false
		} else {
			result[hostName] = true
		}
	}
	return result
}

func (tdc *defaultTiDBControl) getBodyOK(apiURL string) ([]byte, error) {
	res, err := tdc.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		errMsg := fmt.Errorf(fmt.Sprintf("Error response %v", res.StatusCode))
		return nil, errMsg
	}

	defer DeferClose(res.Body, &err)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// FakeTiDBControl is a fake implementation of TiDBControlInterface.
type FakeTiDBControl struct {
	healthInfo map[string]bool
}

// NewFakeTiDBControl returns a FakeTiDBControl instance
func NewFakeTiDBControl() *FakeTiDBControl {
	return &FakeTiDBControl{}
}

// SetHealth set health info for FakeTiDBControl
func (ftd *FakeTiDBControl) SetHealth(healthInfo map[string]bool) {
	ftd.healthInfo = healthInfo
}

func (ftd *FakeTiDBControl) GetHealth(tc *v1alpha1.TidbCluster) map[string]bool {
	return ftd.healthInfo
}