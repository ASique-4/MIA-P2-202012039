package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"proyecto2/analizador"
	"proyecto2/estructuras"
	"strings"
	"time"
)

func LeerEntrada() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

type Comando struct {
	Comando string `json:"comando"`
}

var mensaje estructuras.Mensaje

func handleComando(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var c Comando

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Guardar el comando en una cadena
	comandoString := c.Comando

	// Realizar cualquier otra operación con el comando

	fmt.Println("Comando recibido:", comandoString)
	time.Sleep(5 * time.Second)
	// Responder con un mensaje de éxito
	analizador.Analizar(comandoString, &mensaje)

	fmt.Println("Mensaje enviado: ", mensaje.Accion, mensaje.Mensaje)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mensaje)

	mensaje = estructuras.Mensaje{}

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	// Guardamos el json en una estructura
	var usuario estructuras.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Realizar el login
	comando := "login >user=" + usuario.Username + " >pwd=" + usuario.Password + " >id=" + usuario.Id

	analizador.Analizar(comando, &mensaje)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mensaje)

}

func main() {
	titulo := "Proyecto 1 - MIA"
	nombre := "Angel Francisco Sique Santos"
	codigo := "202012039"
	mensaje := "Ingrese el comando a analizar"

	ancho := 40
	fmt.Printf("+" + strings.Repeat("-", ancho-2) + "+\n")
	fmt.Printf("|" + strings.Repeat(" ", ancho-2) + "|\n")
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(titulo))/2), titulo, strings.Repeat(" ", ((ancho-len(titulo))/2)-2))
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(nombre))/2), nombre, strings.Repeat(" ", ((ancho-len(nombre))/2)-2))
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(codigo))/2), codigo, strings.Repeat(" ", ((ancho-len(codigo))/2)-1))
	fmt.Printf("|" + strings.Repeat(" ", ancho-2) + "|\n")
	fmt.Printf("+" + strings.Repeat("-", ancho-2) + "+\n\n")

	fmt.Println(mensaje)

	// Leer entrada desde la API
	// Crear un endpoint que reciba el comando a analizar
	// y devuelva el resultado del analisis
	http.HandleFunc("/ejecutar-comando", handleComando)
	http.HandleFunc("/login", handleLogin)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Servidor en ejecución. Presiona Ctrl+C para detenerlo.")
	select {}
}
