package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto2/estructuras"
	"strconv"
	"strings"
	"unsafe"
)

type Mkgrp struct {
	Name string
}

func stringToInt(str string) int {
	num := 0
	for _, c := range str {
		num += int(c)
	}
	return num
}

func (mkgrp *Mkgrp) Mkgrp(id string, lista *estructuras.ListaParticionesMontadas) {
	fmt.Println(id)
	// Obtener la partición montada
	particionMontada := lista.ObtenerParticionMontada(id)
	if particionMontada == nil {
		fmt.Println("No se encontró la partición montada")
		return
	}

	// Abrimos el archivo
	filePart, err := os.OpenFile(particionMontada.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer filePart.Close()

	// Leemos el mbr
	mbr := estructuras.MBR{}
	binary.Read(filePart, binary.LittleEndian, &mbr)

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
	binary.Read(filePart, binary.LittleEndian, &superbloque)

	// Nos posicionamos al inicio del archivo users.txt
	filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)

	// Leemos los siguientes 60 bytes
	linea := [64]byte{}
	binary.Read(filePart, binary.LittleEndian, &linea)

	// Copia de linea
	lineaCopia := [64]byte{}
	copy(lineaCopia[:], linea[:])

	// Convertimos a string
	lineaStr := string(linea[:])
	fmt.Println(lineaStr)
	fmt.Println("=====================================")

	ultimoGrupo := "0"
	tamanioTxt := 0

	for _, bytes := range linea {
		if bytes != 0 {
			tamanioTxt++
		}
	}

	// Recorremos el archivo users.txt y buscamos el ultimo grupo
	// txt aceptado grupos -> GID, tipo, grupo
	// txt aceptado usuarios -> UID, tipo, grupo, usuario, password
	lineas := strings.Split(lineaStr, "\n")
	for _, linea := range lineas {
		if linea != "" {
			lineaSplit := strings.Split(linea, ",")
			fmt.Println(lineaSplit, " Len: ", len(lineaSplit))

			if len(lineaSplit) == 1 {
				break
			}
			if len(lineaSplit) < 3 {
				fmt.Println("Error en el archivo users.txt")
				break
			}
			if strings.TrimSpace(lineaSplit[1]) == "G" {
				if mkgrp.Name == strings.TrimSpace(lineaSplit[2]) {
					fmt.Println("Ya existe un grupo con ese nombre")
					return
				}
				// Verificamos si es el ultimo grupo
				if stringToInt(ultimoGrupo) < stringToInt(lineaSplit[0]) {
					ultimoGrupo = lineaSplit[0]
				}

			}
		}
	}

	// Lo agregamos al final del [64]byte
	ultimoGrupoInt, _ := strconv.Atoi(ultimoGrupo)
	ultimoGrupoInt++
	ultimoGrupo = strconv.Itoa(ultimoGrupoInt)
	copy(lineaCopia[tamanioTxt:], []byte(ultimoGrupo+",G,"+mkgrp.Name+"\n"))

	// Nos posicionamos al inicio del archivo users.txt
	filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)

	// Escribimos el archivo
	err = binary.Write(filePart, binary.LittleEndian, &lineaCopia)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Impimimos el archivo
	filePart.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)
	linea2 := [64]byte{}
	binary.Read(filePart, binary.LittleEndian, &linea2)
	fmt.Println(string(linea2[:]))

}
