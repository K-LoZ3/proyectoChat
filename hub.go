// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	//este canal es basicamente para notificar que se crea un cliente.
	register chan *Client

	// Unregister requests from clients.
	//este canal seria para eliminarlo.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		//caso para registrar cliente.
		//almacenamos el cliente en el arreglo.
		case client := <-h.register:
			h.clients[client] = true
			
		//caso para eliminar cliente.
		//Eliminamos el cliente del mapa.
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			
		//caso del mensaje.
		//compartimos el mensaje a cada uni de los cliente.
		//Recorremos el map de clientes y le enviamos por canal el mensaje.
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				//el select es para asegurar que solo se envia mensaje
				//cuando la variable tenga datos que enviar.
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}