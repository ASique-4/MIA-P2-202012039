package estructuras

import (
	"encoding/binary"
	"time"
)

// Estructura para el MBR
type MBR struct {
	Mbr_tamanio        [4]byte
	Mbr_fecha_creacion [16]byte
	Mbr_disk_signature [4]byte
	Dsk_fit            [1]byte
	Mbr_partition_1    Particion
	Mbr_partition_2    Particion
	Mbr_partition_3    Particion
	Mbr_partition_4    Particion
}

// La función `OrdenarParticiones()` es un método de la estructura `MBR` que ordena las particiones en
// el MBR en orden ascendente en función de sus posiciones iniciales. Para ello, crea una porción de
// estructuras `Partición` que contienen las cuatro particiones en el MBR, y luego utiliza un bucle
// anidado para comparar las posiciones iniciales de cada partición e intercambiar sus posiciones si es
// necesario. Finalmente, la función actualiza la estructura MBR con las particiones ordenadas.
func (mbr *MBR) OrdenarParticiones() {
	particiones := []Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}

	for i := 0; i < len(particiones); i++ {
		for j := 0; j < len(particiones)-1; j++ {
			start1 := binary.LittleEndian.Uint32(particiones[j].Part_start[:])
			start2 := binary.LittleEndian.Uint32(particiones[j+1].Part_start[:])
			if start1 > start2 && start2 != 0 {
				particiones[j], particiones[j+1] = particiones[j+1], particiones[j]
			}
		}
	}

	mbr.Mbr_partition_1 = particiones[0]
	mbr.Mbr_partition_2 = particiones[1]
	mbr.Mbr_partition_3 = particiones[2]
	mbr.Mbr_partition_4 = particiones[3]
}

// Estructura para las particiones
type Particion struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  [4]byte
	Part_size   [4]byte
	Part_name   [16]byte
}

// Estrucutra para el EBR
type EBR struct {
	Part_status [1]byte
	Part_fit    [1]byte
	Part_start  [4]byte
	Part_size   [4]byte
	Part_next   [4]byte
	Part_name   [16]byte
}

// Partición montada
type ParticionMontada struct {
	Id            string
	Path          string
	Letra         string
	NumeroDeDisco int
	Name          [16]byte
	Siguiente     *ParticionMontada
	Anterior      *ParticionMontada
	Mount_time    time.Time
	Unmount_time  time.Time
	Mount_count   int
	Montada       bool
}

// Lista de particiones montadas
type ListaParticionesMontadas struct {
	Primero *ParticionMontada
	Ultimo  *ParticionMontada
}

// Función para agregar una partición montada a la lista de particiones montadas
func (lista *ListaParticionesMontadas) AgregarParticionMontada(nuevaParticion *ParticionMontada) {
	if lista.Primero == nil {
		lista.Primero = nuevaParticion
		lista.Ultimo = nuevaParticion
	} else {
		lista.Ultimo.Siguiente = nuevaParticion
		nuevaParticion.Anterior = lista.Ultimo
		lista.Ultimo = nuevaParticion
	}
}

// Función para obtener una partición montada de la lista de particiones montadas
func (lista *ListaParticionesMontadas) ObtenerParticionMontada(id string) *ParticionMontada {
	aux := lista.Primero
	for aux != nil {
		if aux.Id == id {
			return aux
		}
		aux = aux.Siguiente
	}
	return nil
}

// Función para eliminar una partición montada de la lista de particiones montadas
func (lista *ListaParticionesMontadas) EliminarParticionMontada(id string) {
	aux := lista.Primero
	for aux != nil {
		if aux.Id == id {
			if aux.Anterior == nil {
				lista.Primero = aux.Siguiente
				if aux.Siguiente != nil {
					aux.Siguiente.Anterior = nil
				}
			} else if aux.Siguiente == nil {
				lista.Ultimo = aux.Anterior
				aux.Anterior.Siguiente = nil
			} else {
				aux.Anterior.Siguiente = aux.Siguiente
				aux.Siguiente.Anterior = aux.Anterior
			}
			return
		}
		aux = aux.Siguiente
	}
}

// Función para imprimir la lista de particiones montadas
func (lista *ListaParticionesMontadas) ImprimirListaParticionesMontadas() {
	aux := lista.Primero
	for aux != nil {
		println(aux.Id)
		aux = aux.Siguiente
	}
}

// Obtener el número de particiones montadas por número de disco
func (lista *ListaParticionesMontadas) ObtenerUltimaParticionMontadaPorNumeroDeDisco(numeroDeDisco int) *ParticionMontada {
	aux := lista.Primero
	var ultimaParticion *ParticionMontada
	for aux != nil {
		if aux.NumeroDeDisco == numeroDeDisco {
			ultimaParticion = aux
		}
		aux = aux.Siguiente
	}
	return ultimaParticion
}

// Buscar si se repite el path si se repite retorna el número de disco
// si no se repite retorna 1 si no hay particiones montadas
// si ya hay particiones montadas retorna el número de disco más alto + 1
func (lista *ListaParticionesMontadas) ObtenerNumero(path string) int {
	aux := lista.Primero
	var numeroDeDisco int
	for aux != nil {
		if aux.Path == path {
			return aux.NumeroDeDisco
		}
		if aux.NumeroDeDisco > numeroDeDisco {
			numeroDeDisco = aux.NumeroDeDisco
		}
		aux = aux.Siguiente
	}
	if numeroDeDisco == 0 {
		return 1
	}
	return numeroDeDisco + 1
}

// Función para obtener la ultima letra de la partición montada con el mismo número de disco
// si no hay particiones montadas con el mismo número de disco retorna la letra A
// si ya hay particiones montadas con el mismo número de disco retorna la siguiente letra en el abecedario
func (lista *ListaParticionesMontadas) ObtenerLetra(numeroDeDisco int) string {
	aux := lista.Primero
	var ultimaLetra string
	for aux != nil {
		if aux.NumeroDeDisco == numeroDeDisco {
			ultimaLetra = aux.Letra
		}
		aux = aux.Siguiente
	}
	if ultimaLetra == "" {
		return "A"
	}
	return string(ultimaLetra[0] + 1)
}
