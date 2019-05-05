package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
)

// ApropiacionController operations for  Apropiacion
type ApropiacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *ApropiacionController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description create Apropiacion
// @Param	body		body 	models.Apropiacion	true		"body for Apropiacion content"
// @Success 201 {int} models.Apropiacion
// @Failure 403 body is empty
// @router / [post]
func (c *ApropiacionController) Post() {
	var v map[string]interface{}
	var res map[string]interface{}
	var resM map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		try.This(func() {
			mongoData := v["Rubro"].(map[string]interface{})
			v["Rubro"] = map[string]interface{}{"Id": v["Rubro"].(map[string]interface{})["Id"]}
			urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/apropiacion/"
			if err = request.SendJson(urlcrud, "POST", &res, &v); err == nil {
				if res["Type"] != nil && res["Type"].(string) == "success" {
					mongoData["Id"] = res["Body"].(map[string]interface{})["Id"]
					mongoData["ApropiacionInicial"] = v["Valor"]
					urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarApropiacionInicial/" + strconv.Itoa(int(res["Body"].(map[string]interface{})["Vigencia"].(float64)))
					if err = request.SendJson(urlmongo, "POST", &resM, &mongoData); err == nil {
						if resM["Type"].(string) == "success" {
							body := res["Body"].(map[string]interface{})
							body["Rubro"] = mongoData
							res["Body"] = body
							c.Data["json"] = res
						} else {
							resul := res["Body"].(map[string]interface{})
							urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
							if err := request.SendJson(urlcrud, "DELETE", &resM, nil); err == nil {
								fmt.Println("borrado correctamente")
								panic("Mongo API error")
							}else{
								fmt.Println("Error al borrar")
							}

						}
					} else {
						resul := res["Body"].(map[string]interface{})
						urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
						if err := request.SendJson(urlcrud, "DELETE", &resM, nil); err == nil {
							fmt.Println("borrado correctamente")
							panic("Mongo API not Found")
						}else{
							fmt.Println("Error al borrar")
						}

					}
				} else {
					beego.Error(res)
					panic("Financiera CRUD not Found")
				}
			} else {
				panic("Financiera CRUD Service Error")
			}
		}).Catch(func(e try.E) {
			beego.Error("expc ", e)
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
		})
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description Update Apropiacion
// @Param	body		body 	models.Apropiacion	true		"body for Apropiacion content"
// @Success 201 {int} models.Apropiacion
// @Failure 403 body is empty
// @router /:id/:valor/:vigencia [put]
func (c *ApropiacionController) Put() {
	var v map[string]interface{}
	var res map[string]interface{}
	var resM map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		try.This(func() {
			idStr := c.Ctx.Input.Param(":id")
			valStr := c.Ctx.Input.Param(":valor")
			vigStr := c.Ctx.Input.Param(":vigencia")
			valorAnt := v["ApropiacionInicial"].(float64)
			valorAntStr := strconv.Itoa(int(valorAnt))
			mongoData := v
			urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/apropiacion/" + "/UpdateApropiacionValue/" + idStr + "/" + valStr
			if err = request.SendJson(urlcrud, "PUT", &res, nil); err == nil {
				if res["Type"] != nil && res["Type"].(string) == "success" {
					mongoData["ApropiacionInicial"], _ = strconv.Atoi(valStr)
					urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/RegistrarApropiacionInicial/" + vigStr
					if err = request.SendJson(urlmongo, "POST", &resM, &mongoData); err == nil {
						if resM["Type"].(string) == "success" {
							c.Data["json"] = res
						} else {
							resul := res["Body"].(map[string]interface{})
							urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
							if err := request.SendJson(urlcrud, "DELETE", &resM, nil);  err == nil {
								fmt.Println("borrado correctamente")
								panic("Mongo API error")
							}else{
								fmt.Println("Error al borrar")
							}

						}
					} else {
						urlcrud = "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/apropiacion/" + "/UpdateApropiacionValue/" + idStr + "/" + valorAntStr
						if err := request.SendJson(urlcrud, "PUT", &resM, nil);  err == nil {
							fmt.Println("actualizado correctamente")
							panic("Mongo API not Found")
						}else{
							fmt.Println("Error al actualizar")
						}

					}
				} else {
					beego.Error(res)
					panic("Financiera CRUD not Found")
				}
			} else {
				panic("Financiera CRUD Service Error")
			}
		}).Catch(func(e try.E) {
			beego.Error("expc ", e)
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
		})
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// ArbolApropiaciones ...
// @Title ArbolApropiaciones
// @Description Get Arbol Rubros By UE
// @Param	unidadEjecutora		path 	int64	true		"unidad ejecutora a consultar"
// @Param	rama		query 	string	false		"rama a consultar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /ArbolApropiaciones/:unidadEjecutora/:vigencia [get]
func (c *ApropiacionController) ArbolApropiaciones() {

	try.This(func() {
		ueStr := c.Ctx.Input.Param(":unidadEjecutora")
		vigenciaStr := c.Ctx.Input.Param(":vigencia")
		rama := c.GetString("rama")
		urlmongo := ""
		var res []map[string]interface{}
		if rama == "" {
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/RaicesArbolApropiacion/" + ueStr + "/" + vigenciaStr
		} else {
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/ArbolApropiacion/" + rama + "/" + ueStr + "/" + vigenciaStr
		}
		beego.Info("Url ", urlmongo)
		if err := request.GetJson(urlmongo, &res); err != nil {
			panic("Mongo API Service Error")
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}

// SaldoApropiacion ...
// @Title SaldoApropiacion
// @Description Get Arbol Rubros By UE
// @Param	unidadEjecutora		path 	int64	true		"unidad ejecutora a consultar"
// @Param	rama		query 	string	false		"rama a consultar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /SaldoApropiacion/:rubro/:unidadEjecutora/:vigencia [get]
func (c *ApropiacionController) SaldoApropiacion() {
	var (
		rubroParam    string
		unidadEParam  int
		vigenciaParam int
		err           error
	)
	try.This(func() {
		res := make(map[string]float64)
		rubroParam = c.GetString(":rubro")
		if unidadEParam, err = c.GetInt(":unidadEjecutora"); err != nil {
			panic(err.Error())
		}

		if vigenciaParam, err = c.GetInt(":vigencia"); err != nil {
			panic(err.Error())
		}
		res = CalcularSaldoApropiacion(rubroParam, unidadEParam, vigenciaParam)
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}

//CalcularSaldoApropiacion ....
func CalcularSaldoApropiacion(rubroParam string, unidadEParam, vigenciaParam int) (res map[string]float64) {
	var saldo float64
	urlmongo := ""
	urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/SaldoApropiacion/" + rubroParam + "/" + strconv.Itoa(unidadEParam) + "/" + strconv.Itoa(vigenciaParam)

	beego.Info("Url ", urlmongo)
	if err := request.GetJson(urlmongo, &res); err != nil {
		beego.Info(err.Error())
		panic("Mongo API Service Error")
	} else {
		for key, value := range res {

			if !strings.Contains(key, "mes") && !strings.Contains(key, "rp") {
				switch tipoMovimiento := key; tipoMovimiento {
				//rp
				case "Adicion", "Traslado_cuenta_contra_credito", "total_anulado_cdp", "valor_inicial":
					beego.Info("suma ", tipoMovimiento)
					saldo += value
				default:
					beego.Info("resta ", tipoMovimiento)
					saldo -= value
				}
			}
		}

	}

	res["saldo"] = saldo
	return
}
