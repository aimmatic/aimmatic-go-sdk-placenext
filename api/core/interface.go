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
    "github.com/aimmatic/aimmatic-go-sdk-placenext/api/core/v1"
    "github.com/aimmatic/aimmatic-go-sdk-placenext/rest"
    "sync"
)

// RestApi holding all api and api version
type RestApi interface {
    V1() v1.PlaceNext
}

type restApiImpl struct {
    client *rest.Client
    mu     sync.RWMutex
    api    map[uint8]interface{}
}

// NewRespApi return a RestApi
func NewRespApi(client *rest.Client) RestApi {
    return &restApiImpl{client: client, api: make(map[uint8]interface{})}
}

// V1 return a placenext api interface to access placenext api
func (r *restApiImpl) V1() (placenext v1.PlaceNext) {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.api[1] != nil {
        var ok bool
        if placenext, ok = r.api[1].(v1.PlaceNext); !ok {
            panic("error v1 api instance is not a placenext api")
        }
    } else {
        placenext = v1.NewPlaceNext(r.client)
        r.api[1] = placenext
    }
    return
}
