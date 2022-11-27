package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {

	time.Sleep(20 * time.Second)

	conn, _, err := zk.Connect([]string{"zookeeper"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	serverList(conn)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func serverList(conn *zk.Conn) {
	fmt.Println("checking...")

	children, _, _, err := conn.ChildrenW("/servers")
	if err != nil {
		panic(err)
	}

	for _, child := range children {
		fmt.Println(child)
	}

}
