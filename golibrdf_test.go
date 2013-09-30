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

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const rdfxml_content string = `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
xmlns:dc="http://purl.org/dc/elements/1.1/">
<rdf:Description rdf:about="http://www.dajobe.org/">
<dc:title>Dave Beckett's Home Page</dc:title>
<dc:creator>Dave Beckett</dc:creator>
<dc:description>The generic home page of Dave Beckett.</dc:description>
</rdf:Description>
</rdf:RDF>`

const testRemoteUri string = "http://rawgithub.com/PhillP/golibrdf/master/testdata/redland.rdf"
const testRemoteUri2 string = "http://phillp.github.io/"
const testRemoteUri3 string = "http://rawgithub.com/PhillP/golibrdf/master/testdata/planet.rdf"
const testRemoteUri4 string = "https://rawgithub.com/PhillP/golibrdf/master/testdata/dc.rdf"
const testLocalUriFile string = "file:./testdata/dc.rdf"

//Test_ParseAddAndQuery is based on Example1 from the librdf library and tests the following sequence:
//	- Creating a model based on content parsed from a URI. 
//	- Adding a new statement
//	- Querying the model for statements that match a partial statement.
func Test_ParseAddAndQuery(t *testing.T) {

	storageType := "memory"

	var uri *Uri
	var storage *Storage
	var model *Model
	var parser *Parser
	var err error

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", ""); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to create a new model: %s", err.Error())
	}
	defer model.Free()

	// and a parser
	uriString := testRemoteUri
	// construct a URI based on the uri arg
	if uri, err = NewUri(world, uriString); err != nil {
		t.Fatalf("Failed to create a new URI: %s", err.Error())
	}
	defer uri.Free()

	parserName := "rdfxml"
	if parser, err = NewParser(world, parserName, ""); err != nil {
		t.Fatalf("Failed to create a new parser: %s", err.Error())
	}
	defer parser.Free()

	// parse the content at the uri into the model
	if err = parser.ParseIntoModel(uri, nil, model); err != nil {
		t.Fatalf("Failed to parse RDF into the model: %s", err.Error())
	}

	var subject, predicate, object *Node
	var statement *Statement

	// build a statement out of a set of nodes
	subject, err = NewNodeFromUriString(world, "http://www.dajobe.org/")
	if err != nil {
		t.Fatalf("Failed to create subject node from uri: %s", err.Error())
	}
	defer subject.Free()

	predicate, err = NewNodeFromUriString(world, "http://purl.org/dc/elements/1.1/title")
	if err != nil {
		t.Fatalf("Failed to create predicate node from Uri string: %s", err.Error())
	}
	defer predicate.Free()

	object, err = NewNodeFromLiteral(world, "My Home Page")
	if err != nil {
		t.Fatalf("Failed to create node from literal: %s", err.Error())
	}
	defer object.Free()

	statement, err = NewStatementFromNodes(world, subject, predicate, object)
	if err != nil {
		t.Fatalf("Failed to create statement from nodes: %s", err.Error())
	}
	defer statement.Free()

	statementString, err := statement.ToString()
	if err != nil {
		t.Fatalf("Failed to construct statement string: %s", err.Error())
	}
	fmt.Printf("  Statement: %s\n", statementString)

	model.AddStatement(statement)

	// print out the model
	fmt.Printf("%s: Resulting model is:\n%s", os.Args[0], model.ToString())

	var findSubject, findPredicate *Node
	var partialStatement *Statement

	if findSubject, err = NewNodeFromUriString(world, "http://www.dajobe.org/"); err != nil {
		t.Fatalf("Failed to create new node: %s", err.Error())
	}

	if findPredicate, err = NewNodeFromUriString(world, "http://purl.org/dc/elements/1.1/title"); err != nil {
		t.Fatalf("Failed to create findPredicate node: %s", err.Error())
	}

	if partialStatement, err = NewStatement(world); err != nil {
		t.Fatalf("Failed to create partial statement: %s", err.Error())
	}
	defer partialStatement.Free() // this will also free attached nodes

	partialStatement.SetSubject(findSubject)
	partialStatement.SetPredicate(findPredicate)

	fmt.Printf("%s: Trying to find statements", os.Args[0])

	count := 0
	chanStatements := model.FindStatements(partialStatement, 100)

	for statement := range chanStatements {
		matchedStatementString, err := statement.ToString()
		if err != nil {
			t.Fatalf("Failed to create matchedStatement string: %s", err.Error())
		}
		fmt.Printf("  Matched Statement: %s\n", matchedStatementString)
		count = count + 1
		statement.Free()
	}

	fmt.Printf("%s: Got %d matching statement(s)", os.Args[0], count)

	chanTargets := model.FindTargets(subject, predicate, 100)
	count = 0

	for target := range chanTargets {
		targetString, err := target.ToString()
		if err != nil {
			t.Fatalf("Failed to create targetString: %s", err.Error())	
		}
		fmt.Printf("  Matched Target: %s", targetString)
		count = count + 1
		target.Free()
	}

	fmt.Printf("%s: Got %d matching targets", os.Args[0], count)
}

//Test_ParseStringAddCheckAndRemove is based on Example2 from the librdf library and tests the following sequence:
//	- Creating a model based on a string (content from a URI). 
//	- Adding a new statement
//	- Checking that the model contains the statement
//	- Removing the statement
func Test_ParseStringAddCheckAndRemove(t *testing.T) {
	storageType := "memory"

	var err error
	var storage *Storage
	var model *Model
	var parser *Parser
	var uri *Uri

	// create a new world
	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	if uri, err = NewUri(world, testRemoteUri2); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", ""); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	if parser, err = NewParser(world, "rdfxml", ""); err != nil {
		t.Fatalf("Failed to create parser: %s", err.Error())
	}
	defer parser.Free()

	if err = parser.ParseStringIntoModel(rdfxml_content, uri, model); err != nil {
		t.Fatalf("Failed to parse string into model: %s", err.Error())
	}

	var subject, predicate, object *Node
	var statement *Statement

	// build a statement out of a set of nodes
	subject, err = NewNodeFromUriString(world, "http://example.org/subject")
	predicate, err = NewNodeFromUriString(world, "http://example.org/pred1")
	object, err = NewNodeFromLiteral(world, "object")

	statement, err = NewStatementFromNodes(world, subject, predicate, object)
	if err != nil {
		t.Fatalf("Failed to create statement from nodes: %s", err.Error())
	}
	defer statement.Free() // note: this will free the attached nodes

	statementString, err := statement.ToString()
	if err != nil {
		t.Fatalf("Failed to create string representation of statement: %s", err.Error())
	}

	fmt.Printf("  Statement: %s\n", statementString)
	if err = model.AddStatement(statement); err != nil {
		t.Fatalf("Failed to add statement: %s", err.Error())
	}

	// print out the model
	fmt.Printf("%s: Resulting model is:\n%s", os.Args[0], model.ToString())

	if !model.ContainsStatement(statement) {
		fmt.Printf("  Model does not contain the statement: %s\n", statementString)
	} else {
		fmt.Printf("  Model contains the statement: %s\n", statementString)

		fmt.Printf("%s: Removing the statement", os.Args[0])
		if err = model.RemoveStatement(statement); err != nil {
			t.Fatalf("Failed to remove statement: %s", err.Error())
		}

		fmt.Printf("%s: Resulting model is:\n%s", os.Args[0], model.ToString())
	}
}

//Test_NewModelAddAndSerialize is based on Example3 from the librdf library and tests the following sequence:
//	- Creating an empty model 
//	- Adding a new statement
//	- Checking that the model contains the statement
//	- Serializing the model (in this case through ToString)
func Test_NewModelAddAndSerialize(t *testing.T) {
	storageType := "hashes"
	storageOptions := "hash-type='memory',dir='./testdata'"

	var err error
	var storage *Storage
	var model *Model

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", storageOptions); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	var subject, predicate, object *Node
	var statement *Statement

	// build a statement out of a set of nodes
	subject, err = NewNodeFromUriString(world, "http://example.org/subject")
	predicate, err = NewNodeFromUriString(world, "http://example.org/pred1")
	object, err = NewNodeFromLiteral(world, "object")
	statement, err = NewStatementFromNodes(world, subject, predicate, object)
	defer statement.Free()

	statementString, err := statement.ToString() 
	if err != nil {
		t.Fatalf("Failed to create string representation of statement: %s", err.Error())
	}
	fmt.Printf("  Statement: %s\n", statementString)

	if err = model.AddStatement(statement); err != nil {
		t.Fatalf("Failed to add statement: %s", err.Error())
	}

	// print out the model
	fmt.Printf("%s: Resulting model is:\n%s", os.Args[0], model.ToString())
}

//Test_ParseAndSerialize is based on Example4 from the librdf library and tests the following sequence:
//	- Parsing RDFXML into a model 
//	- Serializing the model out to another RDFXML file
func Test_ParseAndSerialize(t *testing.T) {
	var uri *Uri
	var baseUri *Uri

	storageType := "hashes"

	// TODO: hash-type 'bdb'
	storageOptions := "hash-type='memory',dir='.'"

	var err error
	var storage *Storage
	var model *Model
	var parser *Parser
	var serializer *Serializer

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", storageOptions); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	if uri, err = NewUri(world, testLocalUriFile); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	if parser, err = NewParser(world, "rdfxml", "application/rdf+xml"); err != nil {
		t.Fatalf("Failed to create parser: %s", err.Error())
	}
	defer parser.Free()

	if err = parser.ParseIntoModel(uri, uri, model); err != nil {
		t.Fatalf("Failed to parse string into model: %s", err.Error())
	}

	if serializer, err = NewSerializer(world, "rdfxml", "", nil); err != nil {
		t.Fatalf("Failed to create serializer: %s", err.Error())
	}
	defer serializer.Free()

	if baseUri, err = NewUri(world, "http://exampe.org/base.rdf"); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer baseUri.Free()

	var modelString string
	if modelString, err = serializer.SerializeModelToString(model, baseUri); err != nil {
		t.Fatalf("Failed to serialize model: %s", err.Error())
	}

	fmt.Printf("Serialised model: %s", modelString)
}

//Test_ParseAndSparqlQuery is based on Example5 from the librdf library and tests the following sequence:
//	- Parsing RDFXML into a model 
//	- Executing a SPARQL query and iterating the results
func Test_ParseAndSparqlQuery(t *testing.T) {

	storageType := "memory"
	storageOptions := ""

	var err error
	uriString := testLocalUriFile

	// create a new world
	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	var uri *Uri
	if uri, err = NewUri(world, uriString); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	// construct a storage provider
	var storage *Storage
	if storage, err = NewStorage(world, storageType, "test", storageOptions); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	var model *Model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	parserName := "rdfxml"
	parser, err := NewParser(world, parserName, "")
	if err != nil {
		t.Fatalf("Error constructing a parser", err.Error())
	}
	defer parser.Free()

	if err = parser.ParseIntoModel(uri, nil, model); err != nil {
		t.Fatalf("Error parsing uri into model", err.Error())
	}

	queryString := "select ?p ?o where { <http://purl.org/net/dajobe/> ?p ?o}" //"select ?p ?o where (<http://purl.org/net/dajobe/> ?p ?o)"
	query, err := NewQuery(world, "sparql", queryString)
	if err != nil {
		t.Fatalf("Error creating query :%s", err.Error())
	}

	chanQueryResultItems, err := model.ExecuteQueryToResultsChannel(&query, 100)
	if err != nil {
		t.Fatalf("Error executing query: %s", err.Error())
	}

	count := 0
	for queryResultItem := range chanQueryResultItems {
		fmt.Printf("  result: [")

		for _, nameNodePair := range queryResultItem.NameNodePairs {
			nodeString, err := nameNodePair.Node.ToString() 
			if err != nil {
				t.Fatalf("Failed to create string representation of node: %s", err.Error())
			}			
			fmt.Printf("%s=%s", nameNodePair.Name, nodeString)
		}

		fmt.Printf("]")
		count = count + 1
	}
	fmt.Printf("Query returned %d results\n", count)
}

//Test_ModelLoadAndSerialize is based on Example6 from the librdf library and tests the following sequence:
//	- Using the model.Load method to load the contents of a URI into a model 
//	- Serializing the model using ToString()
func Test_ModelLoadAndSerialize(t *testing.T) {
	var uri *Uri

	storageType, storageOptions := "memory", ""

	var err error
	var storage *Storage
	var model *Model

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", storageOptions); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	if uri, err = NewUri(world, testRemoteUri3); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	if err = model.Load(uri); err != nil {
		t.Fatalf("Failed to load model from URI: %s", err.Error())
	}

	modelString := model.ToString()
	fmt.Printf("Serialised model: %s", modelString)
}

//Test_ModelFileStorageAddStatement is based on Example7 from the librdf library and tests the following sequence:
//	- Creating a model backed with file storage 
//	- Adding a statement to the model
func Test_ModelFileStorageAddStatement(t *testing.T) {

	storageType, storageOptions := "file", ""

	var err error
	var storage *Storage
	var model *Model

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "./testoutput/Test_ModelFileStorageAddStatement_file.rdf", storageOptions); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	subject, err := NewNodeFromUriString(world, "http://www.dajobe.org/")
	predicate, err := NewNodeFromUriString(world, "http://purl.org/dc/elements/1.1/title")
	object, err := NewNodeFromLiteral(world, "My Home Page")

	statement, err := NewStatementFromNodes(world, subject, predicate, object)
	defer statement.Free()

	model.AddStatement(statement)
}

//Test_ParseAndQueryToFormattedString is based on Example8 from the librdf library and tests the following sequence:
//	- Parsing the contents of a URI into a model 
//	- Creating a query
//  - Executing the query and retrieving the results as a formatted string (JSON)
func Test_ParseAndQueryToFormattedString(t *testing.T) {
	storageType := "memory"

	var err error
	uriString := testRemoteUri

	// create a new world
	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	var uri *Uri
	if uri, err = NewUri(world, uriString); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	// construct a storage provider
	var storage *Storage
	if storage, err = NewStorage(world, storageType, "test", ""); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	var model *Model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to construct model: %s", err.Error())
	}
	defer model.Free()

	parserName := "rdfxml"
	var parser *Parser
	if parser, err = NewParser(world, parserName, ""); err != nil {
		t.Fatalf("Failed to create parser: %s", err.Error())
	}
	defer parser.Free()
	fmt.Printf("created parser")

	if err = parser.ParseIntoModel(uri, nil, model); err != nil {
		t.Fatalf("Failed to parse uri into model: %s", err.Error())
	}

	queryString := "select * where { ?s ?p ?o}"
	var query Query
	if query, err = NewQuery(world, "sparql", queryString); err != nil {
		t.Fatalf("Failed to create query: %s", err.Error())
	}

	format := "json"
	var resultString string
	if resultString, err = model.ExecuteQueryToFormattedString(&query, format); err != nil {
		t.Fatalf("Failed to execute query to formatted string: %s", err.Error())
	}
	fmt.Printf("results1: |%s|", resultString)
}

//Test_ParseStringIntoModel tests the following sequence:
//	- Loading the contents at a URI into a Go string 
//	- Parsing the string into a model
func Test_ParseStringIntoModel(t *testing.T) {
	uriString := testRemoteUri4
	storageType := "memory"
	parserName := "raptor"

	var uri *Uri
	var storage *Storage
	var model *Model
	var parser *Parser
	var err error

	resp, err := http.Get(uriString)
	if err != nil {
		t.Fatalf("Failed to http.Get URI: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body for URI: %s", err.Error())
	}
	bodyString := fmt.Sprintf("%s", body)

	world := NewWorld()

	if err = world.Open(); err != nil {
		t.Fatalf("World failed to open: %s", err.Error())
	}
	defer world.Close()

	// construct a URI based on the uri arg
	if uri, err = NewUri(world, uriString); err != nil {
		t.Fatalf("Failed to create URI: %s", err.Error())
	}
	defer uri.Free()

	// construct a storage provider
	if storage, err = NewStorage(world, storageType, "test", ""); err != nil {
		t.Fatalf("Failed to create storage: %s", err.Error())
	}
	defer storage.Free()

	// construct a model
	if model, err = NewModel(world, storage, ""); err != nil {
		t.Fatalf("Failed to create model: %s", err.Error())
	}
	defer model.Free()

	// and a parser
	if parser, err = NewParser(world, parserName, ""); err != nil {
		t.Fatalf("Failed to create parser: %s", err.Error())
	}
	defer parser.Free()

	if err = parser.ParseStringIntoModel(bodyString, uri, model); err != nil {
		t.Fatalf("Failed to parser string into model: %s", err.Error())
	}

	// print out the model
	fmt.Printf("Resulting model is:\n%s", model.ToString())
}
