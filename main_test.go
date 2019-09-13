package main

import (
	//"httptest"

	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestPanicDemo(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Errorf("The Error is %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(PanicDemo)

	//handler.ServeHTTP(rr, req)
	defer func() {
		r := recover()
		if r != nil {
			//t.Errorf("Oh LOrd it panicked but i cleaned it Up : %v", r)
			if rr.Code == 200 {
				t.Log("We Passed The Test")
			}
		}
	}()
	handler.ServeHTTP(rr, req)
	t.Errorf("Did Not Panic")
}

func controller() http.Handler {
	srv := mux.NewRouter()

	srv.HandleFunc("/panic", PanicDemo).Methods("GET")
	srv.HandleFunc("/debug", RenderSourceCode)
	srv.HandleFunc("/panic-after", PanicAfterDemo)
	srv.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3001", RecoverMw(srv, true)))
	return srv
}

func TestRenderSrcCode(t *testing.T) {

	req, err := http.NewRequest("GET", "/debug?line=24&path=%2Fusr%2Flocal%2Fgo%2Fsrc%2Fruntime%2Fdebug%2Fstack.go", nil)
	if err != nil {
		t.Errorf("Cant Create new Request : %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(RenderSourceCode)

	handler.ServeHTTP(rr, req)

	if rr.Code == 200 {
		fmt.Println("We Passed the second test")
	}
}
func TestPanicAfterDemo(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic-after", nil)
	if err != nil {
		t.Errorf("The Error is %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(PanicDemo)

	//handler.ServeHTTP(rr, req)
	defer func() {
		r := recover()
		if r != nil {
			//t.Errorf("Oh LOrd it panicked but i cleaned it Up : %v", r)
			if rr.Code == 200 {
				t.Log("We Passed The Test")
			}
		}
	}()
	handler.ServeHTTP(rr, req)
	t.Errorf("Did Not Panic")
}

// func TestMain(m *testing.M) {

// }
func TestFuncThatPanics(t *testing.T) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("We Recovered From teh Panic")
		}
	}()
	FuncThatPanics()
}
