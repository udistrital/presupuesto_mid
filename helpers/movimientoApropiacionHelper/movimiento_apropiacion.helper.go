package movimientoApropiacionHelper

import (
	"math"
	"strconv"
	"strings"

	"github.com/udistrital/presupuesto_mid/helpers/apropiacionHelper"
	"github.com/udistrital/utils_oas/formatdata"
)

const separator = "-"

// movimientoApropiacionAffectCalculation get Apropiation's tree balnce.
func movimientoApropiacionAffectCalculation(afectacion map[string]interface{}, res map[string]float64) {
	var tipo map[string]interface{}
	var idTipo int
	var cuentaCredito map[string]interface{}
	var cuentaContraCredito map[string]interface{}
	var multiplicador float64
	formatdata.FillStructP(afectacion["TipoMovimientoApropiacion"], &tipo)
	formatdata.FillStructP(tipo["Id"], &idTipo)

	formatdata.FillStructP(afectacion["CuentaCredito"], &cuentaCredito)
	formatdata.FillStructP(afectacion["CuentaContraCredito"], &cuentaContraCredito)
	UnidadEjecutora, err := strconv.Atoi(cuentaCredito["UnidadEjecutora"].(string))

	if err != nil {
		panic(err.Error())
	}

	switch cond := idTipo; cond {
	case 3: // Adicion
		multiplicador = 1
	case 4:
		multiplicador = 0
	default:
		multiplicador = -1
	}

	sumValorMovimientoAPropiacion(false, cuentaCredito["Codigo"].(string), UnidadEjecutora, 2018, afectacion["Valor"].(float64)*multiplicador, res)
	if cuentaContraCredito != nil {
		sumValorMovimientoAPropiacion(false, cuentaContraCredito["Codigo"].(string), UnidadEjecutora, 2018, afectacion["Valor"].(float64), res)
	}
}

func sumValorMovimientoAPropiacion(final bool, codigoRubro string, unidadEjecutora, vigencia int, valorMov float64, res map[string]float64) {
	var valorFinal float64
	var saldoObj map[string]float64
	var valorInicial float64
	codigoPadre := strings.Split(codigoRubro, separator)
	if final {
		saldoObj = apropiacionHelper.CalcularSaldoApropiacion(codigoPadre[0], unidadEjecutora, vigencia)
		valorInicial = saldoObj["valor_inicial"]
	} else {
		valorInicial = 0
	}
	valorFinal = valorInicial + valorMov
	res[codigoPadre[0]] = res[codigoPadre[0]] + valorFinal
	return
}

func ComprobacionMovimiento(afectacion []map[string]interface{}, UE, vigencia int) (map[string]float64, float64, bool) {
	res := make(map[string]float64)
	for _, element := range afectacion {
		calcularAfectacionMovimientoApropiacion(element, vigencia, res)
	}
	sumValorMovimientoAPropiacion(true, "3", UE, vigencia, 0, res)
	sumValorMovimientoAPropiacion(true, "2", UE, vigencia, 0, res)
	diff := math.Abs(res["2"] - res["3"])
	if res["2"] != res["3"] {
		return res, diff, false
	}
	return res, diff, true
}

// calcularAfectacionMovimientoApropiacion ... Calcula la afectacion de un movimiento en el arbol
// Antes de realizar la operacion de registro en la db.
func calcularAfectacionMovimientoApropiacion(afectacion map[string]interface{}, vigencia int, res map[string]float64) {

	var tipo map[string]interface{}
	var idTipo int
	var cuentaCredito map[string]interface{}
	var cuentaContraCredito map[string]interface{}
	var multiplicador float64
	formatdata.FillStructP(afectacion["TipoMovimientoApropiacion"], &tipo)
	formatdata.FillStructP(tipo["Id"], &idTipo)

	formatdata.FillStructP(afectacion["CuentaCredito"], &cuentaCredito)
	formatdata.FillStructP(afectacion["CuentaContraCredito"], &cuentaContraCredito)
	UnidadEjecutora, err := strconv.Atoi(cuentaCredito["UnidadEjecutora"].(string))

	if err != nil {
		panic(err.Error())
	}

	switch cond := idTipo; cond {
	case 3: // Adicion
		multiplicador = 1
	case 4:
		multiplicador = 0
	default:
		multiplicador = -1
	}

	sumValorMovimientoAPropiacion(false, cuentaCredito["Codigo"].(string), UnidadEjecutora, vigencia, afectacion["Valor"].(float64)*multiplicador, res)
	if cuentaContraCredito != nil {
		sumValorMovimientoAPropiacion(false, cuentaContraCredito["Codigo"].(string), UnidadEjecutora, vigencia, afectacion["Valor"].(float64), res)
	}

}
