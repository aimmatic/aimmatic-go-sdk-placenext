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
package v1

import (
	"net/http"
	"encoding/json"
	"bytes"

	"github.com/aimmatic/aimmatic-go-sdk-placenext/rest"
)

// insightsEndpoint return endpoint path based on the given name
func placeNextIngestEndpoint(config rest.Config, name string) string {
	return config.GetPlaceNextHost() + apiVersion + "/placeNextIngest/" + name
}

// Geometry a generic standard GeoJSON geometry
type Geometry struct {
	Type        string      `json:"type,omitempty"`
	Coordinates interface{} `json:"coordinates,omitempty"`
}

// GeometryCollection a standard GeoJSON geometry collection
type GeometryCollection struct {
	Type       string      `json:"type,omitempty"`
	Geometries []*Geometry `json:"geometries"`
}

// NewPointGeometry create a new Point GeoJSON geometry
func NewPointGeometry(coordinate []float64) *Geometry {
	return &Geometry{Type: "Point", Coordinates: coordinate}
}

// NewPolygonGeometry create a new Polygon GeoJSON geometry
func NewPolygonGeometry(coordinate [][][]float64) *Geometry {
	return &Geometry{Type: "Polygon", Coordinates: coordinate}
}

// GetNSS send geometry in GeoJSON format to api server
func (p *coreV1) GeometryImport(geometryCollection *GeometryCollection) (resp *Response, err error) {
	var req *http.Request
	buf := bytes.NewBuffer(nil)
	if err = json.NewEncoder(buf).Encode(geometryCollection); err == nil {
		req, err = http.NewRequest(http.MethodPost, placeNextIngestEndpoint(p.client.Config(), "GeometryImport"), buf)
		req.Header.Set(rest.ContentType, rest.MediaGeoJson)
		var httpResp *http.Response
		if httpResp, err = p.client.Do(req); err == nil {
			if httpResp.ContentLength > 0 {
				resp = &Response{}
				err = json.NewDecoder(httpResp.Body).Decode(resp)
			} else {
				resp = &Response{Status: &Status{Code: 0, Message: "OK"}}
			}
		}
	}
	return
}

// PointJSON point data
type PointJSON struct {
	AdvertisingId        string     `json:"advertisingId,omitempty"`
	AdvertisingIdType    string     `json:"advertisingIdType,omitempty"`
	Ipaddress            string     `json:"ipaddress,omitempty"`
	WifiBssid            string     `json:"wifiBssid,omitempty"`
	Coordinates          []*float64 `json:"coordinates,omitempty"`
	EffectiveCreatedDate int64      `json:"effectiveCreatedDate,omitempty"`
	EffectiveUpdatedDate int64      `json:"effectiveUpdatedDate,omitempty"`
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
}

// PointImport send a batch LocationMeasurement to the placenext server
func (p *coreV1) PointImport(lms []*PointJSON) (resp *Response, err error) {
	var req *http.Request
	var buf []byte
	if buf, err = json.Marshal(lms); err == nil {
		req, err = http.NewRequest(http.MethodPost, placeNextIngestEndpoint(p.client.Config(), "PointImport"), bytes.NewReader(buf))
		var httpResp *http.Response
		if httpResp, err = p.client.Do(req); err == nil {
			if httpResp.ContentLength > 0 {
				resp = &Response{}
				err = json.NewDecoder(httpResp.Body).Decode(resp)
			} else {
				resp = &Response{Status: &Status{Code: 0, Message: "OK"}}
			}
		}
	}
	return
}
