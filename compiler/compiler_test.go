package compiler

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const EXPECTED_SIMPLE_COMPILE = `body {
  font-weight: bold; }
`

func TestFindCompilable(t *testing.T) {
	ctx := NewSassContext(NewSassCommand(), "../integration/bad-src", "../integration/out")

	//The following line is expected to employ fileLogCompilationError which will use log to
	//output an error.  By redirecting the output of log temporarily, we can both test that
	//this error takes place and avoid outputing to stderr during a succesful test.

	//Set up error buffer.
	var buf bytes.Buffer
	log.SetOutput(&buf)

	//do the actual call
	compilable := findCompilable(ctx)

	//restore log output to its normal stderr
	log.SetOutput(os.Stderr)

	//and finally make sure we did get the output error we expected.
	if len(buf.String()) == 0 {
		t.Error()
	}

	if len(compilable) != 0 {
		t.Error()
	}

	ctx = NewSassContext(NewSassCommand(), "../integration/src", "../integration/out")
	compilable = findCompilable(ctx)

	if len(compilable) != 5 {
		t.Error()
	}

	if compilable["../integration/src/01.simple.scss"] != "../integration/out/01.simple.css" {
		t.Error()
	}

	if compilable["../integration/src/02.simple-import.scss"] != "../integration/out/02.simple-import.css" {
		t.Error()
	}

	if compilable["../integration/src/03.multiple-imports.scss"] != "../integration/out/03.multiple-imports.css" {
		t.Error()
	}

	if compilable["../integration/src/04.missing.scss"] != "../integration/out/04.missing.css" {
		t.Error()
	}

	if compilable["../integration/src/05.rawcss.scss"] != "../integration/out/05.rawcss.css" {
		t.Error()
	}
}

func TestCompile(t *testing.T) {
	t.Parallel()

	ctx := NewSassContext(NewSassCommand(), "../integration/src", "../integration/out")
	err := compile(ctx, "../integration/src/01.simple.scss", "../integration/out/01.simple.css")

	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadFile("../integration/out/01.simple.css")

	if err != nil {
		t.Error(err)
	}

	if string(b) != EXPECTED_SIMPLE_COMPILE {
		t.Errorf("Unexpected compiled results: %s", string(b))
	}
}
