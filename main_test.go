package main

import (
	"errors"
	"io"
	"runtime/debug"

	"github.com/alecthomas/chroma"
	"github.com/stretchr/testify/assert"

	//"httptest"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevNeg(t *testing.T) {

	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Errorf("Error Occurred while creaing test Request : %v", err)
	}
	rr := httptest.NewRecorder()

	panHandler := http.HandlerFunc(PanicDemo)
	// Just Passing the Dev flag as false here everything else remains the same as the
	// recoverMw testcase
	handler := http.HandlerFunc(RecoverMw(panHandler, false))
	handler.ServeHTTP(rr, req)
}

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

// func controller() http.Handler {
// 	srv := mux.NewRouter()

// 	srv.HandleFunc("/panic", PanicDemo).Methods("GET")
// 	srv.HandleFunc("/debug", RenderSourceCode)
// 	srv.HandleFunc("/panic-after", PanicAfterDemo)
// 	srv.HandleFunc("/", hello)
// 	log.Fatal(http.ListenAndServe(":3001", RecoverMw(srv, true)))
// 	return srv
// }

func TestRenderSrcCode(t *testing.T) {

	req, err := http.NewRequest("GET", "/debug?line=24&path=/usr/local/go/src/runtime/debug/stack.go", nil)
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

func TestMakeLinks(t *testing.T) {
	stack := debug.Stack()

	links := MakeLinks(string(stack))

	assert.NotEqualf(t, "", links, "Expected Links but got nothing")
}

func TestRecoverMw(t *testing.T) {

	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Errorf("Error Occurred while creaing test Request : %v", err)
	}
	rr := httptest.NewRecorder()

	panHandler := http.HandlerFunc(PanicDemo)
	handler := http.HandlerFunc(RecoverMw(panHandler, true))
	handler.ServeHTTP(rr, req)
}

func TestSourceCodeHandler(t *testing.T) {
	tst := []struct {
		ucase string
		link  string
		Code  int
	}{
		{"successful", "/debug?line=24&path=/usr/local/go/src/runtime/debug/stack.go", 200},
		{"lexer_error", "/debug?line=24&path=/usr/local/go/src/runtime/debug/stack.go", 500},
		{"io_error", "/debug?line=24&path=/usr/local/go/src/runtime/debug/stack.go", 500},
		{"got", "/debug?path=/home/wiz/go/src/recover/main.go", 200},
		{"got", "/debug?path=/home/wiz/go/src/recover/main1.go", 500},
		{"got", "/debug", 500},
	}

	tempIoCopy := IOCopy
	tempLexer := LexerTokenise

	for _, item := range tst {
		if item.ucase == "io_error" {
			IOCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return -1, errors.New("IO Error Occured")
			}
		}
		if item.ucase == "lexer_error" {
			LexerTokenise = func(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
				return nil, errors.New("Lexer Tokenisation error")
			}
		}

		req, err := http.NewRequest("GET", item.link, nil)
		if err != nil {
			t.Errorf("Error while Creating Request : %v", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(RenderSourceCode)

		handler.ServeHTTP(rr, req)

		assert.Equalf(t, item.Code, rr.Code, "Status Code expected %v but got %v", item.Code, rr.Code)
		IOCopy = tempIoCopy
		LexerTokenise = tempLexer
	}
}
