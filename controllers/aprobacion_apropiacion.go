package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/golog"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/ruler"
)

// AprobacionController operations for AprobacionController
type AprobacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *AprobacionController) URLMapping() {
	c.Mapping("Aprobar", c.Aprobar)
}

// Aprobar ...
// @Title Aprobar
// @Description Aprobar Apropiacion
// @Param	body		body 	[]models.Apropiacion	true		"body for Apropiacion  content"
// @Success 201 {int} models.InfoSolDisp
// @Failure 403 body is empty
// @router Aprobar/ [post]
func (c *AprobacionController) Aprobar() {

	var predicados []models.Predicado
	//var predicados_apropiacion []models.Predicado
	var alertas []string
	alertas = append(alertas, "success")
	var postdominio string = ""

	if tdominio := c.GetString("tdominio"); tdominio != "" {
		postdominio = postdominio + "&query=Dominio.Id:" + tdominio
	} else {
		c.Data["json"] = "no se especifico el domino del ruler"
		c.ServeJSON()
	}
	var apropiacion []models.Apropiacion
	//var respuesta interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &apropiacion); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0"+postdominio, &predicados); err == nil {
			//var reglas string = ""
			var reglasbase string = ""
			//var reglasinyectadas string = ""
			var arregloReglas = make([]string, len(predicados))
			var comprobacion string
			var estadoAp int
			//var respuesta []models.FormatoPreliqu
			for i := 0; i < len(predicados); i++ {
				arregloReglas[i] = predicados[i].Nombre
			}

			for i := 0; i < len(arregloReglas); i++ {
				reglasbase = reglasbase + arregloReglas[i]
			}
			for i := len(apropiacion) - 1; i >= 0; i-- {

				comprobacion = comprobarApropiacion(apropiacion[i])
				if comprobacion == "" {
					alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" No aprobada, algunas apropiaciones hijo no se encuentran aprobadas")
				} else {
					estadoAp, _ = strconv.Atoi(golog.ComprobarAprobacion(reglasbase, comprobacion))
					estado := models.EstadoApropiacion{Id: estadoAp}
					apropiacion[i].Estado = &estado
					var respuesta interface{}
					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/"+strconv.Itoa(apropiacion[i].Id), "PUT", &respuesta, &apropiacion[i]); err == nil {
						if estadoAp == 1 {
							alertas[0] = "error"
							alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" No aprobada, la suma de las apropiaciones hijo no corresponde al valor de esta apropiacion")
						} else {
							alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" Aprobada")
						}
					} else {
						alertas[0] = "error"
						alertas = append(alertas, "no se pudo cambiar el estado de la apropiacion")
						c.Data["json"] = alertas
						c.ServeJSON()
					}
				}
			}
			c.Data["json"] = alertas
			c.ServeJSON()
		} else {
			fmt.Println("err: ", err)
			alertas[0] = "error"
			alertas = append(alertas, "no se pudo cargar la informacion de la base de datos")
			c.Data["json"] = alertas
			c.ServeJSON()
		}
	} else {
		fmt.Println("err: ", err)
		alertas[0] = "error"
		alertas = append(alertas, "no se resivieron los datos correctamente")
		c.Data["json"] = alertas
		c.ServeJSON()
	}

}

func comprobarApropiacion(padre models.Apropiacion) string {
	var rubroHijo []models.RubroRubro
	var listaValores []string
	var regla string
	var apropiacionHijo []models.Apropiacion
	var hoja int
	err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_rubro?limit=0&query=RubroPadre.Id:"+strconv.Itoa(padre.Rubro.Id)+",RubroPadre.Vigencia:"+strconv.FormatFloat(padre.Vigencia, 'f', -1, 64), &rubroHijo)
	if rubroHijo != nil && err == nil {
		for i := 0; i < len(rubroHijo); i++ {

			err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion?limit=0&query=Vigencia:"+strconv.FormatFloat(padre.Vigencia, 'f', -1, 64)+",Rubro.Id:"+strconv.Itoa(rubroHijo[i].RubroHijo.Id)+"", &apropiacionHijo)
			if apropiacionHijo != nil && err == nil {
				hoja = 0
				for i := 0; i < len(apropiacionHijo); i++ {
					if apropiacionHijo[i].Estado.Id == 2 {
						listaValores = append(listaValores, strconv.FormatFloat(apropiacionHijo[i].Valor, 'f', -1, 64))
						fmt.Println("apro: ", listaValores)
					}

				}

			}
		}
	} else {
		listaValores = append(listaValores, strconv.FormatFloat(padre.Valor, 'f', -1, 64))
		hoja = 1
	}
	fmt.Println("hijo: ", len(apropiacionHijo))
	fmt.Println("valor: ", len(listaValores))
	if listaValores != nil && len(apropiacionHijo) > 0 {
		for i := 0; i < len(listaValores); i++ {
			if len(listaValores) == 1 {
				regla = "verifica_hijos([" + listaValores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
			} else if len(listaValores) > 1 {
				if i == 0 {
					regla = "verifica_hijos([" + listaValores[i]
				} else if i == (len(listaValores) - 1) {
					regla = regla + "," + listaValores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
				} else {
					regla = regla + "," + listaValores[i]
				}
			}
		}
	}
	if hoja == 1 {
		for i := 0; i < len(listaValores); i++ {
			if len(listaValores) == 1 {
				regla = "verifica_hijos([" + listaValores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
			} else if len(listaValores) > 1 {
				if i == 0 {
					regla = "verifica_hijos([" + listaValores[i]
				} else if i == (len(listaValores) - 1) {
					regla = regla + "," + listaValores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
				} else {
					regla = regla + "," + listaValores[i]
				}
			}
		}
	}

	return regla
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
	vigencia, err := c.GetInt("Vigencia")
	if err == nil {
		unidadejecutora, err := c.GetInt("UnidadEjecutora")
		if err == nil {
			fmt.Println(vigencia)
			fmt.Println(unidadejecutora)
			tool := new(ruler.EntornoReglas)
			tool.Agregar_dominio("Presupuesto")
			var res []string
			var infoSaldoInicial []map[string]interface{}
			//saldo := make(map[string]interface{})
			if errorFill := formatdata.FillStruct(tool.Ejecutar_all_result("codigo_rubro_comprobacion_inicial(Y).", "Y"), &res); errorFill == nil {
				for _, rpadre := range res {
					var rubro []map[string]interface{}
					urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/ArbolApropiacion/" + rpadre + "/" + strconv.Itoa(unidadejecutora) + "/" + strconv.Itoa(vigencia)
					//if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro?query=Codigo:"+rpadre, &rubro); err == nil {
					if err = request.GetJson(urlmongo, &rubro); err == nil {
						beego.Info("Rubro ", rubro[0])
						if rubro[0]["Id"] != nil {
							infoSaldoInicial = append(infoSaldoInicial, rubro[0])

							/*if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/SaldoApropiacionPadre/"+strconv.Itoa(rubro[0].Id)+"?Vigencia="+strconv.Itoa(vigencia)+"&UnidadEjecutora="+strconv.Itoa(unidadejecutora), &saldo); err == nil {
								if saldo != nil {
									//infoSaldoInicial = append(infoSaldoInicial, map[string]interface{}{"Id": rubro[0].Id, "Codigo": rpadre, "Nombre": rubro[0].Nombre, "SaldoInicial": saldo["original"]})
								}
							} else {
								fmt.Println(err)
							}*/
						}

					} else {
						fmt.Println(err)
						c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}

				}
			}
			//c.Data["json"] = map[string]interface{}{"Aprobado": "0", "Data": infoSaldoInicial}
			for _, apr := range infoSaldoInicial {
				tool.Agregar_predicado("valor_inicial_rubro(" + fmt.Sprintf("%v", apr["Codigo"]) + "," + fmt.Sprintf("%v", apr["ApropiacionInicial"]) + ").")
			}
			if infoSaldoInicial != nil {
				res := tool.Ejecutar_result("comprobacion_inicial_apropiacion("+fmt.Sprintf("%v", infoSaldoInicial[0]["ApropiacionInicial"])+",Y).", "Y")
				var comp string
				err = formatdata.FillStruct(res, &comp)
				if err == nil {
					c.Data["json"] = map[string]interface{}{"Aprobado": res, "Data": infoSaldoInicial}
				} else {
					fmt.Println("nil2")

					c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}
			}

		} else {
			fmt.Println(err)
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		fmt.Println(err)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}
	c.ServeJSON()
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
	var v []map[string]interface{}
	vigencia, err := c.GetInt("Vigencia")
	if err == nil {
		unidadejecutora, err := c.GetInt("UnidadEjecutora")
		if err == nil {
			if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
				tool := new(ruler.EntornoReglas)
				tool.Agregar_dominio("Presupuesto")
				for _, apr := range v {
					tool.Agregar_predicado("valor_inicial_rubro(" + fmt.Sprintf("%v", apr["Codigo"]) + "," + fmt.Sprintf("%v", apr["SaldoInicial"]) + ").")
				}
				if v != nil {
					res := tool.Ejecutar_result("comprobacion_inicial_apropiacion("+fmt.Sprintf("%v", v[0]["SaldoInicial"])+",Y).", "Y")
					var aprobado string
					err = formatdata.FillStruct(res, &aprobado)
					if err == nil {
						if aprobado == "1" {
							var res interface{}
							if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/AprobacionAsignacionInicial"+"?Vigencia="+strconv.Itoa(vigencia)+"&UnidadEjecutora="+strconv.Itoa(unidadejecutora), &res); err == nil {
								c.Data["json"] = res

							} else {
								c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
							}

						} else {
							c.Data["json"] = models.Alert{Code: "E_AP003", Body: v, Type: "error"}
						}

					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}

				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
				}
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}

	c.ServeJSON()
}
