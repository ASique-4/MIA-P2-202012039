package comandos

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"proyecto2/estructuras"
	"time"
	"unsafe"
)

type MKFS struct {
	Id   string
	Type string
}

func imprimirSuperBloque(superbloque estructuras.SuperBloque) {
	fmt.Println("S_filesystem_type: ", string(superbloque.S_filesystem_type[:]))
	fmt.Println("S_inodes_count: ", byte16ToInt(superbloque.S_inodes_count))
	fmt.Println("S_blocks_count: ", byte16ToInt(superbloque.S_blocks_count))
	fmt.Println("S_free_blocks_count: ", byte16ToInt(superbloque.S_free_blocks_count))
	fmt.Println("S_free_inodes_count: ", byte16ToInt(superbloque.S_free_inodes_count))
	fmt.Println("S_mtime: ", time.Unix(int64(binary.LittleEndian.Uint64(superbloque.S_mtime[:])), 0))
	fmt.Println("S_mnt_count: ", bytesToInt(superbloque.S_mnt_count))
	fmt.Println("S_magic: ", bytes16ToText(superbloque.S_magic))
	fmt.Println("S_inode_size: ", bytesToInt(superbloque.S_inode_size))
	fmt.Println("S_block_size: ", bytesToInt(superbloque.S_block_size))
	fmt.Println("S_first_ino: ", bytesToInt(superbloque.S_first_ino))
	fmt.Println("S_first_blo: ", bytesToInt(superbloque.S_first_blo))
	fmt.Println("S_bm_inode_start: ", byte16ToInt(superbloque.S_bm_inode_start))
	fmt.Println("S_bm_block_start: ", byte16ToInt(superbloque.S_bm_block_start))
	fmt.Println("S_inode_start: ", byte16ToInt(superbloque.S_inode_start))
	fmt.Println("S_block_start: ", byte16ToInt(superbloque.S_block_start))
}

func bytes16ToText(bytes [16]byte) string {
	var text string
	for i := 0; i < 16; i++ {
		text += string(bytes[i])
	}
	return text
}

func crearEXT2(file *os.File, particion estructuras.Particion, particionMontada *estructuras.ParticionMontada, mensaje *estructuras.Mensaje) {
	// Nos posicionamos en el inicio de la partición
	file.Seek(int64(bytesToInt(particion.Part_start)), 0)
	// Creamos el superbloque
	superbloque := estructuras.SuperBloque{}
	superbloque.S_filesystem_type = [4]byte{0x45, 0x58, 0x54, 0x32}
	//n = (tamaño_particion - sizeof(superblock)) / (4 + sizeof(inodos) + 3 * sizeof(block))                                             // EXT2 en ASCII
	inodes_count := math.Floor((float64(bytesToInt(particion.Part_size)) - float64(unsafe.Sizeof(estructuras.SuperBloque{}))) / (4 + float64(unsafe.Sizeof(estructuras.Inodo{})) + 3*float64(unsafe.Sizeof(estructuras.BloqueArchivo{}))))

	superbloque.S_inodes_count = [16]byte{byte(inodes_count)}
	superbloque.S_blocks_count = [16]byte{byte(inodes_count * 3)}
	superbloque.S_free_blocks_count = [16]byte{byte(byte16ToInt(superbloque.S_blocks_count) - 1)}
	superbloque.S_free_inodes_count = [16]byte{byte(inodes_count - 1)}
	// Fecha de creación
	binary.LittleEndian.PutUint64(superbloque.S_mtime[:], uint64(time.Now().Unix()))
	superbloque.S_mnt_count = [4]byte{byte(particionMontada.Mount_count)}
	superbloque.S_magic = [16]byte{0x30, 0x78, 0x45, 0x46, 0x35, 0x33}
	superbloque.S_inode_size = [4]byte{byte(unsafe.Sizeof(estructuras.Inodo{}))}
	superbloque.S_block_size = [4]byte{byte(unsafe.Sizeof(estructuras.BloqueArchivo{}))}
	superbloque.S_first_ino = [4]byte{byte(0)}
	superbloque.S_first_blo = [4]byte{byte(0)}
	superbloque.S_bm_inode_start = [16]byte{byte(bytesToInt(particion.Part_start) + int(unsafe.Sizeof(estructuras.SuperBloque{})))}
	superbloque.S_bm_block_start = [16]byte{byte(byte16ToInt(superbloque.S_bm_inode_start) + int(inodes_count))}
	superbloque.S_inode_start = [16]byte{byte(byte16ToInt(superbloque.S_bm_block_start) + int(inodes_count))}
	superbloque.S_block_start = [16]byte{byte(byte16ToInt(superbloque.S_inode_start) + (int(inodes_count) * int(unsafe.Sizeof(estructuras.Inodo{}))))}

	// Escribimos el superbloque
	binary.Write(file, binary.BigEndian, &superbloque)

	// Creamos el bitmap de inodos
	file.Seek(int64(byte16ToInt(superbloque.S_bm_inode_start)), 0)

	bmInodos := make([]byte, int(inodes_count))
	for i := 0; i < len(bmInodos); i++ {
		if i == 0 {
			bmInodos[i] = 1
		} else {
			bmInodos[i] = 0
		}
	}

	binary.Write(file, binary.BigEndian, &bmInodos)

	// Creamos el bitmap de bloques
	file.Seek(int64(byte16ToInt(superbloque.S_bm_block_start)), 0)

	bmBloques := make([]byte, int(inodes_count*3))
	for i := 0; i < len(bmBloques); i++ {
		if i == 0 {
			bmBloques[i] = 1
		} else {
			bmBloques[i] = 0
		}
	}

	binary.Write(file, binary.BigEndian, &bmBloques)

	// Creamos los inodos
	file.Seek(int64(byte16ToInt(superbloque.S_inode_start)), 0)

	inodos := make([]estructuras.Inodo, int(inodes_count))
	for i := 0; i < len(inodos); i++ {
		if i == 0 {
			inodos[i].I_uid = [4]byte{byte(1)}
			inodos[i].I_gid = [4]byte{byte(1)}
			inodos[i].I_size = [16]byte{byte(unsafe.Sizeof(estructuras.BloqueCarpeta{}))}
			inodos[i].I_atime = [16]byte{byte(0)}
			inodos[i].I_ctime = [16]byte{byte(0)}
			binary.LittleEndian.PutUint64(inodos[i].I_mtime[:], uint64(time.Now().Unix()))
			inodos[i].I_block = [16]byte{byte(0)}
			inodos[i].I_type = [1]byte{byte(1)}
			inodos[i].I_perm = [4]byte{6, 6, 4} // 664
		} else {
			inodos[i].I_uid = [4]byte{byte(0)}
			inodos[i].I_gid = [4]byte{byte(0)}
			inodos[i].I_size = [16]byte{byte(0)}
			inodos[i].I_atime = [16]byte{byte(0)}
			inodos[i].I_ctime = [16]byte{byte(0)}
			inodos[i].I_mtime = [16]byte{byte(0)}
			inodos[i].I_block = [16]byte{byte(0)}
		}
	}

	binary.Write(file, binary.BigEndian, &inodos)

	// Creamos los bloques
	file.Seek(int64(byte16ToInt(superbloque.S_block_start)), 0)

	bloques := make([]estructuras.BloqueArchivo, int(inodes_count*3))
	for i := 0; i < len(bloques); i++ {
		bloques[i].B_content = [64]byte{byte(0)}
	}

	binary.Write(file, binary.LittleEndian, &bloques)

	// Creamos el archivo de usuarios
	file.Seek(int64(byte16ToInt(superbloque.S_block_start))+int64(unsafe.Sizeof(estructuras.BloqueCarpeta{})), 0)
	contenido := [64]byte{}
	copy(contenido[:], "1,G,root\n1,U,root,root,123\n")
	binary.Write(file, binary.LittleEndian, &contenido)
	mensaje.Mensaje = "Se formateó la partición correctamente."

}

func byte16ToInt(b [16]byte) int {
	return int(binary.BigEndian.Uint16(b[:]))
}

func (mkfs *MKFS) FormatearParticion(lista *estructuras.ListaParticionesMontadas, mensaje *estructuras.Mensaje) {
	// Obtenemos la partición
	particion := lista.ObtenerParticionMontada(mkfs.Id)
	// Abrimos el archivo
	file, err := os.OpenFile(particion.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		mensaje.Mensaje = "ERROR: No se pudo abrir el archivo."
		return
	}
	defer file.Close()
	// Leemos el MBR
	mbr := estructuras.MBR{}
	binary.Read(file, binary.BigEndian, &mbr)
	// Buscamos la partición
	Particiones := []estructuras.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for i := 0; i < len(Particiones); i++ {
		if Particiones[i].Part_name == particion.Name {
			// Formateamos la partición
			crearEXT2(file, Particiones[i], particion, mensaje)
			// Escribimos el MBR
			file.Seek(0, 0)
			binary.Write(file, binary.BigEndian, &mbr)
			break
		}
	}
}
