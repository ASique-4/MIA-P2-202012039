package comandos

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"proyecto2/estructuras"
	"strings"
	"time"
)

type Mkdisk struct {
	Size [4]byte
	Path string
	Unit byte
	Fit  [1]byte
	MBR  estructuras.MBR
}

// Función para crear un disco
func CrearDiscos(disco Mkdisk, mensaje *estructuras.Mensaje) {
	//Guardamos el nombre del disco
	nombre := disco.Path[strings.LastIndex(disco.Path, "/")+1:]
	// Verificamos si el path existe y si no lo creamos
	pathSinNombre := strings.TrimSuffix(disco.Path, "/"+disco.Path[strings.LastIndex(disco.Path, "/")+1:]) + "/"
	if err := os.MkdirAll(pathSinNombre, 0777); err != nil {
		fmt.Println("¡Error! Fallé al crear el directorio. Lo siento, parece que no soy tan hábil como pensaba.")
		mensaje.Mensaje = "Error. No se pudo crear el directorio."
		return
	}
	disco.Path = pathSinNombre + "/" + nombre

	//Creamos el archivo
	archivo, err := os.Create(disco.Path)

	if err != nil {
		fmt.Println("¡Error! Fallé al crear el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		mensaje.Mensaje = "Error. No se pudo crear el archivo."

		return
	}

	fmt.Println("Creando disco...")

	//Verificamos si el tamaño es en KB, MB o GB
	var tamanio int64 = int64(binary.LittleEndian.Uint32(disco.Size[:]))
	switch disco.Unit {
	case 'B':
		tamanio *= 1
	case 'K':
		tamanio *= 1024
	case 'M':
		tamanio *= 1024 * 1024
	default:
		tamanio *= 1024 * 1024
	}
	//Creamos el mbr
	//Guardamos el tamaño en el mbr
	binary.LittleEndian.PutUint32(disco.MBR.Mbr_tamanio[:], uint32(tamanio))
	binary.LittleEndian.PutUint64(disco.MBR.Mbr_fecha_creacion[:], uint64(time.Now().Unix()))
	//Generamos un número aleatorio para la firma
	rand.Seed(time.Now().UnixNano())
	binary.LittleEndian.PutUint32(disco.MBR.Mbr_disk_signature[:], uint32(rand.Intn(1000000)))
	//Guardamos el ajuste en el mbr
	disco.MBR.Dsk_fit = disco.Fit
	//Si el ajuste está vacío, lo llenamos con el ajuste por defecto
	if disco.MBR.Dsk_fit[0] == 0 {
		disco.MBR.Dsk_fit[0] = 'F'
	}
	fmt.Println("Fit:", disco.MBR.Dsk_fit[0])

	//Llenamos el archivo con 0
	for i := int64(0); i < tamanio; i++ {
		var c byte = 0
		if err := binary.Write(archivo, binary.LittleEndian, &c); err != nil {
			fmt.Println("¡Error! Fallé al llenar el archivo con 0. Lo siento, parece que no soy tan hábil como pensaba.")
			mensaje.Mensaje = "Error. No se pudo llenar el archivo con 0."
			archivo.Close()
			return
		}
	}

	archivo.Seek(0, 0)

	//Escribimos el mbr al inicio del archivo
	if err := binary.Write(archivo, binary.LittleEndian, &disco.MBR); err != nil {
		fmt.Println("¡Error! Fallé al escribir el MBR. Lo siento, parece que no soy tan hábil como pensaba.")
		mensaje.Mensaje = "Error. No se pudo escribir el MBR."
		archivo.Close()
		return
	}

	archivo.Close()
	fmt.Println("¡Presto! Disco creado correctamente.")
	mensaje.Mensaje = "Disco creado correctamente."

}
