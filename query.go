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

//A Query that can be executed against a model to produce results 
type Query struct {
	world       *World
	name        string
	queryString string
}

//NewQuery constructs a new query given a name indicating the query type and a string containing the query 
func NewQuery(world *World, name string, queryString string) (Query, error) {

	var err error = nil
	query := Query{world: world, name: name, queryString: queryString}
	return query, err
}
