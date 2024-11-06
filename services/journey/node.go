package journey

import (
	"fmt"
	"math"
)

// ============================================================================
// ============================================================================
// all stops are considered a node.
// create Leg based on stops. Add to Cache, map[id]Leg, of Legs
// go through all leg.stops for all legs and
// create cache, map[nodeID] []*Legs, for all stops
// ============================================================================
// ============================================================================

// used for working out if Direction is wrong. eg choosing a leg may be taking
// you in the wrong direction and should not be considered
const DISTANCE_DIRECTION_FACTOR int = 2 // meaning distance + (distance*1/2)

const DISTANCE_CALC_FACTOR int = 4 // meaning distance + (distance*1/4)

const MAX_NESTED_LEVEL int = 3

// ========================================================
// ========================================================
// ========================================================

const (
	UNKOWN_LOCATION LocationType = iota // this makes Unknow as 0
	COACH_STATION
	COACH_AND_BUS_STATION
	BUS_STOP
	TEMP_STOP
)

type GeoLocation struct {
	Lat, Lng float64
}

type StopType int

const (
	UNKNOWN_STOP_TYPE StopType = iota
	DROP_OFF
	PICKUP
	BOTH_DROP_PICKUP
)

type LocationType int

const (
	UNKNOWN_LOCATION_TYPE LocationType = iota
	FROM
	TO
	STOP
)

type Status int

const (
	NOT_PROCESSED          Status = 0
	QUEUED_TO_BE_PROCESSED Status = 1
	BEING_PROCESSED        Status = 2
	PROCESSED              Status = 4
	ON_PATH                Status = 8 // the leg or stop in a leg that takes you to your destination
	PROCESSED_ON_PATH      Status = 12
)

// ========================================================
// ========================================================
// start of journey_2
// ========================================================
// ========================================================
type Location struct {
	ID          string
	Name        string
	GeoLocation GeoLocation
	Type        LocationType
	StopType    StopType
}
type Locations []*Location

type Stops []*Location

type Connection struct {
	ThisStop *Location
	Leg      *Leg
	Previous *Connection
	// Next     *Connection
	Status Status
}

type Connections []*Connection

type Leg struct {
	ID        string
	From      *Location
	To        *Location
	Distance  int   // in KM. but converted to any other unit like miles for display
	TimeTaken int   // approximate time for travel between 'From' and 'To' in minutes
	AllStops  Stops // AllStops include From and To and in order
}
type Legs []*Leg

// master connection/result structure
type JourneyMap struct {
	JourneyStart      *Location
	JourneyEnd        *Location
	JourneyDistance   int // in KM
	ResultConnections Connections
	Processed         *ProcessingStatus
}

type ProcessingStatus struct {
	ConnectionMap       map[string]Status
	LegMap              map[string]Status
	OnPathConnectionMap map[string]*Connection // map[legID-LocationId]Connections
}

// ============================================================================
// ============================================================================
// journey_2
// ============================================================================
// ============================================================================

// ============================================================================
// ============================================================================
// ============================================================================
//
// ============================================================================
// ============================================================================
// create a leg from locations
// ============================================================================
func MakeALeg(legID string, locationIDs ...*Location) (*Leg, error) {
	var _leg Leg = Leg{ID: legID}
	var _lastIndex int = len(locationIDs) - 1
	for _index, _location := range locationIDs {
		if _index == 0 {
			_leg.From = _location
		} else if _index == _lastIndex {
			_leg.To = _location
		}

		_leg.AllStops = append(_leg.AllStops, _location)
	}
	_leg.Distance = _leg.From.Distance(_leg.To)
	_leg.TimeTaken = _leg.Distance
	MICache.AddLeg(&_leg)

	return &_leg, nil
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// Make a connection
// ============================================================================
func (jm *JourneyMap) MakeAConnection(leg *Leg, this *Location, previous *Connection, status Status) *Connection {
	var _connection Connection = Connection{ThisStop: this, Leg: leg, Previous: previous, Status: status}
	return &_connection
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// search for path through nodes
// ============================================================================

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
func (n *Location) Distance(to *Location) int {
	var _distance float64
	_distance = distance(n.GeoLocation.Lat, n.GeoLocation.Lng, to.GeoLocation.Lat, to.GeoLocation.Lng, "K")

	return int(_distance * 1.25)
}

// ============================================================================
// distance between 2 GPS locations
// taken from:  https://www.geodatasource.com/developers/go
// ============================================================================
// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// :::                                                                         :::
// :::  This routine calculates the distance between two points (given the     :::
// :::  latitude/longitude of those points). It is being used to calculate     :::
// :::  the distance between two locations using GeoDataSource (TM) products   :::
// :::                                                                         :::
// :::  Definitions:                                                           :::
// :::    South latitudes are negative, east longitudes are positive           :::
// :::                                                                         :::
// :::  Passed to function:                                                    :::
// :::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
// :::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
// :::    unit = the unit you desire for results                               :::
// :::           where: 'M' is statute miles (default)                         :::
// :::                  'K' is kilometers                                      :::
// :::                  'N' is nautical miles                                  :::
// :::                                                                         :::
// :::  Worldwide cities and other features databases with latitude longitude  :::
// :::  are available at https://www.geodatasource.com                         :::
// :::                                                                         :::
// :::  For enquiries, please contact sales@geodatasource.com                  :::
// :::                                                                         :::
// :::  Official Web site: https://www.geodatasource.com                       :::
// :::                                                                         :::
// :::               GeoDataSource.com (C) All Rights Reserved 2022            :::
// :::                                                                         :::
// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// instance of JourneyMap for a new search
// ============================================================================
// add all conections departing from a, and b
// ============================================================================
func InitialiseJourneyMap(from, to *Location) (*JourneyMap, error) {
	var _jm JourneyMap = JourneyMap{JourneyStart: from, JourneyEnd: to,
		JourneyDistance: from.Distance(to)}

	var _stopsProcessed map[string]Status = make(map[string]Status)
	var _legsProcessed map[string]Status = make(map[string]Status)
	var _connectionMap map[string]*Connection = make(map[string]*Connection)
	var _processed ProcessingStatus = ProcessingStatus{_stopsProcessed, _legsProcessed, _connectionMap}

	_jm.Processed = &_processed
	return &_jm, nil
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// OnPathConnections is used to map all connections with the status of ON_PATH
// if a connection exists in this map then it is ON_PATH
// ============================================================================
func (jm *JourneyMap) getConnectionFromMap(leg *Leg, locationID string) *Connection {
	_connection := jm.Processed.OnPathConnectionMap[leg.ID+"-"+locationID]
	return _connection
}

// ============================================================================
//
// ============================================================================
func (jm *JourneyMap) setConnectionToMap(leg *Leg, connection *Connection) {
	jm.Processed.OnPathConnectionMap[leg.ID+"-"+connection.ThisStop.ID] = connection
}

// ============================================================================
//
// ============================================================================
func (jm *JourneyMap) setLegStatus(leg *Leg, status Status) {
	jm.Processed.LegMap[leg.ID] = status
}

// ============================================================================
//
// ============================================================================
func (jm *JourneyMap) getLegStatus(leg *Leg) Status {
	return jm.Processed.LegMap[leg.ID]
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
func MakeALocation(id, name string, geo GeoLocation, locType LocationType,
	stopType StopType) *Location {
	var _location Location = Location{ID: id, Name: name, GeoLocation: geo,
		Type: locType, StopType: stopType}

	return &_location
}

// ============================================================================
// check if this leg visits the a location. eg used to check if the desired
// location, ie final destination, is part of this journey leg
// ============================================================================
func (leg *Leg) LegPassesThroughThisLocation(startingLocation *Location, destination *Location) bool {
	var _startingLocation bool = false

	for _, _stop := range leg.AllStops {
		if _startingLocation == false {
			// loop through till you get to starting location
			if _stop.ID == startingLocation.ID {
				_startingLocation = true
			}
		} else {
			if _stop.ID == destination.ID {
				return true
			}
		}
	}

	return false
}

// ============================================================================
//
// ============================================================================
func (jm *JourneyMap) AddConnectionToFinalResult(connection *Connection) {
	// we have now found a path all the way to destination from this stop
	// make sure we store this stop in the leg
	// add the connection tree to the result
	jm.ResultConnections = append(jm.ResultConnections, connection)
	jm.SetConnectionTrailStatus(connection, ON_PATH)
}

// ============================================================================
// BINGO: does the leg stop at what is our final destination
// ============================================================================
func (jm *JourneyMap) LegStopsAtOurDestination(leg *Leg, startingStop *Location) *Location {
	var _location *Location = nil
	var _startingLocation bool = false
	for _, _stop := range leg.AllStops {
		if _startingLocation == false {
			if _stop.ID == startingStop.ID {
				_startingLocation = true
			}
		} else {
		}
		if _stop.ID == jm.JourneyEnd.ID {
			_location = _stop
			break
		}
	}
	return _location
}

// ============================================================================
// found a connection that ends at our desired destination.
// ============================================================================
func (jm *JourneyMap) FoundAConnectionThrough(leg *Leg, to *Location, parent *Connection) {
	_newConnection := jm.MakeAConnection(leg, to, parent, ON_PATH)
	jm.SetConnectionTrailStatus(_newConnection, ON_PATH)
	jm.AddConnectionToFinalResult(_newConnection)

}

// ============================================================================
// follow the link back to start and set every connection status to
// ON_PATH
// ============================================================================
func (jm *JourneyMap) SetConnectionTrailStatus(conn *Connection, status Status) {
	if conn != nil {
		jm.SetConnectionProcessingStatus(conn.Leg, conn, status)
		if conn.Previous != nil {
			jm.SetConnectionTrailStatus(conn.Previous, status)
		}
	}
}

// ============================================================================
//
// ============================================================================
func (c *Connection) MakeACopy() *Connection {
	_connection := Connection{c.ThisStop, c.Leg, c.Previous, c.Status}
	return &_connection
}

// ============================================================================
//
// ============================================================================
// func (jm *JourneyMap) copyLinksForward(originalConnection, copyConnection *Connection) {

// 	if originalConnection.Next != nil {
// 		_copy := originalConnection.Next.MakeACopy()
// 		jm.copyLinksForward(originalConnection.Next, _copy)
// 	} else {
// 		jm.ResultConnections = append(jm.ResultConnections, copyConnection)

// 	}

// }

// ============================================================================
// we have come across a leg that is marked as ON_PATH.
// so find the connection and copy the ON_PATH connections to form a new
// set of linked lists
// ============================================================================
// func (jm *JourneyMap) linkToAnExistingTree(leg *Leg, from *Connection) bool {
// 	var _startingLocation bool = false
// 	for _, _stop := range leg.AllStops {
// 		// loop through till you get to current location
// 		if _startingLocation == false {
// 			// loop through stops in this leg this you get to the parent
// 			if _stop.ID == from.ThisStop.ID {
// 				_startingLocation = true
// 			}
// 		} else {
// 			// now look for *connection with ON_PATH. once you find it
// 			// you don't need any more processing on this leg as you
// 			// have found a path to destination. walk through all the
// 			// connections to the final destination and copy them to
// 			// form a new linked list
// 			_connectionOnPath := jm.getConnectionFromMap(leg, _stop.ID)
// 			if _connectionOnPath != nil {
// 				_newConnection := jm.MakeAConnection(leg, _stop, from, nil, ON_PATH)
// 				from.Next = _newConnection
// 				jm.copyLinksForward(_connectionOnPath, _newConnection)
// 				break
// 			}

// 		}
// 	}
// 	return false
// }

// ============================================================================
// add all conections departing from a
// once you have found the connections loop through them as well
// Get legs for this location
// For each leg
//     Create a connection with leg
//     If (legStatus == ON_PATH)
//     {
//         Get [legID+locationID]connection
//         Walk forward
//             Copy each node
//             Add newConnection to result
//             Walk back through the newConnection and Mark all connections and Paths as ON_PATH
//     }
//     Else
//         Does leg meets the destination
//         If yes
//             Create a newConnection for finalDestinatio
//             update connection and new connection with previous, next, …
//             Add newConnection to result
//             Walk back through the newConnection and Mark all connections and Paths as ON_PATH
//     Else
//         For each stop in connection
//             Create a newConnection
//             call findPath (newConnection)
//             If newConnection.status != ON_PATH mark it as processed

// ============================================================================
func (jm *JourneyMap) FindConnectingNodes_v5(level int, parent *Connection) {
	// get all legs through this node
	var _legs Legs = MICache.GetLegFromNode(parent.ThisStop.ID)

	if level < 4 {
		for _, _leg := range _legs {
			legStatus := jm.getLegStatus(_leg)
			if legStatus == PROCESSED || legStatus == BEING_PROCESSED {

			} else {
				// if jm.LegProcessingStatus(_leg) == ON_PATH {
				// 	// link to an existing ON_PATH
				// 	// jm.linkToAnExistingTree(_leg, parent)
				// } else {
				if jm.LegProcessingStatus(_leg) == NOT_PROCESSED {
					jm.SetLegProcessingStatus(_leg, BEING_PROCESSED)

					var _startingLocation bool = false
					for _, _stop := range _leg.AllStops {
						if _startingLocation == false {
							// loop through stops in this leg this you get to the parent
							if _stop.ID == parent.ThisStop.ID {
								_startingLocation = true
							}
						} else {
							// create a connection and link it to parent
							_newConnection := jm.MakeAConnection(_leg, _stop, parent, BEING_PROCESSED)
							// check to see if this leg gets you to the final destination
							var _destinationLocation *Location = jm.LegStopsAtOurDestination(_leg, _stop)
							if _destinationLocation != nil {
								// if so mark the trail as ON_PATH
								jm.FoundAConnectionThrough(_leg, _destinationLocation, _newConnection)
							} else {
								if _stop.ID != parent.ThisStop.ID {
									jm.FindConnectingNodes_v5(level+1, _newConnection)
									if jm.StopProcessingStatus(_leg, _newConnection.ThisStop) != ON_PATH {
										jm.SetConnectionProcessingStatus(_leg, _newConnection, PROCESSED)
									}
								}
							}
						}
					}
					// }
				}
			}
		}

		// jm.ConnectionTree = nil
	}
}

// ============================================================================
// show resulting journies
// ============================================================================
func (jm *JourneyMap) ShowResultingLegs() {
	for _, _connection := range jm.ResultConnections {
		_result := _connection.WalkTheTree()
		fmt.Printf("*** %s\n", _result)
	}
}

// ============================================================================
// show resulting for one plan/route
// ============================================================================
func (conn *Connection) WalkTheTree() string {
	var _result string
	if conn.Leg == nil {
		_result = fmt.Sprintf("From %s -> ", conn.ThisStop.Name)
	} else {
		_result = fmt.Sprintf("Leg %s to %s -> ", conn.Leg.ID, conn.ThisStop.Name)
	}
	if conn.Previous != nil {
		_result = conn.Previous.WalkTheTree() + _result
	}
	return _result
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// mark this leg as
// ============================================================================
func (jm *JourneyMap) SetLegProcessingStatus(leg *Leg, status Status) {
	jm.Processed.LegMap[leg.ID] = status
}

// ============================================================================
// return leg processing status
// ============================================================================
func (jm *JourneyMap) LegProcessingStatus(leg *Leg) Status {
	return jm.Processed.LegMap[leg.ID]
}

// ============================================================================
// mark this leg as visited
// ============================================================================
// func (jm *JourneyMap_2) LegFullyProcessed(leg *Leg_2) bool {
// 	for _, _stop := range leg.AllStops {
// 		if jm.StopProcessingStatus(leg, _stop.Location) != PROCESSED {
// 			return false
// 		}
// 	}

// 	return true
// }

// ============================================================================
// has this stop been processed already?
// ============================================================================
func (jm *JourneyMap) StopProcessingStatus(leg *Leg, location *Location) Status {
	_status, _ := jm.Processed.ConnectionMap[leg.ID+"-"+location.ID]

	return _status
}

// ============================================================================
// mark this stop as visited
// ============================================================================
func (jm *JourneyMap) SetConnectionProcessingStatus(leg *Leg, connection *Connection, status Status) {
	if leg == nil {
		jm.Processed.ConnectionMap["-"+connection.ThisStop.ID] = status
	} else {
		jm.Processed.ConnectionMap[leg.ID+"-"+connection.ThisStop.ID] = status
	}
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
