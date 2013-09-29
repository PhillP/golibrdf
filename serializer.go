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

//A Serializer used to serialize a model into various formats
type Serializer struct {
	librdf_serializer *C.librdf_serializer
}

//NewSerializer construcs a new serializer based on a name defining the type, a mimeType and optional URI
func NewSerializer(world *World, name string, mimeType string, uri *Uri) (*Serializer, error) {

	serializer := Serializer{}

	var uriPtr *C.librdf_uri
	uriPtr = nil

	if uriPtr != nil {
		uriPtr = uri.librdf_uri
	}

	serializer.librdf_serializer = C.librdf_new_serializer(world.librdf_world, C.CString(name), C.CString(mimeType), uriPtr)

	runtime.SetFinalizer(&serializer, (*Serializer).Free)

	return &serializer, nil
}

//Serialize a model to a string in the format appropriate for the serializer
func (serializer *Serializer) SerializeModelToString(model *Model, baseUri *Uri) (string, error) {
	var err error
	var resultString string

	var baseUriPtr *C.librdf_uri
	baseUriPtr = nil

	if baseUri != nil {
		baseUriPtr = baseUri.librdf_uri
	}

	result := C.librdf_serializer_serialize_model_to_string(serializer.librdf_serializer, baseUriPtr, model.librdf_model)

	if result == nil {
		err = errors.New("Unable to serialize model")
	} else {
		resultString = C.GoString((*C.char)(unsafe.Pointer(result)))
	}

	return resultString, err
}

//Free cleans up memory resources held by the Serializer
//	Free will be automatically called when Serializer instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (serializer *Serializer) Free() {

	if serializer.librdf_serializer != nil {
		serializer.librdf_serializer = nil
		C.librdf_free_serializer(serializer.librdf_serializer)
	}

	return
}
