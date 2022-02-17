package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"net/http"
	"strconv"
	"strings"
)

type postBody struct {
	Name        string `json:"Name"`
	HostPort    int    `json:"HostPort"`
	ExposedPort int    `json:"ExposedPort"`
	ImageName   string `json:"ImageName"`
}

func containerList(w http.ResponseWriter, r *http.Request) {

	containers, err := clientApi().ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	jData, err := json.Marshal(containers)
	if err != nil {
		panic(err)
	}

	w.Write(jData)
}

func containerStats(w http.ResponseWriter, r *http.Request) {

	params := handleBody(r)
	containerId := params["containerId"]
	containerId = json.RawMessage(strings.ReplaceAll(string(containerId), "\"", ""))

	stats, err := clientApi().ContainerStats(context.Background(), string(containerId), false)
	if err != nil {
		panic(err)
	}
	jData, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	w.Write(jData)

}

func containerStop(w http.ResponseWriter, r *http.Request) {

	params := handleBody(r)
	containerId := params["containerId"]
	containerId = json.RawMessage(strings.ReplaceAll(string(containerId), "\"", ""))

	err := clientApi().ContainerStop(context.Background(), string(containerId), nil)

	if err != nil {
		panic(err)
	}
	jData, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}
	w.Write(jData)

}

func containerStart(w http.ResponseWriter, r *http.Request) {

	params := handleBody(r)
	containerId := params["containerId"]
	containerId = json.RawMessage(strings.ReplaceAll(string(containerId), "\"", ""))

	err := clientApi().ContainerStart(context.Background(), string(containerId), types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}
	jData, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}
	w.Write(jData)
}

func containerRemove(w http.ResponseWriter, r *http.Request) {

	params := handleBody(r)
	containerId := params["containerId"]
	containerId = json.RawMessage(strings.ReplaceAll(string(containerId), "\"", ""))

	err := clientApi().ContainerRemove(context.Background(), string(containerId), types.ContainerRemoveOptions{})

	if err != nil {
		panic(err)
	}
	jData, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}
	w.Write(jData)
}

func containerRestart(w http.ResponseWriter, r *http.Request) {

	params := handleBody(r)
	containerId := params["containerId"]
	containerId = json.RawMessage(strings.ReplaceAll(string(containerId), "\"", ""))

	err := clientApi().ContainerRestart(context.Background(), string(containerId), nil)

	if err != nil {
		panic(err)
	}
	jData, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}
	w.Write(jData)
}

func containerCreate(w http.ResponseWriter, r *http.Request) {

	//params := handleBody(r)
	//containerName := json.RawMessage(strings.ReplaceAll(string(params["containerName"]), "\"", ""))
	//containerPort := json.RawMessage(strings.ReplaceAll(string(params["containerPort"]), "\"", ""))
	//cpu := json.RawMessage(strings.ReplaceAll(string(params["cpu"]), "\"", ""))
	//hostPort := json.RawMessage(strings.ReplaceAll(string(params["hostPort"]), "\"", ""))
	//memory := json.RawMessage(strings.ReplaceAll(string(params["memory"]), "\"", ""))
	//imageName := json.RawMessage(strings.ReplaceAll(string(params["imageName"]), "\"", ""))

	var postBody postBody
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&postBody)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	newport, err := nat.NewPort("tcp", strconv.Itoa(postBody.HostPort))
	if err != nil {
		errorResponse(w, "Port is not available", http.StatusInternalServerError)
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			newport: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(postBody.HostPort),
				},
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		LogConfig:     container.LogConfig{Type: "json-file", Config: map[string]string{}}}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	gatewayConfig := &network.EndpointSettings{Gateway: "gatewayname"}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	config := &container.Config{Image: postBody.ImageName, ExposedPorts: map[nat.Port]struct{}{
		newport: struct{}{},
	},
	}

	cont, err := clientApi().ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, string(postBody.Name))

	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	clientApi().ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	jData, err := json.Marshal(cont)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jData)
	errorResponse(w, string(jData), http.StatusOK)

}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
