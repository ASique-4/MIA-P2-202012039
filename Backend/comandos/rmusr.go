package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto2/estructuras"
	"strings"
	"unsafe"
)

type Rmusr struct {
	User string
}

func (rmuser *Rmusr) Rmusr(id string, lista *estructuras.ListaParticionesMontadas) {
	//Obtener la partición montada
	particionMontada := lista.ObtenerParticionMontada(id)
	if particionMontada == nil {
		fmt.Println("No se encontró la partición montada")
		return
	}

	//Abrimos el archivo
	filePart, err := os.OpenFile(particionMontada.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer filePart.Close()

	//Leemos el mbr
	mbr := estructuras.MBR{}
	binary.Read(filePart, binary.LittleEndian, &mbr)

	//Buscamos la partición
	var particion estructuras.Particion
	particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for _, particionAux := range particiones {
		if particionAux.Part_name == particionMontada.Name {
			particion = particionAux
			break
		}
	}

	//Nos posicionamos en el inicio de la partición
	filePart.Seek(int64(bytesToInt(particion.Part_start)), 0)

	//Leemos el superbloque
	superbloque := estructuras.SuperBloque{}
	binary.Read(filePart, binary.LittleEndian, &superbloque)

	// Nos posicionamos al inicio del archivo users.txt
	filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)

	// Leemos los siguientes 64 bytes
	linea := [64]byte{}
	binary.Read(filePart, binary.LittleEndian, &linea)

	// Copia de linea
	lineaCopia := [64]byte{}
	copy(lineaCopia[:], linea[:])

	// Convertimos a string
	lineaStr := string(linea[:])

	// Recorremos el archivo users.txt y buscamos el ultimo usuario
	// txt aceptado grupos -> GID, tipo, grupo
	// txt aceptado usuarios -> UID, tipo, grupo, usuario, password
	lineas := strings.Split(lineaStr, "\n")
	posTxt := 0 // Inicializar la posición del texto
	for i, linea := range lineas {
		if linea != "" {
			lineaSplit := strings.Split(linea, ",")

			if len(lineaSplit) == 1 {
				break
			}
			if lineaSplit[1] == "U" {
				if len(lineaSplit) < 5 {
					fmt.Println(lineaSplit)
					fmt.Println("Error en el archivo users.txt")
					break
				}
				// Verificamos si es el usuario a eliminar
				if lineaSplit[3] == rmuser.User {
					// Eliminamos el grupo
					lineaSplit[0] = "0"
					// Actualizamos el valor de lineaSplit en lineas
					lineas[i] = strings.Join(lineaSplit, ",")
					// Actualizamos el archivo users.txt
					// Convertimos lineas a string
					lineaStr = strings.Join(lineas, "\n")
					// Convertimos lineaStr a [64]byte
					copy(lineaCopia[:], lineaStr[:])
					// Escribimos en el archivo
					filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)
					binary.Write(filePart, binary.LittleEndian, &lineaCopia)
					return
				}
			}
			posTxt += len(lineaSplit) // Actualizar la posición del texto
		}
	}

	//Imprimimos el error
	fmt.Println("No se encontró el usuario")

}
