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
	"unsafe"
)

type Node struct {
	librdf_node *C.librdf_node
	world       *World
}

//NewNode constructs a new node
func NewNode(world *World) (*Node, error) {
	node := Node{}
	node.world = world

	return &node, nil
}

//NewNode constructs a new node from a specified URI
func NewNodeFromUri(world *World, uri *Uri) (*Node, error) {
	node, err := NewNode(world)

	if err != nil {
		return nil, err
	}

	node.librdf_node = C.librdf_new_node_from_uri(world.librdf_world, uri.librdf_uri)

	return node, nil
}

//NewNode constructs a new node from a string literal
func NewNodeFromLiteral(world *World, literal string) (*Node, error) {
	node, err := NewNode(world)

	if err != nil {
		return nil, err
	}

	cLiteralString := C.CString(literal)
	defer C.free(unsafe.Pointer(cLiteralString))

	node.librdf_node = C.librdf_new_node_from_literal(world.librdf_world, (*C.uchar)(unsafe.Pointer(cLiteralString)), nil, 0)

	return node, nil
}

//NewNode constructs a new node from an xml literal
func NewNodeFromXmlLiteral(world *World, xmlLiteral string, xmlLanguage string) (*Node, error) {
	node, err := NewNode(world)

	if err != nil {
		return nil, err
	}

	cLiteralString := C.CString(xmlLiteral)
	defer C.free(unsafe.Pointer(cLiteralString))

	cXmlLangString := C.CString(xmlLanguage)
	defer C.free(unsafe.Pointer(cLiteralString))

	node.librdf_node = C.librdf_new_node_from_literal(world.librdf_world, (*C.uchar)(unsafe.Pointer(cLiteralString)), cXmlLangString, 1)

	return node, nil
}

//NewNode constructs a new node from a URI string
func NewNodeFromUriString(world *World, uriString string) (*Node, error) {

	var node *Node
	var err error
	var uri *Uri

	if uri, err = newUriWithoutFinaliser(world, uriString); err == nil {
		node, err = NewNodeFromUri(world, uri)
	}

	return node, err
}

//ToString returns a string representation of the node
func (node *Node) ToString() (string, error) {
	var stringPointer unsafe.Pointer
	var length C.size_t

	raptorWorld := node.world.GetRaptorWorld()
	stream := C.raptor_new_iostream_to_string(raptorWorld, &stringPointer, &length, nil)

	if stream == nil {
		return "", errors.New("Unable to obtain raptor stream")
	}

	defer C.free(unsafe.Pointer(stringPointer))

	C.librdf_node_write(node.librdf_node, stream)
	C.raptor_free_iostream(stream)

	return C.GoString((*C.char)(unsafe.Pointer(stringPointer))), nil
}

//GetUriString returns the URI string for the node.  (Only appropriate if the node is a resource type)
func (node *Node) GetUriString() string {
	var uriString string

	if uri := C.librdf_node_get_uri(node.librdf_node); uri != nil {
		if cString := C.librdf_uri_to_string(uri); cString != nil {
			defer C.free(unsafe.Pointer(cString))

			uriString = C.GoString((*C.char)(unsafe.Pointer(cString)))
		}
	}

	return uriString 
}

//IsLiteral returns true if the node represents a literal value
func (node *Node) IsLiteral() bool {
	isLiteralCode := C.librdf_node_is_literal(node.librdf_node)

	return (isLiteralCode != 0) 
}

//IsResource returns true if the node represents a resource
func (node *Node) IsResource() bool {
	isResourceCode := C.librdf_node_is_resource(node.librdf_node)

	return (isResourceCode != 0) 
}

//IsBlank returns true if the node represents a blank
func (node *Node) IsBlank() bool {
	isBlank := C.librdf_node_is_blank(node.librdf_node)

	return (isBlank != 0) 
}

//GetLiteralValue returns the literal string for the node.  (Only appropriate if the node is a literal type)
func (node *Node) GetLiteralValue() string {
	var literalString string

	value := C.librdf_node_get_literal_value(node.librdf_node)

	if value != nil {
		literalString = C.GoString(((*C.char)(unsafe.Pointer(value))))
	}

	return literalString
}

//GetLiteralValueLanguage returns the language string associated with the literal value for the node
func (node *Node) GetLiteralValueLanguage() string {
	var languageString string

	value := C.librdf_node_get_literal_value_language(node.librdf_node)

	if value != nil {
		languageString = C.GoString(value)
	}

	return languageString
}

//Free cleans up memory resources held by the Node
//	Free will be automatically called when Node instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
//
//An exception to this is that there is no need to explicitly Free Node instances
//that are attached to a statement that has been Freed explicitly
func (node *Node) Free() {
	if node.librdf_node != nil {
		C.librdf_free_node(node.librdf_node)
		node.librdf_node = nil
	}
}
