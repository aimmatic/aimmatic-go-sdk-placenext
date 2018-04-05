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
    "time"
    "strconv"

    "github.com/aimmatic/aimmatic-go-sdk-placenext/rest"
)

// insightsEndpoint return endpoint path based on the given name
func insightsEndpoint(config rest.Config, name string) string {
    return config.GetPlaceNextHost() + apiVersion + "/insights/" + name
}

// NSSResponse response of Net Sentiment Score
type NSSResponse struct {
    *Status
    Score int64 `json:"score"`
}

// GetNSS get Net Sentiment Score of all time
func (p *coreV1) GetNSS() (resp *NSSResponse, err error) {
    return p.GetNSSByRange(time.Time{}, time.Time{})
}

// GetNSSByRange get Net Sentiment Score (NSS) in between the given start and end time
func (p *coreV1) GetNSSByRange(start, end time.Time) (resp *NSSResponse, err error) {
    var req *http.Request
    req, err = http.NewRequest(http.MethodGet, insightsEndpoint(p.client.Config(), "nss"), nil)
    if !start.IsZero() && !end.IsZero() {
        if end.Before(start) {
            return nil, ErrorInvalidDateRange
        }
        query := req.URL.Query()
        query.Set("start", strconv.FormatInt(start.UnixNano(), 10))
        query.Set("end", strconv.FormatInt(end.UnixNano(), 10))
        req.URL.RawQuery = query.Encode()
    } else if start.IsZero() != end.IsZero() {
        return nil, stackAfter(ErrorInvalidDateRange, "start and end time both must be given")
    }
    var httpResp *http.Response
    if httpResp, err = p.client.Do(req); err == nil {
        if httpResp.ContentLength > 0 {
            resp = &NSSResponse{}
            err = json.NewDecoder(httpResp.Body).Decode(resp)
        } else {
            resp = &NSSResponse{Status: &Status{Code: 0, Message: "OK"}}
        }
    }
    return
}
