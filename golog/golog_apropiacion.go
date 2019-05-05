package golog

import (
	"fmt"

	. "github.com/mndrix/golog"
)

func ComprobarAprobacion(reglas string, regla_inyectada string) string {
	//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var res string
	for _, solution := range resultados {
		res = fmt.Sprintf("%s", solution.ByName_("Y"))
	}

	return res

}

//GetBoolean ...
func GetBoolean(reglas string, regla_inyectada string, variable_a_obtener string) bool {
	m := NewMachine().Consult(reglas)
	resultados := m.ProveAll(regla_inyectada)
	var res string
	for _, solution := range resultados {
		res = fmt.Sprintf("%s", solution.ByName_(variable_a_obtener))
	}
	if res == "0" {
		return false
	} else {
		return true
	}

}
