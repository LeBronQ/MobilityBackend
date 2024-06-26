package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"bytes"
	"encoding/json"
	"github.com/LeBronQ/Mobility"
	"testing"

)
/*
//API test
func Benchmark(b *testing.B) {
	b.StopTimer()
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
	
	jsonData, err := json.Marshal(param)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	b.StartTimer()

	for i := 0; i < b.N; i++{
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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	}
	//fmt.Println("Response:", string(body))
}
*/

const (
	n        = 1000
	minSpeed = 10.0
	maxSpeed = 100.0
)

var arr []Mobility.Node

func GenerateNodes() {
	for i := 1; i <= n; i++ {
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
		arr = append(arr, node)
	}
}
func Benchmark_Nodes(b *testing.B) {
	b.StopTimer()
	se := Discovery("Default_Mobility")
	port := se[0].Service.Port
	address := se[0].Service.Address
	request := "http://" + address + ":" + strconv.Itoa(port) + "/mobility"
	GenerateNodes()
	b.StartTimer()
	for i := 0; i < b.N; i++{
		for _, n := range arr{
			param := ReqParams{
				Node: n,
			}
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
		
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
		}
	}
	//fmt.Println("Response:", string(body))
}