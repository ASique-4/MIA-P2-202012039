package estructuras

import (
	"encoding/binary"
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
			if start1 > start2 {
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
