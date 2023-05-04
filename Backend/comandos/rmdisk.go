package comandos

import (
	"fmt"
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
		disco.Path = disco.Path[1 : len(disco.Path)-2]
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
	eliminar(true, disco.Path)
}
