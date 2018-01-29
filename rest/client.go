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

// Package rest provides a help rest http client include config, compute
// authenticate signature and add necessary http header that required by
// placenext api server
package rest

import (
    "net/http"
    "sync"
    "time"
    "fmt"
)

// RESTClient imposes common placenext API conventions on a set of resource paths.
// The baseURL is expected to use Config host where the default value is point to
// placenext api server.
//
// Most consumers should use rest.DefaultClient() or using rest.NewRestClient to provide
// a configuration at runtime or dynamically.
type Client struct {
    *http.Client
    config Config
}

// a default rest client
var defaultRestClient *Client

// singleton for creating default client
var once sync.Once

// DefaultClient return a Client with the default config. The default config
// will contain an apikey and secret key from the variable environment.
func DefaultClient() *Client {
    once.Do(func() {
        defaultRestClient = &Client{
            Client: http.DefaultClient,
            config: DefaultConfig(),
        }
    })
    return defaultRestClient
}

// NewRestClient create a new rest client of default http.DefaultClient with
// the given config. This function is useful when you need to provide a different
// apikey and secret at runtime.
func NewRestClient(config Config) *Client {
    return &Client{
        Client: http.DefaultClient,
        config: config,
    }
}

// Override Do request of the default http client to add custom header and api authorization
func (c *Client) Do(req *http.Request) (*http.Response, error) {
    if err := addHeader(req, c.config); err != nil {
        return nil, err
    }
    return c.Client.Do(req)
}

// Config return a the config of the current client
func (c *Client) Config() Config {
    return c.config
}

// add custom header and api authorization
func addHeader(r *http.Request, config Config) error {
    // add date
    date := time.Now().UTC().Format(time.RFC1123)
    r.Header.Set(Date, date)
    r.Header.Set(XPlacenextDate, date)
    // TODO: remove content type of body is not available
    r.Header.Set(ContentType, MediaJson)
    // calculate signature
    _, _, signatureHashB64, md5hash, err := ComputeSignature(r, config.GetSecretKey())
    if err != nil {
        return err
    }
    r.Header.Set(ContentMD5, md5hash)
    r.Header.Set(Authorization, fmt.Sprintf("AimMatic %s:%s", config.GetApiKey(), signatureHashB64))
    return nil
}
