package rubroHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/helpers"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddRubro ... send finantial information to mongo for some internal operation's
func AddRubro(data models.RubroRubro) (response map[string]interface{}) {
	var (
		urlCrud  string
		padre    bool
		sendData interface{}
	)
	res := make(map[string]interface{})

	defer func() {

		if r := recover(); r != nil {
			go func() {
				defer func() {
					if routineRecv := recover(); routineRecv != nil {
						beego.Error(helpers.RoutineError())
					}
				}()
				if resul, e := res["Body"].(map[string]interface{}); e {
					urlCrud = beego.AppConfig.String("presupuestoApiService")

					if padre {
						urlCrud += "rubro/" + strconv.Itoa(int(resul["Id"].(float64)))
					} else {
						urlCrud += "rubro_rubro/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64)))
					}
					var deleteRes interface{}
					if errorDelete := request.SendJson(urlCrud, "DELETE", &deleteRes, nil); errorDelete == nil {
						beego.Error(helpers.ExternalAPIErrorMessage())
					} else {
						beego.Error("Error delete ", errorDelete)
					}
				} else {
					beego.Error(helpers.InternalErrorMessage())
				}

			}()
			panic(helpers.InternalErrorMessage())
		}

	}()

	urlCrud = beego.AppConfig.String("presupuestoApiService")
	if data.RubroPadre != nil {
		urlCrud += "rubro_rubro"
		sendData = data
		padre = false
	} else if data.RubroHijo != nil {
		sendData = data.RubroHijo
		urlCrud += "rubro"
		padre = true
	} else {
		beego.Error("Null Data")
		panic(helpers.InternalErrorMessage())
	}

	err := request.SendJson(urlCrud, "POST", &res, &sendData)

	if err == nil && res["Type"] != nil && res["Type"].(string) == "success" {

		if !padre {
			sendData = res["Body"].(map[string]interface{})
		} else {
			sendData = map[string]interface{}{"RubroHijo": res["Body"].(map[string]interface{}), "RubroPadre": map[string]interface{}{}}
			res["Body"] = sendData

		}

		urlMongo := beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/registrarRubro"
		var resMongo map[string]interface{}
		err = request.SendJson(urlMongo, "POST", &resMongo, &sendData)
		if err != nil || resMongo["Type"] == nil || resMongo["Type"].(string) == "error" {
			beego.Error("err", err, "res", resMongo)
			panic(helpers.ExternalAPIErrorMessage())
		}

	} else {
		panic(helpers.InternalErrorMessage())
	}

	response = res

	return

}

func DeleteRubro(id int) map[string]interface{} {

	idStr := strconv.Itoa(id)
	urlcrud := beego.AppConfig.String("presupuestoApiService") + "/rubro/" + idStr
	res := make(map[string]interface{})
	defer func() {
		if r := recover(); r != nil {
			beego.Error(r)
			panic(helpers.InternalErrorMessage())

		}
	}()

	if err := request.SendJson(urlcrud, "DELETE", &res, nil); err == nil {
		if res["Type"].(string) == "success" {
			var resMg map[string]interface{}
			urlmongo := beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/eliminarRubro/" + idStr
			if err = request.SendJson(urlmongo, "DELETE", &resMg, nil); err != nil {
				beego.Error(err.Error())
				panic(helpers.ExternalAPIErrorMessage())
			}

		} else {
			beego.Error(err.Error())
			panic(helpers.InternalErrorMessage())
		}
	} else {
		beego.Error(err.Error())
		panic(helpers.InternalErrorMessage())
	}
	return res
}
