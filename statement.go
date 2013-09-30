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

const (
	defaultStatementEncodingBufferSize = 2000
)

const (
	StatementSubject   = 1 << 0
	StatementPredicate = 1 << 1
	StatementObject    = 1 << 2

	/* must be a combination of all of the above */
	StatementAll = (StatementSubject |
		StatementPredicate |
		StatementObject)
)

//A RDF statement
type Statement struct {
	librdf_statement *C.librdf_statement
	world            *World
}

//NewStatementFromNodes constructs a statement given subject, predicate and object nodes
func NewStatementFromNodes(world *World, subject *Node, predicate *Node, object *Node) (*Statement, error) {
	statement := Statement{}
	statement.world = world
	statement.librdf_statement = C.librdf_new_statement_from_nodes(world.librdf_world, subject.librdf_node, predicate.librdf_node, object.librdf_node)

	runtime.SetFinalizer(&statement, (*Statement).Free)

	return &statement, nil
}

//NewStatement constructs a new statement
func NewStatement(world *World) (*Statement, error) {
	statement := Statement{}
	statement.world = world
	statement.librdf_statement = C.librdf_new_statement(world.librdf_world)

	runtime.SetFinalizer(&statement, (*Statement).Free)

	return &statement, nil
}

//DeepClone performs a deep clone of a statement and returns a clone
func (statement *Statement) DeepClone() (*Statement, error) {
	newStatement := Statement{}
	newStatement.world = statement.world
	newStatement.librdf_statement = C.librdf_new_statement_from_statement(statement.librdf_statement)

	runtime.SetFinalizer(&newStatement, (*Statement).Free)

	return &newStatement, nil
}

//ShallowClone performs a shallow clone of a statement and returns a clone
func (statement *Statement) ShallowClone() (*Statement, error) {
	newStatement := Statement{}
	newStatement.world = statement.world
	newStatement.librdf_statement = C.librdf_new_statement_from_statement2(statement.librdf_statement)

	runtime.SetFinalizer(&newStatement, (*Statement).Free)

	return &newStatement, nil
}

//Clear removes the nodes associated with a statement
func (statement *Statement) Clear() {
	if statement.librdf_statement == nil {
		panic(errors.New("Statement can't be cleared as it has already been freed"))
	}
	C.librdf_statement_clear(statement.librdf_statement)

	return
}

//Free cleans up memory resources held by the Statement
//	Free will be automatically called when Statement instances are garbage collected
//  however it is important to explicitly call Free to avoid issues that may result
//  from freeing resources in an unexpected order
func (statement *Statement) Free() {
	if statement.librdf_statement != nil {
		C.librdf_free_statement(statement.librdf_statement)
		statement.librdf_statement = nil
	}

	return
}

//SetSubject associates a node with the statement as a subject
func (statement *Statement) SetSubject(subject *Node) {
	C.librdf_statement_set_subject(statement.librdf_statement, subject.librdf_node)
}

//GetSubject gets the subject node associated with the statement
func (statement *Statement) GetSubject() *Node {
	libRdfNode := C.librdf_statement_get_subject(statement.librdf_statement)

	node := Node{}
	node.librdf_node = libRdfNode
	node.world = statement.world

	return &node
}

//SetPredicate associates a node with the statement as a predicate
func (statement *Statement) SetPredicate(predicate *Node) {
	C.librdf_statement_set_predicate(statement.librdf_statement, predicate.librdf_node)
}

//GetPredicate gets the predicate node associated with the statement
func (statement *Statement) GetPredicate() *Node {
	libRdfNode := C.librdf_statement_get_predicate(statement.librdf_statement)

	node := Node{}
	node.librdf_node = libRdfNode
	node.world = statement.world

	return &node
}

//SetObject associates a node with the statement as object
func (statement *Statement) SetObject(object *Node) {
	C.librdf_statement_set_object(statement.librdf_statement, object.librdf_node)
}

//GetObject gets the object associated with the statement
func (statement *Statement) GetObject() *Node {
	libRdfNode := C.librdf_statement_get_object(statement.librdf_statement)

	node := Node{}
	node.librdf_node = libRdfNode
	node.world = statement.world

	return &node
}

//IsComplete returns true if the statement has subject, predicate and object nodes
func (statement *Statement) IsComplete() bool {
	isComplete := int(C.librdf_statement_is_complete(statement.librdf_statement)) == 0

	return isComplete
}

//IsEqual compares 2 statements and returns true if the statements are equal
func (statement *Statement) IsEqual(other *Statement) bool {
	isEqual := int(C.librdf_statement_equals(statement.librdf_statement, other.librdf_statement)) == 0

	return isEqual
}

//IsMatch compares the statement with a partial statement and returns true if the statement is a match
func (statement *Statement) IsMatch(partial *Statement) bool {
	isMatch := int(C.librdf_statement_match(statement.librdf_statement, partial.librdf_statement)) == 0

	return isMatch
}

//Encode encodes a statement to string
func (statement *Statement) Encode() (string, error) {
	var encodedString string
	var err error

	bufferSize := (C.size_t)(defaultStatementEncodingBufferSize)
	buffer := make([]byte, bufferSize)

	written := C.librdf_statement_encode2(statement.world.librdf_world, statement.librdf_statement, (*C.uchar)(unsafe.Pointer(&buffer[0])), bufferSize)

	if written == 0 {
		// determine the size of the buffer needed by passing a nil buffer
		sizeNeeded := C.librdf_statement_encode2(statement.world.librdf_world, statement.librdf_statement, nil, 0)

		if sizeNeeded == 0 {
			err = errors.New("Unable to determine size of buffer needed for encoding")
		} else {
			bufferSize = sizeNeeded
			buffer = make([]byte, bufferSize)
			written := C.librdf_statement_encode2(statement.world.librdf_world, statement.librdf_statement, (*C.uchar)(unsafe.Pointer(&buffer[0])), bufferSize)

			if written == 0 {
				err = errors.New("Unable to encode, even after explicitly determining required buffer size")
			}
		}
	}

	if written > 0 {
		encodedString = string(buffer[:written])
	}

	return encodedString, err
}

//EncodeParts encodes one or more of the subject,predicate and object parts of a statement to string
func (statement *Statement) EncodeParts(contextNode *Node, parts int) (string, error) {
	var encodedString string
	var err error
	var nodeRef *C.librdf_node

	if contextNode != nil {
		nodeRef = contextNode.librdf_node
	}

	bufferSize := (C.size_t)(defaultStatementEncodingBufferSize)
	buffer := make([]byte, bufferSize)

	partsFields := (C.librdf_statement_part)(parts)

	written := C.librdf_statement_encode_parts2(statement.world.librdf_world, statement.librdf_statement, nodeRef, (*C.uchar)(unsafe.Pointer(&buffer[0])), bufferSize, partsFields)

	if written == 0 {
		// determine the size of the buffer needed by passing a nil buffer
		sizeNeeded := C.librdf_statement_encode_parts2(statement.world.librdf_world, statement.librdf_statement, nodeRef, nil, 0, partsFields)

		if sizeNeeded == 0 {
			err = errors.New("Unable to determine size of buffer needed for encoding")
		} else {
			bufferSize = sizeNeeded
			buffer = make([]byte, bufferSize)
			written := C.librdf_statement_encode_parts2(statement.world.librdf_world, statement.librdf_statement, nodeRef, (*C.uchar)(unsafe.Pointer(&buffer[0])), bufferSize, partsFields)

			if written == 0 {
				err = errors.New("Unable to encode, even after explicitly determining required buffer size")
			}
		}
	}

	if written > 0 {
		encodedString = string(buffer[:written])
	}

	return encodedString, err
}

//Decode decodes a string to a Statement
func (statement *Statement) Decode(world *World, encodedStatement string) error {
	_, err := statement.decodeInner(world, encodedStatement, false)
	return err
}

//Decode decodes a string to a Statement with a context node
func (statement *Statement) DecodeWithContextNode(world *World, encodedStatement string) (*Node, error) {
	return statement.decodeInner(world, encodedStatement, true)
}

//decodeInner decodes a string to a Statement with optional context node
func (statement *Statement) decodeInner(world *World, encodedStatement string, withContextNode bool) (*Node, error) {
	var err error
	var nodeRef *C.librdf_node
	var node *Node
	var read C.size_t

	buffer := (*C.uchar)(unsafe.Pointer(&encodedStatement))
	bufferSize := (C.size_t)(len(encodedStatement))

	if withContextNode {
		read = C.librdf_statement_decode2(world.librdf_world, statement.librdf_statement, &nodeRef, (*C.uchar)(unsafe.Pointer(&buffer)), bufferSize)

		if read > 0 {
			if nodeRef != nil {
				node = &Node{}
				node.librdf_node = nodeRef
				node.world = world
			}
		}
	} else {
		read = C.librdf_statement_decode2(world.librdf_world, statement.librdf_statement, nil, (*C.uchar)(unsafe.Pointer(&buffer)), bufferSize)
	}

	if read == 0 {
		err = errors.New("Unable to determine size of buffer needed for encoding")
	}

	return node, err
}

//ToString serializers a statement to string
func (statement *Statement) ToString() (string, error) {
	var stringPointer unsafe.Pointer
	var length C.size_t
	
	raptorWorld := statement.world.GetRaptorWorld()
	
	stream := C.raptor_new_iostream_to_string(raptorWorld, &stringPointer, &length, nil)
	if stream == nil {
		return "", errors.New("Unable to obtain raptor stream")
	}
	defer C.raptor_free_iostream(stream)

	if result := C.librdf_statement_write(statement.librdf_statement, stream); result != 0 {
		return "", errors.New("Unable to write statement")
	}
	defer C.free(unsafe.Pointer(stringPointer))
	

	return C.GoString((*C.char)(unsafe.Pointer(stringPointer))), nil
}
