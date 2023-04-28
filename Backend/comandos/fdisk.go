package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"proyecto2/estructuras"
)

// The Fdisk type represents a disk partition with its path, size, unit, type, fit, and name.
type Fdisk struct {
	Path string
	Size [4]byte
	Unit byte
	Type [1]byte
	Fit  [1]byte
	Name [16]byte
}

type EspacioLibre struct {
	inicio   int64
	tamaño   int64
	partcion *estructuras.Particion
}

var mensajeTmp *estructuras.Mensaje

func calcularEspaciosLibres(mbr *estructuras.MBR) []EspacioLibre {
	// Creamos un slice para almacenar los espacios libres
	var espaciosLibres []EspacioLibre

	// Creamos un slice de particiones
	particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}

	// Ordenamos las particiones por posición de inicio
	mbr.OrdenarParticiones()

	// Iteramos sobre las particiones para encontrar los espacios libres
	lastEnd := binary.Size(mbr) // Variable para almacenar el final de la última partición
	anyPartitionUsed := false   // Variable para verificar si se ha ocupado alguna partición

	for _, particion := range particiones {
		// Si la partición no tiene tamaño ni start, es porque no se ha creado aún
		if particion.Part_size[0] == 0 && particion.Part_start[0] == 0 && anyPartitionUsed {
			// Agregamos un espacio libre que abarque todo el disco
			espaciosLibres = append(espaciosLibres, EspacioLibre{int64(lastEnd), int64(bytesToInt(mbr.Mbr_tamanio)) - int64(lastEnd), nil})
			break
		}

		// Si la partición tiene un nombre, no es un espacio libre
		if particion.Part_name[0] != 0 {
			// Si la partición comienza más allá del final de la última partición, hay un espacio libre entre ellas
			if bytesToInt(particion.Part_start) > lastEnd {
				espaciosLibres = append(espaciosLibres, EspacioLibre{int64(lastEnd), int64(bytesToInt(particion.Part_start)) - int64(lastEnd), nil})
			}
			// Actualizamos el final de la última partición
			lastEnd = bytesToInt(particion.Part_start) + bytesToInt(particion.Part_size)
			anyPartitionUsed = true
		} else if particion.Part_size[0] != 0 && particion.Part_start[0] != 0 {
			// Si la partición no tiene nombre y tamaño y start no son 0, es un espacio libre
			espaciosLibres = append(espaciosLibres, EspacioLibre{int64(bytesToInt(particion.Part_start)), int64(bytesToInt(particion.Part_size)), nil})
			// Actualizamos el final de la última partición
			lastEnd = bytesToInt(particion.Part_start) + bytesToInt(particion.Part_size)
			anyPartitionUsed = true
		}
	}

	if len(espaciosLibres) == 0 {
		// Agregamos un espacio libre que abarque todo el disco
		espaciosLibres = append(espaciosLibres, EspacioLibre{int64(lastEnd), int64(bytesToInt(mbr.Mbr_tamanio)) - int64(lastEnd), nil})
	}

	// Retornamos el slice de espacios libres
	return espaciosLibres
}

// La función convierte una matriz de bytes de longitud 4 en un número entero utilizando el orden de
// bytes little-endian.
func bytesToInt(bytess [4]byte) int {
	var n uint32
	buf := bytes.NewReader(bytess[:])
	binary.Read(buf, binary.LittleEndian, &n)
	return int(n)
}

// La función busca una partición libre en un MBR y la devuelve si la encuentra.
func buscarParticionLibre(mbr *estructuras.MBR, size int) *estructuras.Particion {
	var particionLibre *estructuras.Particion = nil
	particiones := [4]*estructuras.Particion{&mbr.Mbr_partition_1, &mbr.Mbr_partition_2, &mbr.Mbr_partition_3, &mbr.Mbr_partition_4}

	for _, particion := range particiones {

		if particion.Part_name == [16]byte{} && bytesToInt(particion.Part_start) == size {
			particionLibre = particion
			break
		}
	}

	if particionLibre == nil {
		for i := 3; i >= 0; i-- {
			if particiones[i].Part_size == [4]byte{} && (particiones[i].Part_start) == [4]byte{} && i == 0 {
				particionLibre = particiones[i]
				break
			}
			if particiones[i].Part_size == [4]byte{} && (particiones[i].Part_start) == [4]byte{} &&
				(bytesToInt(particiones[i-1].Part_size)+bytesToInt(particiones[i-1].Part_start)) == (bytesToInt(mbr.Mbr_tamanio)-size) {
				particionLibre = particiones[i]
				break
			}
		}
	}

	return particionLibre
}

// La función busca el índice del espacio libre que mejor se ajusta a una partición determinada en una
// lista de espacios libres.
func buscarIndexMejorAjuste(espaciosLibres []EspacioLibre, particion estructuras.Particion) int {
	// Variable para almacenar el índice del mejor ajuste
	var index int = -1

	// Variable para almacenar la diferencia de tamaño entre el espacio libre y la partición
	var diferencia int64 = 0

	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Iteramos sobre los espacios libres
	for i, espacioLibre := range espaciosLibres {
		// Si el espacio libre es mayor o igual que el tamaño de la partición
		if espacioLibre.tamaño >= tamañoParticion {
			// Si el espacio libre es el primero que cumple con la condición, lo guardamos
			if index == -1 {
				index = i
				diferencia = espacioLibre.tamaño - tamañoParticion
			} else {
				// Si el espacio libre no es el primero que cumple con la condición, comparamos el tamaño de la diferencia
				if espacioLibre.tamaño-tamañoParticion < diferencia {
					index = i
					diferencia = espacioLibre.tamaño - tamañoParticion
				}
			}
		}
	}

	// Retornamos el índice del mejor ajuste
	return index

}

// La función busca el índice del primer espacio libre disponible en una lista de espacios libres lo
// suficientemente grande como para caber en una partición determinada.
func buscarIndexPrimerAjuste(espaciosLibres []EspacioLibre, particion estructuras.Particion) int {
	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Iteramos sobre los espacios libres
	for i, espacioLibre := range espaciosLibres {
		// Si el espacio libre es mayor o igual que el tamaño de la partición, retornamos el índice
		if espacioLibre.tamaño >= tamañoParticion {
			return i
		}
	}

	// Si no se encuentra un espacio libre que cumpla con la condición, retornamos -1
	return -1
}

// La función busca el índice del espacio libre peor ajustado para una partición determinada en una
// lista de espacios libres.
func buscarIndexPeorAjuste(espaciosLibres []EspacioLibre, particion estructuras.Particion) int {
	// Variable para almacenar el índice del peor ajuste
	var index int = -1

	// Variable para almacenar la diferencia de tamaño entre el espacio libre y la partición
	var diferencia int64 = 0

	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Iteramos sobre los espacios libres
	for i, espacioLibre := range espaciosLibres {
		// Si el espacio libre es mayor o igual que el tamaño de la partición
		if espacioLibre.tamaño >= tamañoParticion {
			// Si el espacio libre es el primero que cumple con la condición, lo guardamos
			if index == -1 {
				index = i
				diferencia = espacioLibre.tamaño - tamañoParticion
			} else {
				// Si el espacio libre no es el primero que cumple con la condición, comparamos el tamaño de la diferencia
				if espacioLibre.tamaño-tamañoParticion > diferencia {
					index = i
					diferencia = espacioLibre.tamaño - tamañoParticion
				}
			}
		}
	}

	// Retornamos el índice del peor ajuste
	return index

}

// la función devuelve la partición extendida de un disco
func getParticionExtendida(mbr *estructuras.MBR) *estructuras.Particion {

	//Particiones del disco
	particiones := [4]estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}

	for _, particion := range particiones {
		if particion.Part_type[0] == 'E' || particion.Part_type[0] == 'e' {
			return &particion
		}
	}
	return nil
}

// La función escribe el EBR en el disco
func escribirEBR(ebr estructuras.EBR, particionExtendida estructuras.Particion, path string) {
	// Abrimos el archivo
	if archivo, err := os.OpenFile(path, os.O_RDWR, 0666); err == nil {
		// Nos movemos a la posición del EBR
		archivo.Seek(int64(bytesToInt(ebr.Part_start)), 0)

		// Escribimos el EBR
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, ebr)
		archivo.Write(binario.Bytes())

		// Cerramos el archivo
		archivo.Close()
	}

}

// La función crea una partición en un espacio libre determinado.
func crearParticion(mbr *estructuras.MBR, particion estructuras.Particion, espaciosLibres []EspacioLibre, best_index int, tamanio int64, path string) {

	if particion.Part_type[0] == 'P' || particion.Part_type[0] == 'p' {
		if espaciosLibres[best_index].partcion != nil {
			espaciosLibres[best_index].partcion.Part_status = particion.Part_status
			espaciosLibres[best_index].partcion.Part_type = particion.Part_type
			espaciosLibres[best_index].partcion.Part_fit = particion.Part_fit
			espaciosLibres[best_index].partcion.Part_size = particion.Part_size
			espaciosLibres[best_index].partcion.Part_name = particion.Part_name
			return
		} else {
			// Buscamos la partición libre
			particionLibre := buscarParticionLibre(mbr, int(espaciosLibres[best_index].tamaño))

			// Si la partición no es valida, retornamos
			if particionLibre == nil {
				fmt.Println("No se encontró una partición libre.")
				mensajeTmp.Mensaje = "No se encontró una partición libre."
				return
			}

			// Creamos la partición
			particionLibre.Part_status = particion.Part_status
			particionLibre.Part_type = particion.Part_type
			particionLibre.Part_fit = particion.Part_fit
			particionLibre.Part_size = particion.Part_size
			particionLibre.Part_name = particion.Part_name
			binary.LittleEndian.PutUint32(particionLibre.Part_start[:], uint32(espaciosLibres[best_index].inicio))
		}
	} else if particion.Part_type[0] == 'E' || particion.Part_type[0] == 'e' {

		// Verificamos si existe una partición extendida
		if getParticionExtendida(mbr) != nil {
			fmt.Println("Ya existe una partición extendida.")
			mensajeTmp.Mensaje = "Ya existe una partición extendida."
			return
		}

		// Buscamos la partición libre
		particionLibre := buscarParticionLibre(mbr, int(espaciosLibres[best_index].tamaño))

		// Si la partición no es valida, retornamos
		if particionLibre == nil {
			fmt.Println("No se encontró una partición libre.")
			mensajeTmp.Mensaje = "No se encontró una partición libre."
			return
		}

		// Creamos la partición
		particionLibre.Part_status = particion.Part_status
		particionLibre.Part_type = particion.Part_type
		particionLibre.Part_fit = particion.Part_fit
		particionLibre.Part_size = particion.Part_size
		particionLibre.Part_name = particion.Part_name
		binary.LittleEndian.PutUint32(particionLibre.Part_start[:], uint32(espaciosLibres[best_index].inicio))

	} else if particion.Part_type[0] == 'L' || particion.Part_type[0] == 'l' {
		// Buscamos la partición libre
		particionExtendida := getParticionExtendida(mbr)

		// Si la partición no es valida, retornamos
		if particionExtendida == nil {
			fmt.Println("No se encontró una partición extendida.")
			mensajeTmp.Mensaje = "No se encontró una partición extendida."
			return
		}

		// Abrimos el archivo
		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Error al abrir el archivo.")
			mensajeTmp.Mensaje = "Error al abrir el archivo."
			return
		}

		defer file.Close()

		// Leemos los EBR
		inicio := bytesToInt(particionExtendida.Part_start)
		tamaño := bytesToInt(particionExtendida.Part_size)
		espacioOcupado := 0

		// Variable para almacenar el EBR
		var ebr estructuras.EBR

		for espacioOcupado < tamaño {
			// Nos movemos a la posición del EBR
			file.Seek(int64(inicio), 0)

			// Leemos el EBR
			binary.Read(file, binary.BigEndian, &ebr)

			// Verificamos si el EBR es válido
			if ebr.Part_size == [4]byte{0, 0, 0, 0} || bytesToInt(ebr.Part_size) > tamaño || bytesToInt(ebr.Part_size) < 0 || (ebr.Part_status != [1]byte{'0'} && ebr.Part_status != [1]byte{'1'}) || (ebr.Part_fit != [1]byte{'B'} && ebr.Part_fit != [1]byte{'F'} && ebr.Part_fit != [1]byte{'W'}) || bytesToInt(ebr.Part_start) != inicio {

				break
			}

			// Verificamos si hay espacio para la partición
			inicio += bytesToInt(ebr.Part_size)
			espacioOcupado += bytesToInt(ebr.Part_size)
		}

		// Verificamos si hay espacio para la partición
		if espacioOcupado+int(tamanio) > tamaño {
			fmt.Println("No hay espacio para la partición.")
			mensajeTmp.Mensaje = "No hay espacio para la partición."
			return
		}

		// Actualizamos el next del EBR anterior
		if inicio != bytesToInt(particionExtendida.Part_start) {
			// Nos movemos a la posición del EBR
			file.Seek(int64(inicio)-tamanio, 0)
			binary.Read(file, binary.BigEndian, &ebr)
			binary.LittleEndian.PutUint32(ebr.Part_next[:], uint32(inicio))
			file.Seek(int64(inicio)-tamanio, 0)
			binary.Write(file, binary.BigEndian, &ebr)
		}

		// Creamos el EBR
		nuevoEBR := estructuras.EBR{
			Part_status: [1]byte{'0'},
			Part_fit:    [1]byte{'W'},
			Part_size:   particion.Part_size,
			Part_next:   [4]byte{0, 0, 0, 0},
			Part_name:   particion.Part_name,
		}
		binary.LittleEndian.PutUint32(nuevoEBR.Part_start[:], uint32(inicio))

		// Escribimos el EBR
		file.Seek(int64(inicio), 0)
		binary.Write(file, binary.BigEndian, &nuevoEBR)

	}
}

func peorAjuste(mbr *estructuras.MBR, particion estructuras.Particion, path string) {
	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Creamos una lista de espacios libres
	var espaciosLibres []EspacioLibre = calcularEspaciosLibres(mbr)

	// Buscamos el índice del peor ajuste
	var index int = buscarIndexPeorAjuste(espaciosLibres, particion)

	// Si el índice es -1, no se encontró un espacio libre que cumpla con la condición
	if index == -1 {
		fmt.Println("No se encontró un espacio libre que cumpla con la condición.")
		mensajeTmp.Mensaje = "No se encontró un espacio libre que cumpla con la condición."
	} else {
		// Si el índice es válido, creamos la partición
		crearParticion(mbr, particion, espaciosLibres, index, tamañoParticion, path)
	}

}

func primerAjuste(mbr *estructuras.MBR, particion estructuras.Particion, path string) {
	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Creamos una lista de espacios libres
	var espaciosLibres []EspacioLibre = calcularEspaciosLibres(mbr)

	// Buscamos el índice del primer ajuste
	var index int = buscarIndexPrimerAjuste(espaciosLibres, particion)

	// Si el índice es -1, no se encontró un espacio libre que cumpla con la condición
	if index == -1 {
		fmt.Println("No se encontró un espacio libre que cumpla con la condición.")
		mensajeTmp.Mensaje = "No se encontró un espacio libre que cumpla con la condición."
	} else {
		// Si el índice es válido, creamos la partición
		crearParticion(mbr, particion, espaciosLibres, index, tamañoParticion, path)
	}

}

func mejorAjuste(mbr *estructuras.MBR, particion estructuras.Particion, path string) {
	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Creamos una lista de espacios libres
	var espaciosLibres []EspacioLibre = calcularEspaciosLibres(mbr)

	// Buscamos el índice del mejor ajuste
	var index int = buscarIndexMejorAjuste(espaciosLibres, particion)

	// Si el índice es -1, no se encontró un espacio libre que cumpla con la condición
	if index == -1 {
		fmt.Println("No se encontró un espacio libre que cumpla con la condición.")
		mensajeTmp.Mensaje = "No se encontró un espacio libre que cumpla con la condición."
	} else {
		// Si el índice es válido, creamos la partición
		crearParticion(mbr, particion, espaciosLibres, index, tamañoParticion, path)
	}

}

func existeNombreParticionLogica(particion estructuras.Particion, nombre [16]byte, path string) bool {
	// Abrimos el archivo
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Recorremos las particiones
	var ebr estructuras.EBR
	var inicio int = bytesToInt(particion.Part_start)
	var tamaño int = bytesToInt(particion.Part_size)
	var espacioOcupado int = 0

	for espacioOcupado < tamaño {
		// Nos movemos a la posición del EBR
		file.Seek(int64(inicio), 0)
		binary.Read(file, binary.BigEndian, &ebr)

		// Verificamos si el nombre es igual
		if ebr.Part_name == nombre {
			return true
		}

		// Verificamos si el next es -1
		if ebr.Part_size == [4]byte{0, 0, 0, 0} || bytesToInt(ebr.Part_size) > tamaño || bytesToInt(ebr.Part_size) < 0 || (ebr.Part_status != [1]byte{'0'} && ebr.Part_status != [1]byte{'1'}) || (ebr.Part_fit != [1]byte{'B'} && ebr.Part_fit != [1]byte{'F'} && ebr.Part_fit != [1]byte{'W'}) || ebr.Part_name == [16]byte{0} || bytesToInt(ebr.Part_start) != inicio {
			break
		}

		// Actualizamos el inicio
		inicio = bytesToInt(ebr.Part_next)
		espacioOcupado += bytesToInt(ebr.Part_size)
	}

	return false

}

func existeNombreParticion(mbr estructuras.MBR, nombre [16]byte, path string) bool {
	particiones := [4]estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for i := 0; i < len(particiones); i++ {
		if particiones[i].Part_name == nombre {
			return true
		}
		// Verificamos si la partición es extendida
		if particiones[i].Part_type[0] == 'E' || particiones[i].Part_type[0] == 'e' {
			// Buscamos si existe una partición lógica con el mismo nombre
			if existeNombreParticionLogica(particiones[i], nombre, path) {
				return true
			}
		}
	}
	return false
}

func agregarParticionAlMBR(mbr *estructuras.MBR, particion estructuras.Particion, path string) {
	// Verificamos si existe una partición con el mismo nombre
	if existeNombreParticion(*mbr, particion.Part_name, path) {
		fmt.Println("Ya existe una partición con el mismo nombre.")
		mensajeTmp.Mensaje = "Ya existe una partición con el mismo nombre."
		return
	}
	if particion.Part_fit[0] == 'B' {
		mejorAjuste(mbr, particion, path)
	} else if particion.Part_fit[0] == 'F' {
		primerAjuste(mbr, particion, path)
	} else if particion.Part_fit[0] == 'W' {
		peorAjuste(mbr, particion, path)
	}
}

func ImprimirMBR(mbr estructuras.MBR, path string) {
	fmt.Println("MBR")
	fmt.Printf("MBR_Tamaño: %d\n", bytesToInt(mbr.Mbr_tamanio))
	fmt.Printf("MBR_Fecha_creacion: %s\n", mbr.Mbr_fecha_creacion)
	fmt.Printf("MBR_Disk_signature: %d\n", mbr.Mbr_disk_signature)
	fmt.Println("Particiones:")

	particiones := [4]estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	fmt.Println("---------------------")
	for i := 0; i < len(particiones); i++ {
		if particiones[i].Part_size != [4]byte{0, 0, 0, 0} {
			fmt.Printf("Particion %d%s", i+1, " - ")
			fmt.Printf("Particion_Status: %s\n", particiones[i].Part_status)
			fmt.Printf("Particion_Type: %s\n", particiones[i].Part_type)
			fmt.Printf("Particion_Fit: %s\n", particiones[i].Part_fit)
			fmt.Printf("Particion_Start: %d\n", bytesToInt(particiones[i].Part_start))
			fmt.Printf("Particion_Size: %d\n", bytesToInt(particiones[i].Part_size))
			fmt.Printf("Particion_Name: %s\n", particiones[i].Part_name)
			fmt.Println("---------------------")
			//Si la partición es extendida, imprimimos sus EBR
			if particiones[i].Part_type[0] == 'E' {
				//Abrimos el archivo
				file, err := os.OpenFile(path, os.O_RDWR, 0666)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()

				//Nos movemos a la posición del EBR
				file.Seek(int64(bytesToInt(particiones[i].Part_start)), 0)

				// Leemos los EBR
				inicio := bytesToInt(particiones[i].Part_start)
				tamaño := bytesToInt(particiones[i].Part_size)
				espacioOcupado := 0

				// Variable para almacenar el EBR
				var ebr estructuras.EBR

				for espacioOcupado < tamaño {
					// Nos movemos a la posición del EBR
					file.Seek(int64(inicio), 0)

					// Leemos el EBR
					binary.Read(file, binary.BigEndian, &ebr)
					fmt.Println("inicio: ", inicio)
					fmt.Println("ebr.part_start: ", bytesToInt(ebr.Part_start))

					// Verificamos si el EBR no es válido
					if ebr.Part_size == [4]byte{0, 0, 0, 0} || bytesToInt(ebr.Part_size) > tamaño || bytesToInt(ebr.Part_size) < 0 || (ebr.Part_status != [1]byte{'0'} && ebr.Part_status != [1]byte{'1'}) || (ebr.Part_fit != [1]byte{'B'} && ebr.Part_fit != [1]byte{'F'} && ebr.Part_fit != [1]byte{'W'}) || ebr.Part_name == [16]byte{0} || bytesToInt(ebr.Part_start) != inicio {
						break
					}

					// Imprimimos el EBR
					fmt.Println("EBR")
					fmt.Printf("EBR_Status: %s\n", ebr.Part_status)
					fmt.Printf("EBR_Fit: %s\n", ebr.Part_fit)
					fmt.Printf("EBR_Start: %d\n", bytesToInt(ebr.Part_start))
					fmt.Printf("EBR_Size: %d\n", bytesToInt(ebr.Part_size))
					fmt.Printf("EBR_Name: %s\n", ebr.Part_name)
					fmt.Printf("EBR_Next: %d\n", bytesToInt(ebr.Part_next))
					fmt.Println("---------------------")

					// Obtenemos el siguiente EBR
					inicio = bytesToInt(ebr.Part_next)
					espacioOcupado += bytesToInt(ebr.Part_size)
				}

			}
		}

	}

}

// int64 a [4]byte
func intToBytes(n int) [4]byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(n))
	return b
}

// Función para crear una partición
func CrearParticion(particion Fdisk, mensaje *estructuras.Mensaje) {
	mensajeTmp = mensaje

	//Si el path tiene comillas, las quitamos
	if particion.Path[0] == '"' {
		particion.Path = particion.Path[1 : len(particion.Path)-1]
	}

	//Verificamos si el path existe
	if _, err := os.Stat(particion.Path); os.IsNotExist(err) {
		fmt.Println("¡Error! No existe el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		mensajeTmp.Mensaje = "¡Error! No existe el disco. Lo siento, parece que no soy tan hábil como pensaba."
		return
	}

	//Size
	var tamanio int = bytesToInt(particion.Size)
	switch particion.Unit {
	case 'B':
		tamanio *= 1
	case 'K':
		tamanio *= 1024
	case 'M':
		tamanio *= 1024 * 1024
	default:
		tamanio *= 1024
	}

	//Creamos la partición
	particionNueva := estructuras.Particion{
		Part_status: [1]byte{'0'},
		Part_type:   particion.Type,
		Part_fit:    particion.Fit,
		Part_start:  [4]byte{},
		Part_name:   particion.Name,
	}

	// Guardamos el tamaño de la partición
	binary.LittleEndian.PutUint32(particionNueva.Part_size[:], uint32(tamanio))

	//Abrimos el archivo
	file, err := os.OpenFile(particion.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("¡Error! No se pudo abrir el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		mensajeTmp.Mensaje = "¡Error! No se pudo abrir el archivo. Lo siento, parece que no soy tan hábil como pensaba."
		return
	}

	//Leemos el MBR
	mbr := estructuras.MBR{}
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("¡Error! No se pudo leer el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		mensajeTmp.Mensaje = "¡Error! No se pudo leer el archivo. Lo siento, parece que no soy tan hábil como pensaba."
		return
	}

	agregarParticionAlMBR(&mbr, particionNueva, particion.Path)

	//Ordenamos las particiones
	mbr.OrdenarParticiones()

	//Escribimos el MBR
	file.Seek(0, 0)
	err = binary.Write(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("¡Error! No se pudo escribir el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		mensajeTmp.Mensaje = "¡Error! No se pudo escribir el archivo. Lo siento, parece que no soy tan hábil como pensaba."
		return
	}

	//Cerramos el archivo
	file.Close()

	fmt.Println("¡Partición creada con éxito!")
	mensajeTmp.Mensaje = "¡Partición creada con éxito!"

}
