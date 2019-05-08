package controllers

import (
	"strings"
	"encoding/json"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/astaxie/beego"
)

// IngresoSinSituacionFondosController operations for IngresoSinSituacionFondos
type IngresoSinSituacionFondosController struct {
	beego.Controller
}

type estadoIngreso struct{
	IngresoSinSituacionFondos interface{}
	EstadoIngresoSinSituacionFondos interface{}
	Activo bool
}

// URLMapping ...
func (c *IngresoSinSituacionFondosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Ingreso_sin_situacion_fondos
// @Param	body		body 	interface{}	"body for Ingreso_sin_situacion_fondos content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router / [post]
func (c *IngresoSinSituacionFondosController) Post() {
	var v map[string]interface{}
	var respuesta map[string]interface{}
	var estadoResp []interface{}
	var respuestaCreaEst map[string]interface{}


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		ingresoSinSituacionFondos:=v["IngresoSinSituacionFondos"]
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/ingreso_sin_situacion_fondos", "POST", &respuesta, ingresoSinSituacionFondos); err == nil {
			if (strings.Compare(respuesta["Type"].(string),"success")==0){
				 if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/estado_ingreso_sin_situacion_fondos?query=numeroOrden:1", &estadoResp); err == nil {
					 estadoIng := &estadoIngreso {IngresoSinSituacionFondos:respuesta["Body"],EstadoIngresoSinSituacionFondos:estadoResp[0],Activo:true}
					 if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/ingreso_sin_situacion_fondos_estado", "POST", &respuestaCreaEst, estadoIng); err == nil {
						if (strings.Compare(respuestaCreaEst["Type"].(string),"success")==0){
							alert:=models.Alert{Type:"success",Code:"S_543",Body:respuesta["Body"]}
							c.Data["json"] = alert
							beego.Error(c.Data["json"])
							c.Ctx.Output.SetStatus(201)
						}
					 }else{
						 alert:=models.Alert{Type:"error",Code:"E_0458",Body:err}
						 c.Data["json"] = alert
					 }
				 }else{
					 alert:=models.Alert{Type:"error",Code:"E_0458",Body:err}
					 c.Data["json"] = alert
				 }
			}
		} else {
			alert:=models.Alert{Type:"error",Code:"E_0458",Body:err}
			beego.Error(err.Error());
			c.Data["json"] = alert

		}
	} else {
		beego.Error(err.Error());
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}


// ChangeState ...
// @Title  ChangeState
// @Description Change State from a income without finance situation
// @Param	body	body 	interface{}	"body for Ingreso_sin_situacion_fondos content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router ChangeState/ [post]
func (c *IngresoSinSituacionFondosController) ChangeState() {
	defer c.ServeJSON()
	var v interface{}
	var respuesta map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/ingreso_sin_situacion_fondos_estado/ChangeExistingStates", "POST", &respuesta, v); err == nil {
			if (strings.Compare(respuesta["Type"].(string),"success")==0){
				alert:=models.Alert{Type:"success",Code:"S_543",Body:respuesta["Body"]}
				c.Data["json"] = alert
				beego.Error(c.Data["json"])
				c.Ctx.Output.SetStatus(201)
			}else{
				alert:=models.Alert{Type:"error",Code:"E_0458",Body:err}
				c.Data["json"] = alert
			}
		}
	}
}

// GetOne ...
// @Title GetOne
// @Description get Ingreso_sin_situacion_fondos by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Ingreso_sin_situacion_fondos
// @Failure 403 :id is empty
// @router /:id [get]
func (c *IngresoSinSituacionFondosController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Ingreso_sin_situacion_fondos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Ingreso_sin_situacion_fondos
// @Failure 403
// @router / [get]
func (c *IngresoSinSituacionFondosController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Ingreso_sin_situacion_fondos
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Ingreso_sin_situacion_fondos	true		"body for Ingreso_sin_situacion_fondos content"
// @Success 200 {object} models.Ingreso_sin_situacion_fondos
// @Failure 403 :id is not int
// @router /:id [put]
func (c *IngresoSinSituacionFondosController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Ingreso_sin_situacion_fondos
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *IngresoSinSituacionFondosController) Delete() {

}
