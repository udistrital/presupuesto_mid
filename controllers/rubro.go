package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/presupuesto_mid/helpers/rubroHelper"
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
	var v interface{}
	defer func() {
		if r := recover(); r != nil {
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0458", 500)
		}
	}()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		res := make(map[string]interface{})
		rubroData := v.(map[string]interface{})
		rubroHelper.AddRubro(rubroData)
	} else {
		resposeformat.SetResponseFormat(&c.Controller, nil, "E_0458", 404)
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
	try.This(func() {
		idStr := c.Ctx.Input.Param(":id")
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro/" + idStr
		var res map[string]interface{}
		if err := request.SendJson(urlcrud, "DELETE", &res, nil); err == nil {
			if res["Type"].(string) == "success" {
				var resMg map[string]interface{}
				urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/eliminarRubro/" + idStr
				if err = request.SendJson(urlmongo, "DELETE", &resMg, nil); err != nil {
					fmt.Println("err ", err)
					panic("Mongo Not Found")
				} else if resMg["Type"].(string) == "error" {
					panic("Mongo CRUD Service Error")
				}
			} else if res["Type"].(string) == "error" {
				c.Data["json"] = res
			} else {
				panic("Financiera CRUD Service Error")
			}
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
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
			beego.Info(err.Error())
			panic("Mongo API Service Error")
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}
