package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// IngresoController operations for Ingreso
type IngresoController struct {
	beego.Controller
}

// URLMapping ...
func (c *IngresoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Ingreso
// @Param	body		body 	models.Ingreso	true		"body for Ingreso content"
// @Success 201 {object} models.Ingreso
// @Failure 403 body is empty
// @router / [post]
func (c *IngresoController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Ingreso by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Ingreso
// @Failure 403 :id is empty
// @router /:id [get]
func (c *IngresoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Ingreso
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Ingreso
// @Failure 403
// @router / [get]
func (c *IngresoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Ingreso
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Ingreso	true		"body for Ingreso content"
// @Success 200 {object} models.Ingreso
// @Failure 403 :id is not int
// @router /:id [put]
func (c *IngresoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Ingreso
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *IngresoController) Delete() {

}

// AprobacionPresupuestalIngreso ...
// @Title AprobacionPresupuestalIngreso
// @Description change estate to a income
// @Param	body		body 	interface{}	true		"body for Ingresocontent"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router AprobacionPresupuestalIngreso/ [post]
func (c *IngresoController) AprobacionPresupuestalIngreso() {
	defer c.ServeJSON()
	var v interface{}
	var respuesta interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/ingreso/AprobacionPresupuestalIngreso", "POST", &respuesta, v); err == nil {
			c.Data["json"] = respuesta
		} else {
			beego.Error(err.Error())
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}

	} else {
		beego.Error(err.Error())
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}
