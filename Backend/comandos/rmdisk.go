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

func eliminar(comando string) {
	if err := os.Remove(comando); err != nil {
		mensajeGlobal.Mensaje = "Error. No se pudo eliminar el disco."
		fmt.Println("¡Error! Fallé al eliminar el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		return
	}
	mensajeGlobal.Mensaje = "Disco eliminado con éxito."
	fmt.Println("Disco eliminado con éxito.")
}

// Función para eliminar un disco
func EliminarDiscos(disco Rmdisk, mensaje *estructuras.Mensaje, confirmar bool) {
	mensajeGlobal = *mensaje
	fmt.Println("Confirmar rmdisk: ", confirmar)
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
	if confirmar {
		eliminar(disco.Path)
	} else {
		mensaje.Mensaje = "El disco no se ha eliminado."
	}
}
