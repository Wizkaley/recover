package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	go main()
	dashtest.ControlCoverage(m)
	time.Sleep(1 * time.Second)
}

func TestPanicAfter(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic-after", nil)

	if err != nil {
		t.Errorf("Error while Creating test Request : %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(PanicAfterDemo)
	defer func() {
		r := recover()
		if r != nil {
			log.Println("Recovered from Panic", r)
		}
	}()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		log.Printf("Expected %q but got %q", http.StatusOK, rr.Code)
	}
}
