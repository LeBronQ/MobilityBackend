package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"bytes"
	"encoding/json"
	"github.com/LeBronQ/Mobility"

	consulapi "github.com/hashicorp/consul/api"
)

func Discovery(serviceName string) []*consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("consul client error: %v", err)
	}
	service, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		fmt.Printf("consul client get serviceIp error: %v", err)
	}
	return service
}

type ReqParams struct {
	Node    Mobility.Node    `json:"node"`
}


func main() {
	se := Discovery("Default_Mobility")
	port := se[0].Service.Port
	address := se[0].Service.Address
	request := "http://" + address + ":" + strconv.Itoa(port) + "/mobility"
	node := Mobility.Node{
		Pos:  Mobility.Nbox.RandomPosition3D(),
		Time: 10,
		V: Mobility.Speed{
			X: 10., Y: 10., Z: 10.,
		},
		Model: "RandomWalk",
		Param: Mobility.RandomWalkParam{
			MinSpeed: 0,
			MaxSpeed: 20,
		},
	}
	param := ReqParams{
		Node: node,
	}
	fmt.Println(param)
	jsonData, err := json.Marshal(param)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	

	requestBody := bytes.NewBuffer(jsonData)


	req, err := http.NewRequest("POST", request, requestBody)
    if err != nil {
        fmt.Println(err)
        return
    }
 
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", requestBody.Len()))
 
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(body))
}

