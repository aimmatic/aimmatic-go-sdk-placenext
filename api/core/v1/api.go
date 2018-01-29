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

// Package v1 provides api version 1 access to placenext api
package v1

import "github.com/aimmatic/aimmatic-go-sdk-placenext/rest"

type Response struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    RequestId string `json:"requestId"`
}

const apiVersion = "/v1"

// Interface provides access to all api of placenext version 1
type PlaceNext interface {
    IngestLocationMeasurement([]LocationMeasurement) (*Response, error)
    IngestGeometry(data []*Geometry) (*Response, error)
}

type placeNextImpl struct {
    client *rest.Client
}

// NewPlaceNext returns a new PlaceNext.
func NewPlaceNext(client *rest.Client) PlaceNext {
    return &placeNextImpl{client: client}
}
