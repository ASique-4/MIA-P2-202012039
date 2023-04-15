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
		if particion.Part_name[0] == 0 && bytesToInt(particion.Part_start) == size {
			particionLibre = particion
			break
		}
	}

	if particionLibre == nil {
		for i := 3; i >= 0; i-- {
			if bytesToInt(particiones[i].Part_size) == -1 && bytesToInt(particiones[i].Part_start) == -1 && i == 0 {
				particionLibre = particiones[i]
				break
			}
			if bytesToInt(particiones[i].Part_size) == -1 && bytesToInt(particiones[i].Part_start) == -1 &&
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

// La función crea una partición en un espacio libre determinado.
func crearParticion(mbr *estructuras.MBR, particion estructuras.Particion, espaciosLibres []EspacioLibre, best_index int, tamanio int64) {

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
	}
}

func peorAjuste(mbr *estructuras.MBR, particion estructuras.Particion) {
	// Variable para almacenar el tamaño de la partición
	var tamañoParticion int64 = int64(bytesToInt(particion.Part_size))

	// Creamos una lista de espacios libres
	var espaciosLibres []EspacioLibre = calcularEspaciosLibres(mbr)

	// Buscamos el índice del peor ajuste
	var index int = buscarIndexPeorAjuste(espaciosLibres, particion)

	// Si el índice es -1, no se encontró un espacio libre que cumpla con la condición
	if index == -1 {
		fmt.Println("No se encontró un espacio libre que cumpla con la condición.")
	} else {
		// Si el índice es válido, creamos la partición
		crearParticion(mbr, particion, espaciosLibres, index, tamañoParticion)
	}

}

func primerAjuste(mbr *estructuras.MBR, particion estructuras.Particion) {

}

func mejorAjuste(mbr *estructuras.MBR, particion estructuras.Particion) {

}

func agregarParticionAlMBR(mbr *estructuras.MBR, particion estructuras.Particion) {
	if particion.Part_fit[0] == 'B' {
		mejorAjuste(mbr, particion)
	} else if particion.Part_fit[0] == 'F' {
		primerAjuste(mbr, particion)
	} else if particion.Part_fit[0] == 'W' {
		peorAjuste(mbr, particion)
	}
}

// Función para crear una partición
func CrearParticion(particion Fdisk) {
	fmt.Println("Creando partición...")
	fmt.Println("Size:", particion.Size)
	fmt.Println("Path:", particion.Path)
	fmt.Println("Unit:", particion.Unit)
	fmt.Println("Type:", particion.Type)
	fmt.Println("Fit:", particion.Fit)
	fmt.Println("Name:", particion.Name)

	//Si el path tiene comillas, las quitamos
	if particion.Path[0] == '"' {
		particion.Path = particion.Path[1 : len(particion.Path)-1]
	}

	//Verificamos si el path existe
	if _, err := os.Stat(particion.Path); os.IsNotExist(err) {
		fmt.Println("¡Error! No existe el disco. Lo siento, parece que no soy tan hábil como pensaba.")
		return
	}

	//Size
	var tamanio int64 = int64(binary.LittleEndian.Uint32(particion.Size[:]))
	switch particion.Unit {
	case 'B':
		tamanio *= 1
	case 'K':
		tamanio *= 1024
	case 'M':
		tamanio *= 1024 * 1024
	default:
		tamanio *= 1024 * 1024
	}

	//Creamos la partición
	particionNueva := estructuras.Particion{
		Part_status: [1]byte{0},
		Part_type:   particion.Type,
		Part_fit:    particion.Fit,
		Part_start:  [4]byte{},
		Part_size:   particion.Size,
		Part_name:   particion.Name,
	}

	//Abrimos el archivo
	file, err := os.OpenFile(particion.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("¡Error! No se pudo abrir el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		return
	}

	//Leemos el MBR
	mbr := estructuras.MBR{}
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Println("¡Error! No se pudo leer el archivo. Lo siento, parece que no soy tan hábil como pensaba.")
		return
	}

	agregarParticionAlMBR(&mbr, particionNueva)

}
