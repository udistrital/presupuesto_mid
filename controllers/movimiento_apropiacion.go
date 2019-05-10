package controllers

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// MovimientoApropiacionController operations for MovimientoApropiacion
type MovimientoApropiacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *MovimientoApropiacionController) URLMapping() {

}

const separator = "-"

// AprobarMovimietnoApropiacion ...
// @Title AprobarMovimietnoApropiacion
// @Description create MovimientoApropiacion
// @Param	body		body 	models.MovimientoApropiacion	true		"body for MovimientoApropiacion content"
// @Success 201 {int} models.MovimientoApropiacion
// @Failure 403 body is empty
// @router /AprobarMovimietnoApropiacion/:unidadEjecutora/:vigencia [post]
func (c *MovimientoApropiacionController) AprobarMovimietnoApropiacion() {
	var v map[string]interface{}
	try.This(func() {

		var res interface{}
		var afectacion []map[string]interface{}
		UEStr := c.Ctx.Input.Param(":unidadEjecutora")
		UE, comp := strconv.Atoi(UEStr)
		if comp != nil {
			panic(comp.Error())
		}
		VigStr := c.Ctx.Input.Param(":vigencia")
		vigencia, comp := strconv.Atoi(VigStr)
		if comp != nil {
			panic(comp.Error())
		}
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			formatdata.FillStructP(v["MovimientoApropiacionDisponibilidadApropiacion"], &afectacion)
			for index := 0; index < len(afectacion); index++ {
				afectacion[index]["CuentaCredito"].(map[string]interface{})["Codigo"] = afectacion[index]["CuentaCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["Codigo"]
				afectacion[index]["CuentaCredito"].(map[string]interface{})["UnidadEjecutora"] = strconv.Itoa(int(afectacion[index]["CuentaCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64)))
				if afectacion[index]["CuentaContraCredito"] != nil {
					afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Codigo"] = afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["Codigo"]
					afectacion[index]["CuentaContraCredito"].(map[string]interface{})["UnidadEjecutora"] = strconv.Itoa(int(afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64)))
				}
			}
			_, _, compr := ComprobacionMovimiento(afectacion, UE, vigencia)
			if compr {
				Urlcrud := "http://" + beego.AppConfig.String("crudService")+"movimiento_apropiacion/AprobarMovimietnoApropiacion"
				if err := request.SendJson(Urlcrud, "POST", &res, &v); err == nil {
					//beego.Info("Data to Send ", res)
					c.Data["json"] = res
				} else {
					panic(err.Error())
				}
			} else {
				var alert []models.Alert
				alt := models.Alert{}
				alt.Code = "E_MODP008"
				alt.Body = map[string]interface{}{"Movimiento": v}
				alt.Type = "error"
				alert = append(alert, alt)
				c.Data["json"] = alert
			}

		} else {
			panic(err.Error())
		}
	}).Catch(func(e try.E) {
		beego.Error("catch error registrar valores: ", e)
		var alert []models.Alert
		alt := models.Alert{}
		alt.Code = "E_0458"
		alt.Body = v
		alt.Type = "error"
		alert = append(alert, alt)
		c.Data["json"] = alert
	})
	c.ServeJSON()
}

// ComprobarMovimientoApropiacion ...
// @Title Comprobar Movimiento Apropiacion
// @Description Comprueba si se puede generar un Movimiento en las apropiaciones.
// @Param	body		body 	map[string]string	true		"body for MovimientoApropiacion content"
// @Param	format		query 	string	false		"rama a consultar"
// @Success 200 {object} map[string]string
// @Failure 403
// @router /ComprobarMovimientoApropiacion/:unidadEjecutora/:vigencia [post]
func (c *MovimientoApropiacionController) ComprobarMovimientoApropiacion() {
	try.This(func() {
		var v map[string]interface{}
		var afectacion []map[string]interface{}
		dataSend := make(map[string]interface{})
		UEStr := c.Ctx.Input.Param(":unidadEjecutora")
		UE, comp := strconv.Atoi(UEStr)
		if comp != nil {
			panic(comp.Error())
		}
		format, _ := c.GetInt("format")

		VigStr := c.Ctx.Input.Param(":vigencia")
		vigencia, comp := strconv.Atoi(VigStr)
		if comp != nil {
			panic(comp.Error())
		}
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			formatdata.FillStructP(v["MovimientoApropiacionDisponibilidadApropiacion"], &afectacion)
			if format == 1 {
				for index := 0; index < len(afectacion); index++ {
					afectacion[index]["CuentaCredito"].(map[string]interface{})["Codigo"] = afectacion[index]["CuentaCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["Codigo"]
					afectacion[index]["CuentaCredito"].(map[string]interface{})["UnidadEjecutora"] = strconv.Itoa(int(afectacion[index]["CuentaCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64)))
					if afectacion[index]["CuentaContraCredito"] != nil {
						afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Codigo"] = afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["Codigo"]
						afectacion[index]["CuentaContraCredito"].(map[string]interface{})["UnidadEjecutora"] = strconv.Itoa(int(afectacion[index]["CuentaContraCredito"].(map[string]interface{})["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64)))
					}
				}
			}
			dataSend["Saldo"], dataSend["Diff"], dataSend["Comp"] = ComprobacionMovimiento(afectacion, UE, vigencia)
			alt := models.Alert{}
			alt.Code = "S_CMA001"
			alt.Body = dataSend
			alt.Type = "success"
			c.Data["json"] = alt
		} else {
			panic(err.Error())
		}
	}).Catch(func(e try.E) {
		beego.Error("catch error Comprobar Movimientos: ", e)
		alt := models.Alert{}
		alt.Code = "E_0458"
		alt.Body = e
		alt.Type = "error"
		c.Data["json"] = alt
	})
	c.ServeJSON()
}

// CalcularAfectacionMovimientoApropiacion ... Calcula la afectacion de un movimiento en el arbol
// Antes de realizar la operacion de registro en la db.
func CalcularAfectacionMovimientoApropiacion(afectacion map[string]interface{}, res map[string]float64) {

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
		saldoObj = CalcularSaldoApropiacion(codigoPadre[0], unidadEjecutora, vigencia)
		valorInicial = saldoObj["valor_inicial"]
	} else {
		valorInicial = 0
	}
	valorFinal = valorInicial + valorMov
	res[codigoPadre[0]] = res[codigoPadre[0]] + valorFinal
	return
}

// Mongo function's
// AddMovimientoApropiacionMongo ...
func AddMovimientoApropiacionMongo(parameter ...interface{}) (err interface{}) {
	idMov := 0.0
	var movimientos []map[string]interface{}
	var resM map[string]interface{}
	try.This(func() {
		dataMongo := make(map[string]interface{})
		infoMovimiento := parameter[0].(map[string]interface{})["Movimiento"].(map[string]interface{})
		idMov = infoMovimiento["Id"].(float64)
		var afectacionArr []map[string]interface{}

		dataMongo["FechaMovimiento"] = infoMovimiento["FechaMovimiento"]
		dataMongo["Vigencia"] = infoMovimiento["Vigencia"]
		dataMongo["UnidadEjecutora"] = infoMovimiento["UnidadEjecutora"]
		dataMongo["Id"] = infoMovimiento["Id"]
		err1 := formatdata.FillStruct(infoMovimiento["MovimientoApropiacionDisponibilidadApropiacion"], &movimientos)
		if err1 != nil {
			panic(err1.Error())
		}
		for _, data := range movimientos {
			var CuentaContraCredito string
			var CuentaCredito string
			var Disponibilidad float64
			var Apropiacion float64
			if CuentaContraCreditoInt, e := data["CuentaContraCredito"].(map[string]interface{}); e {
				CuentaContraCredito = CuentaContraCreditoInt["Rubro"].(map[string]interface{})["Codigo"].(string)
			}
			if CuentaCreditoInt, e := data["CuentaCredito"].(map[string]interface{}); e {
				CuentaCredito = CuentaCreditoInt["Rubro"].(map[string]interface{})["Codigo"].(string)
				Apropiacion = CuentaCreditoInt["Id"].(float64)
			}
			if dispo, e := data["Disponibilidad"].(map[string]interface{}); e {
				Disponibilidad = dispo["Id"].(float64)
				var disponibilidadParameter []interface{}
				disponibilidadParameter = append(disponibilidadParameter, dispo)
				beego.Info("Jon Info ", disponibilidadParameter)
				if err := AddDisponibilidadMongo(disponibilidadParameter...); err != nil {
					panic("Mongo api error")
				}
			}

			Valor := data["Valor"]
			beego.Info("data send ", data)

			TipoMovimiento := data["TipoMovimientoApropiacion"].(map[string]interface{})["Nombre"]
			afectacion := map[string]interface{}{
				"CuentaContraCredito": CuentaContraCredito,
				"CuentaCredito":       CuentaCredito,
				"Valor":               Valor,
				"TipoMovimiento":      TipoMovimiento,
				"Disponibilidad":      Disponibilidad,
				"Apropiacion":         Apropiacion,
			}
			afectacionArr = append(afectacionArr, afectacion)
		}
		dataMongo["Afectacion"] = afectacionArr
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/ModificacionApr"
		if err1 := request.SendJson(Urlmongo, "POST", &resM, &dataMongo); err1 == nil {
			if resM["Type"].(string) == "success" {
				err = err1
			} else {
				panic("Mongo api error")
			}
		} else {
			panic("Mongo Not Found")
		}
	}).Catch(func(e try.E) {
		var resC interface{}
		infoMovimiento := parameter[0].(map[string]interface{})["Movimiento"].(map[string]interface{})
		Urlcrud := "http://" + beego.AppConfig.String("crudService")+"movimiento_apropiacion/" + strconv.Itoa(int(idMov))
		estadoMov := infoMovimiento["EstadoMovimientoApropiacion"].(map[string]interface{})
		estadoMov["Id"] = 1
		infoMovimiento["EstadoMovimientoApropiacion"] = estadoMov

		if errorPut := request.SendJson(Urlcrud, "PUT", &resC, &infoMovimiento); errorPut == nil {
			for _, data := range movimientos {
				if dispo, e := data["Disponibilidad"].(map[string]interface{}); e {
					idDisp := dispo["Id"].(float64)
					Urlcrud := "http://" + beego.AppConfig.String("crudService")+"disponibilidad/DeleteDisponibilidadMovimiento/" + strconv.Itoa(int(idDisp))
					if errorDelete := request.SendJson(Urlcrud, "DELETE", &resC, nil); errorDelete == nil {
						beego.Info(resC)
					} else {
						beego.Info(errorDelete)
					}

				}

			}
			beego.Error("error job ", e)
		} else {
			beego.Error("error en put ", errorPut)
		}

	})
	return
}

func ComprobacionMovimiento(afectacion []map[string]interface{}, UE, vigencia int) (map[string]float64, float64, bool) {
	res := make(map[string]float64)
	for _, element := range afectacion {
		CalcularAfectacionMovimientoApropiacion(element, res)
	}
	sumValorMovimientoAPropiacion(true, "3", UE, vigencia, 0, res)
	sumValorMovimientoAPropiacion(true, "2", UE, vigencia, 0, res)
	diff := math.Abs(res["2"] - res["3"])
	if res["2"] != res["3"] {
		return res, diff, false
	}
	return res, diff, true
}
