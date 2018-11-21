package rpc

import (
	"net/http"

	"github.com/fatih/color"
)

func Run() {
	http.HandleFunc("/test", testHander)
	http.ListenAndServe(":3000", nil)
}

func testHander(writer http.ResponseWriter, req *http.Request) {
	color.Red("test!")
}
