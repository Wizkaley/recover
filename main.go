package main

import(
	"fmt"
	"net/http"
	"log"
)

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/panic",panicDemo)
	mux.HandleFunc("/panic-after",panicAfterDemo)
	mux.HandleFunc("/",hello)
	log.Fatal(http.ListenAndServe(":3000",recoverMw(mux)))

}


func recoverMw(app http.Handler)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		defer func(){
			err := recover(); if err!= nil{
				log.Println(err)
				http.Error(w,"Something Went Wrong :| ",http.StatusInternalServerError)
			}
		}()

		app.ServeHTTP(w,r)
	}
}

func panicDemo(w http.ResponseWriter, r * http.Request){
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r * http.Request){
	fmt.Fprint(w,"<html>Helllo!</html>")
	funcThatPanics()
}


func funcThatPanics(){
	panic("Oh no!")	
}

func hello(w http.ResponseWriter , r * http.Request){
	fmt.Fprint(w,"<h1>Hello!</h1>")
}