package analizador

import (
	"strings"
)

type Scrpit struct {
	contenido string
}

func (mkfs *Scrpit) Ejecutar() {
	//Analizar el contenido del script
	contenido := strings.Split(mkfs.contenido, "\n")
	for i := 0; i < len(contenido); i++ {
		if strings.TrimSpace(contenido[i]) == "" {
			continue
		}
		//Si es un comentario
		if strings.TrimSpace(contenido[i])[0] == '#' {
			continue
		}
		Analizar(contenido[i])
	}
}
