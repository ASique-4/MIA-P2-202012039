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
func analizarRmdisk(parametros string, w *estructuras.Mensaje, confirmar bool) {
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
	comandos.EliminarDiscos(disco, w, confirmar)
}

// Función para analizar los parámetros del comando mkdisk
func analizarMkdisk(parametros string, w *estructuras.Mensaje) {
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
			if len(valor) != 2 {
				fmt.Printf("¡Error! El valor de fit debe ser un único carácter: %v\n", valor)
				return
			}

			disco.Fit[0] = valor[0]
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
	comandos.CrearDiscos(disco, w)
}

// Función para analizar los parámetros del comando fdisk
func analizarFdisk(parametros string, w *estructuras.Mensaje) {

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
			if len(valor) != 2 {
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
	comandos.CrearParticion(particion, w)
}

// Lista global de particiones montadas
var particionesMontadas estructuras.ListaParticionesMontadas

// Función para analizar los parámetros del comando mount
func analizarMount(parametros string, w *estructuras.Mensaje) {
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
	particion.MountCommand(&particionesMontadas, w)
	particionesMontadas.ImprimirListaParticionesMontadas()
}

func analizarMKFS(parametros string, w *estructuras.Mensaje) {
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
	particion.FormatearParticion(&particionesMontadas, w)

}

func analizarREP(parametros string, w *estructuras.Mensaje) {
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
	comando.Rep(&particionesMontadas, w)

}

// Usurio actual
var usuarioActual *estructuras.Usuario

func analizarLogin(parametros string, w *estructuras.Mensaje) {
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
	usuarioActual = comando.Login(&particionesMontadas, w)

}

func analizarLogout(mensaje *estructuras.Mensaje) {
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		mensaje.Mensaje = "¡Error! No hay ningún usuario logueado"
		return
	}
	usuarioActual = new(estructuras.Usuario)
	fmt.Println("Se ha cerrado la sesión correctamente")
	mensaje.Mensaje = "Se ha cerrado la sesión correctamente"
}

func analizarMkgrp(parametros string, w *estructuras.Mensaje) {
	// Si no hay sesion iniciada
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		w.Mensaje = "¡Error! No hay ningún usuario logueado"
		return
	}
	// Si no es usuario root
	if usuarioActual.Username != "root" {
		fmt.Println("¡Error! Solo el usuario root puede crear grupos")
		w.Mensaje = "¡Error! Solo el usuario root puede crear grupos"
		return
	}
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
	comando.Mkgrp(usuarioActual.PartID, &particionesMontadas, w)

}

func analizarRmgrp(parametros string, w *estructuras.Mensaje) {
	// Si no hay sesion iniciada
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		w.Mensaje = "¡Error! No hay ningún usuario logueado"
		return
	}
	// Si no es usuario root
	if usuarioActual.Username != "root" {
		fmt.Println("¡Error! Solo el usuario root puede crear grupos")
		w.Mensaje = "¡Error! Solo el usuario root puede crear grupos"
		return
	}
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
	comando.Rmgrp(usuarioActual.PartID, &particionesMontadas, w)
}

func analizarMkuser(parametros string, w *estructuras.Mensaje) {
	// Si no hay sesion iniciada
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		w.Mensaje = "¡Error! No hay ningún usuario logueado"
		return
	}
	// Si no es usuario root
	if usuarioActual.Username != "root" {
		fmt.Println("¡Error! Solo el usuario root puede crear grupos")
		w.Mensaje = "¡Error! Solo el usuario root puede crear grupos"
		return
	}
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Mkuser
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "user":
			comando.Usr = valor
		case "pwd":
			comando.Pwd = valor
		case "grp":
			comando.Grp = valor
		default:
			fmt.Printf("¡Error! mkusr solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.Usr == "" || comando.Pwd == "" || comando.Grp == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'mkusr'")
		return
	}

	//Creamos el reporte
	comando.Mkuser(usuarioActual.PartID, &particionesMontadas, w)
}

func analizarRmusr(parametros string, w *estructuras.Mensaje) {
	// Si no hay sesion iniciada
	if usuarioActual == (nil) {
		fmt.Println("¡Error! No hay ningún usuario logueado")
		w.Mensaje = "¡Error! No hay ningún usuario logueado"
		return
	}
	// Si no es usuario root
	if usuarioActual.Username != "root" {
		fmt.Println("¡Error! Solo el usuario root puede crear grupos")
		w.Mensaje = "¡Error! Solo el usuario root puede crear grupos"
		return
	}
	parametros = strings.TrimSpace(strings.SplitN(parametros, ">", 2)[1])
	var comando comandos.Rmusr
	for parametros != "" {
		tmpParam := parametros
		tipo := getTipoParametro(tmpParam)
		valor := strings.TrimSpace(strings.SplitN(getValorParametro(tmpParam), " ", 2)[0])
		switch tipo {
		case "user":
			comando.User = valor
		default:
			fmt.Printf("¡Error! rmusr solo acepta parámetros válidos, ¿qué intentas hacer con '%v'?\n", valor)
		}
		if index := strings.Index(parametros, ">"); index >= 0 {
			parametros = parametros[index+1:]
		} else {
			parametros = ""
		}

		parametros = strings.TrimSpace(parametros)
	}
	//Verificamos que los parametros obligatorios esten
	if comando.User == "" {
		fmt.Println("¡Error! Parece que alguien olvidó poner los parámetros en 'rmusr'")
		return
	}

	//Creamos el reporte
	comando.Rmusr(usuarioActual.PartID, &particionesMontadas, w)
}

func Analizar(comando string, mensaje *estructuras.Mensaje, confirmar bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Ocurrió un error:", r)
			mensaje.Mensaje = "Ocurrió un error: " + fmt.Sprint(r)
		}
	}()
	//Si es pause
	if comando == "pause" {
		fmt.Println("Pausado...")
		mensaje.Mensaje = "Pausado..."
		mensaje.Accion = "pause"
		return
	}

	// Lógica de análisis del comando aquí
	token := strings.TrimSpace(strings.SplitN(comando, " ", 2)[0])
	parametros := strings.TrimSpace(strings.SplitN(comando, " ", 2)[1])
	if token == "salir" {
		fmt.Println("Saliendo...")
	} else if token == "mkdisk" {
		fmt.Println("Creando disco...")
		mensaje.Accion = "Creando disco..."
		analizarMkdisk(parametros, mensaje)
	} else if token == "rmdisk" {
		fmt.Println("Eliminando disco...")
		mensaje.Accion = "Eliminando disco..."
		analizarRmdisk(parametros, mensaje, confirmar)
	} else if token == "fdisk" {
		fmt.Println("Creando partición...")
		mensaje.Accion = "Creando partición..."
		analizarFdisk(parametros, mensaje)
	} else if token == "mount" {
		fmt.Println("Montando partición...")
		mensaje.Accion = "Montando partición..."
		analizarMount(parametros, mensaje)
	} else if token == "mkfs" {
		fmt.Println("Creando sistema de archivos...")
		mensaje.Accion = "Creando sistema de archivos..."
		analizarMKFS(parametros, mensaje)
	} else if token == "rep" {
		fmt.Println("Creando reporte...")
		mensaje.Accion = "Creando reporte..."
		analizarREP(parametros, mensaje)
	} else if token == "login" {
		fmt.Println("Iniciando sesión...")
		mensaje.Accion = "Iniciando sesión..."
		analizarLogin(parametros, mensaje)
	} else if token == "logout" {
		fmt.Println("Cerrando sesión...")
		mensaje.Accion = "Cerrando sesión..."
		analizarLogout(mensaje)
	} else if token == "mkgrp" {
		fmt.Println("Creando grupo...")
		mensaje.Accion = "Creando grupo..."
		analizarMkgrp(parametros, mensaje)
	} else if token == "rmgrp" {
		fmt.Println("Eliminando grupo...")
		mensaje.Accion = "Eliminando grupo..."
		analizarRmgrp(parametros, mensaje)
	} else if token == "mkusr" {
		fmt.Println("Creando usuario...")
		mensaje.Accion = "Creando usuario..."
		analizarMkuser(parametros, mensaje)
	} else if token == "rmusr" {
		fmt.Println("Eliminando usuario...")
		mensaje.Accion = "Eliminando usuario..."
		analizarRmusr(parametros, mensaje)
	} else {
		fmt.Println("Comando no reconocido")
		mensaje.Accion = "Comando no reconocido"

	}
}
