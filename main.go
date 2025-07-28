//El codigo de chat es del ejemplo de gorilla/websocket
//estare modificando el protecto a partir de aqui.

package main

import (
  "time"
	"log"
	"net/http"
	// Declarada pero aun no la he usado ya que solo se creo el paquete para
	//administrar la base de datos.
	_ "Golang/Practicas/chat/data"
	
	"github.com/go-chi/chi/v5"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
  
  //Archivo a servir
	http.ServeFile(w, r, "home.html")
}

func main() {
  //Como el binario lo ejecuto desde termux en mi celular
  //para que se pueda acceder desde otro dispositivo hay que usar esta ip.
  addr := "0.0.0.0:8080"
  
  //reouter para manejar las rutas
  r :=  chi.NewRouter()
  
  //hub. Se encarga de manejar de administrar los clientes.
  //los almacena, los elimina o les comarte el mensaje que envie uno de ellos.
  //con esto declaramos las variables del hub
	hub := newHub()
	
	//corremos el administrador, esta es la funcion que almacena clientes, elimina clientes o comparte el mensaje entrelos demas
	go hub.run()
	
	//ruta principal. 
	r.Route("/chat", func(r chi.Router) {
	  //con el metodo get para la ruta principal servimos el archivo
	  //Este se encargara de redirigir la pagina por el protocolo ws:// 
	  r.Get("/", serveHome)
	  
	  //esta es la ruta que se maneja desde el protocolo ws.
  	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
  	  //desde el protocolo ws:// en la rita /ws.
  	  //los clientes que entren en esta ruta seran creados y almacenados con esta funcion.
  		serveWs(hub, w, r)
  		
  	})
  	
	})
	
	//creo un servidor con las caracteristicas
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