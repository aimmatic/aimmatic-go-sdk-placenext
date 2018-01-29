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

package rest

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "time"
    "strings"
    "encoding/base64"
    "bytes"
    "io/ioutil"
    "net/http/httputil"
    "net/url"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) {
    now := time.Now().UTC().Format(time.RFC1123)
    if date := r.Header.Get(XPlacenextDate); date != "" && date != now {
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte("missing x-placenext-date"))
        return
    }
    if date := r.Header.Get(Date); date != "" && date != now {
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte("missing date"))
        return
    }
    // TODO: remove verify content type if content body is empty, like GET or DELETE method
    if ctype := r.Header.Get(ContentType); ctype != MediaJson {
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte("missing content-type"))
        return
    }
    var err error
    var apiKey string
    var requestSignature []byte
    var requestBodyMd5 []byte
    // check authorization
    if auth := r.Header.Get(Authorization); auth == "" {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("missing authorization"))
        return
    } else if index := strings.IndexByte(auth, ':'); index > 9 { // AimMatic ....
        apiKey = auth[9:index]
        requestSignature, err = base64.RawStdEncoding.DecodeString(auth[index+1:])
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("decode request signature failed " + err.Error()))
            return
        }
    }
    // check body md5
    if md5B64 := r.Header.Get(ContentMD5); md5B64 == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("missing content-md5"))
        return
    } else if requestBodyMd5, err = base64.RawStdEncoding.DecodeString(md5B64); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("decode request content md5 failed " + err.Error()))
        return
    }
    // compare body md5 and signature
    config, _ := NewConfig(apiKey, secretKey)
    signature, bodyMd5, _, _, err := ComputeSignature(r, config.GetSecretKey())
    if err != nil {
        w.Write([]byte("compute signature failed " + err.Error()))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if !bytes.Equal(bodyMd5, requestBodyMd5) {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("miss matched content md5"))
        return
    }
    if !bytes.Equal(signature, requestSignature) {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("miss matched request signature"))
        return
    }
}

func TestAuthentication(t *testing.T) {
    // create http mock server
    serverHttp := httptest.NewServer(http.HandlerFunc(simpleHandler))
    defer serverHttp.Close()
    client := serverHttp.Client()
    // create client resource and request
    config, _ := NewConfig(apiKey, secretKey)
    req, err := http.NewRequest("GET", serverHttp.URL, nil)
    if err != nil {
        t.Fatal(err)
    }
    addHeader(req, config)
    // process test
    resp, err := client.Do(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            t.Fatal(err)
        }
        t.Error("authentication failed due to", string(b))
    }
    // create https mock server
    serverHttps := httptest.NewTLSServer(http.HandlerFunc(simpleHandler))
    defer serverHttps.Close()
    client = serverHttps.Client()
    // new request with https
    req, err = http.NewRequest("GET", serverHttps.URL, nil)
    if err != nil {
        t.Fatal(err)
    }
    addHeader(req, config)
    // process test
    resp, err = client.Do(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            t.Fatal(err)
        }
        t.Error("authentication failed due to", string(b))
    }
    // test http behind https proxy
    serverHttpInternal := httptest.NewServer(http.HandlerFunc(simpleHandler))
    serverHttpsFrontend := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        proxyUrl, _ := url.Parse(serverHttpInternal.URL)
        proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
        r.Header.Set("X-Forwarded-Proto", "https")
        proxy.ServeHTTP(w, r)
    }))
    defer serverHttpInternal.Close()
    defer serverHttpsFrontend.Close()
    //
    req, err = http.NewRequest("GET", serverHttpsFrontend.URL, nil)
    if err != nil {
        t.Fatal(err)
    }
    addHeader(req, config)
    // process test
    resp, err = serverHttpsFrontend.Client().Do(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            t.Fatal(err)
        }
        t.Error("authentication failed due to", string(b))
    }
}
