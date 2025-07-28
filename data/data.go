//Paquete para manejar la base de datos. Pendiente optimizar las funciones.
//lo primero, que devuelvan el error.
package data

import (
  "database/sql"
  "log"
  
  _ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {
  var err error
  db, err = sql.Open("sqlite", "base.db")
  
  create := `
  CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
  );`
  
  _, err = db.Exec(create)
  if err != nil {
    log.Fatal("Error al crear e uniciar la tabla.", err)
  }
}

//var Close = db.Close
func Close() {
  db.Close()
}

func Register(username string, password string) bool {
  _, err := db.Exec("INSERT INTO users(username, password) VALUES( ?, ?)", username, password)
  if err != nil {
    log.Fatal("Error al guardar usuario y clave", err)
    return false
  }
  return true
}

func FindUser(username string, password string) bool {
  var u, p string
  err := db.QueryRow("SELECT username, password FROM users WHERE username = ? AND password = ?", username, password).Scan(&u, &p)
  if err != nil {
    log.Fatal("Usuario o clave invalida", err)
    return false
  }
  return true
}