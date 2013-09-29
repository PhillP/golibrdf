/*
*
* This file forms part of the golibrdf package containing go language bindings,
* tests and examples for the Redland RDF library.
*
* Please refer to http://librdf.org for copyright and licence information 
* on the Redland libraries that this package wraps 
*
* This golibrdf package is: 
* 	Copyright (C) 2013, Phillip Pettit http://ppettit.net/
* 
* This package is licensed under the following three licenses as alternatives:
* 1. GNU Lesser General Public License (LGPL) V2.1 or any newer version
* 2. GNU General Public License (GPL) V2 or any newer version
* 3. Apache License, V2.0 or any newer version
*
* You may not use this file except in compliance with at least one of
* the above three licenses.
*
*/

package golibrdf

// #cgo linux pkg-config: redland raptor2
// #cgo LDFLAGS: -lrdf
// #include <stdlib.h>
// #include <string.h>
// #include <strings.h>
// #include <librdf.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

//World represents a Redland execution environment 
type World struct {
	librdf_world        *C.librdf_world
	librdf_raptor_world *C.raptor_world
	isOpen              bool
	hasBeenOpen         bool
}

//NewWorld constructs a new World.  The World must be opened before use.
func NewWorld() *World {
	world := World{}

	return &world
}

//Open readies a World for use.  A corresponding Close call must be made to free resources.
func (world *World) Open() error {
	if world.IsOpen() {
		return errors.New("Unable to Open() world.  World is already open.")
	}

	if world.hasBeenOpen {
		return errors.New("Unable to Open() world.  World has previously been opened.")
	}

	world.librdf_world = C.librdf_new_world()
	C.librdf_world_open(world.librdf_world)

	world.isOpen = true
	world.hasBeenOpen = true

	// set the finalizer so that the librdf_world_free call occurs as required
	runtime.SetFinalizer(world, (*World).Close)

	return nil
}

//IsOpen tests whether a world is open and returns true if the world is open
func (world World) IsOpen() bool {
	return world.isOpen
}

//GetRaptorWorld returns a raptor reference associated with the world
func (world *World) GetRaptorWorld() *C.raptor_world {
	return C.librdf_world_get_raptor(world.librdf_world)
}

//SetRaptorWorld associates a raptor world reference with the world
func (world *World) SetRaptorWorld(raptorWorld *C.raptor_world) {
	C.librdf_world_set_raptor(world.librdf_world, raptorWorld)
}

//GetRasqalWorld returns a rasqal reference associated with the world
func (world *World) GetRasqalWorld() *C.rasqal_world {
	return C.librdf_world_get_rasqal(world.librdf_world)
}

//GuessParserName is used to guess the appropriate parser given a URI
func (world *World) GuessParserName(uri *Uri) string {
	var cUriAsString *C.uchar
	if cUriAsString = C.librdf_uri_as_string(uri.librdf_uri); cUriAsString == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cUriAsString))

	var cParserName *C.char
	if cParserName = C.librdf_parser_guess_name2(world.librdf_world, nil, nil, cUriAsString); cParserName == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cParserName))

	parserName := C.GoString(cParserName)

	return parserName
}

//SetRasqalWorld associates a rasqal world reference with the world
func (world *World) SetRasqalWorld(rasqalWorld *C.rasqal_world) {
	C.librdf_world_set_rasqal(world.librdf_world, rasqalWorld)
}

//Close cleans up memory resources held by the World
func (world *World) Close() {
	if world.librdf_world != nil {
		C.librdf_free_world(world.librdf_world)
		world.librdf_world = nil
	}

	world.isOpen = false
}

//SetFeature specifies a value for a world feature (setting)
func (world *World) SetFeature(feature *Uri, value *Node) {
	C.librdf_world_set_feature(world.librdf_world, feature.librdf_uri, value.librdf_node)
}

//GetFeature returns a value node for a world feature
func (world *World) GetFeature(feature *Uri) *Node {
	var node *Node
	var err error
	nodeValue := C.librdf_world_get_feature(world.librdf_world, feature.librdf_uri)

	if nodeValue != nil {
		node, err = NewNode(world)

		if err != nil {
			panic(err)
		}

		node.librdf_node = nodeValue
	}

	return node
}

//SetDigest sets a digest for the world
func (world *World) SetDigest(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.librdf_world_set_digest(world.librdf_world, cName)
}
