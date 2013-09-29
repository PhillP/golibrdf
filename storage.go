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

type Storage struct {
	librdf_storage *C.librdf_storage
}

func NewStorage(world *World, storageName string, name string, options string) (*Storage, error) {

	cStorageName := C.CString(storageName)
	defer C.free(unsafe.Pointer(cStorageName))

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cOptions := C.CString(options)
	defer C.free(unsafe.Pointer(cOptions))

	storage := Storage{}
	storage.librdf_storage = C.librdf_new_storage(world.librdf_world, cStorageName, cName, cOptions)

	if storage.librdf_storage == nil {
		return nil, errors.New("Unable to make new storage.  Call to librdf_storage failed.")
	}

	runtime.SetFinalizer(&storage, (*Storage).Free)

	return &storage, nil
}

//Free cleans up memory resources held by the Storage
//	Free will be automatically called when Storage instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (storage *Storage) Free() {
	if storage.librdf_storage != nil {
		C.librdf_free_storage(storage.librdf_storage)
		storage.librdf_storage = nil
	}
	return
}
