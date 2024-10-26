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

type Location struct {
	ID          string
	Name        string
	GeoLocation GeoLocation
	Type        LocationType
	StopType    StopType
}
type Nodes []*Location

// Stop is specific to a leg
type Stop struct {
	Location     *Location
	PreviousStop *Location // if the Type is from then this is nil
	Distance     int       // in KM. but converted to any other unit like miles for display
	TimeTaken    int       // approximate time for travel between previous Node and this one
	Type         StopType
}
type Stops []*Stop

// not sure if we need this
// type NodesCacheType map[string]Node // node from ID
// var NodesCache NodesCacheType = NodesCacheType(make(map[string]Node))

type Leg struct {
	ID        string
	From      *Location
	To        *Location
	Distance  int   // in KM. but converted to any other unit like miles for display
	TimeTaken int   // approximate time for travel between 'From' and 'To' in minutes
	AllStops  Stops // AllStops include From and To
}
type Legs []*Leg

// create a nested map of connections. this is all the legs servicing all the stops in the leg
type Connection struct {
	Leg         *Leg
	Connections Connections // all conncection from all stops in legs, except from 'from' location,
	Parent      *Connection
	FromNode    *Location
	NestedLevel int
}

// type Visited map[string]Legs

type Connections []*Connection

const (
	NOT_PROCESSED          int = 0
	QUEUED_TO_BE_PROCESSED int = 1
	BEING_PROCESSED        int = 2
	PROCESSED              int = 3
)

type visited struct {
	// Locations map[string]bool
	Stop map[string]int
	Leg  map[string]int
}

// master connection/result structure
type JourneyMap struct {
	JourneyStart      *Location
	JourneyEnd        *Location
	JourneyDistance   int         // in KM
	ConnectionTree    Connections // only the level 0 connections
	ResultConnections Connections
	Processed         *visited
	Level             int
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
//
// ============================================================================
func CreateLeg(id string, from, to *Location, distance, timeTaken int, allStops Stops) Leg {
	var _newLeg Leg = Leg{ID: id, From: from, To: to, AllStops: allStops}
	// var _newLeg Leg = Leg{ID: id, From: from, To: to}
	// var _stops []Stop = []Stop{from}
	// var _allStops Stops = []Stops{from}

	// LegsCache.AddLeg(_newLeg)
	MICache.AddLeg(&_newLeg)
	// NodeLegsCache.LinkLegToAllNode(&_newLeg)

	return _newLeg
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
//
// ============================================================================
func CreateStop(nType LocationType, leg *Leg, location, previous *Location) *Stop {
	var _stop Stop = Stop{Location: location, PreviousStop: previous, Type: BOTH_DROP_PICKUP}
	if previous != nil {
		_stop.Distance = location.Distance(previous)
		_stop.Distance = _stop.Distance + (_stop.Distance / DISTANCE_CALC_FACTOR)
		_stop.TimeTaken = _stop.Distance // assuming time is 1 min/KM
	}

	return &_stop
}

func MakeAStop(location, previous *Location) *Stop {
	var _stop Stop = Stop{Location: location, PreviousStop: previous, Type: BOTH_DROP_PICKUP}
	if previous != nil {
		_stop.Distance = location.Distance(previous)
		_stop.Distance = _stop.Distance + (_stop.Distance / DISTANCE_CALC_FACTOR)
		_stop.TimeTaken = _stop.Distance // assuming time is 1 min/KM
	}

	return &_stop
}

// ============================================================================
//
// ============================================================================
func (ss Stops) Addtop(nType LocationType, leg *Leg, node, previousNode *Location, distance, timeTaken int) Stops {
	var _stop Stop = *CreateStop(nType, leg, node, previousNode)
	ss = append(ss, &_stop)
	return ss
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
// First and easiest is to look for direct no change path beween
// 'from' and 'to'
// ============================================================================
func (mic *miCacheType) NoChangePath(from, to Location) Legs {
	var _fromLegs Legs = mic.NodeIdIndex[from.ID]
	var _resultingLegs Legs = Legs{}

	for _, _leg := range _fromLegs {
		var _fromIndex int = -1
		var _toIndex int = -1
		for _index, _stop := range _leg.AllStops {
			if _stop.Location.ID == from.ID {
				_fromIndex = _index
			} else if _stop.Location.ID == to.ID {
				_toIndex = _index
			}
		}
		if (_toIndex > _fromIndex) && (_fromIndex > -1) && (_toIndex > -1) {
			_resultingLegs = append(_resultingLegs, _leg)
		}
	}
	return _resultingLegs
}

type connectingLeg struct {
	FromLeg *Leg
	ToLeg   *Leg
}
type connectingLegs []connectingLeg

// ============================================================================
// create 2 list of Legs.
// fromLegs is the list of all legs that start 'from'
// toLegs is the list of all legs that finish or stop at 'to'
// this finds 2 connecting legs that takes you from A->B
// loop through fromLegs till you get to the intersection
// ============================================================================
func (mic *miCacheType) AllPaths(from, to *Location) connectingLegs {
	var _fromToDistance int = from.Distance(to)
	var _fromLegs Legs = mic.NodeIdIndex[from.ID]
	var _toLegs Legs = mic.NodeIdIndex[to.ID]
	var _changingLegs connectingLegs = make(connectingLegs, 0)

	for _, _fromLeg := range _fromLegs {
		for _, _fromStop := range _fromLeg.AllStops {
			for _, _toLeg := range _toLegs {
				for _, _toStop := range _toLeg.AllStops {
					var _distance int = _fromStop.Location.Distance(to)
					if _distance > _fromToDistance*DISTANCE_DIRECTION_FACTOR {
						// if the distance between this _fromStop to the final dest
						// has increased by more than factor. break away
						// break
					}
					if _toStop.Location.ID == _fromStop.Location.ID {
						// this is stop in the toLeg that matches a stop in _fromLeg stop
						// It is a possible interchange
						var _connectingLeg connectingLeg = connectingLeg{_fromLeg, _toLeg}
						_changingLegs = append(_changingLegs, _connectingLeg)
					}
				}
			}
		}
	}
	return _changingLegs
}

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
func CreateJourneyMap(from, to *Location) (*JourneyMap, error) {
	var _stopsProcessed map[string]int = make(map[string]int)
	var _legsProcessed map[string]int = make(map[string]int)
	var _processed visited = visited{_stopsProcessed, _legsProcessed}
	var _jm JourneyMap = JourneyMap{JourneyStart: from, JourneyEnd: to,
		JourneyDistance: from.Distance(to), Processed: &_processed}

	return &_jm, nil
}

// ============================================================================
// check if this leg visits the final destination
// ============================================================================
func (leg *Leg) LegPassesThroughThisLocation(startingLocation *Location, destination *Location) bool {
	var _startingLocation bool = false

	for _, _stop := range leg.AllStops {
		if _startingLocation == false {
			// loop through till you get to starting location
			if _stop.Location.ID == startingLocation.ID {
				_startingLocation = true
			}
		} else {
			if _stop.Location.ID == destination.ID {
				return true
			}
		}
	}

	return false
}

func (jm *JourneyMap) AddConnectionToFinalResult(connection *Connection) {
	// we have now found a path all the way to destination from this stop
	// make sure we store this stop in the leg
	// add the connection tree to the result
	jm.ResultConnections = append(jm.ResultConnections, connection)
}

func (jm *JourneyMap) ProcessANewConnection(leg *Leg, parent *Connection, stop *Stop, level int) *Connection {
	_connection := Connection{Leg: leg, Parent: parent, FromNode: stop.Location, NestedLevel: level}
	jm.MarkStopAsProcessed(leg, stop)
	return &_connection
}

// ============================================================================
// add all conections departing from a
// once you have found the connections loop through them as well
// ============================================================================
func (jm *JourneyMap) FindConnectingNodes_v4(level int, parent *Connection) {
	// get all legs through this node
	var _legs Legs = MICache.GetLegFromNode(parent.FromNode.ID)

	if level < 4 {
		for _, _leg := range _legs {
			// if leg.To is the same as where we are coming from,
			// then there is no point in processing, we've already been there
			if _leg.To.ID != parent.FromNode.ID {
				// if the leg is fully processed, then skip this leg
				// This is making the assumption that leg.ID is unique and not shared
				//     by the both inbound and outbound of this leg
				if jm.LegFullyProcessed(_leg) == false {
					if level > 0 && _leg.To.Distance(jm.JourneyEnd) > jm.JourneyDistance {
						// if the destination of this leg is taking us further away, then
						// we are moving in the wrong direction so ignore this leg
						jm.Processed.Leg[_leg.ID] = PROCESSED
					} else {
						var _startingLocation bool = false
						for _, _stop := range _leg.AllStops {
							// any stops from this connection to end matches where
							// we want to go then add this to result
							if jm.StopProcessed(_leg, _stop.Location) == NOT_PROCESSED {

								// we have found the legs passing through this connection.
								// for each leg, step through stops till you get to this
								// stop and start there
								if _startingLocation == false {
									if _stop.Location.ID == parent.FromNode.ID {
										_startingLocation = true
									}
								} else {
									// create a connection
									_connection := jm.ProcessANewConnection(_leg, parent, _stop, level)
									if _stop.Location.ID == jm.JourneyEnd.ID {
										parent.Connections = append(parent.Connections, _connection)
										jm.AddConnectionToFinalResult(_connection)
									} else {
										if _stop.Type == BOTH_DROP_PICKUP {
											parent.Connections = append(parent.Connections, _connection)
										}
									}
								}

							}
						}
					}
				}
			}
		}

		// you have now found all the connections from all the legs from the parent node
		// now recursively go through the same process for these new connections
		for _, _connection := range parent.Connections {
			if _connection.FromNode.ID != jm.JourneyEnd.ID {
				jm.FindConnectingNodes_v4(level+1, _connection)
			}
		}
	}
	jm.ConnectionTree = nil
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
		_result = fmt.Sprintf("From %s -> ", conn.FromNode.Name)
	} else {
		_result = fmt.Sprintf("Leg %s to %s -> ", conn.Leg.ID, conn.FromNode.Name)
	}
	if conn.Parent != nil {
		_result = conn.Parent.WalkTheTree() + _result
	}
	return _result
}

// ============================================================================
// has this stop been processed already?
// ============================================================================
func (jm *JourneyMap) StopProcessed(leg *Leg, location *Location) int {
	// var _visited int

	_visited, _ := jm.Processed.Stop[leg.ID+"-"+location.ID]

	return _visited
}

// ============================================================================
// mark this stop as visited
// ============================================================================
func (jm *JourneyMap) MarkStopAsProcessed(leg *Leg, stop *Stop) {
	jm.Processed.Stop[leg.ID+"-"+stop.Location.ID] = PROCESSED
	// jm.Processed.Stop[stop.Location.ID] = PROCESSED
}

// ============================================================================
// mark this leg as visited
// ============================================================================
func (jm *JourneyMap) LegFullyProcessed(leg *Leg) bool {
	for _, _stop := range leg.AllStops {
		if jm.StopProcessed(leg, _stop.Location) != PROCESSED {
			return false
		}
	}

	return true
}

// ============================================================================
// create a leg from locations
// ============================================================================
func MakeALeg(legID string, locationIDs ...*Location) (*Leg, error) {
	var _previousLocation *Location = nil
	var _stops Stops
	var _distance, _time int
	for _, _locationID := range locationIDs {
		_stop := MakeAStop(_locationID, _previousLocation)
		_previousLocation = _locationID
		_distance = _distance + _stop.Distance
		_time = _time + _stop.Distance
		_stops = append(_stops, _stop)
	}
	_res := CreateLeg(legID, locationIDs[0], locationIDs[len(locationIDs)-1], _distance, _time, _stops)

	return &_res, nil
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// add all conections departing from a, and b
// ============================================================================
// func (jm *JourneyMap) FindConnectingNodes_v6(level int, parent *Connection) {
// 	// var _level int = parent.NestedLevel + 1
// 	var _legs Legs = MICache.GetLegFromNode(parent.FromNode.ID)

// 	if level < 4 {
// 		for _, _leg := range _legs {
// 			// if leg.To is the same as where we are coming from,
// 			// then there is no point in processing, we've already been there
// 			if _leg.To.ID != parent.FromNode.ID {
// 				// if the leg is fully processed, then skip this leg
// 				// This is making the assumption that leg.ID is unique and not shared
// 				//     by the both inbound and outbound of this leg
// 				if jm.LegFullyProcessed(_leg) == false {
// 					if level > 0 && _leg.To.Distance(jm.JourneyEnd) > jm.JourneyDistance {
// 						// if the destination of this leg is taking us further away, then
// 						// we are moving in the wrong direction so ignore this leg
// 						jm.Processed.Leg[_leg.ID] = PROCESSED
// 					} else {
// 						var _startingLocation bool = false
// 						for _, _stop := range _leg.AllStops {
// 							// any stops from this connection to end matches where
// 							// we want to go then add this to result
// 							if jm.StopProcessed(_leg, _stop.Location) == NOT_PROCESSED {
// 								if _leg.LegPassesThroughThisLocation(_stop.Location, jm.JourneyEnd) == true {
// 									_connection := jm.ProcessANewConnection(_leg, parent, _stop, level)
// 									_connection.FromNode = jm.JourneyEnd
// 									jm.AddConnectionToFinalResult(_connection)
// 								} else {
// 									// we have found the legs passing through this connection.
// 									// for each leg, step through stops till you get to this
// 									// stop and start there
// 									if _startingLocation == false {
// 										if _stop.Location.ID == parent.FromNode.ID {
// 											_startingLocation = true
// 										}
// 									}
// 									if _startingLocation == true {
// 										// }
// 										// if _startingLocation == true {
// 										// create a connection and add it to the parent
// 										_connection := jm.ProcessANewConnection(_leg, parent, _stop, level)
// 										if _stop.Location.ID == jm.JourneyEnd.ID {
// 											jm.AddConnectionToFinalResult(_connection)
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 		for _, _connection := range parent.Connections {
// 			if _connection.FromNode.ID != jm.JourneyEnd.ID {
// 				jm.FindConnectingNodes_v6(level+1, _connection)
// 			}
// 		}
// 	}
// 	jm.ConnectionTree = nil
// }

// ============================================================================
// Not working properly. it look like v4 is the best option.
// add all conections departing from a, and b
// ============================================================================
// func (jm *JourneyMap) FindConnectingNodes_v5(level int, parent *Connection) {
// 	// var _level int = parent.NestedLevel + 1
// 	var _legs Legs = MICache.GetLegFromNode(parent.FromNode.ID)

// 	if level < 4 {

// 		for _, _leg := range _legs {
// 			if jm.StopProcessed(_leg, parent.FromNode) == NOT_PROCESSED {
// 				if _leg.To.ID != parent.FromNode.ID {
// 					if jm.LegFullyProcessed(_leg) {
// 						if level > 0 && _leg.To.Distance(jm.JourneyEnd) > jm.JourneyDistance {
// 							// if the destination of this leg is taking us further then
// 							// we are moving in the wrong direction so ignore this leg
// 							jm.Processed.Leg[_leg.ID] = PROCESSED
// 						} else {
// 							var _startingLocation bool = false
// 							for _, _stop := range _leg.AllStops {
// 								if jm.StopProcessed(_leg, _stop.Location) == NOT_PROCESSED {
// 									// we have found the legs passing through this connection.
// 									// for each leg, step through stops till you get to this
// 									// stop and start there
// 									if _startingLocation == false {
// 										if _stop.Location.ID == parent.FromNode.ID {
// 											_startingLocation = true
// 										}
// 									} else {
// 										if _stop.Location.ID == jm.JourneyStart.ID {
// 											// this means we are traversing back in the wrong direction
// 											break
// 										} else {
// 											// create a connection and add it to the parent
// 											_connection := Connection{Leg: _leg, Parent: parent, FromNode: _stop.Location, NestedLevel: level}
// 											parent.Connections = append(parent.Connections, &_connection)
// 											jm.MarkStopAsProcessed(_leg, _stop)
// 											if _stop.Location.ID == jm.JourneyEnd.ID {
// 												// we have now found a path all the way to destination from this stop
// 												// make sure we store this stop in the leg
// 												// _connection.FromNode = _stop.Location
// 												// add the connection tree to the result
// 												jm.ResultConnections = append(jm.ResultConnections, &_connection)
// 												// break
// 											} else {
// 												jm.FindConnectingNodes_v5(level+1, &_connection)
// 											}
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 	}
// 	jm.ConnectionTree = nil
// }
