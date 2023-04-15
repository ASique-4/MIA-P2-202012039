package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"proyecto2/analizador"
)

func LeerEntrada() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func main() {
	titulo := "Proyecto 1 - MIA"
	nombre := "Angel Francisco Sique Santos"
	codigo := "202012039"
	mensaje := "Ingrese el comando a analizar"

	ancho := 40
	fmt.Printf("+" + strings.Repeat("-", ancho-2) + "+\n")
	fmt.Printf("|" + strings.Repeat(" ", ancho-2) + "|\n")
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(titulo))/2), titulo, strings.Repeat(" ", ((ancho-len(titulo))/2)-2))
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(nombre))/2), nombre, strings.Repeat(" ", ((ancho-len(nombre))/2)-2))
	fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", (ancho-len(codigo))/2), codigo, strings.Repeat(" ", ((ancho-len(codigo))/2)-1))
	fmt.Printf("|" + strings.Repeat(" ", ancho-2) + "|\n")
	fmt.Printf("+" + strings.Repeat("-", ancho-2) + "+\n\n")

	fmt.Println(mensaje)

	repetir := true
	for repetir {
		fmt.Print("~ ")
		comando := LeerEntrada()
		analizador.Analizar(comando)
	}
}
