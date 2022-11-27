package main

import (
	"fmt"
	"os"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	time.Sleep(20 * time.Second)

	fmt.Println(os.Getenv("NAME"), "container started!")

	conn, _, err := zk.Connect([]string{"zookeeper"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}

	children, _, _, err := conn.ChildrenW("/")
	if err != nil {
		panic(err)
	}

	if !contains(children, "servers") {
		fmt.Println("Creating servers path...")
		serverPath := "/servers"
		if p, err := conn.Create(serverPath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			panic(err)
		} else if p != serverPath {
			fmt.Printf("Create returned different path '%s' != '%s'", p, serverPath)
		} else {
			createNode(conn)
		}
	} else {
		createNode(conn)
	}

}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func createNode(conn *zk.Conn) {
	path := "/servers/" + os.Getenv("NAME")

	if p, err := conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
		panic(err)
	} else if p != path {
		fmt.Printf("Create returned different path '%s' != '%s'", p, path)
	}
}
