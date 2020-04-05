package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Visteras/vscale/vscale"
)

func main() {

	token := os.Getenv("VSCALE_TOKEN")
	handler := &RegexpHandler{}
	a := action{vs: vscale.NewVScale(token)}

	handler.HandleFunc(http.MethodPost, regexp.MustCompile(`/create/[0-9]+`), a.CreateServers)
	handler.HandleFunc(http.MethodDelete, regexp.MustCompile(`/delete`), a.DeleteServers)

	log.Println(http.ListenAndServe(":3000", handler))
}
