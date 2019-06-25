// Copyright 2019 PingCAP, Inc.
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
	"github.com/pingcap/tidb-operator/pkg/apis/pdapi"
	"github.com/pingcap/tidb-operator/pkg/apis/pingcap.com/v1alpha1"
)

// GetPDClient gets the pd client from the TidbCluster
func GetPDClient(pdi pdapi.PDControlInterface, tc *v1alpha1.TidbCluster) pdapi.PDClient {
	return pdi.GetPDClient(pdapi.Namespace(tc.GetNamespace()), tc.GetName())
}
