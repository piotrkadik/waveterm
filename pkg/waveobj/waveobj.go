// Copyright 2025, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

// Package waveobj defines the core object model for WaveTerm.
// All persistent entities (tabs, blocks, workspaces, etc.) implement the WaveObj interface.
package waveobj

import (
	"fmt"
	"reflect"
	"sort"
)

// OType is the object type identifier string (e.g. "tab", "block", "workspace").
type OType = string

// OID is a unique object identifier (UUID string).
type OID = string

// MetaMapType is a generic string-keyed map used for object metadata.
type MetaMapType = map[string]any

// WaveObj is the interface that all persistent WaveTerm objects must implement.
type WaveObj interface {
	// GetOType returns the object type string (e.g. "tab", "block").
	GetOType() OType

	// GetOID returns the unique identifier of the object.
	GetOID() OID

	// GetVersion returns the current version counter of the object.
	GetVersion() int

	// SetVersion sets the version counter on the object.
	SetVersion(v int)
}

// ORef is a typed reference to a WaveObj by type and OID.
type ORef struct {
	OType OType `json:"otype"`
	OID   OID   `json:"oid"`
}

// String returns a human-readable representation of an ORef.
func (r ORef) String() string {
	return fmt.Sprintf("%s:%s", r.OType, r.OID)
}

// IsEmpty returns true if the ORef has no type or ID set.
func (r ORef) IsEmpty() bool {
	return r.OType == "" || r.OID == ""
}

// MakeORef constructs an ORef from a WaveObj.
func MakeORef(obj WaveObj) ORef {
	return ORef{
		OType: obj.GetOType(),
		OID:   obj.GetOID(),
	}
}

// registry maps OType strings to their reflect.Type for deserialization.
var registry = map[OType]reflect.Type{}

// RegisterType registers a WaveObj implementation type so it can be
// instantiated by OType during deserialization.
// T must be a pointer type that implements WaveObj.
func RegisterType[T WaveObj](otype OType) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("waveobj.RegisterType: %T must be a pointer type", zero))
	}
	registry[otype] = t.Elem()
}

// NewByOType creates a new zero-value instance of the registered type for otype.
// Returns nil and an error if otype is not registered.
func NewByOType(otype OType) (WaveObj, error) {
	t, ok := registry[otype]
	if !ok {
		return nil, fmt.Errorf("waveobj: unknown otype %q", otype)
	}
	return reflect.New(t).Interface().(WaveObj), nil
}

// GetOTypes returns all registered OType strings in sorted order for
// deterministic iteration (useful for debugging and tests).
func GetOTypes() []OType {
	types := make([]OType, 0, len(registry))
	for k := range registry {
		types = append(types, k)
	}
	sort.Strings(types)
	return types
}
