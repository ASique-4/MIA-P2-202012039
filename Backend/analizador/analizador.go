package analizador

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"proyecto2/comandos"
)

// Función para analizar el tipo del parámetro
func getTipoParametro(parametro string) string {
	var tipo string
	for i := 0; i < len(parametro); i++ {
		if parametro[i] == '=' {
			break
		}
		caracter := strings.ToLower(string(parametro[i]))
		tipo += caracter
	}
	return strings.TrimSpace(tipo)
}

// Función para analizar el valor del parámetro
func getValorParametro(parametro string) string {
	var valor string
	var concatenar bool
	for i := 0; i < len(parametro); i++ {
		if parametro[i] == '#' {
			break
		}
		if concatenar {
			valor += string(parametro[i])
		}
		if parametro[i] == '=' {
			concatenar = true
		}
	}
	return strings.TrimSpace(valor)
}

func estaVacia(b [4]byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func estaVaciaName(b [16]byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

// Función para analizar los parámetros del comando rmdisk
func analizarRmdisk(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var disco comandos.Rmdisk
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "path":
			disco.Path = valor
			fmt.Println("Path:", disco.Path)
		default:
			fmt.Printf("¡Error! rmdisk solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", tipo)
			return
		}
		if len(strings.SplitN(parametros, " ", 2)) > 1 {
			parametros = strings.TrimSpace(strings.SplitN(parametros, " ", 2)[1])
		} else {
			parametros = ""
		}
	}
	comandos.EliminarDiscos(disco)
}

// Función para analizar los parámetros del comando mkdisk
func analizarMkdisk(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var disco comandos.Mkdisk
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "size":
			if size, err := strconv.ParseInt(valor, 10, 64); err == nil {
				binary.LittleEndian.PutUint32(disco.Size[:], uint32(size))
			} else {
				fmt.Printf("¡Error! El valor de size no es un número válido: %v\n", valor)
				return
			}
		case "path":
			disco.Path = valor
		case "unit":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de unit debe ser un único carácter: %v\n", valor)
				return
			}
			disco.Unit = valor[0]
		case "fit":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de fit debe ser un único carácter: %v\n", valor)
				return
			}
			binary.LittleEndian.PutUint32(disco.Fit[:], uint32(valor[0]))
		default:
			fmt.Printf("¡Error! mkdisk solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
			return
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if estaVacia(disco.Size) || disco.Path == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'mkdisk'")
		return
	}
	//Creamos el disco
	comandos.CrearDiscos(disco)
}

// Función para analizar los parámetros del comando fdisk
func analizarFdisk(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var particion comandos.Fdisk
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "size":
			if size, err := strconv.ParseInt(valor, 10, 64); err == nil {
				binary.LittleEndian.PutUint32(particion.Size[:], uint32(size))
			} else {
				fmt.Printf("¡Error! El valor de size no es un número válido: %v\n", valor)
				return
			}
		case "path":
			particion.Path = valor
		case "unit":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de unit debe ser un único carácter: %v\n", valor)
				return
			}
			particion.Unit = valor[0]
		case "fit":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de fit debe ser un único carácter: %v\n", valor)
				return
			}
			binary.LittleEndian.PutUint32(particion.Fit[:], uint32(valor[0]))
		case "type":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de type debe ser un único carácter: %v\n", valor)
				return
			}
			binary.LittleEndian.PutUint32(particion.Type[:], uint32(valor[0]))
		case "name":
			if len(valor) > 16 {
				fmt.Printf("¡Error! El valor de name no puede ser mayor a 16 caracteres: %v\n", valor)
				return
			}
			binary.LittleEndian.PutUint32(particion.Name[:], uint32(valor[0]))

		default:
			fmt.Printf("¡Error! fdisk solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
			return
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if estaVacia(particion.Size) || particion.Path == "" || estaVaciaName(particion.Name) {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'fdisk'")
		return
	}

	//Creamos la particion
	comandos.CrearParticion(particion)
}

func Analizar(comando string) {
	// Lógica de análisis del comando aquí
	token := strings.TrimSpace(strings.SplitN(comando, " ", 2)[0])
	parametros := strings.TrimSpace(strings.SplitN(comando, " ", 2)[1])
	if token == "salir" {
		fmt.Println("Saliendo...")
	} else if token == "mkdisk" {
		fmt.Println("Creando disco...")
		analizarMkdisk(parametros)
	} else if token == "rmdisk" {
		fmt.Println("Eliminando disco...")
		analizarRmdisk(parametros)
	} else if token == "fdisk" {
		fmt.Println("Creando partición...")

	} else if token == "mount" {
		fmt.Println("Montando partición...")
	} else {
		fmt.Println("Comando no reconocido")
	}
}
