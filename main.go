package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/styles"

	"github.com/alecthomas/chroma/lexers"

	"github.com/alecthomas/chroma/formatters/html"
)

var (
	// IOCopy is a variable in which the io.Copy signature is stored
	IOCopy   = io.Copy
	lexersgo = lexers.Get("go")
	// LexerTokenise is a variable in which the lexers.Get("go") funtion signature is stored
	LexerTokenise = lexersgo.Tokenise
)

func main() {
	m := Controller()
	log.Fatal(http.ListenAndServe(":3000", RecoverMw(m, true)))
}

// Controller Handles the Requests
func Controller() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/panic", PanicDemo)
	mux.HandleFunc("/debug", RenderSourceCode)
	mux.HandleFunc("/panic-after", PanicAfterDemo)

	return mux

}

// RenderSourceCode ...
func RenderSourceCode(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	//path := "/home/wiz/go/src/github.com/Wizkaley/recover/main.go"
	l := r.FormValue("line")
	line, err := strconv.Atoi(l)
	if err != nil {
		line = -1
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := bytes.NewBuffer(nil)

	_, err = IOCopy(b, file)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//_ = line.

	var lines [][2]int

	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	//exer := LexerTokenise
	iterator, err := LexerTokenise(nil, b.String())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	style := styles.Get("github")
	// if style == nil {
	// 	style = styles.Fallback
	// }

	formatter := html.New(html.TabWidth(2), html.HighlightLines(lines), html.WithLineNumbers())
	w.Header().Set("Content-Type", "text/html")
	formatter.Format(w, style, iterator)

	//err = quick.Highlight(w, b.String(), "go", "html", "monokai")
	// if err != nil {
	// 	log.Fatal(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)

	// 	return
	// }
	io.Copy(w, file)
}

// RecoverMw ...
func RecoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
				stack := debug.Stack()
				//fmt.Println(string(stack))
				if !dev {
					http.Error(w, "Something Went Wrong :| ", http.StatusInternalServerError)
				}

				fmt.Fprintf(w, "<h1>Panic : %v </h1><pre>Trace : %v</pre>", err, MakeLinks(string(stack)))
			}
		}()

		app.ServeHTTP(w, r)
	}
}

// MakeLinks ...
func MakeLinks(stack string) string {
	lines := strings.Split(stack, "\n")

	for li, line := range lines {

		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}
		var lineStr strings.Builder
		for i := len(file) + 2; i < len(line); i++ {

			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineStr.WriteByte(line[i])
		}
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineStr.String())
		lines[li] = "\t<a href=\"/debug?" + v.Encode() + "\">" + file + ":" + lineStr.String() + "</a>" + line[len(file)+2+len(lineStr.String()):]
	}
	return strings.Join(lines, "\n")
}

func PanicDemo(w http.ResponseWriter, r *http.Request) {
	FuncThatPanics()
}

func PanicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html>Hello!</html>")
	FuncThatPanics()
}

func FuncThatPanics() {
	panic("Oh no!")
}
