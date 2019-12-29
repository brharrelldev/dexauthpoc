package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	demo "github.com/brharrelldev/dexauthpoc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"
)

type auth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
}

func authorize(authService string) {

}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	c, err := grpc.DialContext(ctx, ":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	auth := auth{}

	body, err := json.Marshal(auth)
	if err != nil{
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:5555/login", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(resp.Request.URL.String())


	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+"test")

	d := demo.NewDemoServiceClient(c)

	r, err := d.DemoService(ctx, &demo.DemoRequest{Message: "test"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)

}
