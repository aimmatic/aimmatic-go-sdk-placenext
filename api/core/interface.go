/*
Copyright 2018 The AimMatic Authors.
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

// Package core provides a core api client to access placenext api
package core

import (
    "github.com/aimmatic/aimmatic-go-sdk-placenext/rest"
    "github.com/aimmatic/aimmatic-go-sdk-placenext/api/core/v1"
)

// RestApi holding all api and api version
type RestApi interface {
    V1() v1.CoreV1
}

type restApiImpl struct {
    client *rest.Client
    v1     v1.CoreV1
}

// NewRestApi return a RestApi
func NewRestApi(client *rest.Client) RestApi {
    return &restApiImpl{
        client: client,
        v1:     v1.NewCoreV1(client),
    }
}

// V1 return a placenext api interface to access placenext api
func (r *restApiImpl) V1() (v1.CoreV1) {
    return r.v1
}
