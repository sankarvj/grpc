package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	pb "github.com/sankarvj/grpc/pb"
	"google.golang.org/grpc"
)

var (
	grpcServerAddress = "localhost:50051"
	restServerAddress = "http://localhost:8080/check"
)

var grpcReqChannel chan bool
var restReqChannel chan bool

func BenchmarkGRPCClient(b *testing.B) {
	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStatusServiceClient(conn)

	//Rest cannot not handle more concurrent requests
	//Rest api will through "socket: too many open files" err
	//So, Restricting to concurrently process upto 50 requests at a time.
	grpcReqChannel = make(chan bool, 10)

	for i := 0; i < b.N; i++ {
		grpcReqChannel <- true
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go callGRPCServer(ctx, client)
	}
	fmt.Println("Closing GRPC")
	close(grpcReqChannel)
}

func callGRPCServer(ctx context.Context, client pb.StatusServiceClient) {
	_, err := client.CheckStatus(ctx, &pb.StatusRequest{Check: "mcr"})
	if err != nil {
		fmt.Println("Error in grpc  --> ", err)
	}
	//fmt.Println("output ", output)
	<-grpcReqChannel
}

func BenchmarkRestClient(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := http.Client{}
	restReqChannel = make(chan bool, 10)
	for i := 0; i < b.N; i++ {
		restReqChannel <- true
		go callRestServer(ctx, client)
	}
	fmt.Println("Closing REST")
	close(restReqChannel)
}

type status struct {
	message string
	code    int
}

func callRestServer(ctx context.Context, client http.Client) {
	output := &status{}
	req, err := http.NewRequest("GET", restServerAddress, nil)
	if err != nil {
		log.Println("error creating request ", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("error executing request ", err)
		return
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading response body ", err)
		return
	}

	err = json.Unmarshal(bytes, output)
	if err != nil {
		log.Println("error unmarshalling response ", err)
		return
	}
	//fmt.Println("output ", output)
	<-restReqChannel
}
