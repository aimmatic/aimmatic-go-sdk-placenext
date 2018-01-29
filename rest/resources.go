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

// variable environment
const (
    // variable environment to provide a placenext's api key, use to authenticate with
    // placenext backend api.
    PLACENEXT_APIKEY = "PLACENEXT_APIKEY"
    // variable environment to provide a placenext's api secret key, use to create signature
    // authenticate with placenext backend api.
    PLACENEXT_SECRETKEY = "PLACENEXT_SECRETKEY"
    // variable environment to provide address of placenext backend
    // this value can be a host domain or ipaddress pair with listen port.
    // if the port is a standard 80 (HTTP) and 443 (HTTPs) then listen port should not need it.
    // by default if PLACENEXT_ADDRESS is not available from variable environment, the default
    // address api.aimmatic.com is used.
    PLACENEXT_ADDRESS = "PLACENEXT_ADDRESS"
)

// placenext endpoint
const (
    domain       = "api.aimmatic.com"
    scheme       = "https"
    defaultAgent = "placenext 1.0"
)

// required http header to be include in rest api request
const (
    Authorization   = "Authorization"
    ContentType     = "Content-Type"
    ContentMD5      = "Content-MD5"
    Date            = "Date"
    XPlacenextDate  = "X-PlaceNext-Date"
    XForwardedProto = "X-Forwarded-Proto"
)

// Media content type
const (
    MediaJson = "application/json; charset=utf-8"
)
