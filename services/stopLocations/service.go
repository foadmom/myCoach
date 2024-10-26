package main

import (
	"sync"

	nch "github.com/foadmom/common/http"
)

type locationInput struct {
	LocationId string `json:"locationId"`
	StopId     string `json:"stopId"`
	Date       string `json:"date"`
}

type locationOutput struct {
	LocationId     string `json:"locationId"`
	StopId         string `json:"stopId"`
	Date           string `json:"date"`
	LocationName   string `json:"locationName"`
	AdditionalName string `json:"additionalName"`
	StopName       string `json:"stopName"`
}

// ============================================================================
// ============================================================================
func main() {
	var _wg sync.WaitGroup

	// ====================================================================
	// for every different service or url you need to repeat the following
	// 2 lines and write a func for processing the message
	_wg.Add(1)
	go nch.Init("", "8090", "/stopLocations", &_wg, processStopLocationRequest)
	// ====================================================================

	_wg.Wait()
}

// ============================================================================
// This is the actual processing of the request. The comms will pass a json
// to this function and sends the returned string back through comms channel.
// for every separate nch.Init you need one of the functions below
// ============================================================================
func processStopLocationRequest(input string) (string, error) {
	var _dummyResp string = `{"location":"Birminghad", "stop_location":"somewhere"}`
	return _dummyResp, nil
}
