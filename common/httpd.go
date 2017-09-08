package common

import (
	"net/http"
	_ "net/http/pprof"
)

type HttpService struct {
	BaseService

	localnode *LocalNode
}

type Response struct {
	IfaceName   string      `json:"iface"`
	LocalIPAddr string      `json:"local_ip_addr"`
	Peers       interface{} `json:"net_peers"`
}

func (hs *HttpService) Name() string {
	return "http-service"
}

func (hs *HttpService) Init(ln *LocalNode) (err error) {
	hs.localnode = ln
	//hs.iface = ln.Service("iface").(*InterfaceService)
	return nil
}

func (hs *HttpService) Run() error {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO. Lets show something cool here."))
	})
	http.ListenAndServe(":15080", nil)
	return nil
}
