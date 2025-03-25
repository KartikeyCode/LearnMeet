package objects

import (
	"sync"
)

//Thread Safe map of objects with auto incrementing IDs

type SharedCollection[T any] struct {
	objectsMap map[uint64]T
	nextId     uint64
	mapMux     sync.Mutex
}

func NewSharedCollection[T any](capacity ...int) *SharedCollection[T] {
	var newObjMap map[uint64]T
	if len(capacity) > 0 {
		newObjMap = make(map[uint64]T, capacity[0])

	} else {
		newObjMap = make(map[uint64]T)
	}

	return &SharedCollection[T]{
		objectsMap: newObjMap,
		nextId:     1,
	}
}

// Add object to map with Id or next available id (returns id of object added)
func (s *SharedCollection[T]) Add(obj T, id ...uint64) uint64 {
	s.mapMux.Lock()
	defer s.mapMux.Unlock()

	thisId := s.nextId
	if len(id) > 0 {
		thisId = id[0]
	}

	s.objectsMap[thisId] = obj
	s.nextId++
	return thisId
}

// remove object from map of ids (if exists)

func (s *SharedCollection[T]) Remove(id uint64) {
	s.mapMux.Lock()
	defer s.mapMux.Unlock()

	delete(s.objectsMap, id)
}

//call the callback function for each object in the map:

func (s *SharedCollection[T]) ForEach(callback func(uint64, T)) {
	//create a local copy whilst holding the lock:
	s.mapMux.Lock()
	localCopy := make(map[uint64]T, len(s.objectsMap))
	for id, obj := range s.objectsMap {
		localCopy[id] = obj
	}
	s.mapMux.Unlock()

	//iterate over local copy without lock
	for id, obj := range localCopy {
		callback(id, obj)
	}
}

//get an object for given id (if doesn't exist return nil) also returns bool for if object exists

func (s *SharedCollection[T]) Get(id uint64) (T, bool) {
	s.mapMux.Lock()
	defer s.mapMux.Unlock()

	obj, found := s.objectsMap[id]
	return obj, found
}

//get approx number of objects in map:

func (s *SharedCollection[T]) Len() int {
	return len(s.objectsMap)
}
