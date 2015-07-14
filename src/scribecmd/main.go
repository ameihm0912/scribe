// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package main

import (
	"flag"
	"fmt"
	"os"
	"scribe"
)

var flagDebug bool

func failExit(t scribe.TestResult) {
	fmt.Fprintf(os.Stdout, "error: test result for \"%v\" was unexpected, exiting\n", t.Name)
	os.Exit(2)
}

func main() {
	var (
		docpath      string
		expectedExit bool
		testHooks    bool
		showVersion  bool
	)

	err := scribe.Bootstrap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	flag.BoolVar(&flagDebug, "d", false, "enable debugging")
	flag.BoolVar(&expectedExit, "e", false, "exit if result is unexpected")
	flag.StringVar(&docpath, "f", "", "path to document")
	flag.BoolVar(&testHooks, "t", false, "enable test hooks")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Fprintf(os.Stdout, "scribe %v\n", scribe.Version)
		os.Exit(0)
	}

	if flagDebug {
		scribe.SetDebug(true, os.Stderr)
	}

	if docpath == "" {
		fmt.Fprintf(os.Stderr, "error: must specify document path\n")
		os.Exit(1)
	}

	scribe.TestHooks(testHooks)

	fd, err := os.Open(docpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()
	doc, err := scribe.LoadDocument(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// In expectedExit mode, set a callback in the scribe module that will
	// be called immediately during analysis if a test result does not
	// match the boolean expectedresult parameter in the test. The will
	// result in the tool exiting with return code 2.
	if expectedExit {
		scribe.ExpectedCallback(failExit)
	}

	err = scribe.AnalyzeDocument(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, x := range doc.GetTestNames() {
		scribe.GetResults(doc, x)
	}

	os.Exit(0)
}
