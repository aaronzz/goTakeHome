package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// base test case with 3 * 3 matrix
func TestHandleEcho(t *testing.T) {
	input := "1,2,3\n4,5,6\n7,8,9\n"
	expectedOutput := "1,2,3\n4,5,6\n7,8,9\n"
	testHandleAction(t, "/echo", handleEcho, input, expectedOutput)
}

func TestHandleInvert(t *testing.T) {
	input := "1,2,3\n4,5,6\n7,8,9\n"
	expectedOutput := "1,4,7\n2,5,8\n3,6,9\n"
	testHandleAction(t, "/invert", handleInvert, input, expectedOutput)
}

func TestHandleFlatten(t *testing.T) {
	input := "1,2,3\n4,5,6\n7,8,9\n"
	expectedOutput := "1,2,3,4,5,6,7,8,9\n"
	testHandleAction(t, "/flatten", handleFlatten, input, expectedOutput)
}

func TestHandleSum(t *testing.T) {
	input := "1,2,3\n4,5,6\n7,8,9\n"
	expectedOutput := "45\n"
	testHandleAction(t, "/sum", handleSum, input, expectedOutput)
}

func TestHandleMultiply(t *testing.T) {
	input := "1,2,3\n4,5,6\n7,8,9\n"
	expectedOutput := "362880\n"
	testHandleAction(t, "/multiply", handleMultiply, input, expectedOutput)
}

// test case with 4 * 4 matrix
func TestHandleEcho4(t *testing.T) {
	input := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	expectedOutput := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	testHandleAction(t, "/echo", handleEcho, input, expectedOutput)
}

func TestHandleInvert4(t *testing.T) {
	input := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	expectedOutput := "1,5,9,13\n2,6,10,14\n3,7,11,15\n4,8,12,16\n"
	testHandleAction(t, "/invert", handleInvert, input, expectedOutput)
}

func TestHandleFlatten4(t *testing.T) {
	input := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	expectedOutput := "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16\n"
	testHandleAction(t, "/flatten", handleFlatten, input, expectedOutput)
}

func TestHandleSum4(t *testing.T) {
	input := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	expectedOutput := "136\n"
	testHandleAction(t, "/sum", handleSum, input, expectedOutput)
}

func TestHandleMultiply4(t *testing.T) {
	input := "1,2,3,4\n5,6,7,8\n9,10,11,12\n13,14,15,16\n"
	expectedOutput := "20922789888000\n"
	testHandleAction(t, "/multiply", handleMultiply, input, expectedOutput)
}

func TestHandleError(t *testing.T) {
	input := "1,2\n4,5,6\n7,8,9\n"
	expectedOutput := "error record on line 2: wrong number of fields"
	testHandleAction(t, "/echo", handleEcho, input, expectedOutput)
}

func TestHandleErrorEmptyFile(t *testing.T) {
	input := ""
	expectedOutput := "error: Empty CSV file"
	testHandleAction(t, "/echo", handleEcho, input, expectedOutput)
}

func testHandleAction(t *testing.T, urlPath string, operation func(w http.ResponseWriter, r *http.Request), input string, expectedOutput string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Errorf("error creating form file: %s", err.Error())
	}
	_, err = io.Copy(part, strings.NewReader(input))
	if err != nil {
		t.Errorf("error copying input to form file: %s", err.Error())
	}
	writer.Close()

	req, err := http.NewRequest("POST", urlPath, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err != nil {
		t.Errorf("error creating request: %s", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(operation)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expectedOutput {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedOutput)
	}
}
