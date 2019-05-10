package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/ruler"
)

// RegistroPresupuestalController operations for RegistroPresupuestal
type RegistroPresupuestalController struct {
	beego.Controller
}

// URLMapping ...
func (c *RegistroPresupuestalController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetSolicitudesRp", c.GetSolicitudesRp)
	c.Mapping("GetSolicitudesRpByID", c.GetSolicitudesRpByID)
}

func formatoSolicitudRP(solicitudintfc interface{}, params ...interface{}) (res interface{}) {
	//recuperar datos del CDP objetivo de la solicitud
	var rubros []interface{}
	solicitud := models.SolicitudRp{}
	err := formatdata.FillStruct(solicitudintfc, &solicitud)
	fmt.Println(err)
	var afectacionSolicitud []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"disponibilidad_apropiacion_solicitud_rp?limit=-1&query=SolicitudRp:"+strconv.Itoa(solicitud.Id), &afectacionSolicitud); err == nil {
		//consulta de la afectacion presupuestal objetivo.
		fmt.Println(solicitud.Id)
		for _, afect := range afectacionSolicitud {

			var dispAprSol []map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"disponibilidad_apropiacion?limit=1&query=Id:"+fmt.Sprintf("%v", afect["DisponibilidadApropiacion"]), &dispAprSol); err == nil {

				for _, disp_apro := range dispAprSol {
					solicitud.Cdp = int(disp_apro["Disponibilidad"].(map[string]interface{})["Id"].(float64))
					beego.Info("CDP ", solicitud.Cdp)
					disp_apro["ValorAsignado"] = afect["Monto"]
					disp_apro["FuenteFinanciacion"] = disp_apro["FuenteFinanciamiento"]
					rubros = append(rubros, disp_apro)
				}

			} else {
				//si sale mal la consulta de la afectacion del cdp objetivo.
				return
			}
		}
		solicitud.Rubros = rubros
	} else {
		//si sale mal la consulta de la afectacion de la solicitud.
		fmt.Println(err.Error())
		return
	}

	var cdpObjtvo []models.Disponibilidad
	if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.Cdp), &cdpObjtvo); err == nil {
		if cdpObjtvo != nil {
			solicitud.DatosDisponibilidad = &cdpObjtvo[0]
			var necesidadCdp []models.SolicitudDisponibilidad
			if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.DisponibilidadProcesoExterno[0].ProcesoExterno), &necesidadCdp); err == nil {
				if necesidadCdp != nil {
					solicitud.DatosDisponibilidad.DatosNecesidad = necesidadCdp[0].Necesidad
					var depNes []models.DependenciaNecesidad
					if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.DatosNecesidad.Id), &depNes); err == nil {
						if depNes != nil {
							var depSol []models.Dependencia
							var jefeDepSol []models.JefeDependencia
							if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol); err == nil {
								if jefeDepSol != nil {
									if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=1&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol); err == nil {
										if depSol != nil {
											solicitud.DatosDisponibilidad.DatosNecesidad.DatosDependenciaSolicitante = &depSol[0]
										} else {
											//si no hay datos de la dependencia
										}
									} else {
										//si hay error al consultar la dependecia solicitante
									}
								} else {
									//no hay datos jefe dep
								}

							} else {
								//jefe_dep
							}
						} else {
							//si no hay datos en la consulta dependencia_necesidad
						}
					} else {
						//si hay error al consultar dependencia_necesidad
					}
				} else {
					//si no hay datos de la necesidad
				}
			} else {
				//si hay error al encontrar datos de la necesidad
			}
		} else {
			//si no hay datos del CDP objetivo
			fmt.Println("error cdp: no hay datos, id : ", solicitud.Cdp)
		}
	} else {
		//si ocurre error al traer datos del CDP objetivo
		fmt.Println("error cdp: ", err)
	}
	//obtener informacion del contrato del rp
	//var infoContrato []models.ContratoGeneral
	/*var contratista []models.InformacionProveedor
	fmt.Println("prov ", solicitud.Proveedor)
	if solicitud.Proveedor != 0 {
		fmt.Println("http://" + beego.AppConfig.String("AdministrativaAmazonService") + "informacion_proveedor?limit=1&query=Id:" + strconv.Itoa(solicitud.Proveedor))
		if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_proveedor?limit=1&query=Id:"+strconv.Itoa(solicitud.Proveedor), &contratista); err == nil {
			solicitud.DatosProveedor = &contratista[0]
		} else {
			//error consulta proveedor
			fmt.Println(err.Error())
		}

	} else {
		//si no encuentra datos sobre el contrato
		fmt.Println("error contrato: no hay datos, id sol : ", solicitud.Id)
	}*/

	//cargar datos del compromiso de la solicitud de rp
	var compromisoRp []interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"compromiso?limit=1&query=Id:"+strconv.Itoa(solicitud.TipoCompromiso), &compromisoRp); err == nil {
		if compromisoRp != nil {
			solicitud.DatosCompromiso = compromisoRp[0]
		} else {
			//si no encuentra los datos del compromiso

		}
	} else {
		//si hay error al cargar el compromiso del rp
		fmt.Println(err.Error())
	}
	return solicitud
}

//funcion para recopilar datos externos de los rp a listar
func formatoListaRP(rpintfc interface{}, params ...interface{}) (res interface{}) {
	rp := rpintfc.(map[string]interface{})
	try.This(func() {
		idSolicitudDisponibilidad := int(rp["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{})[0].(map[string]interface{})["DisponibilidadApropiacion"].(map[string]interface{})["Disponibilidad"].(map[string]interface{})["DisponibilidadProcesoExterno"].([]interface{})[0].(map[string]interface{})["ProcesoExterno"].(float64))
		solicituddisp, err := DetalleSolicitudDisponibilidadByID(strconv.Itoa(idSolicitudDisponibilidad))

		if err == nil {
			rp["InfoSolicitudDisponibilidad"] = solicituddisp
		}
	}).Catch(func(e try.E) {
		// Print crash
		//fmt.Println("expc ",e)
	})
	return rp
}

// ListaRp ...
// @Title ListaRp
// @Description get RP by vigencia
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado para la lista de los cdp"
// @Success 200 {object} models.RegistroPresupuestal
// @Failure 403
// @router ListaRp/:vigencia [get]
func (c *RegistroPresupuestalController) ListaRp() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err1 := strconv.Atoi(vigenciaStr)
	var rpresupuestal []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var querybase string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if r := c.GetString("query"); r != "" {
		querybase = r

	}
	if startrange != "" && endrange != "" {
		query = querybase + ",FechaRegistro__gte:" + startrange + ",FechaRegistro__lte:" + endrange

	} else if querybase != "" {
		query = "," + querybase
	}
	UnidadEjecutora, err2 := c.GetInt("UnidadEjecutora")
	if err1 == nil && err2 == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"registro_presupuestal?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora:"+strconv.Itoa(UnidadEjecutora)+",Vigencia:"+strconv.Itoa(vigencia)+query, &rpresupuestal); err == nil {
			if rpresupuestal != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(rpresupuestal...)
				chrpresupuestal := optimize.Digest(done, formatoListaRP, resch, nil)
				for rp := range chrpresupuestal {
					if rp != nil {
						respuesta = append(respuesta, rp.(map[string]interface{}))
					}

				}
				c.Data["json"] = respuesta
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}

	c.ServeJSON()
}

// GetSolicitudesRp ...
// @Title GetSolicitudesRp
// @Description get saldo rp by apropiacion
// @Param	vigencia	query	string	false	"vigencia de las solicitudes a consultar"
// @Param	UnidadEjecutora	query	string	false	"unidad ejecutora de las solicitudes a consultar"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado para la lista de los cdp"
// @Success 200 {object} models.SolicitudRp
// @Failure 403 :vigencia is empty
// @router /GetSolicitudesRp/:vigencia [get]
func (c *RegistroPresupuestalController) GetSolicitudesRp() {

	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var querybase string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("query"); r != "" {
		querybase = r

	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if startrange != "" && endrange != "" {
		query = querybase + ",FechaSolicitud__gte:" + startrange + ",FechaSolicitud__lte:" + endrange

	} else if querybase != "" {
		query = "," + querybase
	}
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err1 := strconv.Atoi(vigenciaStr)
	UnidadEjecutora, err2 := c.GetInt("UnidadEjecutora")
	fmt.Println(UnidadEjecutora)
	if err1 == nil && err2 == nil {
		var solicitudesRp []interface{}
		var respuesta []models.SolicitudRp
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Masivo:false,Expedida:false"+query+",Vigencia:"+strconv.Itoa(vigencia)+"&sortby=Id&order=desc", &solicitudesRp); err == nil {
			if solicitudesRp != nil {
				//encontrar datos del CDP objetivo del RP Solicitado

				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(solicitudesRp...)
				chsolicitud := optimize.Digest(done, formatoSolicitudRP, resch, nil)
				for solicitud := range chsolicitud {
					if data, e := solicitud.(models.SolicitudRp); e {
						respuesta = append(respuesta, data)
					}
				}
				c.Data["json"] = respuesta
			} else {
				//si no hay datos de solicitudes
				c.Data["json"] = nil
			}
		} else {
			//si ocurre error al traer las solicitudes
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			fmt.Println(err.Error())
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

// GetSolicitudesRpByID ...
// @Title GetSolicitudesRpByID
// @Description get GetSolicitudesRpById by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SolicitudRp
// @Failure 403 :id is empty
// @router /GetSolicitudesRpById/:id [get]
func (c *RegistroPresupuestalController) GetSolicitudesRpByID() {
	var solicitudesRp []models.SolicitudRp
	var respuesta []models.SolicitudRp
	idStr := c.Ctx.Input.Param(":id")
	if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=1&query=Id:"+idStr, &solicitudesRp); err == nil {
		if solicitudesRp != nil {
			//encontrar datos del CDP objetivo del RP Solicitado
			for _, solicitud := range solicitudesRp {
				//recuperar datos del CDP objetivo de la solicitud

				var afectacionSolicitud []map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"disponibilidad_apropiacion_solicitud_rp?limit=0&query=SolicitudRp:"+strconv.Itoa(solicitud.Id), &afectacionSolicitud); err == nil {
					//consulta de la afectacion presupuestal objetivo.
					for _, afect := range afectacionSolicitud {
						var dispAprSol []map[string]interface{}
						if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"disponibilidad_apropiacion?limit=1&query=Id:"+fmt.Sprintf("%v", afect["DisponibilidadApropiacion"]), &dispAprSol); err == nil {
							for _, disp_apro := range dispAprSol {
								disp_apro["ValorAsignado"] = afect["Monto"]
								disp_apro["FuenteFinanciacion"] = disp_apro["FuenteFinanciamiento"]
							}
						} else {
							//si sale mal la consulta de la afectacion del cdp objetivo.

						}
					}
				} else {
					//si sale mal la consulta de la afectacion de la solicitud.
				}

				var cdpObjtvo []models.Disponibilidad
				if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.Cdp), &cdpObjtvo); err == nil {
					if cdpObjtvo != nil {
						solicitud.DatosDisponibilidad = &cdpObjtvo[0]
						var necesidadCdp []models.SolicitudDisponibilidad
						if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.Solicitud), &necesidadCdp); err == nil {
							if necesidadCdp != nil {
								solicitud.DatosDisponibilidad.DatosNecesidad = necesidadCdp[0].Necesidad
								var depNes []models.DependenciaNecesidad
								if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.DatosNecesidad.Id), &depNes); err == nil {
									if depNes != nil {
										var depSol []models.Dependencia
										var jefeDepSol []models.JefeDependencia
										if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol); err == nil {
											if jefeDepSol != nil {
												if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=1&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol); err == nil {
													if depSol != nil {
														solicitud.DatosDisponibilidad.DatosNecesidad.DatosDependenciaSolicitante = &depSol[0]
													} else {
														//si no hay datos de la dependencia
													}
												} else {
													//si hay error al consultar la dependecia solicitante
												}
											} else {
												//no hay datos jefe dep
											}

										} else {
											//jefe_dep
										}
									} else {
										//si no hay datos en la consulta dependencia_necesidad
									}
								} else {
									//si hay error al consultar dependencia_necesidad
								}
							} else {
								//si no hay datos de la necesidad
							}
						} else {
							//si hay error al encontrar datos de la necesidad
						}
					} else {
						//si no hay datos del CDP objetivo
						fmt.Println("error cdp: no hay datos, id : ", solicitud.Cdp)
					}
				} else {
					//si ocurre error al traer datos del CDP objetivo
					fmt.Println("error cdp: ", err)
				}
				//obtener informacion del contrato del rp
				var infoContrato []models.ContratoGeneral
				var contratista []models.InformacionProveedor
				if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"contrato_general?limit=1&query=Id:"+solicitud.NumeroContrato, &infoContrato); err == nil {
					if infoContrato != nil {
						if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_proveedor?limit=1&query=Id:"+strconv.Itoa(infoContrato[0].Contratista), &contratista); err == nil {
							solicitud.DatosProveedor = &contratista[0]
						} else {
							//error consulta proveedor
							fmt.Println("ss", err.Error())
						}
					} else {
						//si no encuentra datos sobre el contrato
						fmt.Println("error contrato: no hay datos, id : ", solicitud.NumeroContrato)
					}
				} else {
					//si ocurre error al obtener los datos del contrato
				}
				//cargar datos del compromiso de la solicitud de rp
				var compromisoRp []models.Compromiso
				if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"compromiso?limit=1&query=Id:"+strconv.Itoa(solicitud.TipoCompromiso), &compromisoRp); err == nil {
					if compromisoRp != nil {
						solicitud.DatosCompromiso = &compromisoRp[0]
					} else {
						//si no encuentra los datos del compromiso
					}
				} else {
					//si hay error al cargar el compromiso del rp
				}
				respuesta = append(respuesta, solicitud)
			}
			if respuesta != nil {
				c.Data["json"] = respuesta[0]
			} else {
				c.Data["json"] = respuesta
			}
		} else {
			//si no hay datos de solicitudes
			c.Data["json"] = nil
		}
	} else {
		//si ocurre error al traer las solicitudes
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// CargueMasivoPr ...
// @Title CargueMasivoPr
// @Description create RegistroPresupuestal
// @Param	body		body 	models.DatosRegistroPresupuestal	true		"body for DatosRegistroPresupuestal content"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /CargueMasivoPr [post]
func (c *RegistroPresupuestalController) CargueMasivoPr() {
	//variables a usar globales en el proceso.
	var dataRpRegistro []models.DatosRegistroPresupuestal
	var dataAlertas []models.Alert //array con las alertas generadas en aprobacion masiva de solicitudes
	var saldoCDP map[string]float64
	VigActual := time.Now().Year()
	tool := new(ruler.EntornoReglas)
	var respuestaServices interface{}
	//------------------------------------------------------
	tool.Agregar_dominio("Presupuesto")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataRpRegistro); err == nil {
		for _, rp_a_registrar := range dataRpRegistro { //recorrer el array de solicitudes cargadas
			for _, rubros_a_comprobar := range rp_a_registrar.Rubros { //recorrer la afectacion de la solicitud para inyeccion de reglas.
				datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
					Apropiacion: rubros_a_comprobar.Apropiacion, FuenteFinanciacion: rubros_a_comprobar.FuenteFinanciacion,
				}
				if err := request.SendJson("http://"+beego.AppConfig.String("crudService")+"disponibilidad/SaldoCdp", "POST", &saldoCDP, &datos); err == nil {
					fmt.Println(rubros_a_comprobar.FuenteFinanciacion)
					if rubros_a_comprobar.FuenteFinanciacion == nil {
						tool.Agregar_predicado("rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) +
							"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ").")
					} else {
						tool.Agregar_predicado("rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) +
							"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ").")
					}

				} else {
					dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
				}
				if rubros_a_comprobar.FuenteFinanciacion == nil {
					tool.Agregar_predicado("valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ").")

				} else {
					tool.Agregar_predicado("valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ").")

				}
			}
			var res string
			fmt.Println(rp_a_registrar)
			if rp_a_registrar.Rubros != nil {
				err := formatdata.FillStruct(tool.Ejecutar_result("aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y"), &res)
				if err == nil { //
					if int(rp_a_registrar.Rp.Vigencia) == VigActual {
						if res == "1" { // si se aprueba la solicitud
							rp_a_registrar.Rp.FechaRegistro = time.Now().Local()
							var comprobacion models.DatosRegistroPresupuestal
							if err := request.SendJson("http://"+beego.AppConfig.String("crudService")+"registro_presupuestal", "POST", &comprobacion, &rp_a_registrar); err == nil {
								dataAlertas = append(dataAlertas, models.Alert{Code: "S_543", Body: comprobacion, Type: "success"})
								rp_a_registrar.Rp.DatosSolicitud.Expedida = true
								if err := request.SendJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp/"+strconv.Itoa(rp_a_registrar.Rp.Solicitud), "PUT", &respuestaServices, &rp_a_registrar.Rp.DatosSolicitud); err == nil {
									//dataAlertas = append(dataAlertas, models.Alert{Code: "S_RP002", Body: respuestaServices, Type: "success"})

								} else {
									dataAlertas = append(dataAlertas, models.Alert{Code: "E_RP002", Body: rp_a_registrar, Type: "success"})
								}
							} else {
								dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
							}
						} else {
							dataAlertas = append(dataAlertas, models.Alert{Code: "E_RP001", Body: rp_a_registrar, Type: "error"})

						}
					} else {

						dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
					}
				} else {
					dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
				}
			} else {
				dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
			}
			//res := golog.GetBoolean(reglas, "aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y")
		}
	} else {
		fmt.Println("err 2 ", err.Error())
		dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"})
	}
	c.Data["json"] = dataAlertas //respuesta de las alertas generadas durante el proceso.
	c.ServeJSON()

}

func listaNecesidadesByRp(solicitudintfc interface{}, params ...interface{}) (res interface{}) {
	solicitud, e := solicitudintfc.(map[string]interface{})
	var rp []map[string]interface{}
	if e {
		idSol, e := solicitud["Id"].(float64)
		if e {
			if err := request.GetJson("http://"+beego.AppConfig.String("crudService")+"registro_presupuestal?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Disponibilidad.Solicitud:"+strconv.Itoa(int(idSol)), &rp); err == nil {
				if rp != nil {
					solicitud["InfoRp"] = rp[0]
					return solicitud
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
}

// ListaNecesidadesByRp ...
// @Title ListaNecesidadesByRp
// @Description Lsta de las necesidades origen de los rp registrados
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	tipoNecesidad	query	string	false	"tipo de la necesidad origen del rp"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /ListaNecesidadesByRp/:vigencia [get]
func (c *RegistroPresupuestalController) ListaNecesidadesByRp() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	var solicitudNecesidad []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var tipoNecesidad string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if startrange != "" && endrange != "" {
		query = ",FechaSolicitud__gte:" + startrange + ",FechaSolicitud__lte:" + endrange

	}
	if r := c.GetString("tipoNecesidad"); r != "" {
		tipoNecesidad = r
		//peticion a los rp expedidos
		if err = request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Vigencia:"+strconv.Itoa(vigencia)+",Necesidad.TipoNecesidad.CodigoAbreviacion:"+tipoNecesidad+query, &solicitudNecesidad); err == nil {
			if solicitudNecesidad != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(solicitudNecesidad...)
				chsolicitud := optimize.Digest(done, listaNecesidadesByRp, resch, nil)
				for solicitud := range chsolicitud {
					if solicitud != nil {
						respuesta = append(respuesta, solicitud.(map[string]interface{}))
					}
				}
				c.Data["json"] = respuesta
			} else {
				//si no encuentra solicitudes...
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			//si ocurre error al traer las solicitudes
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			fmt.Println(err.Error())
		}

	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "TipoNecesidad lost", Type: "error"}
	}

	c.ServeJSON()

}

// SolicitudesRpByDependencia ...
// @Title SolicitudesRpByDependencia
// @Description Lsta de las necesidades origen de los rp registrados
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	idDependencia	query	string	false	"tipo de la necesidad origen del rp"
// @Param	tipoNecesidad	query	string	false	"tipo de la necesidad origen del rp"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /SolicitudesRpByDependencia/:vigencia [get]
func (c *RegistroPresupuestalController) SolicitudesRpByDependencia() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	//var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var queryF string
	var idDependencia int64
	var tipoNecesidad string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
		fmt.Println(limit)
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
		fmt.Println(offset)
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if startrange != "" && endrange != "" {
		queryF = ",FechaSolicitud__gte:" + startrange + ",FechaSolicitud__lte:" + endrange

	}
	if tipoNecesidad = c.GetString("tipoNecesidad"); tipoNecesidad != "" {
		if idDependencia, err = c.GetInt64("idDependencia"); err == nil {
			var jefeDependenciaInfo []interface{}

			//consultar por rangos los jefes de dependencia asociados a la dependencia objetivo por toda la vigencia
			//si no se establece un rango de fechas.
			fInicio := time.Date(vigencia, time.January, 1, 0, 0, 0, 0, time.Local)
			fFin := time.Date(vigencia+1, time.January, 1, 0, 0, 0, 0, time.Local)
			query = "FechaInicio__lte:" + fInicio.Format("2006-01-02") + ",FechaFin__lte:" + fFin.Format("2006-01-02")

			fmt.Println(query)
			//consulta del servicio para determinar el filtro de la necesidad.
			//fmt.Println("http://" + beego.AppConfig.String("coreService") + "jefe_dependencia?fields=Id&limit=-1&query=DependenciaId:" + strconv.FormatInt(idDependencia, 10) + "," + query)
			if err = request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?fields=Id&limit=-1&query=DependenciaId:"+strconv.FormatInt(idDependencia, 10)+","+query, &jefeDependenciaInfo); err == nil && jefeDependenciaInfo != nil {
				//consulta sobre necesidades expedidas
				//solicitud_disponibilidad?query=Necesidad.DependenciaReversa.JefeDependenciaSolicitante:
				inQuery := ""
				for i, jefeDep := range jefeDependenciaInfo {
					jefeDepmap, e := jefeDep.(map[string]interface{})
					if e {
						if idJefeDep, e := jefeDepmap["Id"].(float64); e {
							if i == len(jefeDependenciaInfo)-1 {
								inQuery = inQuery + strconv.Itoa(int(idJefeDep))
							} else {
								inQuery = inQuery + strconv.Itoa(int(idJefeDep)) + "|"
							}

						}

					}

				}
				var solicitudNecesidad []map[string]interface{}
				//fmt.Println("http://" + beego.AppConfig.String("argoService") + "solicitud_disponibilidad?limit=-1&query=Necesidad.DependenciaReversa.JefeDependenciaSolicitante__in:" + inQuery)
				if err = request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=-1&query=Expedida:true,JustificacionRechazo:,Necesidad.TipoNecesidad.CodigoAbreviacion:"+tipoNecesidad+",Necesidad.DependenciaReversa.JefeDependenciaSolicitante__in:"+inQuery, &solicitudNecesidad); err == nil && solicitudNecesidad != nil {
					inQuery = ""
					for i, solicitud := range solicitudNecesidad {
						if idSolicitud, e := solicitud["Id"].(float64); e {
							if i == len(solicitudNecesidad)-1 {
								inQuery = inQuery + strconv.Itoa(int(idSolicitud))
							} else {
								inQuery = inQuery + strconv.Itoa(int(idSolicitud)) + "|"
							}
						}
					}
					//consulta de los CDP expedidos a las necesidades filtradas
					var Disponibilidades []map[string]interface{}
					if err = request.GetJson("http://"+beego.AppConfig.String("kronosService")+"disponibilidad?limit=-1&query=Solicitud__in:"+inQuery+",Vigencia:"+vigenciaStr, &Disponibilidades); err == nil && Disponibilidades != nil {
						inQuery = ""
						for i, disponibilidad := range Disponibilidades {
							if idDisponibilidad, e := disponibilidad["Id"].(float64); e {
								if i == len(Disponibilidades)-1 {
									inQuery = inQuery + strconv.Itoa(int(idDisponibilidad))
								} else {
									inQuery = inQuery + strconv.Itoa(int(idDisponibilidad)) + "|"
								}
							}
						}
						var solicitudRp []interface{}
						var respuesta []interface{}
						//fmt.Println("http://" + beego.AppConfig.String("argoService") + "solicitud_rp?limit=-1&query=Masivo:true,Expedida:false,JustificacionRechazo:,Cdp__in:" + inQuery + queryF)
						if err = request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Masivo:true,Expedida:false,JustificacionRechazo:,Cdp__in:"+inQuery+queryF+",Vigencia:"+vigenciaStr, &solicitudRp); err == nil && solicitudRp != nil {
							done := make(chan interface{})
							defer close(done)
							resch := optimize.GenChanInterface(solicitudRp...)
							chsolicitud := optimize.Digest(done, formatoSolicitudRP, resch, nil)
							var rubrosAfectados []map[string]interface{}
							for solicitud := range chsolicitud {
								if solicitudmap, e := solicitud.(models.SolicitudRp); e {

									rubrosint := solicitudmap.Rubros
									for _, afectacion := range rubrosint {
										row := make(map[string]interface{})
										if afectacionmap, e := afectacion.(map[string]interface{}); e {

											existe := false
											for _, rubro := range rubrosAfectados {
												if rubro["Apropiacion"].(map[string]interface{})["Id"].(float64) == afectacionmap["Apropiacion"].(map[string]interface{})["Id"].(float64) && rubro["FuenteFinanciamiento"].(map[string]interface{})["Id"].(float64) == afectacionmap["FuenteFinanciamiento"].(map[string]interface{})["Id"].(float64) {
													rubro["Valor"] = rubro["Valor"].(float64) + afectacionmap["ValorAsignado"].(float64)
													existe = true
												}

											}
											if !existe {
												row["Apropiacion"] = afectacionmap["Apropiacion"]
												row["FuenteFinanciamiento"] = afectacionmap["FuenteFinanciamiento"]
												row["Valor"] = afectacionmap["ValorAsignado"]
												rubrosAfectados = append(rubrosAfectados, row)
											}

										}
									}

								}
								respuesta = append(respuesta, solicitud)
							}
							c.Data["json"] = map[string]interface{}{"InformacionRp": respuesta, "ResumenCargueRp": rubrosAfectados}
						} else {
							c.Data["json"] = models.Alert{Code: "E_0458", Body: "Solicitud RP lost", Type: "error"}
						}
					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: "Disponibilidad lost", Type: "error"}

					}
				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: "Solicitud Necesidad Lost", Type: "error"}
				}

			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: "jefe dependencia lost", Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: "idDependencia lost", Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "tipoNecesidad lost", Type: "error"}
	}

	c.ServeJSON()
}

// AprobarAnulacion ...
// @Title AprobarAnulacion
// @Description create RegistroPresupuestal
// @Param	body		body 	models.AnulacionRegistroPresupuestal	true		"body for AnulacionRegistroPresupuestal content"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /AprobarAnulacion [post]
func (c *RegistroPresupuestalController) AprobarAnulacion() {
	try.This(func() {
		var v models.AnulacionRegistroPresupuestal
		var res models.Alert
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			Urlcrud := "http://" + beego.AppConfig.String("crudService")+"registro_presupuestal/AprobarAnulacion/"
			if err := request.SendJson(Urlcrud, "POST", &res, &v); err == nil {
				c.Data["json"] = res
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = err.Error()
		}
	}).Catch(func(e try.E) {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: e, Type: "error"}
	})

	c.ServeJSON()
}

// AddRpMongo ...
func AddRpMongo(parameter ...interface{}) (err interface{}) {
	try.This(func() {
		infoRp := parameter[0].(models.DatosRegistroPresupuestal)
		dataSend := make(map[string]interface{})
		dataSend["Vigencia"] = strconv.FormatFloat(infoRp.Rp.Vigencia, 'f', 0, 64)
		dataSend["Id"] = infoRp.Rp.Id
		dataSend["MesRegistro"] = strconv.Itoa(int(infoRp.Rp.FechaRegistro.Month()))
		var afectacion []interface{}
		var resM map[string]interface{}
		aux := make(map[string]interface{})

		for _, data := range infoRp.Rubros {
			aux["Rubro"] = data.Apropiacion.Rubro.Codigo
			aux["UnidadEjecutora"] = strconv.Itoa(int(data.Apropiacion.Rubro.UnidadEjecutora))
			aux["Valor"] = data.Valor
			aux["FuenteCodigo"] = data.FuenteFinanciacion.Codigo
			aux["FuenteNombre"] = data.FuenteFinanciacion.Nombre
			dataSend["Disponibilidad"] = data.Disponibilidad.Id
			afectacion = append(afectacion, aux)
		}

		dataSend["Afectacion"] = afectacion
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/Rp"
		if err1 := request.SendJson(Urlmongo, "POST", &resM, &dataSend); err1 == nil {
			if resM["Type"].(string) == "success" {
				err = err1
			} else {
				panic("Mongo api error")
			}
		} else {
			panic("Mongo Not Found")
		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
		infoRp := parameter[0].(models.DatosRegistroPresupuestal)
		var resC interface{}
		Urlcrud := "http://" + beego.AppConfig.String("crudService")+"registro_presupuestal/DeleteRpData/" + strconv.Itoa(infoRp.Rp.Id)
		if errorDelete := request.SendJson(Urlcrud, "DELETE", &resC, nil); errorDelete == nil {
			beego.Info("Data ", resC)
		} else {
			beego.Info("Error ", errorDelete)
		}

		err = e
	})
	return
}

// AddOpMongo ...
func AddOpMongo(parameter ...interface{}) (err interface{}) {
	try.This(func() {
		infoOp := models.OrdenPago{}
		infoOpData := parameter[0].(map[string]interface{})["OrdenPago"].([]interface{})[0]
		//beego.Info("infoOpData", infoOpData)
		err := formatdata.FillStruct(infoOpData, &infoOp)
		if err != nil {
			panic(err.Error())
		}
		dataSend := make(map[string]interface{})
		dataSend["Vigencia"] = strconv.FormatFloat(infoOp.Vigencia, 'f', 0, 64)
		dataSend["Id"] = infoOp.Id
		dataSend["MesRegistro"] = strconv.Itoa(int(time.Now().Month()))
		dataSend["TipoMovimiento"] = "OP-Proveedor"
		var afectacion []interface{}
		var resM map[string]interface{}
		dispApropiacion := []models.RegistroPresupuestalDisponibilidadApropiacion{}
		aux := make(map[string]interface{})
		urlcrud := "http://" + beego.AppConfig.String("crudService") 
		for _, data := range infoOp.OrdenPagoRegistroPresupuestal {
			aux["RP"] = data.RegistroPresupuestal.Id
			aux["FuenteNombre"] = "" //data.FuenteFinanciacion.Nombre
			//beego.Info(urlcrud + "/registro_presupuestal_disponibilidad_apropiacion?limit=-1" + "&query=RegistroPresupuestal.Id:" + strconv.Itoa(data.RegistroPresupuestal.Id) + "&fields=DisponibilidadApropiacion")
			if err := request.GetJson(urlcrud+"/registro_presupuestal_disponibilidad_apropiacion?limit=-1"+"&query=RegistroPresupuestal.Id:"+strconv.Itoa(data.RegistroPresupuestal.Id)+"&fields=DisponibilidadApropiacion", &dispApropiacion); err == nil {
				aux["Rubro"] = dispApropiacion[0].DisponibilidadApropiacion.Apropiacion.Rubro.Codigo
			}
			aux["UnidadEjecutora"] = infoOp.UnidadEjecutora
			aux["Valor"] = infoOp.ValorBase
			aux["FuenteCodigo"] = "" //data.FuenteFinanciacion.Codigo
			afectacion = append(afectacion, aux)
		}
		dataSend["Afectacion"] = afectacion
		//beego.Info("dataSend", dataSend)
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/Op"
		if err1 := request.SendJson(Urlmongo, "POST", &resM, &dataSend); err1 == nil {
			if resM["Type"].(string) == "success" {
				err = err1
			} else {
				panic("Mongo api error")
			}
		} else {
			panic("Mongo Not Found")
		}

	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
		infoOp := parameter[0].(models.OrdenPagoEstadoOrdenPago)
		var resC interface{}
		Urlcrud := "http://" + beego.AppConfig.String("crudService")+"orden_pago_estado_orden_pago/DeleteOpData/" + strconv.Itoa(infoOp.OrdenPago.Id)
		if errorDelete := request.SendJson(Urlcrud, "DELETE", &resC, nil); errorDelete == nil {
			beego.Info("Data ", resC)
		} else {
			beego.Info("Error ", errorDelete)
		}

		err = e
	})
	return
}

//AddAnulacionRpMongo ...
func AddAnulacionRpMongo(parameter ...interface{}) (err interface{}) {
	infoAnulacion := models.AnulacionRegistroPresupuestal{}
	try.This(func() {
		infoAnulacionint := parameter[0].(map[string]interface{})
		error := formatdata.FillStruct(infoAnulacionint, &infoAnulacion)
		if err != nil {
			panic(error.Error())
		}
		dataSend := make(map[string]interface{})
		dataSend["Id"] = infoAnulacion.Id
		dataSend["MesRegistro"] = strconv.Itoa(int(infoAnulacion.FechaRegistro.Month()))
		dataSend["Vigencia"] = strconv.FormatFloat(infoAnulacion.AnulacionRegistroPresupuestalDisponibilidadApropiacion[0].RegistroPresupuestalDisponibilidadApropiacion.RegistroPresupuestal.Vigencia, 'f', 0, 64)
		var afectacion []interface{}
		var resM map[string]interface{}
		aux := make(map[string]interface{})
		for _, data := range infoAnulacion.AnulacionRegistroPresupuestalDisponibilidadApropiacion {
			aux["Rubro"] = data.RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Apropiacion.Rubro.Codigo
			aux["UnidadEjecutora"] = strconv.Itoa(int(data.RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora))
			aux["Valor"] = data.Valor
			dataSend["Disponibilidad"] = data.RegistroPresupuestalDisponibilidadApropiacion.RegistroPresupuestal.Id
			afectacion = append(afectacion, aux)
		}
		dataSend["Afectacion"] = afectacion
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/AnulacionRp"
		beego.Info("Data to send ", dataSend)
		if err1 := request.SendJson(Urlmongo, "POST", &resM, &dataSend); err1 == nil {
			if resM["Type"].(string) == "success" {
				err = err1
			} else {
				panic("Mongo api error")
			}
		} else {
			panic("Mongo Not Found")
		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
		var resC interface{}
		Urlcrud := "http://" + beego.AppConfig.String("crudService")+"anulacion_registro_presupuestal/" + strconv.Itoa(infoAnulacion.Id)
		infoAnulacion.EstadoAnulacion["Id"] = 2
		if errorPut := request.SendJson(Urlcrud, "PUT", &resC, &infoAnulacion); errorPut == nil {
			beego.Info("Data ", resC)
		} else {
			beego.Info("Error ", errorPut)
		}

	})
	return
}

// SaldoRp ...
// @Title SaldoRp
// @Description create RegistroPresupuestal
// @Param	idPsql		path 	int	true		"idPsql del documento"
// @Param	rubro		path 	string	true		"código del rubro"
// @Param	fuente		query	string false		"fuente de financiamiento"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /SaldoRp/:idPsql/:rubro [get]
func (c *RegistroPresupuestalController) SaldoRp() {
	try.This(func() {
		var (
			cdpID     int
			err       error
			infoSaldo map[string]float64
		)

		cdpID, err = c.GetInt(":idPsql") // id psql del cdp
		if err != nil {
			panic(err.Error())
		}
		rubro := c.GetString(":rubro")
		fuente := c.GetString("fuente")
		valoresSuman := map[string]bool{"Valor": true, "TotalAnuladoOp": true}
		infoSaldo = CalcularSaldoMovimiento(rubro, fuente, "Rp", cdpID, valoresSuman)
		c.Data["json"] = infoSaldo

	}).Catch(func(e try.E) {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: e, Type: "error"}
	})

	c.ServeJSON()
}
