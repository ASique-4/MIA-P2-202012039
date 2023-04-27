package analizador

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"proyecto2/comandos"
	"proyecto2/estructuras"
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
			particion.Fit = [1]byte{valor[0]}
		case "type":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de type debe ser un único carácter: %v\n", valor)
				return
			}
			particion.Type = [1]byte{valor[0]}
		case "name":
			if len(valor) > 16 {
				fmt.Printf("¡Error! El valor de name no puede ser mayor a 16 caracteres: %v\n", valor)
				return
			}
			for i := 0; i < len(valor); i++ {
				particion.Name[i] = valor[i]
			}

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

	//Si el fit esta vacio, lo ponemos por defecto
	if particion.Fit[0] == 0 {
		particion.Fit = [1]byte{'W'}
	}

	//Si el type esta vacio, lo ponemos por defecto
	if particion.Type[0] == 0 {
		particion.Type = [1]byte{'P'}
	}

	//Creamos la particion
	comandos.CrearParticion(particion)
}

// Lista global de particiones montadas
var particionesMontadas estructuras.ListaParticionesMontadas

// Función para analizar los parámetros del comando mount
func analizarMount(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var particion comandos.Mount
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "path":
			particion.Path = valor
		case "name":
			if len(valor) > 16 {
				fmt.Printf("¡Error! El valor de name no puede ser mayor a 16 caracteres: %v\n", valor)
				return
			}
			for i := 0; i < len(valor); i++ {
				particion.Name[i] = valor[i]
			}
		default:
			fmt.Printf("¡Error! mount solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
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
	if particion.Path == "" || estaVaciaName(particion.Name) {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'mount'")
		return
	}

	//Montamos la particion
	particion.MountCommand(&particionesMontadas)
	particionesMontadas.ImprimirListaParticionesMontadas()
}

func analizarMKFS(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var particion comandos.MKFS
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "id":
			if len(valor) > 16 {
				fmt.Printf("¡Error! El valor de id no puede ser mayor a 16 caracteres: %v\n", valor)
				return
			}
			particion.Id = valor
		case "type":
			if len(valor) != 1 {
				fmt.Printf("¡Error! El valor de type debe ser un único carácter: %v\n", valor)
				return
			}
			particion.Type = valor
		default:
			fmt.Printf("¡Error! mkfs solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
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
	if particion.Id == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'mkfs'")
		return
	}

	//Creamos el sistema de archivos
	particion.FormatearParticion(&particionesMontadas)

}

func analizarREP(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Rep
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "name":
			comando.Name = valor
		case "path":
			comando.Path = valor
		case "id":
			comando.Id = valor
		case "ruta":
			comando.Ruta = valor
		default:
			fmt.Printf("¡Error! rep solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.Name == "" || comando.Path == "" || comando.Id == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'rep'")
		return
	}

	//Creamos el reporte
	comando.Rep(&particionesMontadas)

}

// Usurio actual
var usuarioActual *estructuras.Usuario

func analizarLogin(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Login
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "user":
			comando.Usuario = valor
		case "pwd":
			comando.Pass = valor
		case "id":
			comando.Id = valor
		default:
			fmt.Printf("¡Error! login solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.Usuario == "" || comando.Pass == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'login'")
		return
	}

	//Guaramos el usuario
	usuarioActual = comando.Login(&particionesMontadas)

}

func analizarLogout() {
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		return
	}
	usuarioActual = new(estructuras.Usuario)
	fmt.Println("Se ha cerrado la sesión correctamente")
}

func analizarMkgrp(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Mkgrp
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "name":
			comando.Name = valor
		default:
			fmt.Printf("¡Error! mkgrp solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.Name == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'mkgrp'")
		return
	}

	//Creamos el reporte
	comando.Mkgrp(usuarioActual.PartID, &particionesMontadas)

}

func analizarRmgrp(parametros string) {
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Rmgrp
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "name":
			comando.Name = valor
		default:
			fmt.Printf("¡Error! rmgrp solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.Name == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'rmgrp'")
		return
	}

	//Creamos el reporte
	comando.Rmgrp(usuarioActual.PartID, &particionesMontadas)
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
		analizarFdisk(parametros)
	} else if token == "mount" {
		fmt.Println("Montando partición...")
		analizarMount(parametros)
	} else if token == "mkfs" {
		fmt.Println("Creando sistema de archivos...")
		analizarMKFS(parametros)
	} else if token == "rep" {
		fmt.Println("Creando reporte...")
		analizarREP(parametros)
	} else if token == "login" {
		fmt.Println("Iniciando sesión...")
		analizarLogin(parametros)
	} else if token == "logout" {
		fmt.Println("Cerrando sesión...")
		analizarLogout()
	} else if token == "mkgrp" {
		fmt.Println("Creando grupo...")
		analizarMkgrp(parametros)
	} else if token == "rmgrp" {
		fmt.Println("Eliminando grupo...")
		analizarRmgrp(parametros)
	} else {
		fmt.Println("Comando no reconocido")
	}
}
