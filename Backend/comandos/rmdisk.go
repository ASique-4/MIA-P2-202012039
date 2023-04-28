package comandos

import (
	"fmt"
	"os"
	"proyecto2/estructuras"
)

type Rmdisk struct {
	Path string
}

// Función para eliminar un disco
func EliminarDiscos(disco Rmdisk, mensaje *estructuras.Mensaje) {
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
	comando := disco.Path
	if err := os.Remove(comando); err != nil {
		fmt.Println("¡Error! Fallé al eliminar el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		mensaje.Mensaje = "Error. No se pudo eliminar el disco."
		return
	}
	fmt.Println("Disco eliminado con éxito.")
	mensaje.Mensaje = "Disco eliminado con éxito."
}
