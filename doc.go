/*
Package golibrdf provides go language bindings, tests and examples for the
Redland RDF library (see http://librdf.org)

Please refer to the tests within golibrdf_test.go which also serve as examples
of various usage scenarios.  These tests are based on corresponding examples
within Redland RDF itself.

Prerequisites:
	- Redland libRDF, Rasqal and Raptor libraries must be installed first.
		- Refer to instructions at http://librdf.org

Once the prequisite libraries have been installed, golibrdf can be installed with one of the following:
Linux:
	go get github.com/PhillP/golibrdf
	(this relies on pkg-config for Redland library locations)
Windows / OSX / Explicit library locations:
	Inform Go of the Redland library locations using CGO_CFLAGS and CGO_LDFLAGS by modifying the paths in the following example
	CGO_CFLAGS="-I/usr/local/include/ -I/usr/local/include/raptor2 -I/usr/local/include/rasqal" CGO_LDFLAGS="-L/usr/local/lib" go get github.com/PhillP/golibrdf

Refer to LICENSE.txt for license information.
*/
package golibrdf
