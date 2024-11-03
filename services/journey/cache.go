package journey

// ============================================================================
// ============================================================================
// this is an attempt to create a simple cache by using multi-index map
// This simple test does not use semaphore to sync to safeguard against
// multiple instance accessing and changing the data. NOT THREAD/GOROUTINE safe.
// ============================================================================
// ============================================================================
type miCacheType struct {
	NodeIdIndex map[string][]*Leg_2 // return all legs through a locationID
	LegIdIndex  map[string]*Leg_2   // return leg for a legID
}

const _LegIdCacheInitialSize = 100
const _NodeIdCacheInitialSize = _LegIdCacheInitialSize * 4

var MICache miCacheType = miCacheType{NodeIdIndex: make(map[string][]*Leg_2, _NodeIdCacheInitialSize),
	LegIdIndex: make(map[string]*Leg_2, _LegIdCacheInitialSize)}

// ============================================================================
//
// ============================================================================
func (mic *miCacheType) AddLeg(leg *Leg_2) error {
	var _err error
	// Leg is unique and therefore we can add to the cache
	MICache.LegIdIndex[leg.ID] = leg

	// one node can be part of many legs. so the mapping is not one to one
	// so for a node.ID we have an array of legs. []*legs
	for _, _stop := range leg.AllStops {
		_legs := MICache.NodeIdIndex[_stop.ID]
		// if _legs == nil {
		// 	_legs = make([]*Leg, 0)
		// }
		_legs = append(_legs, leg)
		MICache.NodeIdIndex[_stop.ID] = _legs
	}
	return _err
}

// ============================================================================
//
// ============================================================================
func (mic *miCacheType) GetLegFromId(id string) *Leg_2 {
	return mic.LegIdIndex[id]

}

// ============================================================================
//
// ============================================================================
func (mic *miCacheType) CreateAndAddLeg(id string, from, to *Location_2, distance int,
	timeTaken int, allStops Stops_2) {

	// var _newLeg Leg = Leg{ID: id, From: from, To: to, Distance: distance, TimeTaken: timeTaken, AllStops: allStops}
	var _newLeg Leg_2 = Leg_2{ID: id, Distance: distance, TimeTaken: timeTaken, AllStops: allStops}
	mic.AddLeg(&_newLeg)
}

// ============================================================================
//
// ============================================================================
func (mic *miCacheType) GetLegFromNode(locationId string) Legs_2 {
	return mic.NodeIdIndex[locationId]
}

// ============================================================================
//
// ============================================================================

// ============================================================================
//
// ============================================================================
