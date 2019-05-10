package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
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
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		try.This(func() {
			res := make(map[string]interface{})
			rubroData := v.(map[string]interface{})
			if rubroData["RubroPadre"] != nil { //Si se registra Un rubro con padre
				urlcrud := "http://" + beego.AppConfig.String("crudService")+"rubro_rubro"
				if err := request.SendJson(urlcrud, "POST", &res, &rubroData); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.
					//Se debe comprobar si se pudo registrar el rubro y la relacion rubro_rubro en postgres.
					if res["Type"] != nil && res["Type"].(string) == "success" {
						urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro/registrarRubro"
						var data map[string]interface{}
						sendData := res["Body"].(map[string]interface{})
						err := request.SendJson(urlmongo, "POST", &data, &sendData)
						beego.Info("data: ", sendData)
						if data["Type"] != nil && err == nil {
							if data["Type"].(string) == "error" {
								resul := res["Body"].(map[string]interface{})
								ue := resul["RubroHijo"].(map[string]interface{})["UnidadEjecutora"].(float64)
								urlcrud = urlcrud + "/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64))) + "/" + strconv.Itoa(int(ue))
								if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
									beego.Info("Data ", data)
									panic("Mongo API Error")
								} else {
									beego.Info("Error delete ", errorDelete)
									panic("Delete API Error")
								}

							} else if data["Type"].(string) == "success" {
								beego.Info("Registrado en Mongo")
							}
						} else {
							resul := res["Body"].(map[string]interface{})
							ue := resul["RubroHijo"].(map[string]interface{})["UnidadEjecutora"].(float64)
							urlcrud = urlcrud + "/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64))) + "/" + strconv.Itoa(int(ue))
							if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
								beego.Info("Data ", data)
								panic("Mongo API not Found")
							} else {
								beego.Info("Error delete ", errorDelete)
								panic("Delete API Error")
							}

						}
					} else if res["Type"] == nil {
						panic("Financiera Crud Service Error")
					}

					c.Data["json"] = res
				} else {
					panic("Financiera Crud Service Error")
				}
			} else if rubroData["RubroHijo"] != nil { //Si se registra un rubro Padre
				rubro := rubroData["RubroHijo"]
				urlcrud := "http://" + beego.AppConfig.String("crudService")+"rubro"
				if err := request.SendJson(urlcrud, "POST", &res, &rubro); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.
					//Se debe comprobar si se pudo registrar el rubro en postgres.
					if res["Type"] != nil && res["Type"].(string) == "success" {
						urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/registrarRubro"
						var data map[string]interface{}
						res["Body"] = map[string]interface{}{"RubroHijo": res["Body"].(map[string]interface{}), "RubroPadre": map[string]interface{}{}}
						rubroRow := res["Body"].(map[string]interface{})
						errorPost := request.SendJson(urlmongo, "POST", &data, &rubroRow)
						beego.Info("data: ", urlmongo)
						if data["Type"] != nil && errorPost == nil {
							if data["Type"].(string) == "error" {
								beego.Info("Error en mongo")
								resul := res["Body"].(map[string]interface{})["RubroHijo"].(map[string]interface{})
								beego.Info("Send Data: ", resul)
								urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
								if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
									beego.Info("Data ", data)
									panic("Mongo API Error")
								} else {
									beego.Info("Error ", errorDelete)
									panic("delete API Error")
								}

							} else if data["Type"].(string) == "success" {
								beego.Info("Registrado en Mongo")
							}
						} else {
							resul := res["Body"].(map[string]interface{})["RubroHijo"].(map[string]interface{})
							urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
							if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
								beego.Info("Data ", data)
								panic("Mongo API not Found")
							} else {
								beego.Info("Error ", errorDelete)
								panic("delete API Error")
							}

						}

					} else if res["Type"] == nil {
						panic("Financiera Crud Service Error")
					}

					c.Data["json"] = res
				} else {
					panic("Service Error")
				}
			} else {
				panic("Data Undefined")
			}

		}).Catch(func(e try.E) {
			// Print crash
			fmt.Println("expc ", e)
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
		})
	} else {
		c.Data["json"] = nil
	}
	c.ServeJSON()
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
		urlcrud := "http://" + beego.AppConfig.String("crudService")+"rubro/" + idStr
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
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/RaicesArbol/" + ueStr
		} else {
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/ArbolRubro/" + rama + "/" + ueStr
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
