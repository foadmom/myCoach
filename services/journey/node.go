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
type Location_2 struct {
	ID          string
	Name        string
	GeoLocation GeoLocation
	Type        LocationType
	StopType    StopType
}
type Locations []*Location_2

type Stops_2 []*Location_2

type Connection_2 struct {
	ThisStop *Location_2
	Leg      *Leg_2
	Previous *Connection_2
	Next     *Connection_2
	Status   Status
}

type Connections_2 []*Connection_2

type Leg_2 struct {
	ID        string
	From      *Location_2
	To        *Location_2
	Distance  int     // in KM. but converted to any other unit like miles for display
	TimeTaken int     // approximate time for travel between 'From' and 'To' in minutes
	AllStops  Stops_2 // AllStops include From and To and in order
}
type Legs_2 []*Leg_2

// master connection/result structure
type JourneyMap_2 struct {
	JourneyStart      *Location_2
	JourneyEnd        *Location_2
	JourneyDistance   int // in KM
	ResultConnections Connections_2
	Processed         *visited_2
}

type visited_2 struct {
	connection map[string]Status
	Leg        map[string]Status
}

// ============================================================================
// ============================================================================
// journey_2
// ============================================================================
// ============================================================================

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
func MakeALocation(id, name string, geo GeoLocation, locType LocationType,
	stopType StopType) *Location_2 {
	var _location Location_2 = Location_2{ID: id, Name: name, GeoLocation: geo,
		Type: locType, StopType: stopType}

	return &_location
}

// ============================================================================
// ============================================================================
// ============================================================================
//
// ============================================================================
// ============================================================================
// create a leg from locations
// ============================================================================
func MakeALeg(legID string, locationIDs ...*Location_2) (*Leg_2, error) {
	var _leg Leg_2 = Leg_2{ID: legID}
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
func MakeAConnection(leg *Leg_2, this *Location_2, previous,
	next *Connection_2, status Status) *Connection_2 {
	var _connection Connection_2 = Connection_2{ThisStop: this, Leg: leg,
		Previous: previous, Next: next, Status: status}

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
func (n *Location_2) Distance(to *Location_2) int {
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
func InitialiseJourneyMap(from, to *Location_2) (*JourneyMap_2, error) {
	var _stopsProcessed map[string]Status = make(map[string]Status)
	var _legsProcessed map[string]Status = make(map[string]Status)
	// var _processed visited = visited{_stopsProcessed, _legsProcessed}

	var _processed visited_2 = visited_2{_stopsProcessed, _legsProcessed}
	var _jm JourneyMap_2 = JourneyMap_2{JourneyStart: from, JourneyEnd: to,
		JourneyDistance: from.Distance(to), Processed: &_processed}

	return &_jm, nil
}

// ============================================================================
// check if this leg visits the final destination
// ============================================================================
func (leg *Leg_2) LegPassesThroughThisLocation(startingLocation *Location_2, destination *Location_2) bool {
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
func (jm *JourneyMap_2) AddConnectionToFinalResult(connection *Connection_2) {
	// we have now found a path all the way to destination from this stop
	// make sure we store this stop in the leg
	// add the connection tree to the result
	jm.ResultConnections = append(jm.ResultConnections, connection)
	jm.SetConnectionTrailStatus(connection, ON_PATH)
}

// ============================================================================
// BINGO: does the leg stop at what is our final destination
// ============================================================================
func (jm *JourneyMap_2) LegStopsAtOurDestination(leg *Leg_2, startingStop *Location_2) *Location_2 {
	var _location *Location_2 = nil
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
//
// ============================================================================
func (jm *JourneyMap_2) FoundAConnectionThrough_2(leg *Leg_2, stop *Location_2, parent *Connection_2) {
	_newConnection := MakeAConnection(leg, stop, parent, nil, ON_PATH)
	jm.SetConnectionTrailStatus(_newConnection, ON_PATH)
	jm.AddConnectionToFinalResult(_newConnection)
}

// ============================================================================
// follow the link back to start and set every connection status to
// ON_PATH
// ============================================================================
func (jm *JourneyMap_2) SetConnectionTrailStatus(conn *Connection_2, status Status) {
	if conn != nil {
		jm.SetConnectionProcessingStatus(conn.Leg, conn, status)
		if conn.Previous != nil {
			jm.SetConnectionTrailStatus(conn.Previous, status)
		}
	}
}

// ============================================================================
// add all conections departing from a
// once you have found the connections loop through them as well
// ============================================================================
func (jm *JourneyMap_2) FindConnectingNodes_v5(level int, parent *Connection_2) {
	// get all legs through this node
	var _legs Legs_2 = MICache.GetLegFromNode(parent.ThisStop.ID)

	if level < 4 {
		for _, _leg := range _legs {
			if jm.LegProcessingStatus(_leg) == PROCESSED {
				break
			}
			if jm.LegProcessingStatus(_leg) == BEING_PROCESSED {
				break
			}
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
						// check to see if this leg gets you to the final destination
						var _destinationLocation *Location_2 = jm.LegStopsAtOurDestination(_leg, _stop)
						if _destinationLocation != nil {
							// if so mark the trail as ON_PATH
							jm.FoundAConnectionThrough_2(_leg, _stop, parent)
						} else {
							if _stop.ID != parent.ThisStop.ID {
								_connection := Connection_2{ThisStop: _stop, Leg: _leg, Previous: parent,
									Next: nil, Status: BEING_PROCESSED}
								jm.SetConnectionProcessingStatus(_leg, &_connection, BEING_PROCESSED)
								// fmt.Printf("        line 479: level %d - leg %s, _stop %s\n", level, _leg.ID, _stop.Location.ID)
								jm.FindConnectingNodes_v5(level+1, &_connection)
								if jm.StopProcessingStatus(_leg, _connection.ThisStop) != ON_PATH {
									jm.SetConnectionProcessingStatus(_leg, &_connection, PROCESSED)
								}
							}
						}
					}
				}
			}
		}

		// jm.ConnectionTree = nil
	}
}

// ============================================================================
// show resulting journies
// ============================================================================
func (jm *JourneyMap_2) ShowResultingLegs() {
	for _, _connection := range jm.ResultConnections {
		_result := _connection.WalkTheTree()
		fmt.Printf("*** %s\n", _result)
	}
}

// ============================================================================
// show resulting for one plan/route
// ============================================================================
func (conn *Connection_2) WalkTheTree() string {
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
func (jm *JourneyMap_2) SetLegProcessingStatus(leg *Leg_2, status Status) {
	jm.Processed.Leg[leg.ID] = status
}

// ============================================================================
// return leg processing status
// ============================================================================
func (jm *JourneyMap_2) LegProcessingStatus(leg *Leg_2) Status {
	return jm.Processed.Leg[leg.ID]
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
func (jm *JourneyMap_2) StopProcessingStatus(leg *Leg_2, location *Location_2) Status {
	_status, _ := jm.Processed.connection[leg.ID+"-"+location.ID]

	return _status
}

// ============================================================================
// mark this stop as visited
// ============================================================================
func (jm *JourneyMap_2) SetConnectionProcessingStatus(leg *Leg_2, connection *Connection_2, status Status) {
	jm.Processed.connection[leg.ID+"-"+connection.ThisStop.ID] = status
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
