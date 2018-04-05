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
    "strconv"
    "strings"
)

const (
    InvalidErrorCode = 0
    InvalidDateRange = 1
)

type errorStack struct {
    code    int
    message string
}

func (i *errorStack) Error() string {
    return "code [" + strconv.FormatInt(int64(i.code), 10) + "]: " + i.message
}

func (i *errorStack) StackAfter(message string) {
    i.message += " " + message
}

func (i *errorStack) StackBefore(message string) {
    i.message = message + " " + i.message
}

func stackAfter(err error, message string) error {
    if es, ok := err.(*errorStack); ok {
        es.StackAfter(message)
    } else {
        es = fromError(err)
        es.StackAfter(message)
        err = es
    }
    return err
}

func stackBefore(err error, message string) error {
    if es, ok := err.(*errorStack); ok {
        es.StackAfter(message)
    } else {
        es = fromError(err)
        es.StackBefore(message)
        err = es
    }
    return err
}

func fromError(err error) *errorStack {
    message := err.Error()
    startInd := strings.IndexByte(message, '[') + 1
    endInd := strings.IndexByte(message, ']')
    if startInd > endInd {
        return &errorStack{code: InvalidErrorCode, message: message}
    }
    if code, err := strconv.ParseInt(message[startInd:endInd], 10, 64); err != nil {
        return &errorStack{code: InvalidErrorCode, message: message}
    } else {
        return &errorStack{code: int(code), message: message}
    }
}

func newError(code int, message string) *errorStack {
    return &errorStack{code: code, message: message}
}

// the given date range is invalid probably the end date is set as before start date
var ErrorInvalidDateRange = newError(InvalidDateRange, "invalid date range")
