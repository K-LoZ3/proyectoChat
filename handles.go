//el archivo separa los handle para las funciones
//y tenerlas ordenadas.
package main

import (
  "encoding/json"
  "net/http"
  "log"
  
  "Golang/Practicas/chat/data"
)

//struct para administrar la base de datos con los usuarios
type users struct {
  Id int
  Username string
  Password string
}

//handleRegistro maneja la ruta para registrar un usuario pasado en formato json
func handleRegistro(w http.ResponseWriter, r *http.Request) {
  var u users
  
  //decodificamos el json del usuario.
  err := json.NewDecoder(r.Body).Decode(&u)
  if err != nil {
    http.Error(w, "Error en el usuario", http.StatusBadRequest)
  }
  
  //si el string del usuario no cumple con los carecteres retornamos el error
  if !regexpUsuario(u.Username) {
    http.Error(w, "Error, el usuario solo debe contener numeros letras o '_'", http.StatusBadRequest)
    return
  }
  
  //retornamos el error si la contraseña no cumple con los carecteres.
  if !regexpPassword(u.Password) {
    http.Error(w, "Error, la contraseña debe tener mayuscula, minuscula, numero y caracter especial.", http.StatusBadRequest)
    return
  }
  
  //Almacenamos el usuario en la base de datos.
  err = data.Register(u.Username, u.Password)
  if err != nil {
    http.Error(w, "Error al guardar usuario en la base de datos", http.StatusBadRequest)
    return
  }
  
  //devolvemos el estado de creado.
  w.WriteHeader(http.StatusCreated)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
  var u users
  
  err := json.NewDecoder(r.Body).Decode(&u)
  if err != nil {
    http.Error(w, "Error al leer el usuario del json", http.StatusBadRequest)
    return
  }
  
    //si el string del usuario no cumple con los carecteres retornamos el error
  if !regexpUsuario(u.Username) {
    http.Error(w, "Error, el usuario solo debe contener numeros letras o '_'", http.StatusBadRequest)
    return
  }
  
  //retornamos el error si la contraseña no cumple con los carecteres.
  if !regexpPassword(u.Password) {
    http.Error(w, "Error, la contraseña debe tener mayuscula, minuscula, numero y caracter especial.", http.StatusBadRequest)
    return
  }
  
  err = data.FindUser(u.Username, u.Password)
  if err != nil {
    http.Error(w, "Error, usuario o contraseña incorrecta", http.StatusBadRequest)
    return
  }
  
  token, err := crearJWT(u.Username)
  if err != nil {
    http.Error(w, "Error al crear el token", http.StatusBadRequest)
    return
  }

  http.SetCookie(w, &http.Cookie{
      Name:     "auth_token",
      Value:    token,
      Path:     "/",
      HttpOnly: true,
      Secure:   false,
  })
  http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

//funcion para servir el archivo html
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
  
  //Archivo a servir
	http.ServeFile(w, r, "home.html")
}