package mongoProcess

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/manucorporat/try"
	"github.com/udistrital/presupuesto_mid/controllers"
	"github.com/udistrital/presupuesto_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/optimize"
)

func sendDisponibilidadInfoToMongo(ctx *context.Context) {
	try.This(func() {
		serviceResponse := ctx.Input.Data()["json"].([]models.Alert)
		var params []interface{}
		for _, data := range serviceResponse {
			if data.Type == "success" {
				info := data.Body.(map[string]interface{})
				params = append(params, info)
				work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddDisponibilidadMongo)}
				// Push the work onto the queue.
				optimize.WorkQueue <- work
				beego.Info("Job Queued!")
			}

		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

func sendAnulacionRpInfoToMongo(ctx *context.Context) {
	try.This(func() {
		data := ctx.Input.Data()["json"].(models.Alert)
		var params []interface{}
		if data.Type == "success" {
			info := data.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddAnulacionRpMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		}

	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

func sendAnulacionCdpInfoToMongo(ctx *context.Context) {
	try.This(func() {
		data := ctx.Input.Data()["json"].(models.Alert)
		beego.Info("Data job ", data)
		var params []interface{}
		if data.Type == "success" {
			info := data.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddAnulacionCdpMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		}

	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

func sendRpInfoToMongo(ctx *context.Context) {
	try.This(func() {
		serviceResponse := ctx.Input.Data()["json"].([]models.Alert)
		var params []interface{}
		for _, data := range serviceResponse {
			if data.Type == "success" {
				info := data.Body
				params = append(params, info)
				work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddRpMongo)}
				// Push the work onto the queue.
				optimize.WorkQueue <- work
				beego.Info("Job Queued!")
			}

		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})
}

func sendOpInfoToMongo(ctx *context.Context) {
	try.This(func() {
		data := ctx.Input.Data()["json"].(models.Alert)
		//beego.Info(data.Body)
		dataaux := data.Body.(map[string]interface{})
		var params []interface{}
		if data.Type == "success" && dataaux["NuevoEstado"].(map[string]interface{})["Id"].(float64) == 7 {
			info := data.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddOpMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		} else {
			beego.Info("Job Not Work! Id !=7")
		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})
}

func sendMovimientoInfoToMongo(ctx *context.Context) {
	try.This(func() {
		var serviceResponse []models.Alert
		err := formatdata.FillStruct(ctx.Input.Data()["json"], &serviceResponse)
		if err != nil {
			panic(err.Error())
		}
		var params []interface{}
		data := serviceResponse[len(serviceResponse)-1]
		// for _, data := range serviceResponse {
		if data.Type == "success" {
			info := data.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddMovimientoApropiacionMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		}

		// }
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})
}
func sendFuenteFinanciamientoInfoToMongo(ctx *context.Context) {
	try.This(func() {
		var serviceResponse models.Alert
		err := formatdata.FillStruct(ctx.Input.Data()["json"], &serviceResponse)
		if err != nil || serviceResponse.Body == nil {
			panic(err.Error())
		}
		beego.Info("Job Fuente start!! ", ctx.Input.Data()["json"])
		var params []interface{}
		if serviceResponse.Type == "success" {
			info := serviceResponse.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddFuenteFinanciamientoMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		}

	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

func sendModificacionFuenteFinanciamientoInfoToMongo(ctx *context.Context) {
	try.This(func() {
		var serviceResponse models.Alert
		err := formatdata.FillStruct(ctx.Input.Data()["json"], &serviceResponse)
		if err != nil || serviceResponse.Body == nil {
			panic(err.Error())
		}
		beego.Info("Job Fuente start!! ", ctx.Input.Data()["json"])
		var params []interface{}
		if serviceResponse.Type == "success" {
			info := serviceResponse.Body
			params = append(params, info)
			work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddModificacionFuenteFinanciamientoMongo)}
			// Push the work onto the queue.
			optimize.WorkQueue <- work
			beego.Info("Job Queued!")
		}

	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

// PresupuestoMongoJobInit ...
func PresupuestoMongoJobInit() { //inicia los escuchadores de los procesos que deben guardarse simultaneamente en postgres y mongo
	optimize.StartDispatcher(1, 200)

	beego.InsertFilter("/v1/disponibilidad/ExpedirDisponibilidad", beego.AfterExec, sendDisponibilidadInfoToMongo, false)
	beego.InsertFilter("/v1/registro_presupuestal/CargueMasivoPr", beego.AfterExec, sendRpInfoToMongo, false)
	beego.InsertFilter("/v1/registro_presupuestal/AprobarAnulacion", beego.AfterExec, sendAnulacionRpInfoToMongo, false)
	beego.InsertFilter("/v1/disponibilidad/AprobarAnulacion", beego.AfterExec, sendAnulacionCdpInfoToMongo, false)
	beego.InsertFilter("/v1/movimiento_apropiacion/AprobarMovimietnoApropiacion/*", beego.AfterExec, sendMovimientoInfoToMongo, false)
	beego.InsertFilter("/v1/fuente_financiamiento/RegistrarFuente", beego.AfterExec, sendFuenteFinanciamientoInfoToMongo, false)
	beego.InsertFilter("/v1/fuente_financiamiento/RegistrarModificacionFuente", beego.AfterExec, sendModificacionFuenteFinanciamientoInfoToMongo, false)
	beego.InsertFilter("/v1/orden_pago/WorkFlowOrdenPago", beego.AfterExec, sendOpInfoToMongo, false)
}
