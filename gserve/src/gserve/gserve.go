package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	time.Sleep(15 * time.Second)

	fmt.Println(os.Getenv("NAME"), "container started!")

	//connect to zookeeper
	conn, _, err := zk.Connect([]string{"zookeeper"}, 5*time.Second)
	if err != nil {
		panic(err)
	}

	//create servers path/node
	_, err = conn.Create("/servers", []byte{}, 0, zk.WorldACL(zk.PermAll))

	//create node for gserve
	_, err = conn.Create("/servers/"+os.Getenv("NAME"), nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getenv("NAME"), "node created")

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
		err = json.Unmarshal(unencodedJSON, &unencodedRows)

		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}

		encodedRows := unencodedRows.encode()

		encodedJSON, err := json.Marshal(encodedRows)

		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}

		if unencodedRows.Row == nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}

		client := &http.Client{}

		key := unencodedRows.Row[0].Key

		req, err := http.NewRequest(http.MethodPut, "http://hbase:8080/se2:library/"+key, bytes.NewBuffer(encodedJSON))

		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		fmt.Println(resp.StatusCode)

		if resp.StatusCode != 200 {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}

	case "GET":
		fmt.Println("Get request received")

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPut, "http://hbase:8080/se2:library/scanner/", bytes.NewBuffer([]byte(`<Scanner batch="10"/>`)))

		if err != nil {
			http.Error(w, "500 something went wrong", http.StatusInternalServerError)
			return
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
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Accept", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		var encodedRows EncRowsType
		json.Unmarshal(responseBody, &encodedRows)

		decodedRows, err := encodedRows.decode()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(decodedRows)

		tpl, err := template.New("template.html").Funcs(template.FuncMap{
			"isDocument":    isDocument,
			"isMetadata":    isMetadata,
			"getCleanValue": getCleanValue,
			"getServerName": getServerName,
		}).ParseFiles("template.html")
		if err != nil {
			fmt.Println(err)
		}
		tpl.Execute(w, decodedRows)

	}
}

func isDocument(value string) bool {
	if strings.HasPrefix(value, "document:") {
		return true
	} else {
		return false
	}
}

func isMetadata(value string) bool {
	if strings.HasPrefix(value, "metadata:") {
		return true
	} else {
		return false
	}
}

func getCleanValue(value string) string {
	index := strings.Index(value, ":")
	cleanVal := value[(index + 1):]
	return cleanVal
}

func getServerName() string {
	return os.Getenv("NAME")
}
