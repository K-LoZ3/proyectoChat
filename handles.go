//el archivo separa los handle para las funciones
//y tenerlas ordenadas.
package main

import (
  "encoding/json"
  "net/http"
  
  "Golang/Practicas/chat/data"
)

type users struct {
  Id int
  Username string
  Password string
}

func handleRegistro(w http.ResponseWriter, r *http.Request) {
  var u users
  
  err := json.NewDecoder(r.Body).Decode(&u)
  if err != nil {
    http.Error(w, "Error en el usuario", http.StatusBadRequest)
  }
  
  if !regexpUsuario(u.Username) {
    http.Error(w, "Error, el usuario solo debe contener numeros letras o '_'", http.StatusBadRequest)
    return
  }
  
  if !regexpPassword(u.Password) {
    http.Error(w, "Error, la contraseña debe tener mayuscula, minuscula, numero y caracter especial.", http.StatusBadRequest)
    return
  }
  
  err = data.Register(u.Username, u.Password)
  if err != nil {
    http.Error(w, "Error al guardar usuario en la base de datos", http.StatusBadRequest)
    return
  }
  
  w.WriteHeader(http.StatusCreated)
}