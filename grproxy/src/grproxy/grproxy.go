package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	var err error

	conn, _, err = zk.Connect([]string{"zookeeper"}, time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("Grproxy started...")

	http.HandleFunc("/", handler)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

var conn *zk.Conn
var currServer string
var roundRobinList []string

func handler(w http.ResponseWriter, r *http.Request) {

	var server *url.URL
	var err error

	fmt.Println("Request received...")

	fmt.Println("checking available servers...")

	if r.URL.Path == "/library" {

		children, _, _, err := conn.ChildrenW("/servers")
		if err != nil {
			panic(err)
		}

		if children == nil {
			fmt.Println("No servers found!!")
			return
		}
		serverSelected := false

		for _, child := range children {
			if contains(roundRobinList, child) {
				continue
			} else {
				roundRobinList = append(roundRobinList, child)
				currServer = child
				serverSelected = true
				break
			}
		}

		if !serverSelected {
			roundRobinList = nil
			roundRobinList = append(roundRobinList, children[0])
			currServer = children[0]
		}

		fmt.Println("Sending request to", currServer)
		server, err = url.Parse("http://" + currServer + ":80")
		if err != nil {
			panic(err)
		}
	} else {
		server, err = url.Parse("http://nginx:80")
		if err != nil {
			panic(err)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(server)

	proxy.ServeHTTP(w, r)

	fmt.Println(r.URL.Path)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
