package controllers

import (
	"github.com/udistrital/presupuesto_mid/helpers"
	"github.com/udistrital/presupuesto_mid/helpers/apropiacionHelper"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/resposeformat"
)

// AprobacionController operations for AprobacionController
type AprobacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *AprobacionController) URLMapping() {
}

// InformacionAsignacionInicial ...
// @Title InformacionAsignacionInicial
// @Description Devuelve saldos iniciales antes de aprobar
// @Param	Vigencia		query 	string	true		"vigencia a comprobar"
// @Param	UnidadEjecutora		query 	string	true		"unidad ejecutora de los rubros a comprobar"
// @Success 200 {string} resultado
// @Failure 403
// @router /InformacionAsignacionInicial/ [get]
func (c *AprobacionController) InformacionAsignacionInicial() {

	asignationInfo := map[string]float64{"2": 0.0, "3": 0.0}

	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0459", 500)
		}
	}()
	vigencia, err := c.GetInt("Vigencia")
	if err != nil {
		panic(helpers.InternalErrorMessage())
	}
	unidadejecutora, err := c.GetInt("UnidadEjecutora")
	if err != nil {
		panic(helpers.InternalErrorMessage())
	}

	compareFlag := apropiacionHelper.CompareApropiationNodes(&asignationInfo, unidadejecutora, vigencia)

	beego.Debug(compareFlag, asignationInfo)
	response := make(map[string]interface{})
	response["InfoApropiacion"] = asignationInfo
	response["Aprobado"] = compareFlag

	resposeformat.SetResponseFormat(&c.Controller, response, "", 200)
}

// AprobacionAsignacionInicial ...
// @Title AprobacionAsignacionInicial
// @Description aprueba la asignacion inicial de presupuesto
// @Param	Vigencia		query 	string	true		"vigencia a comprobar"
// @Param	UnidadEjecutora		query 	string	true		"unidad ejecutora de los rubros a comprobar"
// @Success 200 {string} resultado
// @Failure 403
// @router /AprobacionAsignacionInicial/ [post]
func (c *AprobacionController) AprobacionAsignacionInicial() {

	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0459", 500)
		}
	}()

	vigencia, err := c.GetInt("Vigencia")

	if err != nil {
		panic(helpers.InternalErrorMessage())
	}

	unidadejecutora, err := c.GetInt("UnidadEjecutora")

	if err != nil {
		panic(helpers.InternalErrorMessage())
	}

	response := apropiacionHelper.AprobarPresupuesto(vigencia, unidadejecutora)
	resposeformat.SetResponseFormat(&c.Controller, response, "", 200)
}

// Aprobado ...
// @Title Aprobado
// @Description aprueba la asignacion inicial de presupuesto
// @Param	Vigencia		query 	string	true		"vigencia a comprobar"
// @Param	UnidadEjecutora		query 	string	true		"unidad ejecutora de los rubros a comprobar"
// @Success 200 {string} resultado
// @Failure 403
// @router /Aprobado [get]
func (c *AprobacionController) Aprobado() {

	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0459", 500)
		}
	}()

	vigencia, err := c.GetInt("Vigencia")

	if err != nil {
		panic(helpers.InternalErrorMessage())
	}

	unidadejecutora, err := c.GetInt("UnidadEjecutora")

	if err != nil {
		panic(helpers.InternalErrorMessage())
	}

	response := apropiacionHelper.PresupuestoAprobado(vigencia, unidadejecutora)
	resposeformat.SetResponseFormat(&c.Controller, response, "", 200)
}
