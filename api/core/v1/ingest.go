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
    "net/http"
    "bytes"
    "encoding/json"
)

// Geometry a generic standard GeoJSON      geometry
type Geometry struct {
    Type        string      `json:"type,omitempty"`
    Coordinates interface{} `json:"coordinates,omitempty"`
}

// NewPointGeometry create a new Point GeoJSON geometry
func NewPointGeometry(coordinate []float64) *Geometry {
    return &Geometry{Type: "Point", Coordinates: coordinate}
}

// NewPolygonGeometry create a new Polygon GeoJSON geometry
func NewPolygonGeometry(coordinate [][][]float64) *Geometry {
    return &Geometry{Type: "Polygon", Coordinates: coordinate}
}

// IngestGeometry send geometry in GeoJSON format to api server
func (p *placeNextImpl) IngestGeometry(data []*Geometry) (resp *Response, err error) {
    var req *http.Request
    buf := bytes.NewBuffer(nil)
    url := p.client.Config().GetPlaceNextHost() + apiVersion + "/ingest/geometries"
    if err = json.NewEncoder(buf).Encode(data); err == nil {
        req, err = http.NewRequest(http.MethodPost, url, buf)
        var httpResp *http.Response
        if httpResp, err = p.client.Do(req); err == nil {
            if httpResp.ContentLength > 0 {
                resp = &Response{}
                err = json.NewDecoder(httpResp.Body).Decode(resp)
            } else {
                resp = &Response{Code: 0, Message: "OK"}
            }
        }
    }
    return
}
