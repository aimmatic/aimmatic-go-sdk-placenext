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
    "os"
    "errors"
)

// Config holds the common attributes that can be passed to rest client
type Config interface {
    GetApiKey() string
    GetSecretKey() []byte
    GetPlaceNextHost() string
    GetUserAgent() string
}

// implement Config
type configImpl struct {
    // placenext apikey raw encoding base64 (no padding). By default, the default
    // config will read apikey from variable environment name PLACENEXT_APIKEY
    apiKey string

    // placenext secretkey raw encoding base64 (no padding). By default, the default
    // config will read apikey from variable environment name PLACENEXT_SECRETKEY
    secretKey string

    // placenext secretkey byte array, an actual key use with hash hmac-sha256 function
    // to compute signature. This key is decoded from base64 of SecretKey.
    rawSecretKey []byte

    // Host must be a host string, a host:port pair, or a URL to the base of the placenext
    // api server. If not provided it will use our default host placenext api server
    // api.aimmatic.com.
    host string

    // userAgent is an optional field that specifies the caller of this request.
    userAgent string
}

// GetApiKey return the api key
func (c *configImpl) GetApiKey() string {
    return c.apiKey
}

// GetSecretKey return a byte array of base64 decoded from secret key
func (c *configImpl) GetSecretKey() []byte {
    return c.rawSecretKey
}

// GetPlaceNextHost return a host domain or ip address of placenext api server
func (c *configImpl) GetPlaceNextHost() string {
    return c.host
}

// GetUserAgent return a placenext user-agent
func (c *configImpl) GetUserAgent() string {
    return c.userAgent
}

// Default config that initial when application start or a config form SetConfig
var defConf Config

// initialize placenext configuration with apikey and secret
// the apikey and secret must available via variable environment
func init() {
    if os.Getenv(PLACENEXT_SECRETKEY) == "" || os.Getenv(PLACENEXT_APIKEY) == "" {
        return
    }
    var err error
    if defConf, err = NewConfig(os.Getenv(PLACENEXT_APIKEY), os.Getenv(PLACENEXT_SECRETKEY)); err != nil {
        panic("invalid secret key " + err.Error())
    }
}

// DefaultConfig return a default config where api key and secret key is reading from variable environment.
// if variable environment is not available then default config will return nil otherwise a config will return
// Note: variable environment must export before application running if not default will be nil
func DefaultConfig() Config {
    return defConf
}

// NewConfig create new configure based on the given api key and secret key.
func NewConfig(apiKey, secretKey string) (Config, error) {
    rawSecretKey, err := GetSecretKeyAsByte(secretKey)
    if err != nil {
        err = errors.New("invalid secret " + err.Error())
        return nil, err
    }
    var placenextAddr string
    // if PLACENEXT_ADDRESS not given then default is set
    if placenextAddr = os.Getenv(PLACENEXT_ADDRESS); placenextAddr == "" {
        placenextAddr = scheme + "://" + domain
    }
    return &configImpl{
        apiKey:       apiKey,
        secretKey:    secretKey,
        rawSecretKey: rawSecretKey,
        host:         placenextAddr,
        userAgent:    defaultAgent,
    }, nil
}

// SetConfig set the given configuration globally as well as default Client
// If you need a different api key and secret key for each request to placenext api
// you must create new Client with your new config and then use the new Client with our function.
func SetConfig(config Config) {
    defConf = config
    DefaultClient().config = config
}
