// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"sync"
)

// PropertyCollection provides access to service properties
type PropertyCollection struct {
	lock        sync.RWMutex
	props       map[PropertyID]string
	propsByName map[string]string
}

// NewPropertyCollection creates a new empty property collection
func NewPropertyCollection() *PropertyCollection {
	return &PropertyCollection{
		props:       make(map[PropertyID]string),
		propsByName: make(map[string]string),
	}
}

// SetProperty sets a property value by PropertyID
func (pc *PropertyCollection) SetProperty(id PropertyID, value string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.props[id] = value
}

// SetPropertyByName sets a property value by name
func (pc *PropertyCollection) SetPropertyByName(name string, value string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.propsByName[name] = value
}

// GetProperty retrieves a property value by PropertyID
func (pc *PropertyCollection) GetProperty(id PropertyID, defaultValue ...string) string {
	pc.lock.RLock()
	defer pc.lock.RUnlock()

	if val, ok := pc.props[id]; ok {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetPropertyByName retrieves a property value by name
func (pc *PropertyCollection) GetPropertyByName(name string, defaultValue ...string) string {
	pc.lock.RLock()
	defer pc.lock.RUnlock()

	if val, ok := pc.propsByName[name]; ok {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// SetProperties sets multiple properties by ID
func (pc *PropertyCollection) SetProperties(properties map[PropertyID]string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	for id, val := range properties {
		pc.props[id] = val
	}
}

// SetPropertiesByName sets multiple properties by name
func (pc *PropertyCollection) SetPropertiesByName(properties map[string]string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	for name, val := range properties {
		pc.propsByName[name] = val
	}
}
