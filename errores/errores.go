package errores

import (
  "net/http"
  "fmt"
)

//Pueba de manejo de errores con una arreglo de estructuras de errores.

const (
  ErrorBadRequest = iota
  ErrorInternalDataBase
  ErrorUserNotFound
  ErrorBadFormatUser
  ErrorBadFormatPassword
  ErrorInvalidToken
  ErrorRegister
  ErrorLogin
)


type err struct {
  title: string
  msg: string
  status int
}

var errors = []err{
  err{ "Bad request", "Error en la informacion compartida por el usuario", http.StatusBadRequest },
  err{ "Internal data base", "Error en la base de datos", http.StatusInternalServerError },
  err{ "User not found", "Usuario incorrecto", http.StatusNotFound },
  err{ "Bad format user", "Solo se permiten mayusculas, minusculas y caracteres especiales, @ # ? _ - .", http.StatusBadRequest },
  err{ "Invalid format password", "Debe contener al menos una mayuscula, una minuscula y un caracter especial @ # ? _ - .", http.StatusBadRequest },
  err{ "Invalid token", "Token invalido o expirado.", http.StatusUnauthorized},
  err{ "Register", "Registro incompleto.", http.StatusInternalServerError },
  err{ "Login", "Error al hacer login", http.StatusInternalServerError },
}

func ErrorCode(code int) error {
  if code < 0 {
    code *= -1
  }
  if code > len(errors) - 1 {
    //manejamos el error de codigo si no lo existe codigo de error
    return fmt.Errorf("Error no identificado.")
  }
  
  e := errors[code]
  
  return fmt.Errorf("Error: %s, descripcion: %s, status: %v", e.title, e.msg, e.status)
}

//Manejo de errores con errores personalizados.
type ErrorStruct struct {
  Title string
  Msg string
  Status int
}

func (e ErrorStruct) Error() string {
  return fmt.Sprintf("Error: %s, descripcion: %s, status: %v", e.Title, e.Msg, e.Status)
}

/*
func New(title string, msg string, status int) error {
  if title == "" {
    title = "Error no definido."
  }
  return ErrorStruct{title, msg, status}
}*/

func NewStruct(title string, msg string, status int) ErrorStruct {
  if title == "" {
    title = "Error no definido."
  }
  return ErrorStruct{title, msg, status}
}

func NewCode(i int) ErrorStruct {
  if i < 0 || i > len(errors) - 1 {
    //manejamos el error de codigo si no lo existe codigo de error
    return NewStruct("", "", -1)
  }
  
  e := errors[i]
  
  return NewStruct(e.title, e.msg, e.status)
}

func WriteHTTP(w http.ResponseWriter, title string, msg string, status int) {
  e := NewStruct(title, msg, status)
  
  http.Error(w, e.Error(), e.status)
}

func WriteError(w http.ResponseWriter, e ErrorStruct) {
  http.Error(w, e.Error(), e.Status)
}