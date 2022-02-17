package main

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"net/http"
)

func imagesList(w http.ResponseWriter, r *http.Request) {

	containers, err := clientApi().ImageList(context.Background(), types.ImageListOptions{All: true})
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
