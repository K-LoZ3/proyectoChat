//Paquete para manejar la base de datos. Pendiente optimizar las funciones.
//lo primero, que devuelvan el error.
package data

import (
  "database/sql"
  
  "golang.org/x/crypto/bcrypt"
  _ "modernc.org/sqlite"
)

//Declaro la bariable de esta manera ya que quiero que el codigo del mamejo
//de la base de datos desde el paquete main sea muy sensillo.
var db *sql.DB

//Inicia la base de datos con una tabla para usuarios y conyraseñas.
func InitDB() error {
  var err error
  
  //abrimos a creamos el archivo para la base de dstos.
  db, err = sql.Open("sqlite", "base.db")
  if err != nil {
    return err
  }
  
  //Intruccion para la tabla.
  create := `
  CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
  );`
  
  //creamos la tabla.
  _, err = db.Exec(create)
  if err != nil {
    return err
  }
  
  return nil
}

func Close() {
  db.Close()
}

func Register(username string, password string) error {
  //Se encripta la contraseña entes de guardarla
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return err
  }
  //sobreescribimos la contraseña
  password = string(hash)
  
  //almacenamos en la base de datos
  _, err = db.Exec("INSERT INTO users(username, password) VALUES( ?, ?)", username, password)
  if err != nil {
    return err
  }
  
  return nil
}

func FindUser(username string, password string) error {
  var p string
  //buscamos en la base de datos el usuario seleccionamos la contraseña
  //para luego compararla
  err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&p)
  if err != nil {
    return err
  }
  
  //conparamos la contraseña con la almacenada.
  return bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
}