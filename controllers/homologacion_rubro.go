package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/helpers/homologacionRubroHelper"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/resposeformat"
)

// HomologacionRubroController operations for HomologacionRubroController
type HomologacionRubroController struct {
	beego.Controller
}

// URLMapping ...
func (c *HomologacionRubroController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Homologacion_rubro
// @Param	body		body 	models.Homologacion_rubro	true		"body for Homologacion_rubro content"
// @Success 201 {object} models.Homologacion_rubro
// @Failure 403 body is empty
// @router / [post]
func (c *HomologacionRubroController) Post() {

}

// CreateRubroHomologado ...
// @Title CreateRubroHomologado
// @Description create homologate category for an organization
// @Param	body		body 	interface	true		"body for Homologacion_rubro content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router /CreateRubroHomologado [post]
func (c *HomologacionRubroController) CreateRubroHomologado() {
	var rubro interface{}
	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			resposeformat.SetResponseFormat(&c.Controller, r, "E_0458", 500)
		}
	}()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rubro); err == nil {
		response := homologacionRubroHelper.CreateRubroHomologado(rubro)
		beego.Debug(response)
		resposeformat.SetResponseFormat(&c.Controller, response["Body"], response["Code"].(string), 200)
	} else {
		resposeformat.SetResponseFormat(&c.Controller, nil, "E_0458", 500)
	}

}

// CreateHomologacion ...
// @Title CreateHomologacion
// @Description create homologate category for an organization
// @Param	body		body 	interface	true		"body for Homologacion_rubro content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router /CreateHomologacion [post]
func (c *HomologacionRubroController) CreateHomologacion() {
	defer c.ServeJSON()
	var rubroHomologadorubro interface{}
	var response map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rubroHomologadorubro); err == nil {
		beego.Info(rubroHomologadorubro)
		if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado_rubro", "POST", &response, rubroHomologadorubro); err == nil {
			beego.Error("rubro homologado rubro", response)
			if strings.Compare(response["Type"].(string), "success") == 0 {
				c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: response["Body"]}
				c.Ctx.Output.SetStatus(201)
			} else {
				c.Data["json"] = models.Alert{Type: response["Type"].(string), Code: response["Code"].(string), Body: response["Body"]}
			}
		} else {
			beego.Error("Error", err)
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}
	} else {
		beego.Error("Error", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}

}

// GetOne ...
// @Title GetOne
// @Description get Homologacion_rubro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Homologacion_rubro
// @Failure 403 :id is empty
// @router /:id [get]
func (c *HomologacionRubroController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Homologacion_rubro
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Homologacion_rubro
// @Failure 403
// @router / [get]
func (c *HomologacionRubroController) GetAll() {
	defer c.ServeJSON()
	var rubrosHomol []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var urlCrud string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("query"); r != "" {
		query = r
	}
	urlCrud = beego.AppConfig.String("presupuestoApiService") + "/rubro_homologado"
	if err := request.GetJson(urlCrud+"?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &rubrosHomol); err == nil {
		if rubrosHomol != nil {
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(rubrosHomol...)
			chRubros := optimize.Digest(done, models.GetOrganizacionRubroHomologado, resch, nil)
			for rubroHomologado := range chRubros {
				beego.Error("rubro homologado", rubroHomologado)
				respuesta = append(respuesta, rubroHomologado.(map[string]interface{}))
			}
			c.Data["json"] = respuesta
		}
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

// Put ...
// @Title Put
// @Description update the Homologacion_rubro
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Homologacion_rubro	true		"body for Homologacion_rubro content"
// @Success 200 {object} models.Homologacion_rubro
// @Failure 403 :id is not int
// @router /:id [put]
func (c *HomologacionRubroController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Homologacion_rubro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *HomologacionRubroController) Delete() {

}

// GetAllRubrosHomologado ...
// @Title GetAllRubrosHomologado
// @Description get rubro_Homologacion_rubro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} interface{}
// @Failure 403 :id is empty
// @router /GetAllRubrosHomologado/:id [get]
func (c *HomologacionRubroController) GetAllRubrosHomologado() {
	defer c.ServeJSON()
	var respuesta []map[string]interface{}

	idStr := c.Ctx.Input.Param(":id")
	var rubrosHomolRubro []interface{}

	beego.Error("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro_homologado_rubro/?query=Rubro.Id:" + idStr)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado_rubro/?query=Rubro.Id:"+idStr, &rubrosHomolRubro); err == nil {
		beego.Error(rubrosHomolRubro)

		if rubrosHomolRubro != nil {
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(rubrosHomolRubro...)
			chRubros := optimize.Digest(done, models.GetOrganizacionRubroHomologado, resch, nil)
			for rubroHomologado := range chRubros {
				beego.Error("rubro homologado", rubroHomologado)
				respuesta = append(respuesta, rubroHomologado.(map[string]interface{}))
			}
			c.Data["json"] = respuesta
		}
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

// GetHomologationNumberRubro ...
// @Title Get Homologation Number Rubro
// @Description get the number of record in crud from RubroHOmologacionRubro for a rubro
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} interface{}
// @Failure 403 :id is empty
// @router /GetHomologationNumberRubro/:id [get]
func (c *HomologacionRubroController) GetHomologationNumberRubro() {
	defer c.ServeJSON()
	var respuesta map[string]interface{}
	idStr := c.Ctx.Input.Param(":id")

	beego.Error("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro_homologado/" + idStr)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado/GetRecordsNumberRubroHomologadoRubroById/"+idStr, &respuesta); err == nil {
		c.Data["json"] = respuesta
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

// GetArbolRubrosHomologado ...
// @Title Get Arbol Rubros Homologado
// @Description get tree from homologated items
// @Param	idEntidad	path 	string	true		"Id entidad para la cual se quieren consultar rubros"
// @Param	idPadre		path 	string	true		"Id del padre para armar la rama default todos"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GetArbolRubrosHomologado [get]
func (c *HomologacionRubroController) GetArbolRubrosHomologado() {
	defer c.ServeJSON()
	var respuesta []map[string]interface{}

	idEntidad, err := c.GetInt("idEntidad")
	idPadre, _ := c.GetInt("idPadre")
	beego.Error("entidad ", idEntidad, "padre ", idPadre)

	if err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado/ArbolRubros/?idEntidad="+strconv.Itoa(idEntidad)+"&idPadre="+strconv.Itoa(idPadre), &respuesta); err == nil {
			c.Data["json"] = respuesta
		} else {
			beego.Error("error", err)
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}
	} else {
		beego.Error("error", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err.Error()}
	}
}

// GetHomologationNumberEntity ...
// @Title Get Homologation Number getting entity
// @Description get the number of record given a Entity
// @Param	idEntidad		path 	string	true		"The key for staticblock"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GetHomologationNumberEntity [get]
func (c *HomologacionRubroController) GetHomologationNumberEntity() {
	defer c.ServeJSON()
	var respuesta map[string]interface{}
	idEntidad, err := c.GetInt("idEntidad")
	beego.Error("entidad ", idEntidad, "error", err)
	beego.Error("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro_homologado/?idEntidad=" + strconv.Itoa(idEntidad))
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado/GetRecordsNumberByEntity/?idEntidad="+strconv.Itoa(idEntidad), &respuesta); err == nil {
		c.Data["json"] = respuesta
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}
