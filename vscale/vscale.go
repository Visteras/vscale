package vscale

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type VScale struct {
	token string
}

func NewVScale(token string) *VScale {
	return &VScale{
		token: token,
	}
}

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	COPY    HttpMethod = "COPY"
	PATCH   HttpMethod = "PATCH"
	OPTIONS HttpMethod = "OPTIONS"
	HEAD    HttpMethod = "HEAD"
)

func (v *VScale) prepareBody(params interface{}) (io.Reader, error) {
	if params == nil {
		return http.NoBody, nil
	}
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(reqBody)
	return buffer, nil
}

func (v *VScale) fetch(method HttpMethod, url string, params interface{}) ([]byte, error) {
	client := http.Client{}

	buffer, err := v.prepareBody(params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(string(method), url, buffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Token", v.token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 429:
		log.Println("too many requests, sleep 1 second...")
		time.Sleep(time.Second)
		return v.fetch(method, url, params)
	case 403:
		return nil, BadRequest
	case 401:
		return nil, Unauthorized
	}

	if response.StatusCode > 299 || response.StatusCode < 200 {
		e := fmt.Errorf("code: %d, info: %s", response.StatusCode, response.Header.Get("Vscale-Error-Message"))
		return nil, e
	}

	return ioutil.ReadAll(response.Body)
}

func (v *VScale) CreateServers(servers []NewServer) ([]*Server, error) {
	n := len(servers)
	err := make(chan error, n)
	res := make(chan *Server, n)
	var wg sync.WaitGroup

	wg.Add(n)
	work := func(server NewServer, cErr chan error, cServer chan *Server) {
		defer wg.Done()
		res, err := v.CreateServer(&server)
		if err != nil {
			cErr <- err
		} else {
			cServer <- res
		}
	}

	for _, server := range servers {
		go work(server, err, res)
	}
	wg.Wait()
	close(err)
	close(res)

	hasError := false
	for {
		if e, opened := <-err; opened {
			log.Println(fmt.Errorf("[CreateServers][ERROR]: %w", e))
			hasError = true
		} else {
			break
		}
	}

	var newServers []*Server
	for {
		if srv, opened := <-res; opened {
			newServers = append(newServers, srv)
		} else {
			break
		}
	}

	if hasError {
		var ctids []int
		for _, newServer := range newServers {
			ctids = append(ctids, newServer.CTID)
		}
		_, err := v.DeleteServers(ctids)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("created servers not successfully")
	}

	return newServers, nil
}

func (v *VScale) CreateServer(server *NewServer) (*Server, error) {
	body, err := v.fetch(POST, "https://api.vscale.io/v1/scalets", server)
	if err != nil {
		return nil, fmt.Errorf("[VSCALE][API]: %s", err)
	}

	var result Server
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("[JSON][UNMARSHAL]: %s", err)
	}

	return &result, nil
}

func (v *VScale) DeleteServer(ctid int) (*Server, error) {
	body, err := v.fetch(DELETE, fmt.Sprintf("https://api.vscale.io/v1/scalets/%d", ctid), nil)
	if err != nil {
		return nil, fmt.Errorf("[VSCALE][API]: %s", err)
	}
	var result Server
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("[JSON][UNMARSHAL]: %s", err)
	}

	return &result, nil
}

func (v *VScale) GetAccointInfo() (*Account, error) {
	body, err := v.fetch(GET, "https://api.vscale.io/v1/account", nil)
	if err != nil {
		return nil, fmt.Errorf("[VSCALE][API]: %s", err)
	}

	var account Account
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, fmt.Errorf("[JSON][UNMARSHAL]: %s", err)
	}

	return &account, nil
}

func (v *VScale) DeleteServers(ctids []int) ([]*Server, error) {
	n := len(ctids)
	err := make(chan error, n)
	res := make(chan *Server, n)
	var wg sync.WaitGroup
	wg.Add(n)
	work := func(ctid int, cErr chan error, cServer chan *Server) {
		defer wg.Done()
		res, err := v.DeleteServer(ctid)
		if err != nil {
			cErr <- err
		} else {
			cServer <- res
		}
	}

	for _, ctid := range ctids {
		go work(ctid, err, res)
	}
	wg.Wait()
	close(err)
	close(res)

	hasError := false
	for {
		if e, opened := <-err; opened {
			log.Println(fmt.Errorf("[DeleteServers][ERROR]: %w", e))
			hasError = true
		} else {
			break
		}
	}

	var servers []*Server
	for {
		if srv, opened := <-res; opened {
			servers = append(servers, srv)
		} else {
			break
		}
	}

	if hasError {
		return nil, errors.New("deleted servers not successfully")
	}

	return servers, nil
}

func (v *VScale) GetAllServers() ([]*Server, error) {
	body, err := v.fetch(GET, "https://api.vscale.io/v1/scalets", nil)
	if err != nil {
		return nil, fmt.Errorf("[VSCALE][API]: %w", err)
	}
	var result []*Server
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("[JSON][UNMARSHAL]: %s", err)
	}
	return result, nil
}
