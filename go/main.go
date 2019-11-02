package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grpc-poc/registration"
	"io/ioutil"
	"net/http"
	"time"
)

const numOfIterations int64 = 10
const numOfCalls int64 = 10000

func main() {
	var grpc int64 = 0
	var rest int64 = 0
	var encGrpc int64 = 0
	var compressedGrpc int64 = 0
	var encCompressedGrpc int64 = 0
	var traefilEncGrpc int64 = 0
	var traefilEncCompressedGrpc int64 = 0
	for i := 0; i < int(numOfIterations); i++ {
		grpc += readGrpc(false, false, false, "2302")
		encGrpc += readGrpc(true, false, false, "2303")
		compressedGrpc += readGrpc(false, true, false, "2304")
		encCompressedGrpc += readGrpc(true, true, false, "2305")
		traefilEncGrpc += readGrpc(false, true, true, "80")
		//traefilEncCompressedGrpc += readGrpc(false, true, true, "80")
		rest += readRest()
	}

	fmt.Printf("Averages %d iterations of %d calls - GRPC %d microseconds\n"+
		"Encrypted GRPC %d microseconds\n"+
		"Compressed GRPC %d microseconds\n"+
		"Encrypted compressed GRPC %d microseconds\n"+
		"Traefik GRPC %d microseconds\n"+
		"Traefik compressed GRPC %d microseconds\n"+
		"REST %d microseconds", numOfIterations,
		numOfCalls, grpc/numOfIterations, encGrpc/numOfIterations, compressedGrpc/numOfIterations,
		encCompressedGrpc/numOfIterations, traefilEncGrpc/numOfIterations, traefilEncCompressedGrpc/numOfIterations, rest/numOfIterations)
}

func readRest() int64 {
	url := "http://localhost:8007"
	var sum int64 = 0
	for i := 0; i < int(numOfCalls); i++ {
		start := time.Now()
		doRest(url)
		t := time.Now()
		sum += t.Sub(start).Nanoseconds()
	}

	fmt.Printf("REST - sum %d average per call %d nano, %d micros\n", sum, sum/numOfCalls, sum/numOfCalls/1e3)
	return sum / numOfCalls / 1e3
}

func doRest(url string) {
	jsonBytes, _ := json.Marshal(&registration.JoinToken{NodeId: "abc", ServiceId: "def"})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConnsPerHost: 0,
		MaxIdleConns: 0,
		MaxConnsPerHost: 1,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
}

func readGrpc(encryption bool, compressed bool, traefik bool, port string) int64 {
	client, conn, err := NewClient("localhost:"+port, encryption, compressed)
	defer conn.Close()

	if err != nil {
		panic(err.Error())
	}

	var sum int64 = 0
	for i := 0; i < int(numOfCalls); i++ {
		start := time.Now()
		_, err := client.Join(context.Background(), &registration.JoinToken{NodeId: "abc", ServiceId: "def"})
		t := time.Now()
		sum += t.Sub(start).Nanoseconds()

		if err != nil {
			panic(fmt.Sprintf("err! %+v", err.Error()))
		}
	}

	fmt.Printf("%t %t %t GRPC - sum %d average per call %d nano, %d micros\n",
		encryption, compressed, traefik,
		sum, sum/numOfCalls, sum/numOfCalls/1e3)
	return sum / numOfCalls / 1e3
}
