package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	. "github.com/mndrix/golog"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/ruler"
)

// AprobacionFuenteController operations for AprobacionFuenteController
type AprobacionFuenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *AprobacionFuenteController) URLMapping() {
	c.Mapping("AprobacionFuente", c.ValorMovimientoFuente)
}

//http://localhost:8080/v1/aprobacion_fuente/ValorMovimientoFuenteTraslado?idfuente=37&idapropiacion=256&iddependencia=122&traslado=40000000
//http://localhost:8080/v1/aprobacion_fuente/ValorMovimientoFuente?idfuente=37&idapropiacion=256&iddependencia=122

// ValorMovimientoFuente ...
// @Title ValorMovimientoFuente
// @Description muestra el valor comprometido de una fuente especifica por dependencia y apropiacion
// @Param	idfuente	query	int	false	"id de la fuente a consultar"
// @Param	iddependencia	query	int	false	"id de la dependencia a consultar"
// @Param	idapropiacion	query	int	false	"id de la apropiacion a consultar"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ValorMovimientoFuente [get]
func (c *AprobacionFuenteController) ValorMovimientoFuente() {
	var res []interface{}
	var resfuente []interface{}
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		fmt.Println(idfuente)
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {

				var Movimiento []map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?query=FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente)+",FuenteFinanciamientoApropiacion.Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamientoApropiacion.Dependencia:"+strconv.Itoa(iddependencia), &Movimiento); err == nil {
					if Movimiento != nil {

						for _, Movimientos := range Movimiento {
							resfuente = append(resfuente, Movimientos)
						}
						if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
							beego.Error(valorGastado)
							if valorGastado != nil {
								for _, valores := range valorGastado {
									res = append(res, valores)
								}

								if resfuente != nil {
									if res != nil {
										var valor float64
										valor = 0
										var valorcdp float64
										valorcdp = 0
										valorcdp = res[0].(float64)
										for _, rowfuente := range resfuente {
											valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
										}

										var valordisponible float64
										valordisponible = 0
										valordisponible = valor - valorcdp

										c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "ValorDisponible": valordisponible, "ValorGastado": valorcdp, "ValorTotal": valor}
									}
								}
							} else {

								var valor float64
								valor = 0
								for _, rowfuente := range resfuente {
									valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
								}

								c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "ValorDisponible": valor, "ValorGastado": 0, "ValorTotal": valor}

							}
						} else {

							var valor float64
							valor = 0
							for _, rowfuente := range resfuente {
								valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
							}
							c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "ValorDisponible": valor, "ValorGastado": 0, "ValorTotal": valor}
						}
					} else {
						c.Data["json"] = nil
					}
				} else {
					fmt.Println("err4 ", err.Error())
					c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}
			} else {
				fmt.Println("err3 ", err.Error())
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			fmt.Println("err2 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		fmt.Println("err1 ", err.Error())
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}

	c.ServeJSON()
}

// ValorMovimientoFuenteTraslado ...
// @Title ValorMovimientoFuenteTraslado
// @Description retorna si se puede realizar un traslado o no en una fuente de financiamiento
// @Param	idfuente	query	int	false	"id de la fuente a consultar"
// @Param	iddependencia	query	int	false	"id de la dependencia a consultar"
// @Param	idapropiacion	query	int	false	"id de la apropiacion a consultar"
// @Param	valortraslado	query	int	false	"id de la apropiacion a consultar"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ValorMovimientoFuenteTraslado [get]
func (c *AprobacionFuenteController) ValorMovimientoFuenteTraslado() {
	var res []interface{}
	var resfuente []interface{}
	var predicados []ruler.Predicado
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {
				if valortraslado, err := c.GetFloat("traslado"); err == nil {

					var Movimiento []map[string]interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?query=FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente)+",FuenteFinanciamientoApropiacion.Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamientoApropiacion.Dependencia:"+strconv.Itoa(iddependencia), &Movimiento); err == nil {
						if Movimiento != nil {

							for _, Movimientos := range Movimiento {
								resfuente = append(resfuente, Movimientos)
							}

							//var valorGastado map[string]interface{}
							if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
								beego.Error("Valor Fuente CDP: ", valorGastado)
								if valorGastado != nil {
									for _, valores := range valorGastado {
										res = append(res, valores)
									}

									if resfuente != nil {
										if res != nil {
											var valor float64
											valor = 0
											var valorcdp float64
											valorcdp = 0
											valorcdp = res[0].(float64)

											for _, rowfuente := range resfuente {
												valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
											}
											//reglas
											reglasBase := ruler.CargarReglasBase("FuenteFinanciamiento")

											predicados = append(predicados, ruler.Predicado{Nombre: "movimientofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) + "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valor, 'f', -1, 64) + ")."})
											predicados = append(predicados, ruler.Predicado{Nombre: "saldofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) + "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valorcdp, 'f', -1, 64) + ")."})
											//valor que se va a transladar
											predicados = append(predicados, ruler.Predicado{Nombre: "saldofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) + "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valortraslado, 'f', -1, 64) + ")."})
											reglas := ruler.FormatoReglas(predicados) + reglasBase

											m := NewMachine().Consult(reglas)
											resultados := m.ProveAll("validacion_total_fuente_dependencia_apropiacion_saldo(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) + "," + strconv.Itoa(idfuente) + ",Y).")
											var resp string

											for _, solution := range resultados {
												resp = fmt.Sprintf("%s", solution.ByName_("Y"))
											}
											f, _ := strconv.ParseBool(resp)

											var trasladar bool
											trasladar = f

											c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "Trasladar": trasladar}
										}
									}
								} else {

									var valor float64
									valor = 0
									for _, rowfuente := range resfuente {
										valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
									}
									valor = valor - valortraslado

									var trasladar bool
									if valor >= 0 {
										trasladar = true
									} else {
										trasladar = false
									}

									c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "Trasladar": trasladar}

								}

							} else {
								var valor float64
								valor = 0
								for _, rowfuente := range resfuente {
									valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
								}
								valor = valor - valortraslado

								var trasladar bool
								if valor >= 0 {
									trasladar = true
								} else {
									trasladar = false
								}

								c.Data["json"] = map[string]interface{}{"Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "Trasladar": trasladar}
							}
						} else {
							c.Data["json"] = nil
						}
					} else {
						fmt.Println("err4 ", err.Error())
						c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}
				} else {
					fmt.Println("err3.5 ", err.Error())
					c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}
			} else {
				fmt.Println("err3 ", err.Error())
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			fmt.Println("err2 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		fmt.Println("err1 ", err.Error())
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}

	c.ServeJSON()
}

// ValorMovimientoFuenteLista ...
// @Title ValorMovimientoFuenteLista
// @Description devuelve una lista de las fuentes, dependecias y apropiaciones con los valores de la fuente (comprometido, disponible y el valor total)
// @Param	idfuente	query	int	false	"id de la fuente a consultar"
// @Param	iddependencia	query	int	false	"id de la dependencia a consultar"
// @Param	idapropiacion	query	int	false	"id de la apropiacion a consultar"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ValorMovimientoFuenteLista [get]
func (c *AprobacionFuenteController) ValorMovimientoFuenteLista() {
	var res []interface{}
	var resfuente []models.MovimientoFuenteFinanciamientoApropiacion

	var Movimiento []models.MovimientoFuenteFinanciamientoApropiacion

	if idfuente, err := c.GetInt("idfuente"); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?limit=-1&query=FuenteFinanciamientoApropiacion.Apropiacion.Rubro.Codigo__startswith:3-3-001-15-01-08-0119-,FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente), &Movimiento); err == nil {
			if Movimiento != nil {

				for _, Movimientos := range Movimiento {
					var idfuente = Movimientos.FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id
					var idapropiacion = Movimientos.FuenteFinanciamientoApropiacion.Apropiacion.Id
					var iddependencia = Movimientos.FuenteFinanciamientoApropiacion.Dependencia
					if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
						if valorGastado != nil {
							for _, valores := range valorGastado {
								res = append(res, valores)
							}

							if res != nil {

								var valorcdp float64
								valorcdp = 0
								valorcdp = res[0].(float64)

								Movimientos.ValorGastado = valorcdp
								Movimientos.ValorDisponible = Movimientos.Valor - valorcdp
							}

						} else {

							Movimientos.ValorGastado = 0
							Movimientos.ValorDisponible = Movimientos.Valor

						}

					} else {
						Movimientos.ValorGastado = 0
						Movimientos.ValorDisponible = Movimientos.Valor
					}

					resfuente = append(resfuente, Movimientos)

				}

				c.Data["json"] = resfuente

			} else {
				c.Data["json"] = nil
			}
		} else {
			fmt.Println("err4 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

	} else {

		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?limit=-1&query=FuenteFinanciamientoApropiacion.Apropiacion.Rubro.Codigo__startswith:3-3-001-15-01-08-0119-", &Movimiento); err == nil {
			if Movimiento != nil {

				for _, Movimientos := range Movimiento {
					var idfuente = Movimientos.FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id
					var idapropiacion = Movimientos.FuenteFinanciamientoApropiacion.Apropiacion.Id
					var iddependencia = Movimientos.FuenteFinanciamientoApropiacion.Dependencia
					if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
						if valorGastado != nil {
							for _, valores := range valorGastado {
								res = append(res, valores)
							}

							if res != nil {

								var valorcdp float64
								valorcdp = 0
								valorcdp = res[0].(float64)

								Movimientos.ValorGastado = valorcdp
								Movimientos.ValorDisponible = Movimientos.Valor - valorcdp
							}

						} else {

							Movimientos.ValorGastado = 0
							Movimientos.ValorDisponible = Movimientos.Valor

						}

					} else {
						Movimientos.ValorGastado = 0
						Movimientos.ValorDisponible = Movimientos.Valor
					}

					resfuente = append(resfuente, Movimientos)

				}

				c.Data["json"] = resfuente

			} else {
				c.Data["json"] = nil
			}
		} else {
			fmt.Println("err4 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	}

	c.ServeJSON()
}

// ValorMovimientoFuenteListaFunc ...
// @Title ValorMovimientoFuenteListaFunc
// @Description devuelve una lista de las fuentes, dependecias y apropiaciones con los valores de la fuente (comprometido, disponible y el valor total)
// @Param	idfuente	query	int	false	"id de la fuente a consultar"
// @Param	iddependencia	query	int	false	"id de la dependencia a consultar"
// @Param	idapropiacion	query	int	false	"id de la apropiacion a consultar"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ValorMovimientoFuenteListaFunc [get]
func (c *AprobacionFuenteController) ValorMovimientoFuenteListaFunc() {
	var res []interface{}
	var resfuente []models.MovimientoFuenteFinanciamientoApropiacion

	var Movimiento []models.MovimientoFuenteFinanciamientoApropiacion

	if idfuente, err := c.GetInt("idfuente"); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?limit=-1&query=FuenteFinanciamientoApropiacion.Apropiacion.Rubro.Codigo__startswith:3-1-,FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente), &Movimiento); err == nil {
			if Movimiento != nil {

				for _, Movimientos := range Movimiento {
					var idfuente = Movimientos.FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id
					var idapropiacion = Movimientos.FuenteFinanciamientoApropiacion.Apropiacion.Id
					var iddependencia = Movimientos.FuenteFinanciamientoApropiacion.Dependencia
					if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
						if valorGastado != nil {
							for _, valores := range valorGastado {
								res = append(res, valores)
							}

							if res != nil {

								var valorcdp float64
								valorcdp = 0
								valorcdp = res[0].(float64)

								Movimientos.ValorGastado = valorcdp
								Movimientos.ValorDisponible = Movimientos.Valor - valorcdp
							}

						} else {

							Movimientos.ValorGastado = 0
							Movimientos.ValorDisponible = Movimientos.Valor

						}

					} else {
						Movimientos.ValorGastado = 0
						Movimientos.ValorDisponible = Movimientos.Valor
					}

					resfuente = append(resfuente, Movimientos)

				}

				c.Data["json"] = resfuente

			} else {
				c.Data["json"] = nil
			}
		} else {
			fmt.Println("err4 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

	} else {

		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_fuente_financiamiento_apropiacion?limit=-1&query=FuenteFinanciamientoApropiacion.Apropiacion.Rubro.Codigo__startswith:3-1-", &Movimiento); err == nil {
			if Movimiento != nil {

				for _, Movimientos := range Movimiento {

					var idfuente = Movimientos.FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id
					var idapropiacion = Movimientos.FuenteFinanciamientoApropiacion.Apropiacion.Id
					var iddependencia = Movimientos.FuenteFinanciamientoApropiacion.Dependencia
					if err, valorGastado := ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion); err == nil {
						if valorGastado != nil {
							for _, valores := range valorGastado {
								res = append(res, valores)
							}

							if res != nil {

								var valorcdp float64
								valorcdp = 0
								valorcdp = res[0].(float64)

								Movimientos.ValorGastado = valorcdp
								Movimientos.ValorDisponible = Movimientos.Valor - valorcdp
							}

						} else {

							Movimientos.ValorGastado = 0
							Movimientos.ValorDisponible = Movimientos.Valor

						}

					} else {
						Movimientos.ValorGastado = 0
						Movimientos.ValorDisponible = Movimientos.Valor
					}

					resfuente = append(resfuente, Movimientos)

				}

				c.Data["json"] = resfuente

			} else {
				c.Data["json"] = nil
			}
		} else {
			fmt.Println("err4 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	}

	c.ServeJSON()
}
