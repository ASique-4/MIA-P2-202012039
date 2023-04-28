package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"proyecto2/estructuras"
	"strconv"
	"strings"
	"time"
)

type Rep struct {
	Path string
	Name string
	Id   string
	Ruta string
}

func calcularPorcentajeOcupado(tamanioParticion int, mbr estructuras.MBR) float64 {
	// Calculamos el tamaño del disco
	tamanioDisco := bytesToInt(mbr.Mbr_tamanio)

	// Calculamos el porcentaje
	porcentaje := float64(tamanioParticion) / float64(tamanioDisco) * 100

	return porcentaje
}

func calcularEspacioLibreMBR(mbr estructuras.MBR) float64 {
	// Calculamos el tamaño del disco
	tamanioDisco := bytesToInt(mbr.Mbr_tamanio)

	// Calculamos el tamaño de las particiones
	tamanioParticiones := 0
	particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for _, particion := range particiones {
		tamanioParticiones += bytesToInt(particion.Part_size)
	}

	// Calculamos el espacio libre
	espacioLibre := float64(tamanioDisco-tamanioParticiones) / float64(tamanioDisco) * 100

	return espacioLibre
}

func ReporteDisk(rep *Rep, lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {

	// Obtener la partición montada
	particionMontada := lista.ObtenerParticionMontada(rep.Id)
	if particionMontada == nil {
		fmt.Println("No se encontró la partición montada")
		mensaje.Mensaje = "No se encontró la partición montada"
		return
	}

	// Abrimos el archivo
	filePart, err := os.Open(particionMontada.Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		mensaje.Mensaje = "Error al abrir el archivo"
		return
	}
	defer filePart.Close()

	// Leemos el mbr
	mbr := estructuras.MBR{}
	binary.Read(filePart, binary.BigEndian, &mbr)

	// Quitamos el nombre del archivo del path
	directorio := path.Dir(rep.Path)

	// Guardamos el nombre del archivo
	nombreDelArchivo := strings.Split(path.Base(particionMontada.Path), ".")[0]

	// Creamos el directorio si no existe
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		os.MkdirAll(directorio, 0777)
	}

	// path con extensión .dot
	dot := strings.Split(rep.Path, ".")[0] + ".dot"
	// Creamos el archivo
	fileDot, err := os.Create(dot)
	if err != nil {
		fmt.Println("Error al crear el archivo")
		mensaje.Mensaje = "Error al crear el archivo"
		return
	}
	defer fileDot.Close()

	// String para las EBR
	ebrString := ""
	rowspan := 0

	// Escribimos el archivo
	// Escribimos el disco
	fileDot.WriteString("digraph G {\n")
	fileDot.WriteString("labelloc=\"t\";\n")
	fileDot.WriteString("label=\"" + nombreDelArchivo + "\";\n")
	fileDot.WriteString("parent [\n")
	fileDot.WriteString("shape=plaintext\n")
	fileDot.WriteString("label=<\n")
	fileDot.WriteString("<table border=\"1\" cellborder=\"1\">\n")
	fileDot.WriteString("<tr><td rowspan=\"3\" bgcolor='#0E8388'>MBR</td>\n")

	// Recorremos las particiones
	particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for _, particion := range particiones {
		// Si es primaria
		if particion.Part_type[0] == 'p' || particion.Part_type[0] == 'P' {
			// Escribimos la partición
			fileDot.WriteString("<td rowspan=\"3\" bgcolor='#0E8388'>PRIMARIA<br /> " + floatToString(calcularPorcentajeOcupado(bytesToInt(particion.Part_size), mbr)) + "%</td>\n")
		} else if particion.Part_type[0] == 'e' || particion.Part_type[0] == 'E' {
			// Revisamos si tiene particiones lógicas
			filePart.Seek(int64(bytesToInt(particion.Part_start)), 0)
			// Variable para leer el ebr
			var ebr estructuras.EBR

			// leer el ebr
			binary.Read(filePart, binary.BigEndian, &ebr)

			inicio := bytesToInt(particion.Part_start)
			tamanio := bytesToInt(particion.Part_size)
			espacioOcupado := 0

			for espacioOcupado < tamanio {
				// Leemos el ebr
				filePart.Seek(int64(inicio), 0)
				binary.Read(filePart, binary.BigEndian, &ebr)

				// Si el status es 0, ya no hay más particiones
				if (bytesToInt(ebr.Part_start)) != inicio {
					break
				}

				// Guardamos el ebr
				ebrString += "<td rowspan=\"2\" bgcolor='#0E8388'>EBR</td>\n"
				ebrString += "<td rowspan=\"2\" bgcolor='#0E8388'>LOGICA<br /> " + floatToString(calcularPorcentajeOcupado(bytesToInt(ebr.Part_size), mbr)) + "%</td>\n"

				// Sumamos 2 al rowspan
				rowspan += 2

				// Calculamos el espacio ocupado
				espacioOcupado += (bytesToInt(ebr.Part_size))

				// Calculamos el inicio
				inicio = (bytesToInt(ebr.Part_next))
			}
		}
	}

	// Espacio libre
	fileDot.WriteString("<td rowspan=\"3\" bgcolor='#0E8388'>LIBRE<br /> " + floatToString(calcularEspacioLibreMBR(mbr)) + "%</td>\n")

	// Si hubieron particiones lógicas
	if rowspan != 0 && ebrString != "" {
		fileDot.WriteString("<td rowspan=\"1\" colspan=\"")
		fileDot.WriteString(intToString(rowspan))
		fileDot.WriteString("\" bgcolor='#0E8388'>EXTENDIDA</td>\n")
		fileDot.WriteString("</tr>\n")
		fileDot.WriteString("<tr>\n")
		fileDot.WriteString(ebrString)
	}

	// Cerramos el archivo
	fileDot.WriteString("</tr>\n")
	fileDot.WriteString("</table>>\n")
	fileDot.WriteString("];\n")
	fileDot.WriteString("}")

	// Generamos el reporte
	// Obtenemos la extensión
	extension := strings.Split(rep.Path, ".")[1]
	if extension == "png" {
		// Generamos el png
		cmd := exec.Command("dot", "-Tpng", dot, "-o", rep.Path)
		cmd.Run()
	} else if extension == "pdf" {
		// Generamos el pdf
		cmd := exec.Command("dot", "-Tpdf", dot, "-o", rep.Path)
		cmd.Run()
	} else if extension == "jpg" {
		// Generamos el jpg
		cmd := exec.Command("dot", "-Tjpg", dot, "-o", rep.Path)
		cmd.Run()
	} else {
		fmt.Println("No se reconoce la extensión")
		mensaje.Mensaje = "No se reconoce la extensión"
	}

	mensaje.Mensaje = "Reporte generado con éxito"
}

func reorteSP(rep *Rep, lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {
	// Abrimos el archivo
	filePart, err := os.Open(lista.ObtenerParticionMontada(rep.Id).Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		mensaje.Mensaje = "Error al abrir el archivo"
		return
	}
	defer filePart.Close()

	// Leemos el superbloque
	var sb estructuras.SuperBloque
	filePart.Seek(0, 0)
	binary.Read(filePart, binary.BigEndian, &sb)

	// Nombre del archivo
	nombreDelArchivo := strings.Split(lista.ObtenerParticionMontada(rep.Id).Path, "/")[len(strings.Split(lista.ObtenerParticionMontada(rep.Id).Path, "/"))-1]

	// Directorio
	directorio := strings.Split(rep.Path, "/")[0]
	for i := 1; i < len(strings.Split(rep.Path, "/"))-1; i++ {
		directorio += "/" + strings.Split(rep.Path, "/")[i]
	}

	// Verificamos si existe el directorio
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		os.MkdirAll(directorio, 0777)
	}

	// path con extensión .dot
	dot := strings.Split(rep.Path, ".")[0] + ".dot"
	// Creamos el archivo
	fileDot, err := os.Create(dot)
	if err != nil {
		fmt.Println("Error al crear el archivo")
		mensaje.Mensaje = "Error al crear el archivo"
		return
	}
	defer fileDot.Close()

	// Escribimos el archivo
	// Escribimos el disco
	fileDot.WriteString("digraph G {\n")
	fileDot.WriteString("labelloc=\"t\";\n")
	fileDot.WriteString("label=\"" + nombreDelArchivo + "\";\n")
	fileDot.WriteString("parent [\n")
	fileDot.WriteString("shape=plaintext\n")
	fileDot.WriteString("label=<\n")
	fileDot.WriteString("<table border=\"1\" cellborder=\"1\">\n")
	fileDot.WriteString("<tr><td bgcolor='#EA5455'>SUPERBLOQUE</td>\n")

	// Escribimos el superbloque
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_nombre: " + nombreDelArchivo + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_inodos_count: " + intToString(byte16ToInt(sb.S_inodes_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_blocks_count: " + intToString(byte16ToInt(sb.S_blocks_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_inodos_free: " + intToString(byte16ToInt(sb.S_free_inodes_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_blocks_free: " + intToString(byte16ToInt(sb.S_free_blocks_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_date_creacion: " + time.Unix(int64(binary.LittleEndian.Uint32(sb.S_mtime[:])), 0).String() + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_mount_count: " + intToString(bytesToInt(sb.S_mnt_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_magic_num: " + intToString(byte16ToInt(sb.S_magic)) + "</td></tr>\n")

}

// int to string
func intToString(input_num int) string {
	return strconv.Itoa(input_num)
}

// float64 to string
func floatToString(input_num float64) string {
	return fmt.Sprintf("%f", input_num)
}

func (rep *Rep) Rep(lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {
	switch strings.ToLower(rep.Name) {
	case "disk":
		ReporteDisk(rep, lista, mensaje)
	default:
		fmt.Println("No se reconoce el reporte")
		mensaje.Mensaje = "No se reconoce el reporte"
	}
}
