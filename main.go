package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/cocopc/pcutil/handler"
	"github.com/gorilla/mux"
	"github.com/sdbaiguanghe/glog"
)

var (
	hostname string
	port     int
	tmpDir string
	dstHostDir string
)

/* register command line options */
func init() {
	flag.StringVar(&hostname, "hostname", "0.0.0.0", "监听的主机名")
	flag.IntVar(&port, "port", 8111, "监听的端口号")
	flag.StringVar(&tmpDir, "tmpDir", "/tmp/", "上传文件写入的的目录")
	// flag.StringVar(&dstHostDir, "dstHostDir", "", "目标主机目录")
}

func main() {

	defer glog.Flush()

	flag.Parse()

	var address = fmt.Sprintf("%s:%d", hostname, port)
	
	glog.Infoln("REST service listening on", address)

	base := &handler.Base{tmpDir}

	// register router
	router := mux.NewRouter().StrictSlash(true)
	router.
		HandleFunc("/api/htcp", base.HSCPHandler)

	// start server listening
	err := http.ListenAndServe(address, router)
	if err != nil {
		glog.Fatalln("ListenAndServe err:", err)
	}

	glog.Infoln("Server end")
}
