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

//An RDF Model
type Model struct {
	librdf_model *C.librdf_model
	world        *World
}

//NewModel constructs a new model backed by the provided storage
//Refer to librdf_new_model documentation for available options
func NewModel(world *World, storage *Storage, options string) (*Model, error) {
	cOptions := C.CString(options)
	defer C.free(unsafe.Pointer(cOptions))

	model := Model{}
	model.librdf_model = C.librdf_new_model(world.librdf_world, storage.librdf_storage, cOptions)

	if model.librdf_model == nil {
		return nil, errors.New("Unable to make new model.  Call to librdf_new_model failed.")
	}

	model.world = world

	// set the finalizer so that free call occurs as required
	runtime.SetFinalizer(&model, (*Model).Free)

	return &model, nil
}

//AddStatement adds the specified statement to the model
func (model *Model) AddStatement(statement *Statement) (err error) {
	C.librdf_model_add_statement(model.librdf_model, statement.librdf_statement)
	return nil
}

//ToString serializes the model to a string representation (RDFXML)
func (model *Model) ToString() string {
	cModelString := C.librdf_model_to_string(model.librdf_model, nil, nil, nil, nil)
	defer C.free(unsafe.Pointer(cModelString))

	return C.GoString((*C.char)(unsafe.Pointer(cModelString)))
}

//FindTargets returns a channel used to iterate through a set of matched targets give a subject + predicate pair to match
func (model *Model) FindTargets(subject *Node, predicate *Node, bufferSize int) chan *Node {
	chanNode := make(chan *Node, bufferSize)

	go func() {
		iterator := C.librdf_model_get_targets(model.librdf_model, subject.librdf_node, predicate.librdf_node)
		var librdfNode *C.librdf_node
		var node *Node
		var err error
		var atEnd C.int

		if iterator != nil {
			atEnd = C.librdf_iterator_end(iterator)
			for atEnd == 0 {

				librdfNode = (*C.librdf_node)(unsafe.Pointer(C.librdf_iterator_get_object(iterator)))

				if librdfNode == nil {
					panic(errors.New("librdf returned null node"))
				}

				if node, err = NewNode(model.world); err != nil {
					panic(errors.New("Unable to create node"))
				}
				node.librdf_node = librdfNode
				chanNode <- node

				C.librdf_iterator_next(iterator)
				atEnd = C.librdf_iterator_end(iterator)
			}
			C.librdf_free_iterator(iterator)
		}

		close(chanNode)
	}()

	return chanNode
}

//FindStatements creates a channel used to iterate the set of statements in the model that matched the given partial statement
//	bufferSize indicates how many statements can be on the channel at one time
func (model *Model) FindStatements(partialStatement *Statement, bufferSize int) chan *Statement {
	chanStatement := make(chan *Statement, bufferSize)

	go func() {
		stream := C.librdf_model_find_statements(model.librdf_model, partialStatement.librdf_statement)
		var librdfStatement *C.librdf_statement
		var statement *Statement
		var err error
		var endOfStream C.int

		if stream == nil {
			panic(errors.New("librdf returned null stream"))
		} else {
			endOfStream = C.librdf_stream_end(stream)
			for endOfStream == 0 {

				if librdfStatement = C.librdf_stream_get_object(stream); librdfStatement == nil {
					panic(errors.New("librdf returned null statement"))
				}

				if statement, err = NewStatement(model.world); err != nil {
					panic(errors.New("Unable to create statement"))
				}
				statement.librdf_statement = librdfStatement
				chanStatement <- statement

				C.librdf_stream_next(stream)
				endOfStream = C.librdf_stream_end(stream)
			}
			C.librdf_free_stream(stream)
		}

		close(chanStatement)
	}()

	return chanStatement
}

//ContainsStatement returns true if the model contains the given statement
func (model *Model) ContainsStatement(statement *Statement) bool {
	var contains bool = false

	if retCode := C.librdf_model_contains_statement(model.librdf_model, statement.librdf_statement); retCode != 0 {
		contains = true
	}

	return contains
}

//RemoveStatement removes the specified statement from the model
func (model *Model) RemoveStatement(statement *Statement) error {
	if retCode := C.librdf_model_remove_statement(model.librdf_model, statement.librdf_statement); retCode != 0 {
		return errors.New("Statement could not be removed")
	}
	return nil
}

func (model *Model) Load(uri *Uri) error {
	if retCode := C.librdf_model_load(model.librdf_model, uri.librdf_uri, nil, nil, nil); retCode != 0 {
		return errors.New("Failed to load model")
	}
	return nil
}

//ExecuteQueryToResultsChannel executes the given query and returns a channel used to read the results
func (model *Model) ExecuteQueryToResultsChannel(query *Query, bufferSize int) (chan *QueryResultItem, error) {
	var err error = nil
	chanQueryResultItem := make(chan *QueryResultItem, bufferSize)

	cQueryString := C.CString(query.queryString)
	if cQueryString != nil {
		defer C.free(unsafe.Pointer(cQueryString))
	}

	cName := C.CString(query.name)
	if cName != nil {
		defer C.free(unsafe.Pointer(cName))
	}

	librdf_query := C.librdf_new_query(query.world.librdf_world, (*C.char)(unsafe.Pointer(cName)), nil, (*C.uchar)(unsafe.Pointer(cQueryString)), nil)
	if librdf_query == nil {
		return nil, errors.New("Failed to create new query")
	}

	results := C.librdf_model_query_execute(model.librdf_model, librdf_query)

	if results == nil {
		err = errors.New("Error executing query")
	}

	if C.librdf_query_results_finished(results) != 0 {
		return nil, errors.New("Query returned no results")
	}

	go func() {

		for C.librdf_query_results_finished(results) == 0 {
			item := new(QueryResultItem)

			cBindingCount := C.librdf_query_results_get_bindings_count(results)

			bindingCount := int(cBindingCount)
			if bindingCount == 0 {
				continue
			}

			item.NameNodePairs = make([]NameNodePair, bindingCount, bindingCount)

			for i := 0; i < bindingCount; i++ {
				cName := C.librdf_query_results_get_binding_name(results, C.int(i))
				librdf_node := C.librdf_query_results_get_binding_value(results, C.int(i))

				item.NameNodePairs[i].Name = C.GoString(cName)
				node := new(Node)
				node.world = query.world
				node.librdf_node = librdf_node
				item.NameNodePairs[i].Node = node
			}

			chanQueryResultItem <- item
			C.librdf_query_results_next(results)
		}
		close(chanQueryResultItem)
	}()

	return chanQueryResultItem, err
}

//ExecuteQueryToFormattedString executes a query and serializes the results to a string in the format provided
func (model *Model) ExecuteQueryToFormattedString(query *Query, format string) (string, error) {
	var librdf_query *C.librdf_query

	cQueryString := C.CString(query.queryString)
	if cQueryString != nil {
		defer C.free(unsafe.Pointer(cQueryString))
	}

	cName := C.CString(query.name)
	if cName != nil {
		defer C.free(unsafe.Pointer(cName))
	}

	cFormat := C.CString(format)
	if cFormat != nil {
		defer C.free(unsafe.Pointer(cFormat))
	}

	librdf_query = C.librdf_new_query(query.world.librdf_world, (*C.char)(unsafe.Pointer(cName)), nil, (*C.uchar)(unsafe.Pointer(cQueryString)), nil)

	if librdf_query == nil {
		return "", errors.New("Unable to create query for execution")
	}

	results := C.librdf_model_query_execute(model.librdf_model, librdf_query)

	if results == nil {
		return "", errors.New("Failed to execute query")
	}

	var cFormattedString *C.uchar
	isBindings := C.librdf_query_results_is_bindings(results)
	isBoolean := C.librdf_query_results_is_boolean(results)

	if isBindings == 0 || isBoolean == 0 {
		cFormattedString = C.librdf_query_results_to_string2(results, (*C.char)(unsafe.Pointer(cFormat)), nil, nil, nil)
	} else {
		var serializer *C.librdf_serializer
		if serializer = C.librdf_new_serializer(query.world.librdf_world, (*C.char)(unsafe.Pointer(cFormat)), nil, nil); serializer == nil {
			return "", errors.New("Failed to build serializer")
		}
		defer C.librdf_free_serializer(serializer)

		var stream *C.librdf_stream
		if stream = C.librdf_query_results_as_stream(results); stream == nil {
			return "", errors.New("Failed to build stream from results")
		}
		defer C.librdf_free_stream(stream)

		cFormattedString = C.librdf_serializer_serialize_stream_to_string(serializer, nil, stream)
	}

	if cFormattedString == nil {
		return "", errors.New("Failed to build formatted string from results")
	}
	defer C.free(unsafe.Pointer(cFormattedString))

	formattedString := C.GoString((*C.char)(unsafe.Pointer(cFormattedString)))

	return formattedString, nil
}

//Free cleans up memory resources held by the model
//	Free will be automatically called when Model instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (model *Model) Free() {
	if model.librdf_model != nil {
		C.librdf_free_model(model.librdf_model)
		model.librdf_model = nil
	}
	return
}
