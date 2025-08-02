//el archivo separa los handle para las funciones
//y tenerlas ordenadas.
package main

import (
  "encoding/json"
  "context"
  "net/http"
  "log"
  "fmt"
  
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
    errorStr := fmt.Sprintf("Error al guardar usuario en la base de datos, %v", err)
    http.Error(w, errorStr, http.StatusBadRequest)
    return
  }
  
  //devolvemos el estado de creado.
  http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//handleLogin Se encarga de manejar la ruta para logear a un usuario.
//Valida que la informacion recibida del usuario y la compara con la
//base de datos para luego devolver un jwt firmado.
func handleLogin(w http.ResponseWriter, r *http.Request) {
  var u users
  //recibimos los datos de usuario desde un json
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
  
  //Buscamos el usuario en la base de datos, y obtenemos el hash para
  //comparar las contraseñas.
  err = data.FindUser(u.Username, u.Password)
  if err != nil {
    http.Error(w, "Error, usuario o contraseña incorrecta", http.StatusBadRequest)
    return
  }
  
  //Creamos un JWT firmado con el nombre de usuario y tiempo de expiracion.
  //el tiempo esta harcodeado en 2 horas de momento para pruebas en la funcion.
  token, err := crearJWT(u.Username)
  if err != nil {
    http.Error(w, "Error al crear el token", http.StatusBadRequest)
    return
  }
  
  //agregamos el token a las cookies del navegador del usuario.
  //esta opcion sera para este proyecto ya que asi el usuario no tiene
  //que agregar el token manualmente porque no es una api.
  http.SetCookie(w, &http.Cookie{
      Name:     "auth_token",
      Value:    token,
      Path:     "/",
      HttpOnly: true,
      Secure:   false,
  })
  //redireccionamos a la ruta del chat.
  http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

//funcion para servir el archivo html
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
  
  //Archivo a servir
	http.ServeFile(w, r, "home.html")
}

func authMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("auth_token")
    if err != nil {
      http.Error(w, "Error al leer token", http.StatusBadRequest)
      http.Redirect(w, r, "/login", http.StatusSeeOther)
      return
    }
    
    username, err := validarJWT(cookie.Value)
    if err != nil {
      errorStr := fmt.Sprintf("Error al validar token %v", err)
      http.Error(w, errorStr, http.StatusBadRequest)
      http.Redirect(w, r, "/login", http.StatusSeeOther)
      return
    }
    
    //Lo convertimos a contexto para pasarlo a al siguiente HandlerFunc
    ctx := context.WithValue(r.Context(), "username", username)
    
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}