package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto2/estructuras"
	"strings"
	"unsafe"
)

type Login struct {
	Usuario string
	Pass    string
	Id      string
}

func (login *Login) Login(lista *estructuras.ListaParticionesMontadas, w *estructuras.Mensaje) *estructuras.Usuario {
	// Obtener la partición montada
	particionMontada := lista.ObtenerParticionMontada(login.Id)
	if particionMontada == nil {
		fmt.Println("No se encontró la partición montada")
		w.Mensaje = "No se encontró la partición montada"
		return nil
	}

	// Abrimos el archivo
	filePart, err := os.Open(particionMontada.Path)
	if err != nil {
		fmt.Println(err)
		w.Mensaje = "Error al abrir el archivo"
		return nil
	}
	defer filePart.Close()

	// Leemos el mbr
	mbr := estructuras.MBR{}
	binary.Read(filePart, binary.BigEndian, &mbr)

	// Buscamos la partición
	var particion estructuras.Particion
	particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for _, particionAux := range particiones {
		if particionAux.Part_name == particionMontada.Name {
			particion = particionAux
			break
		}
	}

	// Nos posicionamos en el inicio de la partición
	filePart.Seek(int64(bytesToInt(particion.Part_start)), 0)

	// Leemos el superbloque
	superbloque := estructuras.SuperBloque{}
	binary.Read(filePart, binary.BigEndian, &superbloque)

	// Nos posicionamos al inicio del archivo users.txt
	filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)

	// Leemos los siguientes 60 bytes
	linea := [64]byte{}
	binary.Read(filePart, binary.BigEndian, &linea)

	// Convertimos la linea a string
	lineaStr := string(linea[:])

	// Recorremos el archivo users.txt
	// txt aceptado grupos -> GID, tipo, grupo
	// txt aceptado usuarios -> UID, tipo, grupo, usuario, password
	for linea[0] != 0 {
		// Separamos la linea por comas
		lineaSplit := strings.Split(lineaStr, ",")

		// Verificamos si es un usuario
		if lineaSplit[1] == "U" {
			// Verificamos si el usuario es el que buscamos
			if lineaSplit[3] == login.Usuario {
				// Verificamos si la contraseña es la que buscamos
				if lineaSplit[4] == login.Pass {
					// Creamos el usuario
					usuario := estructuras.Usuario{
						UID:      lineaSplit[0],
						Type:     lineaSplit[1],
						Group:    lineaSplit[2],
						Username: lineaSplit[3],
						Password: lineaSplit[4],
						PartID:   login.Id,
					}
					fmt.Println("Login correcto")
					w.Mensaje = "Login correcto"
					return &usuario
				}
			}
		}

		// Sacamos la linea del string
		linea = [64]byte{}
		copy(linea[:], lineaStr)

		// Leemos la siguiente línea
		lineaStr = strings.Split(lineaStr, "\n")[1]

	}

	fmt.Println("Login incorrecto")
	w.Mensaje = "Login incorrecto"
	return nil

}
