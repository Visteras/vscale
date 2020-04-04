package vscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
)

type TestVScale struct {
	*VScale
}

var server = NewServer{
	MakeFrom: "ubuntu_18.04_64_001_master",
	RPlan:    "small",
	DoStart:  true,
	Name:     "TmpSrv",
	password: "000",
	Keys:     nil,
	Location: "spb0",
}
var createServer = []byte(`{
    "status": "defined",
    "deleted": null,
    "public_address": {},
    "active": false,
    "location": "spb0",
    "locked": true,
    "hostname": "cs11533.vscale.io",
    "created": "20.08.2015 14:57:04",
    "keys": [],
    "private_address": {},
    "made_from": "ubuntu_18.04_64_001_master",
    "name": "TmpSrv",
    "ctid": 11,
    "rplan": "small"
}`)

func testFetch(token string, method HttpMethod, url string, params interface{}) ([]byte, error) {
	switch method {
	case POST:
		switch url {
		case "https://api.vscale.io/v1/scalets":
			//CreateServer
			return testCreateServer(params)
		}
	case DELETE:
		if strings.Contains(url, "https://api.vscale.io/v1/scalets/") {
			return testDeleteServer(url, params)
		}
	}
	return nil, errors.New("test")
}

func testDeleteServer(url string, params interface{}) ([]byte, error) {
	strId := strings.TrimLeft(url, "https://api.vscale.io/v1/scalets/")
	log.Println(strId)

	return nil, nil
}

func testCreateServer(params interface{}) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("error from vscale, code: %d, info: %s", 400, "Incorrect json received")

	}
	ns := params.(*NewServer)
	if ns.password == "" && len(ns.Keys) == 0 {
		return nil, fmt.Errorf("error from vscale, code: %d, info: %s", 400, "Password or ssh keys should be passed")
	}

	return createServer, nil
}

func NewTestVScale() *TestVScale {
	vs := &TestVScale{
		VScale: NewVScale(""),
	}
	vs.fetch = testFetch

	return vs
}

func TestVScale_CreateServer(t *testing.T) {
	vs := NewTestVScale()
	res, err := vs.CreateServer(&server)
	log.Println(res, err)

}

func TestVScale_DeleteServer(t *testing.T) {
	vs := NewTestVScale()

	res, err := vs.DeleteServer(1)
	log.Println(res, err)
}
