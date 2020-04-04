package vscale

import "errors"

var toManyRequests = errors.New("to many requests")
var badToken = errors.New("bad token")
