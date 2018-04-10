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
    "log"
    "strings"
    "sort"
    "net/http"
    "io/ioutil"
    "bytes"
    "errors"
    "encoding/base64"
    "crypto/md5"
    "crypto/hmac"
    "crypto/sha256"
)

// ErrMissingContentType an error indicate the request does not have content-type header
var ErrMissingContentType = errors.New("content-type header is not available")

// ErrMissingDate an error indicate the request does not have date or x-placenext-date header
var ErrMissingDate = errors.New("x-placenext-date or date header is not available")

// temporary struct to store key of map use to sort key order
type keyOrder struct {
    key string
    val string
}

// slide of key order
type keyOrderSlice []*keyOrder

func (ks keyOrderSlice) Len() int           { return len(ks) }
func (ks keyOrderSlice) Less(i, j int) bool { return ks[i].key < ks[j].key }
func (ks keyOrderSlice) Swap(i, j int)      { ks[i], ks[j] = ks[j], ks[i] }

// create a sorted and concatenate header for all X-Placenext header
func concatenateHeader(r *http.Request) string {
    allHeader := map[string][]string(r.Header)
    size := len(allHeader)
    keyOrderSlice := make(keyOrderSlice, 0, size)
    for k := range allHeader {
        if strings.HasPrefix(k, "X-Placenext") {
            lk := strings.ToLower(k)
            allVal := ""
            if len(allHeader[k]) > 1 {
                sort.Strings(allHeader[k])
                allVal = strings.Join(allHeader[k], ",")
            } else {
                allVal = allHeader[k][0]
            }
            keyOrderSlice = append(keyOrderSlice, &keyOrder{
                key: lk,
                val: lk + ":" + allVal,
            })
        }
    }
    sort.Sort(keyOrderSlice)
    concat := ""
    for i := range keyOrderSlice {
        concat += keyOrderSlice[i].val
    }
    return concat
}

// GetSecretKeyAsByte will decode secret key from string base 64 (no padding)
// to a byte array. The byte array key will be cache with given secret key
func GetSecretKeyAsByte(secret string) (secretByte []byte, err error) {
    secretByte, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(secret)
    return
}

// ComputeBodyMd5Base64 read the whole body of the request and calculate hash md5
// if the content body is empty then the md5 of empty byte array will return
func ComputeBodyMd5Base64(r *http.Request) ([]byte) {
    // calculate body md5
    if r.ContentLength > 0 {
        var buf bytes.Buffer
        buf.ReadFrom(r.Body)
        r.Body = ioutil.NopCloser(&buf)
        bmd5 := md5.Sum(buf.Bytes())
        return bmd5[:]
    } else {
        return nil
    }
}

// ComputeSignature calculate hash result from the given request based on the secret key
// The ComputeSignature use Hmac with sha256 standard hash function to calculate signature of the current request.
func ComputeSignature(r *http.Request, secret []byte) (signature, contentMD5 []byte, signatureB64, contentMD5B64 string, err error) {
    // calculate body md5
    buf := bytes.NewBuffer(nil)
    contentMD5 = ComputeBodyMd5Base64(r)
    if contentMD5 != nil {
        contentMD5B64 = base64.RawStdEncoding.EncodeToString(contentMD5)
        // write content md5
        buf.WriteString(contentMD5B64)
        buf.WriteByte('\n')
    }
    // write content type
    if r.Header.Get("Content-Type") != "" {
        buf.WriteString(r.Header.Get("Content-Type"))
        buf.WriteByte('\n')
    }
    // write date
    if r.Header.Get("X-Placenext-Date") != "" {
        buf.WriteString(r.Header.Get("X-Placenext-Date"))
        buf.WriteByte('\n')
    } else if r.Header.Get("Date") != "" {
        buf.WriteString(r.Header.Get("Date"))
        buf.WriteByte('\n')
    } else {
        err = ErrMissingDate
        return
    }
    // write header concat
    buf.WriteString(concatenateHeader(r))
    buf.WriteByte('\n')
    // write url
    if r.URL.Scheme != "" {
        buf.WriteString(r.URL.String())
    } else {
        url := r.Host + r.URL.String()
        if strings.HasSuffix(url, "/") {
            url = url[:len(url)-1]
        }
        var scheme string
        if scheme = r.Header.Get(XForwardedProto); scheme == "" {
            if r.TLS != nil {
                scheme = "https"
            } else {
                scheme = "http"
            }
        }
        buf.WriteString(scheme + "://" + url)
    }
    // encode signature
    hm := hmac.New(sha256.New, secret)
    if _, err = hm.Write(buf.Bytes()); err == nil {
        signature = hm.Sum(nil)
        signatureB64 = base64.RawStdEncoding.EncodeToString(signature)
    } else {
        log.Println("Warning hmac.Write produce err which is strange behavior")
    }
    return
}
