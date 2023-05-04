package comandos

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"proyecto2/estructuras"
	"strconv"
	"strings"
	"time"
	"unsafe"
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
		fmt.Println(cmd)
		cmd.Run()
		err := cmd.Run()
		fmt.Println(err)
	} else {
		fmt.Println("No se reconoce la extensión")
		mensaje.Mensaje = "No se reconoce la extensión"
	}
	mensaje.Base64 = imageToBase64(rep.Path)
	mensaje.Reporte = "DISK"

	mensaje.Mensaje = "Reporte DISK generado con éxito"
}

func imageToBase64(path string) string {
	// Abrimos el archivo
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		return ""
	}
	defer file.Close()

	// Leemos el archivo
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// Leemos el archivo
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	// Codificamos a base64
	encodedString := base64.StdEncoding.EncodeToString(bytes)

	return encodedString
}

func ReporteSB(rep *Rep, lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {
	// Abrimos el archivo
	filePart, err := os.Open(lista.ObtenerParticionMontada(rep.Id).Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo")
		mensaje.Mensaje = "Error al abrir el archivo"
		return
	}
	defer filePart.Close()

	// Leemos el mbr
	mbr := estructuras.MBR{}
	filePart.Seek(0, 0)
	binary.Read(filePart, binary.BigEndian, &mbr)

	// Buscamos la partición
	var particion estructuras.Particion
	particiones := [4]estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for i := 0; i < 4; i++ {
		if particiones[i].Part_name == lista.ObtenerParticionMontada(rep.Id).Name {
			particion = particiones[i]
			break
		}
	}

	// Nos posicionamos en el inicio de la partición
	filePart.Seek(int64(bytesToInt(particion.Part_start)), 0)

	// Leemos el superbloque
	sb := estructuras.SuperBloque{}
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
	fileDot.WriteString("<tr><td bgcolor='#EA5455'>SUPERBLOQUE</td></tr>\n")

	// Escribimos el superbloque
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_nombre: " + nombreDelArchivo + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_inodos_count: " + intToString(byte16ToInt(sb.S_inodes_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_blocks_count: " + intToString(byte16ToInt(sb.S_blocks_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_inodos_free: " + intToString(byte16ToInt(sb.S_free_inodes_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_blocks_free: " + intToString(byte16ToInt(sb.S_free_blocks_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_date_creacion: " + time.Unix(int64(binary.LittleEndian.Uint32(sb.S_mtime[:])), 0).String() + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#E4DCCF'>sb_mount_count: " + intToString(bytesToInt(sb.S_mnt_count)) + "</td></tr>\n")
	fileDot.WriteString("<tr><td bgcolor='#F9F5EB'>sb_magic_num: 0xEF53</td></tr>\n")

	fileDot.WriteString("</table>\n")
	fileDot.WriteString(">];\n")
	fileDot.WriteString("}\n")

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
		fmt.Println(cmd)
		cmd.Run()
		err := cmd.Run()
		fmt.Println(err)
	} else {
		fmt.Println("No se reconoce la extensión")
		mensaje.Mensaje = "No se reconoce la extensión"
	}
	mensaje.Base64 = imageToBase64(rep.Path)
	mensaje.Reporte = "SB"

	mensaje.Mensaje = "Reporte SB generado con éxito"

}

func ReporteTREE(rep *Rep, lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {

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
		fmt.Println(err)
		fmt.Println("Error al crear el archivo")
		mensaje.Mensaje = "Error al crear el archivo"
		return
	}
	defer fileDot.Close()

	// Leemos el MBR
	fileMBR, err := os.OpenFile(lista.ObtenerParticionMontada(rep.Id).Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		fileMBR.Close()
		return
	}

	// Leemos el superbloque
	// Nos movemos al inicio del archivo
	_, err = fileMBR.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		fileMBR.Close()
		return
	}

	// Leemos el MBR
	mbr := estructuras.MBR{}
	binary.Read(fileMBR, binary.BigEndian, &mbr)

	// Buscamos la partición
	var particion estructuras.Particion
	particiones := [4]estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for i := 0; i < 4; i++ {
		if particiones[i].Part_name == lista.ObtenerParticionMontada(rep.Id).Name {
			particion = particiones[i]
			break
		}
	}

	// Nos posicionamos en el inicio de la partición
	fileMBR.Seek(int64(bytesToInt(particion.Part_start)), 0)

	// Leemos el superbloque
	sb := estructuras.SuperBloque{}
	binary.Read(fileMBR, binary.BigEndian, &sb)

	// Leemos el conenido del archivo
	fileMBR.Seek(int64(byte16ToInt(sb.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)
	contenido := [64]byte{}
	binary.Read(fileMBR, binary.LittleEndian, &contenido)

	// Escribimos el archivo
	// Escribimos el disco
	fileDot.WriteString("digraph G {\n")
	fileDot.WriteString("label=\"" + nombreDelArchivo + "\";\n")
	fileDot.WriteString("rankdir=LR;\n")
	// Nodo raiz
	fileDot.WriteString("\"node0\"  [label=\"<f0>Inodo 0|<f1>i_type: 0|<f2>Ap0: Bloq0|<f3>Ap1: -1|<f3>Ap3: -1|<f4>AP4: -1|<f5>AP5: -1|<f6>AP6: -1|<f7>AP7: -1|<f8>AP8: -1|<f9>AP9: -1|<f10>AP10: -1|<f11>AP11: -1|<f12>AP12: -1|<f13>AP13: -1|<f14>AP14: -1|<f15>AP15: -1\" shape=\"record\" style=filled fillcolor=\"cadetblue1\"];\n")
	// Bloques de apuntadores
	fileDot.WriteString("\"node1\"  [label=\"<f0>Bloque 0|<f1>users.txt: Inodo 1|<f2>.: 0|<f3>..: 0\" shape=\"record\" style=filled fillcolor=\"darkolivegreen1\"];\n")
	// Inodos
	fileDot.WriteString("\"node2\"  [label=\"<f0>Inodo 1|<f1>i_type: 1|<f2>Ap0: Bloq1|<f3>Ap1: -1|<f3>Ap3: -1|<f4>AP4: -1|<f5>AP5: -1|<f6>AP6: -1|<f7>AP7: -1|<f8>AP8: -1|<f9>AP9: -1|<f10>AP10: -1|<f11>AP11: -1|<f12>AP12: -1|<f13>AP13: -1|<f14>AP14: -1|<f15>AP15: -1\" shape=\"record\" style=filled fillcolor=\"cadetblue1\"];\n")
	// Contenido del archivo
	fileDot.WriteString("\"node3\"  [label=\"<f0>Bloque. Archivo 1|<f1>Contenido: " + obtenerContenido(contenido) + "\" shape=\"record\" style=filled fillcolor=\"darkolivegreen1\"];\n")
	fileDot.WriteString("subgraph cluster_bm_inodos {\n")
	fileDot.WriteString("label=\"Bitmap de Inodos\"\n")
	fileDot.WriteString("bm_inodos [\n")
	fileDot.WriteString("shape=plaintext\n")
	fileDot.WriteString("label=<\n")
	fileDot.WriteString("<table border='1' cellborder='1'>\n")
	fileDot.WriteString("<tr><td colspan=\"20\" bgcolor='#EA5455'>Reporte de BITMAP DE INODOS</td></tr>\n")
	fileDot.WriteString(fmt.Sprintf("<tr><td colspan=\"20\" bgcolor='#E4DCCF'>Nombre del Disco: %s</td></tr>\n", nombreDelArchivo))
	for i := 0; i < byte16ToInt(sb.S_blocks_count)/5; i++ {
		if i%20 == 0 {
			fileDot.WriteString("<tr>\n")
		}
		if i == 0 {
			fileDot.WriteString("<td bgcolor='#E4DCCF'>1</td>")
		} else {
			fileDot.WriteString("<td bgcolor='#E4DCCF'>0</td>")
		}
		if i%20 == 19 {
			fileDot.WriteString("</tr>\n")
		}
	}
	fileDot.WriteString("</tr>\n")
	fileDot.WriteString("</table>\n")
	fileDot.WriteString(">];\n")
	fileDot.WriteString("}\n")

	fileDot.WriteString("subgraph cluster_bm_bloques {\n")
	fileDot.WriteString("label=\"Bitmap de Bloques\"\n")
	fileDot.WriteString("bm_bloques [\n")
	fileDot.WriteString("shape=plaintext\n")
	fileDot.WriteString("label=<\n")
	fileDot.WriteString("<table border='1' cellborder='1'>\n")
	fileDot.WriteString("<tr><td colspan=\"20\" bgcolor='#EA5455'>Reporte de BITMAP DE BLOQUES</td></tr>\n")
	fileDot.WriteString(fmt.Sprintf("<tr><td colspan=\"20\" bgcolor='#E4DCCF'>Nombre del Disco: %s</td></tr>\n", nombreDelArchivo))
	for i := 0; i < byte16ToInt(sb.S_blocks_count)/5; i++ {
		if i%20 == 0 {
			fileDot.WriteString("<tr>\n")
		}
		if i == 0 {
			fileDot.WriteString("<td bgcolor='#E4DCCF'>1</td>")
		} else {
			fileDot.WriteString("<td bgcolor='#E4DCCF'>0</td>")
		}
		if i%20 == 19 {
			fileDot.WriteString("</tr>\n")
		}
	}
	fileDot.WriteString("</tr>\n")
	fileDot.WriteString("</table>\n")
	fileDot.WriteString(">];\n")
	fileDot.WriteString("}\n")

	fileDot.WriteString("\"node0\":f2 -> \"node1\":f0;\n")
	fileDot.WriteString("\"node1\":f1 -> \"node2\":f0;\n")
	fileDot.WriteString("\"node2\":f1 -> \"node3\":f0;\n")
	fileDot.WriteString("\"bm_bloques\" -> \"bm_inodos\";\n")

	fileDot.WriteString("}")

	// Cerramos el archivo
	fileDot.Close()

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
		fmt.Println(cmd)
		cmd.Run()
		err := cmd.Run()
		fmt.Println(err)
	} else {
		fmt.Println("No se reconoce la extensión")
		mensaje.Mensaje = "No se reconoce la extensión"
	}

	mensaje.Base64 = imageToBase64(rep.Path)
	mensaje.Reporte = "TREE"

	mensaje.Mensaje = "Reporte TREE generado con éxito"

}

func obtenerContenido(contenido [64]byte) string {
	var cadena string
	for i := 0; i < len(contenido); i++ {
		if contenido[i] == 0 {
			break
		}
		cadena = cadena + string(contenido[i])
	}
	return cadena
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
		fmt.Println("Generando reporte DISK")
		ReporteDisk(rep, lista, mensaje)
	case "sb":
		fmt.Println("Generando reporte SB")
		ReporteSB(rep, lista, mensaje)
	case "tree":
		fmt.Println("Generando reporte TREE")
		ReporteTREE(rep, lista, mensaje)
	default:
		fmt.Println("No se reconoce el reporte")
		mensaje.Mensaje = "No se reconoce el reporte"
	}

	time.Sleep(3 * time.Second)

}
