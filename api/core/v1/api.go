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

import (
    "github.com/aimmatic/aimmatic-go-sdk-placenext/rest"
    "time"
)

type Status struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    RequestId string `json:"requestId"`
}

type Response struct {
    Status *Status `json:"status"`
}

const apiVersion = "/v1"

// CoreV1 provides access to all api of aimmatic service
type CoreV1 interface {
    placeNext
    insights
}

type placeNext interface {
    PointImport([]*PointJSON) (*Response, error)
    GeometryImport(*GeometryCollection) (*Response, error)
}

type insights interface {
    GetNSS() (*NSSResponse, error)
    GetNSSByRange(start, end time.Time) (*NSSResponse, error)
}

type coreV1 struct {
    client *rest.Client
}

// NewCoreV1 returns a new CoreV1 object.
func NewCoreV1(client *rest.Client) CoreV1 {
    return &coreV1{client: client}
}
