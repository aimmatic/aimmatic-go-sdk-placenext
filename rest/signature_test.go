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
    "encoding/base64"
    "bytes"
)

const (
    apiKey    = "UOCMBvhRFLwxDhUFdDeK2QpfvV80Og"
    secretKey = "dMAMNw6HE60xDhV0SWZNsVZSVW91culvEXBFLE76ij62wsZXXqI+aQ"
)

var (
    body = `{ "message":"I am body" }`
    // TODO: change UTC to GMT for standard
    date               = "Tue, 10 Nov 2009 23:00:00 UTC"
    bodyEmptyMD5       = "1B2M2Y8AsgTpgAmY7PhCfg"
    signatureEmptyBody = "7JQ0hgzALiCM7cOgWGD9t30ZpjS441z6G2wJLtyRqag"
    bodyMD5            = "5Eu4uX4ASRWSsmTQRzOllw"
    signatureBody      = "VeAyJzCoYYoscwbsXx8SIacyJVtuGjJXiPHMv1CLCAA"
)

func TestComputeSignature(t *testing.T) {
    config, _ := NewConfig(apiKey, secretKey)
    // test client request build
    // 1. empty body
    req, _ := http.NewRequest("GET", "http://api.aimmatic.com", nil)
    req.Header.Set(Date, date)
    req.Header.Set(XPlacenextDate, date)
    req.Header.Set(ContentType, MediaJson)
    signByte, md5Byte, signB64, md5B64, err := ComputeSignature(req, config.GetSecretKey())

    if err != nil {
        t.Error(err)
    }
    if md5B64 != bodyEmptyMD5 {
        t.Error("wrong body md5 expect base64", bodyEmptyMD5, "got", md5B64)
    }
    if signB64 != signatureEmptyBody {
        t.Error("wrong signature expect base64", signatureEmptyBody, "got", signB64)
    }
    b, _ := base64.RawStdEncoding.DecodeString(bodyEmptyMD5)
    if !bytes.Equal(md5Byte, b) {
        t.Error("wrong md5 byte array hashing expect", b, "got", md5Byte)
    }
    b, _ = base64.RawStdEncoding.DecodeString(signatureEmptyBody)
    if !bytes.Equal(signByte, b) {
        t.Error("wrong signature byte array hashing expect", b, "got", signByte)
    }
    // 1. not empty body
    req, _ = http.NewRequest("GET", "http://api.aimmatic.com", bytes.NewBufferString(body))
    req.Header.Set(Date, date)
    req.Header.Set(XPlacenextDate, date)
    req.Header.Set(ContentType, MediaJson)
    signByte, md5Byte, signB64, md5B64, err = ComputeSignature(req, config.GetSecretKey())

    if err != nil {
        t.Error(err)
    }
    if md5B64 != bodyMD5 {
        t.Error("wrong body md5 expect base64", bodyMD5, "got", md5B64)
    }
    if signB64 != signatureBody {
        t.Error("wrong signature expect base64", signatureBody, "got", signB64)
    }
    b, _ = base64.RawStdEncoding.DecodeString(bodyMD5)
    if !bytes.Equal(md5Byte, b) {
        t.Error("wrong md5 byte array hashing expect", b, "got", md5Byte)
    }
    b, _ = base64.RawStdEncoding.DecodeString(signatureBody)
    if !bytes.Equal(signByte, b) {
        t.Error("wrong signature byte array hashing expect", b, "got", signByte)
    }
}
