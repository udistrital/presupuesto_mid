package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/presupuesto_mid/helpers/rubroHelper"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/resposeformat"
)

// RubroController operations for Rubro
type RubroController struct {
	beego.Controller
}

// RegistrarRubro ...
// @Title RegistrarRubro
// @Description Registra Rubro en postgres y mongo
// @Param       body            body    models.Rubro    true            "body for Rubro content"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /RegistrarRubro/ [post]
func (c *RubroController) RegistrarRubro() {
	var v models.RubroRubro
	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0458", 500)
		}
	}()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		response := rubroHelper.AddRubro(v)
		beego.Debug(response)
		resposeformat.SetResponseFormat(&c.Controller, response["Body"], response["Code"].(string), 200)
	} else {
		resposeformat.SetResponseFormat(&c.Controller, nil, "E_0458", 500)
	}
}

// EliminarRubro ...
// @Title EliminarRubro
// @Description delete the Rubro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /EliminarRubro/:id [delete]
func (c *RubroController) EliminarRubro() {

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0458", 500)
		}
	}()

	if err != nil {
		beego.Error(err.Error())
		panic(err.Error())
	}

	response := rubroHelper.DeleteRubro(id)
	beego.Debug("response", response)
	resposeformat.SetResponseFormat(&c.Controller, response["Body"], response["Code"].(string), 200)

}

// ArbolRubros ...
// @Title ArbolRubros
// @Description Get Arbol Rubros By UE
// @Param	unidadEjecutora		path 	int64	true		"unidad ejecutora a consultar"
// @Param	rama		query 	string	false		"rama a consultar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /ArbolRubros/:unidadEjecutora [get]
func (c *RubroController) ArbolRubros() {

	try.This(func() {
		ueStr := c.Ctx.Input.Param(":unidadEjecutora")
		rama := c.GetString("rama")
		urlmongo := ""
		var res []map[string]interface{}
		if rama == "" {
			urlmongo = beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/RaicesArbol/" + ueStr
		} else {
			urlmongo = beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/ArbolRubro/" + rama + "/" + ueStr
		}
		beego.Info("Url ", urlmongo)
		if err := request.GetJson(urlmongo, &res); err != nil {
			beego.Error(err.Error())
			panic("Mongo API Service Error")
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}
