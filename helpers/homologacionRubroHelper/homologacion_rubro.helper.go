package homologacionRubroHelper

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/presupuesto_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// CreateRubroHomologado ... send finantial information to presupuestoApiService for some internal operation's
func CreateRubroHomologado(rubro interface{}) (response map[string]interface{}) {

	var (
		urlCrud String
	)
	res := make(map[string]interface{})
	// rubroHomologado := make(map[string]interface{})
	urlCrud = beego.AppConfig.String("presupuestoApiService") + "/rubro_homologado"
	err := request.SendJson(urlCrud, "POST", &res, rubro)

	if err == nil && res["Type"] != nil && res["Type"].(string) == "success" {
		res["Body"] = rubro
	} else {
		panic(helpers.InternalErrorMessage())
	}

	response = res

	return

}
