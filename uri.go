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

//Represents a URI
type Uri struct {
	librdf_uri *C.librdf_uri
}

//NewUri constructs a new URI given a string
func NewUri(world *World, uriString string) (*Uri, error) {
	uri, err := newUriWithoutFinaliser(world, uriString)

	if uri != nil {
		runtime.SetFinalizer(uri, (*Uri).Free)
	}

	return uri, err
}

//newUriWithoutFinaliser constructs a new URI given a string, but does not associate a finalizer for automatic free
func newUriWithoutFinaliser(world *World, uriString string) (*Uri, error) {
	uri := new(Uri)

	cUriString := C.CString(uriString)
	defer C.free(unsafe.Pointer(cUriString))

	uri.librdf_uri = C.librdf_new_uri(world.librdf_world, (*C.uchar)(unsafe.Pointer(cUriString)))

	if uri.librdf_uri == nil {
		return nil, errors.New("Unable to create URI for uri string")
	}

	return uri, nil
}

//Free cleans up memory resources held by the Uri
//	Free will be automatically called when Uri instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (uri *Uri) Free() {
	if uri.librdf_uri != nil {
		C.librdf_free_uri(uri.librdf_uri)
		uri.librdf_uri = nil
	}

	return
}

//ToString serializers a URI to string
func (uri Uri) ToString() string {
	cUriString := C.librdf_uri_as_string(uri.librdf_uri)
	defer C.free(unsafe.Pointer(cUriString))

	return C.GoString((*C.char)(unsafe.Pointer(cUriString)))
}

//NewUriFromUri constructs a new URI given an existing URI
func NewUriFromUri(fromUri *Uri) (*Uri, error) {
	var err error

	uri := new(Uri)
	uri.librdf_uri = C.librdf_new_uri_from_uri(fromUri.librdf_uri)

	runtime.SetFinalizer(uri, (*Uri).Free)

	return uri, err
}

//NewUriFromUri constructs a new URI given an existing URI and a localName
func NewUriFromUriLocalName(fromUri *Uri, localName string) (*Uri, error) {
	var err error

	cLocalName := C.CString(localName)
	defer C.free(unsafe.Pointer(cLocalName))

	uri := new(Uri)
	uri.librdf_uri = C.librdf_new_uri_from_uri_local_name(fromUri.librdf_uri, (*C.uchar)(unsafe.Pointer(cLocalName)))

	runtime.SetFinalizer(uri, (*Uri).Free)

	return uri, err
}

//NewUriFromUri constructs a new URI given an existing URI normalised to the specified baseUri
func NewUriNormalisedBase(uriString string, sourceUri *Uri, baseUri *Uri) (*Uri, error) {
	var err error

	cUriString := C.CString(uriString)
	defer C.free(unsafe.Pointer(cUriString))

	uri := new(Uri)
	uri.librdf_uri = C.librdf_new_uri_normalised_to_base((*C.uchar)(unsafe.Pointer(cUriString)), sourceUri.librdf_uri, baseUri.librdf_uri)

	runtime.SetFinalizer(uri, (*Uri).Free)

	return uri, err
}

//NewUriRelativeToBase constructs a new URI given a URI string made relative to the specified baseUri
func NewUriRelativeToBase(baseUri *Uri, uriString string) (*Uri, error) {
	var err error

	cUriString := C.CString(uriString)
	defer C.free(unsafe.Pointer(cUriString))

	uri := new(Uri)
	uri.librdf_uri = C.librdf_new_uri_relative_to_base(baseUri.librdf_uri, (*C.uchar)(unsafe.Pointer(cUriString)))

	runtime.SetFinalizer(uri, (*Uri).Free)

	return uri, err
}

//NewUriFromFileName constructs a new URI for a file given a filename
func NewUriFromFileName(world *World, fileName string) (*Uri, error) {
	var err error

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	uri := new(Uri)
	uri.librdf_uri = C.librdf_new_uri_from_filename(world.librdf_world, (*C.char)(unsafe.Pointer(cFileName)))

	runtime.SetFinalizer(uri, (*Uri).Free)

	return uri, err
}

//ToFileName converts a URI representing a file to a filename
func (uri *Uri) ToFileName() (string, error) {
	var err error

	cFileName := C.librdf_uri_to_filename(uri.librdf_uri)
	defer C.free(unsafe.Pointer(cFileName))

	fileName := C.GoString((*C.char)(unsafe.Pointer(cFileName)))

	return fileName, err
}

//IsFileUri tests whether a URI represents a file or not.
func (uri *Uri) IsFileUri() bool {
	cIsFileUri := int(C.librdf_uri_is_file_uri(uri.librdf_uri))
	return cIsFileUri == 0
}

//Equals compares 2 URIs and returns true if they are equal
func (uri *Uri) Equals(other *Uri) bool {
	cEquals := int(C.librdf_uri_equals(uri.librdf_uri, other.librdf_uri))
	return cEquals == 0
}

//Compare compares 2 URIs
// Returns <0 if the URI instance is less than other
// Returns >0 if the URI instance is greater than other
// Returns 0 if the URIs are equal
func (uri *Uri) Compare(other *Uri) int {
	return int(C.librdf_uri_compare(uri.librdf_uri, other.librdf_uri))
}
