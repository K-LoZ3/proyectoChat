package main

import (
  "regexp"
  "os"
  "time"
  
  "Golang/Practicas/chat/errores"
  
  "github.com/golang-jwt/jwt"
)

//regexpUsuario comprueba que el nombre de usuario tengan un formato expecifico.
func regexpUsuario(s string) bool {
  //Caracteres que se pueden usar [A-Za-z\d_]
  //Longitud minima 4 y maxima 9 {4,9}
  regPass := regexp.MustCompile(`^[A-Za-z\d_]{4,9}$`)
  
  //comparamos la clave con la expresion regular
  return regPass.MatchString(s)
}

//regexpPassword conprueba que la contraseña tenga el formato expecifico
func regexpPassword(s string) bool {
  //al menos una minuscula
  minuscula := regexp.MustCompile(`[a-z]+`)
  
  //al menks una mayuscula
  mayuscula := regexp.MustCompile(`[A-Z]+`)
  
  //al menos un digito
  digito  := regexp.MustCompile(`\d+`)
  
  //debe iniciar con una letra
  inicio  := regexp.MustCompile(`^[A-Za-z]+`)
  
  //al menos una caracter especial
  especial  := regexp.MustCompile(`[@\_\-.#?]+`)
  
  //solo los permitido con un minimo de 4 a 9 caracteres.
  caracteres  := regexp.MustCompile(`^[a-zA-Z0-9@#?\_\-.]{5,12}$`)
  
  return mayuscula.MatchString(s) && minuscula.MatchString(s) && digito.MatchString(s) && inicio.MatchString(s) && especial.MatchString(s) && caracteres.MatchString(s)
}

//crearJWT devuelve un jwt firmado con la variable de entorno y el
//nombre de usuario con 2 hora de vencimiento
func crearJWT(n string) (string, error) {
  //buscamos en .env la frase para firmar el jwt
  firma := os.Getenv("FRASE")
  
  //preparamos los datos para el jwt: el nombre de usuario y el tiempo maximo.
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": n,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	})
	
	//firmamos e token.
	tokenString, err := token.SignedString([]byte(firma))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

//validarJWT verifica el el token este correcto.
func validarJWT(t string) (string, error) {
  firma := os.Getenv("FRASE") //obtenemos la firma del .env
  //paseamos el token.
  token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
      return []byte(firma), nil
  })
  if err != nil || !token.Valid {
    return "", err
  }
    
  //convertimos el token a un map claims para extraer los datos que
  //incluimos en el antes de firmarlo.
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return "", errores.NewStruct("Error claim", "Error al convertir a mapClaim el token", errores.ErrorInvalidToken)
  }
  
  //validamos qie no este expirado
  exp := claims["exp"].(float64)
  if time.Now().Unix() > int64(exp) {
    return "", errores.NewStruct("Expired token", "Toke expirado.", errores.ErrorInvalidToken)
  }
    
  //extraemos el nombre de usuario y lo casteamos a string
  return claims["username"].(string), nil
}