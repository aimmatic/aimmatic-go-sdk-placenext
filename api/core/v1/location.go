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

// LocationMeasurement data, all below field is required
type LocationMeasurement struct {
    Ipaddress     string    `json:"ipaddress,omitempty"`
    AdvertisingId string    `json:"advertisingId,omitempty"`
    CampaignId    string    `json:"campaignId,omitempty"`
    Label         string    `json:"label,omitempty"`
    Location      *Location `json:"location,omitempty"`
}

// Location a location of latitude and longitude of device
type Location struct {
    Lat float64 `json:"lat,omitempty"`
    Lng float64 `json:"lon,omitempty"`
}

// IngestLocationMeasurement send a batch LocationMeasurement to the placenext server
func (p *placeNextImpl) IngestLocationMeasurement(lms []LocationMeasurement) (resp *Response, err error) {
    var req *http.Request
    var buf []byte
    url := p.client.Config().GetPlaceNextHost() + apiVersion + "/location/measurement"
    if buf, err = json.Marshal(lms); err == nil {
        req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
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
