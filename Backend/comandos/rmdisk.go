package comandos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"proyecto2/estructuras"
)

type Rmdisk struct {
	Path string
}

type Confirmar struct {
	Aceptar bool
}

var mensajeGlobal estructuras.Mensaje
var comandoGlobal string

// Esta función maneja una solicitud de confirmación y devuelve un valor booleano en función de si el
// usuario acepta o cancela la operación.
func handleConfirmar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	// Json a enviar
	var jsonConfirmar Confirmar

	// Decodificamos el mensaje
	err := json.NewDecoder(r.Body).Decode(&jsonConfirmar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		mensajeGlobal.Mensaje = "Error al decodificar el mensaje."
		eliminar(false, comandoGlobal)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mensajeGlobal)
		return
	}

	// Si el usuario acepta, eliminamos el disco
	if jsonConfirmar.Aceptar {
		fmt.Println("¡Operación confirmada!")
		eliminar(true, comandoGlobal)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mensajeGlobal)
		return
	} else {
		fmt.Println("¡Operación cancelada!")
		mensajeGlobal.Mensaje = "Operación cancelada."
		eliminar(false, comandoGlobal)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mensajeGlobal)
		return
	}
}

func eliminar(aceptar bool, comando string) {
	if err := os.Remove(comando); err != nil {
		fmt.Println("¡Error! Fallé al eliminar el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		mensajeGlobal.Mensaje = "Error. No se pudo eliminar el disco."
		return
	}
	fmt.Println("Disco eliminado con éxito.")
	mensajeGlobal.Mensaje = "Disco eliminado con éxito."
}

// Función para eliminar un disco
func EliminarDiscos(disco Rmdisk, mensaje *estructuras.Mensaje) {
	mensajeGlobal = *mensaje
	//Verificamos si el path tiene comillas
	if disco.Path[0] == '"' {
		disco.Path = disco.Path[1 : len(disco.Path)-1]
	}
	//Verificamos si el path existe
	if _, err := os.Stat(disco.Path); os.IsNotExist(err) {
		fmt.Println("¡Error! No existe el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		mensaje.Mensaje = "Error. No existe el disco."
		return
	}
	//Eliminamos el disco
	comandoGlobal = disco.Path
	// Enviamos un endpoint para confirmar que se eliminara el disco
	// Si el usuario acepta, se elimina el disco
	// Si el usuario cancela, se cancela la operación
	http.HandleFunc("/confirmar", handleConfirmar)
	go func() {
		if err := http.ListenAndServe(":3030", nil); err != nil {
			log.Fatal(err)
		}
	}()
}
