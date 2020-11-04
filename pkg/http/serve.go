package serve

import (
	"fmt"
	"log"
	"net/http"
)

// HandleRequests 定义创建服务。
// Handler 为自定义处理器；
// Port 为端口，格式如 [:5000]；
// APIPath 为自定义处理器路径，可以用来加密，格式如：[/abc/123]。
func HandleRequests(Handler func(http.ResponseWriter, *http.Request), Port string, APIPath string) {
	http.HandleFunc(APIPath, Handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "OK") })
	log.Printf("Running at %s, and API path:\n    status: /health\n    customAPI: %s", Port, APIPath)
	log.Fatal(http.ListenAndServe(Port, nil))
}
