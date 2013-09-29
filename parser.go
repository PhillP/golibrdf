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
	"fmt"
	"runtime"
	"unsafe"
)

//Parser used to read and transform data in various formats into a model
type Parser struct {
	librdf_parser *C.librdf_parser
	Name          string
	world         *World
	mimeType      string
}

//NewParser constructs a new parser given a parserName and mimeType
//mimeType may be left empty
func NewParser(world *World, parserName string, mimeType string) (*Parser, error) {

	parser := Parser{Name: parserName, mimeType: mimeType}
	parser.world = world

	cParserName := C.CString(parser.Name)
	defer C.free(unsafe.Pointer(cParserName))

	cMimeType := C.CString(parser.mimeType)
	defer C.free(unsafe.Pointer(cMimeType))

	parser.librdf_parser = C.librdf_new_parser(world.librdf_world, cParserName, cMimeType, nil)
	runtime.SetFinalizer(&parser, (*Parser).Free)

	return &parser, nil
}

//Parse a string containing RDF dat into a model
func (parser *Parser) ParseStringIntoModel(rdfString string, baseUri *Uri, model *Model) error {

	var err error

	var baseUriPtr *C.librdf_uri
	baseUriPtr = nil

	cRdfString := C.CString(rdfString)
	defer C.free(unsafe.Pointer(cRdfString))

	if baseUri != nil {
		baseUriPtr = baseUri.librdf_uri
	}

	result := C.librdf_parser_parse_string_into_model(parser.librdf_parser, (*C.uchar)(unsafe.Pointer(cRdfString)), baseUriPtr, model.librdf_model)

	if result != 0 {
		err = errors.New("Unable to parse string into model")
	}

	return err
}

// Parse data at a specified URI into a model
func (parser *Parser) ParseIntoModel(uri *Uri, baseUri *Uri, model *Model) error {

	var err error

	var baseUriPtr *C.librdf_uri
	baseUriPtr = nil

	if baseUri != nil {
		baseUriPtr = baseUri.librdf_uri
	}

	result := C.librdf_parser_parse_into_model(parser.librdf_parser, uri.librdf_uri, baseUriPtr, model.librdf_model)

	fmt.Printf("%s", "parsed")

	if result != 0 {
		err = errors.New("Unable to parse URI into model")
	}

	return err
}

//Free cleans up memory resources held by the Parser
//	Free will be automatically called when Parser instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (parser *Parser) Free() {
	if parser.librdf_parser != nil {
		C.librdf_free_parser(parser.librdf_parser)
		parser.librdf_parser = nil
	}
}
