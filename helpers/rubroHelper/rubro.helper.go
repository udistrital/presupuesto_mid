package rubroHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// AddRubro ... send finantial information to mongo for some internal operation's
func AddRubro(data map[string]interface{}) {
	var (
		urlCrud string
		rubro   interface{}
	)
	res := make(map[string]interface{})

	defer func() {
		resul := res["Body"].(map[string]interface{})
		ue := resul["RubroHijo"].(map[string]interface{})["UnidadEjecutora"].(float64)
		urlCrud = urlCrud + "/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64))) + "/" + strconv.Itoa(int(ue))
		if errorDelete := request.SendJson(urlCrud, "DELETE", &data, nil); errorDelete == nil {
			panic(helpers.APINotFoundErrorMessage())
		} else {
			beego.Error("Error delete ", errorDelete)
			panic(helpers.ExternalAPIErrorMessage())
		}
	}()

	urlCrud = beego.AppConfig.String("presupuestoApiService")
	if data["RubroPadre"] != nil {
		urlCrud += "rubro_rubro"
		rubro = data
	} else if data["RubroHijo"] != nil {
		rubro = data["RubroHijo"]
		urlCrud += "rubro"
	} else {
		panic(helpers.InternalErrorMessage())
	}

	err := request.SendJson(urlCrud, "POST", &res, &rubro)
	if err == nil && res["Type"] != nil && res["Type"].(string) == "success" {
		urlMongo := beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro/registrarRubro"
		var resMongo map[string]interface{}
		sendData := res["Body"].(map[string]interface{})
		err = request.SendJson(urlMongo, "POST", &resMongo, &sendData)
		if err != nil || res["Type"] == nil || res["Type"].(string) == "error" {
			panic(helpers.ExternalAPIErrorMessage())
		}
	} else {
		panic(helpers.InternalErrorMessage())
	}

}
