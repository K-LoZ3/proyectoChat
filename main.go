//El codigo de chat es del ejemplo de gorilla/websocket
//estare modificando el protecto a partir de aqui.

package main

import (
  "time"
	"log"
	"net/http"
	
	"github.com/go-chi/chi/v5"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	http.ServeFile(w, r, "./home.html")
}

func main() {
  addr := "10.254.97.246:8080"
  
  r :=  chi.NewRouter()
  
	hub := newHub()
	go hub.run()
	
	r.Route("/", func(r chi.Router) {
	  r.Get("/", serveHome)
	  
  	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
  		serveWs(hub, w, r)
  		
  	})
  	
	})
	
	
	server := http.Server{
    Addr: addr,
    Handler: r,
    WriteTimeout: 10 * time.Second,
    ReadTimeout: 10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  
  log.Println("Listening in http://" + addr + "/...")
  log.Fatal(server.ListenAndServe())
}