// ============================================================================
// ============================================================================
// Request:
// [
//         {
//             "locationId": "42205",
//             "stopId": "A",
//             "date": "2022-01-20 00:00:00"
//         }
// ]
// _response:
// [
//     {
//         "locationId": "42205",
//         "stopId": "A",
//         "date": "2022-01-20 00:00:00",
//         "locationName": "Taunton",
//         "additionalName": "",
//         "stopName": "Park Street, opposite County Hall"
//     }
// ]
// ============================================================================
// "cp_Get_Location_Description":
//                 locationCode      IN
//                 stop              IN
//                 travelDate        IN
//                 Location_Name     OUT
//                 Additional_Name   OUT
//                 Stop_Name         OUT
// ============================================================================
// ============================================================================

package stopLocation

// import (
// 	"database/sql"

// 	nl "github.com/foadmom/common/logger"
// 	ms "github.com/foadmom/common/sql"
// )

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
// This is the actual processing of the request. The comms will pass a json
// to this function and sends the returned string back through comms channel.
// for every separate nch.Init you need one of the functions below
// ============================================================================
// func ProcessStopLocationRequest(input string) (string, error) {
// 	nl.Instance().Debug().Msg("Entering processStopLocationRequest(" + input + ")")
// 	var _response string
// 	var _err error
// 	_response, _err = getStopLocationData(input)

// 	defer nl.Instance().Debug().Msg("exiting  processStopLocationRequest ()" + _response)

// 	// _response = `{"location":"Birminghad", "stop_location":"somewhere"}`
// 	return _response, _err
// }

// ============================================================================
// This is the actual processing of the request. The comms will pass a json
// to this function and sends the returned string back through comms channel.
// for every separate nch.Init you need one of the functions below
// var localServer ms.DBServer = ms.DBServer{"local test server", "localhost", "1433",
//
//	"sa", "Pa55w0rd", "NEX", "msserver"}
// var localServer ms.InstanceData = ms.InstanceData{"local test server", "localhost", "1433",
// 	"sa", "Pa55w0rd", "NEX", "msserver"}
// var localMSSqlServer ms.MsDbServer = ms.MsDbServer{Server: localServer}
// var procName string = "seats.ja_cp_Get_Location_Description"

// var inputJson string = `[{"locationId": "42205", "stopId": "A", "date": "2022-01-20 00:00:00"}]`
// ============================================================================
// func getStopLocationData(inputJson string) (string, error) {
// 	var _conn *sql.DB
// 	var _err error
// 	var _outputJson string

// 	localMSSqlServer.Init()

// 	_conn, _err = localMSSqlServer.Connection()

// 	if _err == nil {
// 		_outputJson, _err = localMSSqlServer.CallStoredProc(_conn, procName, inputJson)
// 	}
// 	if _err != nil {
// 		nl.Instance().Error().Msg("error getting json payload from Body. Error = " + _err.Error())
// 	}
// 	return _outputJson, _err
// }
