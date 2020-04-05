package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Visteras/vscale/vscale"
	"github.com/sethvargo/go-password/password"
)

type route struct {
	pattern *regexp.Regexp
	method  string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(method string, pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, method, handler})
}

func (h *RegexpHandler) HandleFunc(method string, pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, method, http.HandlerFunc(handler)})
}

func RemoteAddr(r *http.Request) string {

	addr := r.Header.Get("X-Real-IP")
	if len(addr) == 0 {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i > -1 {
				addr = addr[:i]
			}
		}
	}
	return addr
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.method == r.Method {
			log.Printf("[ServeHTTP][REQUEST][%s] %s from %s\n", r.Method, r.URL.Path, RemoteAddr(r))
			route.handler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

type action struct {
	vs *vscale.VScale
}

func PasswordGenerator(length uint) string {
	res, err := password.Generate(int(length), int(length)/3, 0, false, false)
	if err != nil {
		log.Println(fmt.Errorf("[PASSWORD][ERROR]: %w", err))
		res = passwordGenerator(length)
	}
	return res
}

func passwordGenerator(length uint) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := uint(2); i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	for i := len(buf) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

type jsonRes struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg,omitempty"`
	Error error       `json:"error,omitempty"`
}

func (j jsonRes) MarshalJSON() ([]byte, error) {
	st := struct {
		Code  int         `json:"code"`
		Msg   interface{} `json:"msg,omitempty"`
		Error string      `json:"error,omitempty"`
	}{
		Code: j.Code,
		Msg:  j.Msg,
	}
	if j.Error != nil {
		st.Error = j.Error.Error()
	}
	r, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (j *jsonRes) ToJSON() []byte {
	res, err := json.Marshal(j)
	if err != nil {
		log.Println(fmt.Errorf("[ToJSON][MARSHAL]: %w", err))
		return nil
	}
	return res
}

func (a *action) CreateServers(w http.ResponseWriter, r *http.Request) {
	strCount := strings.Split(r.URL.Path, "/")[2]
	count, err := strconv.Atoi(strCount)
	out := jsonRes{
		Code: 200,
	}
	if err != nil {
		out.Code, out.Error = checkError(err)
		w.WriteHeader(out.Code)
		_, _ = w.Write(out.ToJSON())
		return
	}
	if count == 0 {
		out.Code = 441
		out.Error = fmt.Errorf(`variable 'count' cannot be 0`)
		w.WriteHeader(out.Code)
		_, _ = w.Write(out.ToJSON())
		return
	}

	var servers []vscale.NewServer
	for i := 0; i < count; i++ {
		server := vscale.NewServer{
			MakeFrom: "ubuntu_18.04_64_001_master",
			RPlan:    "small",
			DoStart:  true,
			Name:     "TmpSrv" + strconv.Itoa(i),
			Keys:     nil,
			Location: "spb0",
		}

		server.SetPassword(PasswordGenerator(16))

		servers = append(servers, server)
	}

	res, err := a.vs.CreateServers(servers)
	if err != nil {
		out.Code, out.Error = checkError(err)
		w.WriteHeader(out.Code)
		_, _ = w.Write(out.ToJSON())
		return
	}

	out.Msg = res
	w.WriteHeader(out.Code)
	_, _ = w.Write(out.ToJSON())
	return
}

func (a *action) DeleteServers(w http.ResponseWriter, _ *http.Request) {
	out := jsonRes{
		Code: 200,
	}
	srvs, err := a.vs.GetAllServers()
	if err != nil {
		out.Code, out.Error = checkError(err)
		w.WriteHeader(out.Code)
		_, _ = w.Write(out.ToJSON())
		return
	}
	var ctids []int
	for _, srv := range srvs {
		ctids = append(ctids, srv.CTID)
	}
	res, err := a.vs.DeleteServers(ctids)
	if err != nil {
		out.Code, out.Error = checkError(err)
		w.WriteHeader(out.Code)
		_, _ = w.Write(out.ToJSON())
		return
	}

	out.Msg = res
	w.WriteHeader(out.Code)
	_, _ = w.Write(out.ToJSON())
	return
}

func checkError(err error) (int, error) {
	if errors.Is(err, vscale.Unauthorized) {
		return 500, fmt.Errorf("internal error")
	}
	return 500, err
}
