package main

import (
	"log"
	"strconv"

	"github.com/Visteras/vscale/vscale"
)

const token = ""
const passwd = "secret-password"

func main() {
	vs := vscale.NewVScale(token)
	//account, err := vs.GetAccointInfo()

	var servers []vscale.NewServer
	for i := 0; i < 11; i++ {
		server := vscale.NewServer{
			MakeFrom: "ubuntu_18.04_64_001_master",
			RPlan:    "small",
			DoStart:  true,
			Name:     "TmpSrv" + strconv.Itoa(i),
			Keys:     nil,
			Location: "spb0",
		}
		server.SetPassword(passwd)
		servers = append(servers, server)
	}

	res, err := vs.CreateServers(servers)
	if err != nil {
		log.Printf("create error: %s", err)
	}
	log.Println(res)

	//srvs, err := vs.GetAllServers()
	//if err != nil {
	// log.Printf("get all error: %s",err)
	//}
	//var ctids []int
	//for _, srv := range srvs {
	// ctids = append(ctids, srv.CTID)
	//}
	//res, err := vs.DeleteServers(ctids)
	//if err != nil {
	// log.Printf("del all error: %s",err)
	//}
	//for _, srv := range res {
	// log.Println(srv.CTID)
	//}

	log.Println("vscale")
}
