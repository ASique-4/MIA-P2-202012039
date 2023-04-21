package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"proyecto2/estructuras"
	"time"
)

type Mount struct {
	Path          string
	Name          [16]byte
	ListaMontadas *estructuras.ListaParticionesMontadas
}

// Crear el id
func (mount *Mount) CrearId(lista *estructuras.ListaParticionesMontadas, particion Mount) string {
	numero := lista.ObtenerNumero(particion.Path)
	letra := lista.ObtenerLetra(numero)
	id := "39" + numero + letra
	return id
}

// Función para reescribir el status de una partición
func (mount *Mount) ReescribirStatus() {
	// Abrimos el archivo
	file, err := os.OpenFile(mount.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// Leemos el MBR
	mbr := estructuras.MBR{}
	binary.Read(file, binary.BigEndian, &mbr)
	// Buscamos la partición
	Particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for i := 0; i < len(Particiones); i++ {
		if Particiones[i].Part_name == mount.Name {
			Particiones[i].Part_status[0] = 1
			// Escribimos el MBR
			file.Seek(0, 0)
			binary.Write(file, binary.BigEndian, &mbr)
			break
		}
	}
}

func (mount *Mount) MountCommand(lista *estructuras.ListaParticionesMontadas) {
	mount.ListaMontadas = lista
	id := mount.CrearId(lista, *mount)
	if mount.ListaMontadas.ObtenerParticionMontada(id) != nil {
		fmt.Println("La partición ya está montada")
	} else {
		nuevaParticionMontada := &estructuras.ParticionMontada{
			Id:            id,
			Path:          mount.Path,
			NumeroDeDisco: lista.ObtenerNumero(mount.Path),
			Letra:         lista.ObtenerLetra(lista.ObtenerNumero(mount.Path)),
			Name:          mount.Name,
			Siguiente:     nil,
			Anterior:      nil,
			Mount_time:    time.Now(),
			Mount_count:   1,
			Montada:       true,
		}
		mount.ListaMontadas.AgregarParticionMontada(nuevaParticionMontada)
		mount.ReescribirStatus()
		fmt.Println("Partición montada con éxito")
	}
}
