package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	time.Sleep(20 * time.Second)

	fmt.Println(os.Getenv("NAME"), "container started!")

	conn, _, err := zk.Connect([]string{"zookeeper"}, time.Second)
	if err != nil {
		panic(err)
	}
	_, err = conn.Create("/servers", []byte{}, 0, zk.WorldACL(zk.PermAll))
	_, err = conn.Create("/servers/"+os.Getenv("NAME"), nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	} else {
		fmt.Println(os.Getenv("name"), "node created")
	}

	http.HandleFunc("/", handler)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received in ", os.Getenv("NAME"))

	switch r.Method {
	case "POST":
		unencodedJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		var unencodedRows RowsType
		json.Unmarshal(unencodedJSON, &unencodedRows)

		encodedRows := unencodedRows.encode()

		encodedJSON, _ := json.Marshal(encodedRows)

		println("unencoded:", string(unencodedJSON))
		println("encoded:", string(encodedJSON))

		client := &http.Client{}

		key := unencodedRows.Row[0].Key

		req, err := http.NewRequest(http.MethodPut, "http://hbase:8080/se2:library/"+key, bytes.NewBuffer(encodedJSON))

		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		fmt.Println(resp.StatusCode)

		if resp.StatusCode != 200 {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

	case "GET":
		fmt.Println("Get request received")

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPut, "http://hbase:8080/se2:library/scanner/", bytes.NewBuffer([]byte(`<Scanner batch="10"/>`)))

		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

		req.Header.Set("Accept", "text/plain")
		req.Header.Set("Content-Type", "text/xml")
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}
		scanner, _ := resp.Location()

		req, err = http.NewRequest(http.MethodGet, scanner.String(), nil)
		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

		req.Header.Set("Accept", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			log.Fatalln(err)
		}

		var encodedRows EncRowsType
		json.Unmarshal(responseBody, &encodedRows)

		decodedRows, err := encodedRows.decode()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(decodedRows)

	}
}
