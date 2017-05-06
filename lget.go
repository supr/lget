package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Artist struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LovedTrack struct {
	Artist Artist `json:"artist"`
	Name   string `json:"name"`
	Url    string `json:"url"`
}

type Attr struct {
	Page    string `json:"page"`
	PerPage string `json:"perPage"`
	Total   string `json:"total"`
}

type LovedTracks struct {
	Attr   Attr         `json:"@attr"`
	Tracks []LovedTrack `json:"track"`
}

type LovedTrackResponse struct {
	LovedTracks LovedTracks `json:"lovedtracks"`
}

const API_KEY = "BLAH"

func getLovedTracks(user string) ([]LovedTrack, error) {
	u := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=user.getlovedtracks&user=%s&api_key=%s&format=json", user, API_KEY)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("invalid response from the server")
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lovedTracks := &LovedTrackResponse{}
	err = json.Unmarshal(data, &lovedTracks)
	if err != nil {
		return nil, err
	}
	return lovedTracks.LovedTracks.Tracks, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: %s <user>", args[0])
	}
	lts, err := getLovedTracks(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	for _, t := range lts {
		fmt.Printf("%s,%s\n", t.Name, t.Url)
	}
}
