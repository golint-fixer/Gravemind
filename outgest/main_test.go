package main

import (
	"net/http/httptest"
	"testing"
)

// TODO: Write better tests. This is ugly!
func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	handler(w, nil)

	v1 := w.Code
	e1 := 200
	if v1 != e1 {
		t.Error("For handler body: expected", e1, "got", v1)
	}

	v2 := w.Body.String()
	e2 := "Hello World!\n"
	if v2 != e2 {
		t.Error("For handler body: expected", e2, "got", v2)
	}
}
