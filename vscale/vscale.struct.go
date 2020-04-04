package vscale

import (
	"encoding/json"
	"strings"
	"time"
)

type NewServer struct {
	MakeFrom string `json:"make_from"`
	RPlan    string `json:"rplan"`
	DoStart  bool   `json:"do_start"`
	Name     string `json:"name"`
	Keys     []int  `json:"keys,omitempty"`
	password string `json:"password,omitempty"`
	Location string `json:"location"`
}

func (s *NewServer) SetPassword(passwd string) {
	s.password = passwd
}

func (s NewServer) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		MakeFrom string `json:"make_from"`
		RPlan    string `json:"rplan"`
		DoStart  bool   `json:"do_start"`
		Name     string `json:"name"`
		Keys     []int  `json:"keys,omitempty"`
		Password string `json:"password,omitempty"`
		Location string `json:"location"`
	}{
		MakeFrom: s.MakeFrom,
		RPlan:    s.RPlan,
		DoStart:  s.DoStart,
		Name:     s.Name,
		Keys:     s.Keys,
		Password: s.password,
		Location: s.Location,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

type Account struct {
	Info   AccountInfo `json:"info"`
	Status string      `json:"status"`
}

type AccountInfo struct {
	AcceptCookies string   `json:"accept_cookies"`
	ActDate       JSONTime `json:"actdate"`
	Country       string   `json:"country"`
	Email         string   `json:"email"`
	EU            bool     `json:"eu"`
	FaceID        string   `json:"face_id"`
	ID            string   `json:"id"`
	IsBlocked     bool     `json:"is_blocked"`
	Locale        string   `json:"locale"`
	Middlename    string   `json:"middlename"`
	Mobile        string   `json:"mobile"`
	Name          string   `json:"name"`
	State         string   `json:"state"`
	Surname       string   `json:"surname"`
}

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if strInput == "null" {
		t.Time = time.Time{}
		return nil
	}
	tParse, err := time.Parse("2006-01-02 15:04:05.99", strInput)
	if err != nil {
		return err
	}
	t.Time = tParse
	return nil
}

type ServerKey struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ServerAddress struct {
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Address string `json:"address"`
}
type Server struct {
	Status         string         `json:"status"`
	Deleted        string         `json:"deleted"`
	PublicAddress  *ServerAddress `json:"public_address"`
	Active         bool           `json:"active"`
	Location       string         `json:"location"`
	Locked         bool           `json:"locked"`
	Hostname       string         `json:"hostname"`
	Created        string         `json:"created"`
	Keys           []ServerKey    `json:"keys"`
	PrivateAddress *ServerAddress `json:"private_address"`
	MadeFrom       string         `json:"made_from"`
	Name           string         `json:"name"`
	CTID           int            `json:"ctid"`
	RPlan          string         `json:"rplan"`
}
