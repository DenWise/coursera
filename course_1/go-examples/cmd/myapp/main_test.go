package main

import (
	"bufio"
	"bytes"
	//"os"
	"strings"
	"testing"
)

var testOk = `1
2
3
3
3
4
5
5
6`

var testOkResult = `1
2
3
4
5
6
`

func TestOk(t *testing.T) {

	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)

	err := uniq(in,out)
	if err != nil {
		t.Errorf("Test for OK is Failed.")
	}

	result := out.String()
	if result != testOkResult {
		t.Errorf("Test for OK failed - results not matched\n %v %v\n",result, testOkResult)
	}
}

var notSorted = `1
4
3
6
2
5`

func TestErrNotSorted(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(notSorted))
	out := new(bytes.Buffer)

	err := uniq(in, out)
	if err == nil {
		t.Errorf("Test for OK failed - error\n %v\n",err)
	}
}