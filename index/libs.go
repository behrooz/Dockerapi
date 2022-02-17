package main

import (
	"encoding/json"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"net/http"
)

func handleBody(r *http.Request) map[string]json.RawMessage {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var objmap map[string]json.RawMessage

	json.Unmarshal(body, &objmap)

	return objmap
}

func clientApi() *client.Client {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return cli
}
