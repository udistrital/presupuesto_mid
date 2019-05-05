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

//DisponibilidadController operations for Disponibilidad
type DisponibilidadController struct {
	beego.Controller
}

// URLMapping ...
func (c *DisponibilidadController) URLMapping() {
	c.Mapping("Solicitudes", c.InfoSolicitudDisponibilidad)
	c.Mapping("SolicitudById", c.InfoSolicitudDisponibilidadByID)
}
func formatoListaCDPDispatcher(tipo int) (f func(data map[string]interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case 1:
		return formatoListaCDPConSolicitud
	case 2:
		return formatoListaCDPMovimiento
	default:
		return nil
	}
}

func expedirDisponibilidadDispatcher(tipo int) (f func(data []map[string]interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case 1:
		return expedirDisponibilidadConNecesidad
	case 2:
		return expedirDisponibilidadConNecesidad
	default:
		return nil
	}
}
func formatoListaCDP(disponibilidad interface{}, params ...interface{}) (res interface{}) {

	dispMap := disponibilidad.(map[string]interface{})
	if dispPeInt, e := dispMap["DisponibilidadProcesoExterno"].([]interface{}); e && dispPeInt[0].(map[string]interface{})["TipoDisponibilidad"] != nil {
		dispPe := dispPeInt[0].(map[string]interface{})
		var params []interface{}
		params = append(params, dispPe["ProcesoExterno"].(float64))
		if f := formatoListaCDPDispatcher(int(dispPe["TipoDisponibilidad"].(map[string]interface{})["Id"].(float64))); f != nil {

			res = f(dispMap, params)
		} else {

			return dispMap
		}

	} else {
		return dispMap
	}

	return
}

func formatoListaCDPMovimiento(disponibilidad map[string]interface{}, params ...interface{}) (res interface{}) {
	var solicitud []map[string]interface{}
	solicitudArr, ee := disponibilidad["DisponibilidadProcesoExterno"].([]interface{})
	solicitudMap, ee := solicitudArr[0].(map[string]interface{})
	solicitudNo, ee := solicitudMap["ProcesoExterno"].(float64)
	if params != nil && ee {
		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"movimiento_apropiacion?limit=0&query=Id:"+strconv.FormatFloat(solicitudNo, 'f', -1, 64), &solicitud); err == nil {
			disponibilidad["Solicitud"] = solicitud[0]
		} else {
			return map[string]interface{}{"Code": "E_0458", "Body": "kronos Service", "Type": "error"}
		}
	} else {
		return map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in Disponibilidad Procses", "Type": "error"}
	}
	return disponibilidad
}

func formatoListaCDPConSolicitud(disponibilidad map[string]interface{}, params ...interface{}) (res interface{}) {
	var solicitud []models.SolicitudDisponibilidad
	solicitudArr, ee := disponibilidad["DisponibilidadProcesoExterno"].([]interface{})
	solicitudMap, ee := solicitudArr[0].(map[string]interface{})
	solicitudNo, ee := solicitudMap["ProcesoExterno"].(float64)
	if params != nil && ee {
		beego.Info("Peticion: ", "http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=0&query=Id:"+strconv.FormatFloat(solicitudNo, 'f', -1, 64))
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=0&query=Id:"+strconv.FormatFloat(solicitudNo, 'f', -1, 64), &solicitud); err == nil {
			beego.Info("Respuesta: ", solicitud)
			for _, resultado := range solicitud {

				var depNes []models.DependenciaNecesidad
				var jefeDepSol []models.JefeDependencia
				var depSol []models.Dependencia
				var depDest []models.Dependencia
				var necesidad []models.Necesidad
				//var temp models.InfoSolDisp
				//temp.SolicitudDisponibilidad = &resultado
				if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(resultado.Necesidad.Id), &necesidad); err == nil {
					if necesidad != nil {
						necesidadaux := necesidad[0]
						resultado.Necesidad = &necesidadaux
						fmt.Println(necesidadaux)

						beego.Info("http://" + beego.AppConfig.String("argoService") + "dependencia_necesidad?limit=0&query=Necesidad.Id:" + strconv.Itoa(necesidad[0].Id))
						if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes); err == nil {
							//fmt.Println("http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=0&query=Id:" + strconv.Itoa(depNes[0].JefeDependenciaSolicitante))
							if depNes != nil {
								if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol); err == nil {
									//temp.DependenciaSolicitante = &depSol[0]
									if jefeDepSol != nil {
										if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol); err == nil {
											//temp.DependenciaDestino = &depDest[0]
										} else {

										}
										if depSol != nil {
											if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(jefeDepSol[0].TerceroId), &depSol[0].InfoJefeDependencia); err == nil {
												//temp.DependenciaSolicitante = &depSol[0]

											} else {

											}
											fmt.Println(depNes[0].OrdenadorGasto)
											if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); err == nil {
												//temp.DependenciaSolicitante = &depSol[0]

											} else {

											}
											if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depDest); err == nil {
												//temp.DependenciaDestino = &depDest[0]
											} else {

											}
										}
									} else {

									}
								} else {
									fmt.Println(err)
								}
							}

							if depSol == nil {
								depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
							}
							if depDest == nil {
								depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
							}
							temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: depSol[0], DependenciaDestino: depDest[0]}
							disponibilidad["Solicitud"] = temp

						} else {
							// return map[string]interface{}{"Code": "E_0458", "Body": "argo Service1", "Type": "error"}
							temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: models.Dependencia{Nombre: "Indefinida"}, DependenciaDestino: models.Dependencia{Nombre: "Indefinida"}}
							disponibilidad["Solicitud"] = temp
						}
					} //else {
					// 	return map[string]interface{}{"Code": "E_0458", "Body": "argo Service2", "Type": "error"}
					// }
				} //else {
				// 	return map[string]interface{}{"Code": "E_0458", "Body": "argo Service3", "Type": "error"}
				// }
			}

		} // else {
		// 	beego.Info("Error argo: ", err)
		// 	return map[string]interface{}{"Code": "E_0458", "Body": err, "Type": "error"}
		// }
	} //else {
	// 	return map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in Disponibilidad Procses", "Type": "error"}
	// }

	return disponibilidad
}

// ListaDisponibilidades ...
// @Title ListaDisponibilidades
// @Description get Disponibilidad by vigencia
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	UnidadEjecutora	query	string	false	"unidad ejecutora de las solicitudes a consultar"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado para la lista de los cdp"
// @Success 200 {object} models.Disponibilidad
// @Failure 403
// @router ListaDisponibilidades/:vigencia [get]
func (c *DisponibilidadController) ListaDisponibilidades() {
	var disponibilidades []interface{}
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
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err1 := strconv.Atoi(vigenciaStr)
	UnidadEjecutora, err2 := c.GetInt("UnidadEjecutora")
	beego.Info("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/disponibilidad?limit=" + strconv.FormatInt(limit, 10) + "&offset=" + strconv.FormatInt(offset, 10) + "&query=Vigencia:" + strconv.Itoa(vigencia) + ",DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora:" + strconv.Itoa(UnidadEjecutora) + query)
	if err1 == nil && err2 == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Vigencia:"+strconv.Itoa(vigencia)+",DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora:"+strconv.Itoa(UnidadEjecutora)+query, &disponibilidades); err == nil {
			if disponibilidades != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(disponibilidades...)
				chdisponibilidades := optimize.Digest(done, formatoListaCDP, resch, nil)
				for disponibilidad := range chdisponibilidades {
					respuesta = append(respuesta, disponibilidad.(map[string]interface{}))
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

// DisponibilidadByNecesidad ...
// @Title DisponibilidadByNecesidad
// @Description get Disponibilidad by id Necesidad
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} interface{}
// @Failure 403
// @router DisponibilidadByNecesidad/:id [get]
func (c *DisponibilidadController) DisponibilidadByNecesidad() {
	var resdisponibilidad []map[string]interface{}
	var solicitud []map[string]interface{}
	idStr := c.Ctx.Input.Param(":id")
	try.This(func() {
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Expedida:true,Necesidad.Id:"+idStr, &solicitud); err == nil {
			var id float64
			//err = formatdata.FillStruct(solicitud[0]["Id"], &id)
			id, _ = solicitud[0]["ProcesoExterno"].(float64)
			if err != nil {
				c.Data["json"] = models.Alert{Code: "", Type: "error", Body: err.Error()}
				c.ServeJSON()
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad?limit=1&query=Solicitud:"+strconv.FormatFloat(id, 'f', -1, 64), &resdisponibilidad); err == nil {
				err = formatdata.FillStruct(resdisponibilidad[0]["Id"], &id)
				if err != nil {
					c.Data["json"] = models.Alert{Code: "", Type: "error", Body: err.Error()}
					c.ServeJSON()
				}
				var rp []map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Disponibilidad.Id:"+strconv.FormatFloat(id, 'f', -1, 64), &rp); err == nil {
					resdisponibilidad[0]["registro_presupuestal"] = rp
					c.Data["json"] = resdisponibilidad
				} else {
					fmt.Println(err)
					c.Data["json"] = models.Alert{Code: "", Type: "error", Body: err.Error()}
				}

			} else {
				fmt.Println(err)
				c.Data["json"] = models.Alert{Code: "", Type: "error", Body: err.Error()}
			}

		} else {
			c.Data["json"] = err.Error()
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = resdisponibilidad
	})
	c.ServeJSON()
}

func formatoSolicitudCDP(solicitudint interface{}, params ...interface{}) (res interface{}) {
	var depNes []models.DependenciaNecesidad
	var jefeDepSol []models.JefeDependencia
	var jefeDepDest []models.JefeDependencia
	var depSol []models.Dependencia
	var depDest []models.Dependencia
	var necesidad []models.Necesidad
	solicitud, e := solicitudint.(map[string]interface{})
	if e {
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(int(solicitud["Necesidad"].(map[string]interface{})["Id"].(float64))), &necesidad); err == nil {
			necesidadaux := necesidad[0]
			solicitud["Necesidad"] = &necesidadaux
			err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes)
			if depNes != nil && err == nil {
				err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol)
				if jefeDepSol != nil && err == nil {
					err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol)
					if depSol != nil && err == nil {
						if errJefeDep := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(jefeDepSol[0].TerceroId), &depSol[0].InfoJefeDependencia); errJefeDep == nil {
							fmt.Println("Información del jefe de dependencia consultada correctamente")
						} else {
							fmt.Println("Error al consultar Información del jefe de dependencia")
						}
						if errorInfoOrdenador := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); errorInfoOrdenador == nil {
							fmt.Println("Información del ordenador consultada correctamente")
						} else {
							fmt.Println("Error al consultar el ordenador")
						}
					}
				}
				err1 := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaDestino), &jefeDepDest)
				if jefeDepDest != nil && err1 == nil {
					err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepDest[0].DependenciaId), &depDest)
					if depDest != nil && err == nil {
						if errJefeDep := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(jefeDepDest[0].TerceroId), &depDest[0].InfoJefeDependencia); errJefeDep == nil {
							fmt.Println("Información del jefe de dependencia consultada correctamente")
						} else {
							fmt.Println("Error al consultar Información del jefe de dependencia")
						}
						if errorInfoOrdenador := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depDest[0].InfoOrdenador); errorInfoOrdenador == nil {
							fmt.Println("Información del ordenador consultada correctamente")
						} else {
							fmt.Println("Error al consultar el ordenador")
						}
					}
				}
			}

			if depSol == nil {
				depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
			}
			if depDest == nil {
				depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
			}
			temp := map[string]interface{}{"SolicitudDisponibilidad": solicitud, "DependenciaSolicitante": depSol[0], "DependenciaDestino": depDest[0]}
			res = temp
			return
		} else {
			fmt.Println("error3: ", err)
			//return nil
		}
	}
	return res

}

// InfoSolicitudDisponibilidad ...
// @Title InfoSolicitudDisponibilidad
// @Description get InfoSolicitudDisponibilidad
// @Param	vigencia	query	string	false	"vigencia de las solicitudes a consultar"
// @Param	UnidadEjecutora	query	string	false	"unidad ejecutora de las solicitudes a consultar"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	query	query	string	false	"query de filtrado para la lista de los cdp"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Success 200 {object} models.InfoSolDisp
// @Failure 403
// @router Solicitudes/:vigencia [get]
func (c *DisponibilidadController) InfoSolicitudDisponibilidad() {
	var solicitud []interface{}
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
	if err1 == nil && err2 == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Necesidad.UnidadEjecutora:"+strconv.Itoa(UnidadEjecutora)+",Vigencia:"+strconv.Itoa(vigencia)+",Expedida:false"+query+",JustificacionRechazo:&sortby=Id&order=desc", &solicitud); err == nil {
			if solicitud != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(solicitud...)
				chsolicitud := optimize.Digest(done, formatoSolicitudCDP, resch, nil)
				for solicitud := range chsolicitud {
					aux, e := solicitud.(map[string]interface{})
					if e {
						respuesta = append(respuesta, aux)
					}

				}
				c.Data["json"] = respuesta
			} else {

			}

		} else {
			c.Data["json"] = err.Error()
			fmt.Println("error1: ", err)
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}

	c.ServeJSON()
}

//InfoSolicitudDisponibilidadByID ...
// @Title InfoSolicitudDisponibilidadByID
// @Description get InfoSolicitudDisponibilidadByID
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.InfoSolDisp
// @Failure 403
// @router SolicitudById/:id [get]
func (c *DisponibilidadController) InfoSolicitudDisponibilidadByID() {
	var solicitud []models.SolicitudDisponibilidad
	idStr := c.Ctx.Input.Param(":id")
	if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=0&query=Id:"+idStr, &solicitud); err == nil {

		for _, resultado := range solicitud {

			var depNes []models.DependenciaNecesidad
			var jefeDepSol []models.JefeDependencia
			var depSol []models.Dependencia
			var depDest []models.Dependencia
			var necesidad []models.Necesidad
			//var temp models.InfoSolDisp
			//temp.SolicitudDisponibilidad = &resultado
			if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(resultado.Necesidad.Id), &necesidad); err == nil {
				if necesidad != nil {
					necesidadaux := necesidad[0]
					resultado.Necesidad = &necesidadaux
					fmt.Println(necesidadaux)
					if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes); err == nil {
						//fmt.Println("http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=0&query=Id:" + strconv.Itoa(depNes[0].JefeDependenciaSolicitante))
						if depNes != nil {
							if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol); err == nil {
								//temp.DependenciaSolicitante = &depSol[0]
								if jefeDepSol != nil {
									if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol); err == nil {
										//temp.DependenciaDestino = &depDest[0]
									} else {
										c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
									}
									if depSol != nil {
										if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(jefeDepSol[0].TerceroId), &depSol[0].InfoJefeDependencia); err == nil {
											//temp.DependenciaSolicitante = &depSol[0]

										} else {
											c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

										}
										fmt.Println(depNes[0].OrdenadorGasto)
										if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); err == nil {
											//temp.DependenciaSolicitante = &depSol[0]

										} else {
											c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

										}
										if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depDest); err == nil {
											//temp.DependenciaDestino = &depDest[0]
										} else {
											c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

										}
									}
								} else {
									c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
								}
							} else {
								fmt.Println(err)
							}
						} else {
							c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}

						}

						if depSol == nil {
							depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
						}
						if depDest == nil {
							depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
						}
						temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: depSol[0], DependenciaDestino: depDest[0]}
						c.Data["json"] = temp

					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

					}
				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}

				}
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

			}
		}

	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}

	}
	c.ServeJSON()
}

//DetalleSolicitudDisponibilidadByID ....
func DetalleSolicitudDisponibilidadByID(idStr string) (res models.InfoSolDisp, err error) {
	var resultado models.SolicitudDisponibilidad
	if err = request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/"+idStr, &resultado); err == nil {

		var depNes []models.DependenciaNecesidad
		var jefeDepSol []models.JefeDependencia
		var depSol []models.Dependencia
		var depDest []models.Dependencia
		var necesidad []models.Necesidad
		//var temp models.InfoSolDisp
		//temp.SolicitudDisponibilidad = &resultado
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(resultado.Necesidad.Id), &necesidad); err == nil {
			if necesidad != nil {
				necesidadaux := necesidad[0]
				resultado.Necesidad = &necesidadaux
				fmt.Println(necesidadaux)
				if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes); err == nil {
					//fmt.Println("http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=0&query=Id:" + strconv.Itoa(depNes[0].JefeDependenciaSolicitante))
					if depNes != nil {
						if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefeDepSol); err == nil {
							//temp.DependenciaSolicitante = &depSol[0]
							if jefeDepSol != nil {
								if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depSol); err == nil {
									//temp.DependenciaDestino = &depDest[0]
								} else {

								}
								if depSol != nil {
									if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(jefeDepSol[0].TerceroId), &depSol[0].InfoJefeDependencia); err == nil {
										//temp.DependenciaSolicitante = &depSol[0]

									} else {

									}
									fmt.Println(depNes[0].OrdenadorGasto)
									if err := request.GetJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); err == nil {
										//temp.DependenciaSolicitante = &depSol[0]

									} else {

									}
									if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(jefeDepSol[0].DependenciaId), &depDest); err == nil {
										//temp.DependenciaDestino = &depDest[0]
									} else {

									}
								}
							} else {

							}
						} else {
							fmt.Println(err)
						}
					} else {

					}

					if depSol == nil {
						depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
					}
					if depDest == nil {
						depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
					}
					temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: depSol[0], DependenciaDestino: depDest[0]}
					return temp, nil

				} else {

				}
			} else {

			}
		} else {

		}

	} else {
		return
	}
	return
}

// AprobarAnulacion ...
// @Title AprobarAnulacion
// @Description create AprobarAnulacion
// @Param	body		body 	map[string]interface{}	true		"body for InfoSolDisp content"
// @Success 201 {int} models.InfoSolDisp
// @Failure 403 body is empty
// @router /AprobarAnulacion [post]
func (c *DisponibilidadController) AprobarAnulacion() {
	var solicitudAnulacion map[string]interface{}
	var res []map[string]interface{} //objeto usado para guardar resultado de la peticion del servicio
	var id int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitudAnulacion); err == nil {
		//carga responsable de presupuesto
		var responsablePres []models.JefeDependencia //responsable
		if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=DependenciaId:102,FechaInicio__lte:"+time.Now().Format("2006-01-02")+",FechaFin__gte:"+time.Now().Format("2006-01-02"), &responsablePres); err == nil {
			//cargar necesidad que dio origen al cdp

			err = formatdata.FillStructDeep(solicitudAnulacion, "AnulacionDisponibilidadApropiacion", &res)
			//fmt.Println("err ", solicitudAnulacion)
			if err != nil {
				//err fill idsolicitud
				c.Data["json"] = err.Error()
				c.ServeJSON()
			}
			err = formatdata.FillStructDeep(res[0], "DisponibilidadApropiacion.Disponibilidad.Solicitud", &id)
			res[0] = make(map[string]interface{})
			//fmt.Println("err ", solicitudAnulacion)
			if err != nil {
				//err fill idsolicitud
				c.Data["json"] = "error 1 " + err.Error()
				c.ServeJSON()
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Id:"+strconv.Itoa(id), &res); err == nil { //traer solicitudes por id de necesidad

				err = formatdata.FillStructDeep(res[0], "Necesidad.Id", &id)
				res[0] = make(map[string]interface{})
				if err != nil {
					//err fill idsolicitud
					c.Data["json"] = err.Error()
					c.ServeJSON()
				}
				if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(id), &res); err == nil {

					err = formatdata.FillStructDeep(res[0], "JefeDependenciaDestino", &id)
					res[0] = make(map[string]interface{})
					if err != nil {
						//err fill idsolicitud
						c.Data["json"] = err.Error()
						c.ServeJSON()
					}
					if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=0&query=Id:"+strconv.Itoa(id), &res); err == nil {
						err = formatdata.FillStructDeep(res[0], "DependenciaId", &id)
						res[0] = make(map[string]interface{})
						fmt.Println("res ", id)
						if err != nil {
							//err fill idsolicitud
							c.Data["json"] = err.Error()
							c.ServeJSON()
						}
						if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=DependenciaId:"+strconv.Itoa(id)+",FechaInicio__lte:"+time.Now().Format("2006-01-02")+",FechaFin__gte:"+time.Now().Format("2006-01-02"), &res); err == nil {
							fmt.Println("http://" + beego.AppConfig.String("coreService") + "jefe_dependencia?limit=0&query=Id:" + strconv.Itoa(id))
							err = formatdata.FillStructDeep(res[0], "TerceroId", &id)
							res[0] = make(map[string]interface{})
							fmt.Println("res ", id)
							if err != nil {
								//err fill idsolicitud
								c.Data["json"] = err.Error()
								c.ServeJSON()
							}
							c.Data["json"] = id
							c.ServeJSON()
						}
					} else {
						//error json dependencia
						fmt.Println("res ", err)
						fmt.Println("http://" + beego.AppConfig.String("coreService") + "jefe_dependencia?limit=0&query=Id:" + strconv.Itoa(id))
						c.Data["json"] = "error json dependencia " + err.Error()
						c.ServeJSON()
					}
				} else {
					//error json dependencia_necesidad
					fmt.Println("res ", err)
					c.Data["json"] = "error json dependencia_necesidad " + err.Error()
					c.ServeJSON()
				}

			} else {
				//error json solicitud_disponibilidad
				fmt.Println("res ", err)
				c.Data["json"] = "error json solicitud_disponibilidad " + err.Error()
				c.ServeJSON()
			}

		}
	} else {
		//error json
		fmt.Println("res ", err)
		c.Data["json"] = "error json dependencia " + err.Error()
		c.ServeJSON()
	}
	/*var v map[string]interface{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		var id interface{}
		err = formatdata.FillStructDeep(v, "nivel1.nivel2.value", &id)
		if err != nil {
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = id
		}

		this.ServeJSON()
	}*/
}

// ExpedirDisponibilidad ...
// @Title ExpedirDisponibilidad
// @Description create Disponibilidad
// @Param	body		body 	map[string]string	true		"body for InfoSolDisp content"
// @Param	tipoDisponibilidad		query	string	true	"tipo de disponibilidad a expedir"
// @Success 201 {int} map[string]string
// @Failure 403 body is empty
// @router /ExpedirDisponibilidad [post]
func (c *DisponibilidadController) ExpedirDisponibilidad() {
	var infoSolicitudes []map[string]interface{}
	var alertas []models.Alert
	var res interface{}
	tipoDisponibilidad, err1 := c.GetInt("tipoDisponibilidad")
	if err1 == nil {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &infoSolicitudes); err == nil {
			//recorrer las solicitudes enviadas desde el cliente.
			if f := expedirDisponibilidadDispatcher(tipoDisponibilidad); f != nil {
				res = f(infoSolicitudes, nil)
			} else {
				alertas = append(alertas, models.Alert{Code: "E_0458", Body: "No function Found!", Type: "error"})
				res = alertas
			}
		} else {
			//no se recibieron los datos del cliente correctamente. c.Data["json"] = alertas
			alertas = append(alertas, models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"})
			res = alertas

		}
	} else {
		alertas = append(alertas, models.Alert{Code: "E_0458", Body: err1.Error(), Type: "error"})
		res = alertas
	}
	c.Data["json"] = res
	c.ServeJSON()
}

func expedirDisponibilidadConNecesidad(infoSolicitudes []map[string]interface{}, params ...interface{}) (res interface{}) {
	tool := new(ruler.EntornoReglas)
	var alertas []models.Alert
	var rubrosSolicitud []map[string]interface{}
	var mapSaldoApropiacion map[string]float64
	var apropiacion []map[string]interface{}
	VigActual := time.Now().Year()
	disponibilidad := make(map[string]interface{})
	infoDisponibilidad := make(map[string]interface{})
	aprobada := true
	tool.Agregar_dominio("Presupuesto")
	for _, solicitud := range infoSolicitudes {
		var afectacion []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"fuente_financiacion_rubro_necesidad?limit=-1&query=Necesidad.Id:"+strconv.Itoa(int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"].(float64))), &rubrosSolicitud); err == nil {
			//recorrer los rubros y/o fuentes solicitados
			for _, infoRubro := range rubrosSolicitud {
				//Solicitar el saldo de la apropiacion objetivo.
				//Se debe consultar la informacion del rubro para poder consumir saldo desde mongo
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion?query=Id:"+strconv.Itoa(int(infoRubro["Apropiacion"].(float64))), &apropiacion); err == nil {
					mapSaldoApropiacion = CalcularSaldoApropiacion(apropiacion[0]["Rubro"].(map[string]interface{})["Codigo"].(string), int(apropiacion[0]["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64)), int(apropiacion[0]["Vigencia"].(float64)))
					beego.Info("Saldo Apr: ", mapSaldoApropiacion)
					if mapSaldoApropiacion["saldo"] >= 0 {
						tool.Agregar_predicado("rubro_apropiacion(" + strconv.Itoa(int(infoRubro["Apropiacion"].(float64))) + "," + strconv.Itoa(int(infoRubro["FuenteFinanciamiento"].(float64))) + "," + strconv.FormatFloat(mapSaldoApropiacion["saldo"], 'f', -1, 64) + ").")
						tool.Agregar_predicado("valor_rubro_cdp(" + strconv.Itoa(int(infoRubro["Apropiacion"].(float64))) + "," + strconv.Itoa(int(infoRubro["FuenteFinanciamiento"].(float64))) + "," + strconv.FormatFloat(infoRubro["MontoParcial"].(float64), 'f', -1, 64) + ").")
					} else {
						tool.Agregar_predicado("rubro_apropiacion(" + strconv.Itoa(int(infoRubro["Apropiacion"].(float64))) + "," + strconv.Itoa(int(infoRubro["FuenteFinanciamiento"].(float64))) + ",0).")
						tool.Agregar_predicado("valor_rubro_cdp(" + strconv.Itoa(int(infoRubro["Apropiacion"].(float64))) + "," + strconv.Itoa(int(infoRubro["FuenteFinanciamiento"].(float64))) + "," + strconv.FormatFloat(infoRubro["MontoParcial"].(float64), 'f', -1, 64) + ").")
					}

				} else {
					alertas = append(alertas, models.Alert{Code: "E_CDP002", Body: solicitud, Type: "error"})
				}
				var res string
				//aqui condicion saldos fuentes.
				err := formatdata.FillStruct(tool.Ejecutar_result("aprobacion_cdp("+strconv.Itoa(int(infoRubro["Apropiacion"].(float64)))+",Y).", "Y"), &res)
				if err == nil {
					if res == "1" {
						//-----
						disponibilidadApropiacion := make(map[string]interface{})
						disponibilidadApropiacion["Apropiacion"] = map[string]interface{}{"Id": infoRubro["Apropiacion"]}
						disponibilidadApropiacion["disponibilidad"] = disponibilidad
						disponibilidadApropiacion["Valor"] = infoRubro["MontoParcial"].(float64)
						disponibilidadApropiacion["FuenteFinanciamiento"] = map[string]interface{}{"Id": infoRubro["FuenteFinanciamiento"]}
						afectacion = append(afectacion, disponibilidadApropiacion)

					} else {
						aprobada = false
						alertas = append(alertas, models.Alert{Code: "E_CDP001", Body: solicitud, Type: "error"})
					}
				} else {
					//si hay error al consultar las reglas de negocio. protocolo dentro de la peticion
					aprobada = false
					alertas = append(alertas, models.Alert{Code: "E_CDP002", Body: solicitud, Type: "error"})
				}
				if VigActual != int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Necesidad"].(map[string]interface{})["Vigencia"].(float64)) {
					aprobada = false
					alertas = append(alertas, models.Alert{Code: "E_CDP003", Body: solicitud, Type: "error"})
				}
			}
			if aprobada {
				disponibilidad["Vigencia"] = int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Necesidad"].(map[string]interface{})["Vigencia"].(float64))
				disponibilidad["FechaRegistro"] = time.Now().Local()
				disponibilidad["Estado"] = map[string]interface{}{"Id": 1}
				//disponibilidad["Solicitud"] = int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Id"].(float64))
				disponibilidad["Responsable"] = int(solicitud["Responsable"].(float64))
				disponibilidad["UnidadEjecutora"] = int(solicitud["Afectacion"].([]interface{})[0].(map[string]interface{})["Apropiacion"].(map[string]interface{})["Rubro"].(map[string]interface{})["UnidadEjecutora"].(float64))
				//----------------
				infoDisponibilidad["Disponibilidad"] = disponibilidad
				infoDisponibilidad["DisponibilidadApropiacion"] = afectacion
				DisponibilidadProcesoExterno := map[string]interface{}{"ProcesoExterno": int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Id"].(float64))}
				TipoDisponibilidad := map[string]interface{}{"Id": 1}
				DisponibilidadProcesoExterno["TipoDisponibilidad"] = TipoDisponibilidad
				infoDisponibilidad["DisponibilidadProcesoExterno"] = DisponibilidadProcesoExterno
				var respuesta models.Alert
				err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad", "POST", &respuesta, &infoDisponibilidad)
				if err == nil && respuesta.Type != "error" {
					var respuestaMod interface{}
					modsol := solicitud["SolicitudDisponibilidad"].(map[string]interface{})
					modsol["Expedida"] = true
					if err := request.SendJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/"+strconv.Itoa(int(solicitud["SolicitudDisponibilidad"].(map[string]interface{})["Id"].(float64))), "PUT", &respuestaMod, &modsol); err == nil {
						fmt.Println("solicitud de disponibilidad actualizada correctamente")
					} else {

						fmt.Println("error al actualizar solicitud de disponiblidad")
					}
					alertas = append(alertas, respuesta)
				} else {
					alertas = append(alertas, models.Alert{Code: "E_0458", Body: solicitud, Type: "error"})
				}

			}
		} else {
			//error al consumir los datos de la afectacion presupuestal definida en la solicitud.
			alertas = append(alertas, models.Alert{Code: "E_CDP002", Body: solicitud, Type: "error"})
		}
		tool.Quitar_predicados()
	}
	return alertas
}

// ValorDisponibilidadesFuenteRubroDependencia ...
// @Title ValorDisponibilidadesFuenteRubroDependencia
// @Description Obtener el valor total de las disponibilidades expediadas a una determinada dependencia por una fuente y apropiacion especifica
// @Param	idfuente	query	int	false	"id de la fuente a consultar"
// @Param	iddependencia	query	int	false	"id de la dependencia a consultar"
// @Param	idapropiacion	query	int	false	"id de la apropiacion a consultar"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ValorDisponibilidadesFuenteRubroDependencia [get]
func (c *DisponibilidadController) ValorDisponibilidadesFuenteRubroDependencia() {
	var res []interface{}
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {
				var apropiacion map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/"+strconv.Itoa(idapropiacion), &apropiacion); err == nil {
					if apropiacion != nil {
						var dependencias []map[string]interface{}
						fechaconsulta := time.Date(int(apropiacion["Vigencia"].(float64)), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")

						if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?query=DependenciaId:"+strconv.Itoa(iddependencia)+",FechaInicio__lte:"+fechaconsulta+",FechaFin__gt:"+fechaconsulta+"&sortby=Id&order=desc", &dependencias); err == nil {

							for _, dependencia := range dependencias {
								peticion := "solicitud_disponibilidad?"
								peticion = peticion + "limit=-1&query=Necesidad.FuenteReversa.FuenteFinanciamiento:" + strconv.Itoa(idfuente) + ","
								peticion = peticion + "Necesidad.FuenteReversa.Apropiacion:" + strconv.Itoa(idapropiacion) + ","
								peticion = peticion + "Necesidad.DependenciaReversa.JefeDependenciaDestino:" + strconv.Itoa(int(dependencia["Id"].(float64))) + ","
								peticion = peticion + "Expedida:true"
								var solicitudDisponibilidades []map[string]interface{}
								if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+peticion, &solicitudDisponibilidades); err == nil {
									fmt.Println(solicitudDisponibilidades)
									for _, solicitud_disponibilidad := range solicitudDisponibilidades {
										var disponibilidades []map[string]interface{}
										fmt.Println("/disponibilidad_apropiacion?query=Disponibilidad.Solicitud:" + strconv.Itoa(int(solicitud_disponibilidad["Id"].(float64))) + ",Apropiacion.Id:" + strconv.Itoa(idapropiacion))
										if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?query=Disponibilidad.Solicitud:"+strconv.Itoa(int(solicitud_disponibilidad["Id"].(float64)))+",Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamiento:"+strconv.Itoa(idfuente), &disponibilidades); err == nil {
											if disponibilidades != nil {
												for _, disponibilidad := range disponibilidades {
													res = append(res, disponibilidad)
												}
											}
										} else {
											fmt.Println("err7  ", err.Error())
											c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
										}
									}
								} else {
									fmt.Println("err6 ", err.Error())
									c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
								}
							}
							if res != nil {
								var valorcdp float64
								valorcdp = 0
								for _, row := range res {
									valorcdp = valorcdp + row.(map[string]interface{})["Valor"].(float64)
								}
								c.Data["json"] = map[string]interface{}{"valor": valorcdp}
							}
						} else {
							fmt.Println("err5 ", err.Error())
							c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
						}
					} else {
						fmt.Println("aqui")
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

// ValorDisponibilidadPorFuenteDependencia ....
func ValorDisponibilidadPorFuenteDependencia(idfuente, iddependencia, idapropiacion int) (error, map[string]interface{}) {
	var apropiacion map[string]interface{}
	var res []interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/"+strconv.Itoa(idapropiacion), &apropiacion); err == nil {
		if apropiacion != nil {
			var dependencias []map[string]interface{}
			fechaconsulta := time.Date(int(apropiacion["Vigencia"].(float64)), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")

			if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?query=DependenciaId:"+strconv.Itoa(iddependencia)+",FechaInicio__lte:"+fechaconsulta+",FechaFin__gt:"+fechaconsulta+"&sortby=Id&order=desc", &dependencias); err == nil {

				for _, dependencia := range dependencias {
					peticion := "solicitud_disponibilidad?"
					peticion = peticion + "limit=-1&query=Necesidad.FuenteReversa.FuenteFinanciamiento:" + strconv.Itoa(idfuente) + ","
					peticion = peticion + "Necesidad.FuenteReversa.Apropiacion:" + strconv.Itoa(idapropiacion) + ","
					peticion = peticion + "Necesidad.DependenciaReversa.JefeDependenciaDestino:" + strconv.Itoa(int(dependencia["Id"].(float64))) + ","
					peticion = peticion + "Expedida:true"
					var solicitudDisponibilidades []map[string]interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+peticion, &solicitudDisponibilidades); err == nil {
						for _, solicitud_disponibilidad := range solicitudDisponibilidades {
							var disponibilidades []map[string]interface{}
							if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?query=Disponibilidad.Solicitud:"+strconv.Itoa(int(solicitud_disponibilidad["Id"].(float64)))+",Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamiento:"+strconv.Itoa(idfuente), &disponibilidades); err == nil {
								if disponibilidades != nil {
									for _, disponibilidad := range disponibilidades {
										res = append(res, disponibilidad)
									}
								}
							} else {
								fmt.Println("err7  ", err.Error())
								return err, nil
							}
						}
					} else {
						fmt.Println("err6 ", err.Error())
						return err, nil
					}
				}
				if res != nil {
					var valorcdp float64
					valorcdp = 0
					for _, row := range res {
						valorcdp = valorcdp + row.(map[string]interface{})["Valor"].(float64)
					}
					return nil, map[string]interface{}{"valor": valorcdp}
				}
			} else {
				return err, nil
			}
		} else {
			return nil, nil
		}
	} else {
		return err, nil
	}
	return nil, nil
}

// AprobarAnulacionDisponibilidad ...
// @Title AprobarAnulacionDisponibilidad
// @Description aprueba la anulacion de un cdp ya sea total o parcial
// @Param	body		body 	models.AnulacionDisponibilidad	true		"body for AnulacionDisponibilidad content"
// @Success 201 {int} models.AnulacionDisponibilidad
// @Failure 403 body is empty
// @router /AprobarAnulacion [post]
func (c *DisponibilidadController) AprobarAnulacionDisponibilidad() {
	try.This(func() {
		var v models.AnulacionDisponibilidad
		var res models.Alert
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/disponibilidad/AprobarAnulacion"
			beego.Info("URL ", Urlcrud)
			if err := request.SendJson(Urlcrud, "POST", &res, &v); err == nil {
				c.Data["json"] = res
			} else {
				panic(err.Error())
			}
		} else {
			c.Data["json"] = err.Error()
		}
	}).Catch(func(e try.E) {
		beego.Info("excep: ", e)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: e, Type: "error"}
	})

	c.ServeJSON()
}

// SaldoCdp ...
// @Title SaldoCdp
// @Description create RegistroPresupuestal
// @Param	idPsql		path 	int	true		"idPsql del documento"
// @Param	rubro		path 	string	true		"código del rubro"
// @Param	fuente		query	string false		"fuente de financiamiento"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /SaldoCdp/:idPsql/:rubro [get]
func (c *DisponibilidadController) SaldoCdp() {
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
		valoresSuman := map[string]bool{"Valor": true, "TotalAnuladoRp": true}
		infoSaldo = CalcularSaldoMovimiento(rubro, fuente, "Cdp", cdpID, valoresSuman)
		c.Data["json"] = infoSaldo

	}).Catch(func(e try.E) {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: e, Type: "error"}
	})

	c.ServeJSON()
}

//AddDisponibilidadMongo ....
func AddDisponibilidadMongo(parameter ...interface{}) (err interface{}) {
	try.This(func() {
		infoDisp := parameter[0].(map[string]interface{})
		infoDisp["Vigencia"] = strconv.Itoa(int(infoDisp["Vigencia"].(float64)))
		var afectacion []map[string]interface{}
		idDisp := int(infoDisp["Id"].(float64))
		Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/disponibilidad/GetPrincDisponibilidadInfo/" + strconv.Itoa(idDisp)
		if err1 := request.GetJson(Urlcrud, &afectacion); err1 == nil {
			infoDisp["Afectacion"] = afectacion
			dateStr := infoDisp["FechaRegistro"].(string)
			t, err1 := time.Parse(time.RFC3339, dateStr)
			if err1 != nil {
				panic(err1.Error())
			}
			var resM map[string]interface{}
			infoDisp["MesRegistro"] = strconv.Itoa(int(t.Month()))
			beego.Info("Data send ", infoDisp)
			Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/Cdp"
			if err1 = request.SendJson(Urlmongo, "POST", &resM, &infoDisp); err1 == nil {
				if resM["Type"].(string) == "success" {
					err = err1
				} else {
					panic("Mongo api error")
				}
			} else {
				panic("Mongo Not Found")
			}
			//beego.Info("infoDisp ", infoDisp)
		} else {
			panic(err1.Error())
		}
	}).Catch(func(e try.E) {
		infoDisp := parameter[0].(map[string]interface{})
		idDisp := int(infoDisp["Id"].(float64))
		beego.Info("Exepc ", e)
		var resC interface{}
		Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/disponibilidad/DeleteDisponibilidadData/" + strconv.Itoa(idDisp)
		if errDelete := request.SendJson(Urlcrud, "DELETE", &resC, nil); errDelete == nil {
			fmt.Println("registro borrado correctamente")
			beego.Info("Data ", resC)
			err = e
		} else {
			fmt.Println("registro borrado correctamente")
		}
	})
	return
}

// AddAnulacionCdpMongo ....
func AddAnulacionCdpMongo(parameter ...interface{}) (err interface{}) {
	infoAnulacion := models.AnulacionDisponibilidad{}
	try.This(func() {
		infoAnulacionint := parameter[0].(map[string]interface{})
		error := formatdata.FillStruct(infoAnulacionint, &infoAnulacion)
		if err != nil {
			panic(error.Error())
		}
		dataSend := make(map[string]interface{})
		dataSend["Id"] = infoAnulacion.Id
		dataSend["MesRegistro"] = strconv.Itoa(int(infoAnulacion.FechaRegistro.Month()))
		dataSend["Vigencia"] = strconv.FormatFloat(infoAnulacion.AnulacionDisponibilidadApropiacion[0].DisponibilidadApropiacion.Disponibilidad.Vigencia, 'f', 0, 64)
		var afectacion []interface{}
		var resM map[string]interface{}
		aux := make(map[string]interface{})
		for _, data := range infoAnulacion.AnulacionDisponibilidadApropiacion {
			aux["Rubro"] = data.DisponibilidadApropiacion.Apropiacion.Rubro.Codigo
			aux["UnidadEjecutora"] = strconv.Itoa(int(data.DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora))
			aux["Valor"] = data.Valor
			dataSend["Disponibilidad"] = data.DisponibilidadApropiacion.Disponibilidad.Id
			afectacion = append(afectacion, aux)
		}
		dataSend["Afectacion"] = afectacion
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarMovimiento/AnulacionCdp"
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
		Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/anulacion_disponibilidad/" + strconv.Itoa(infoAnulacion.Id)
		infoAnulacion.EstadoAnulacion["Id"] = 2
		if errorPut := request.SendJson(Urlcrud, "PUT", &resC, &infoAnulacion); errorPut == nil {
			fmt.Println("actualizado correctamente")
			beego.Info("Data ", resC)
		} else {
			beego.Info("Error ", errorPut)
		}

	})
	return
}

//CalcularSaldoMovimiento ....
func CalcularSaldoMovimiento(rubro, fuente, tipo string, cdpID int, valoresSuman map[string]bool) (res map[string]float64) {
	var saldo float64
	urlMongo := ""
	urlMongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro_apropiaciones/SaldoMovimiento/" + strconv.Itoa(cdpID) + "/" + rubro + "/" + tipo
	if fuente != "" {
		urlMongo += "?fuente=" + fuente
	}
	if err := request.GetJson(urlMongo, &res); err != nil {
		beego.Info(err.Error())
		panic("Mongo API Service Error")
	} else {
		for key, value := range res {
			if valoresSuman[key] {
				beego.Info("suma ", key)
				saldo += value
			} else {
				beego.Info("resta ", key)
				saldo -= value
			}

			// beego.Info("Saldo ", saldo)
			// switch tipoMovimiento := key; tipoMovimiento {
			//
			// //rp
			// case valoresSuman:
			// 	if tipo == "Cdp" {
			//
			// 	} else {
			//
			// 	}
			// 	beego.Info("suma ", tipoMovimiento)
			// 	saldo += value
			// default:
			//
			// 	beego.Info("resta ", tipoMovimiento)
			// 	saldo -= value
			// }

		}

	}

	res["saldo"] = saldo
	return
}
