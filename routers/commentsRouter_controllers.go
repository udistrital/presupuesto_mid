package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"],
        beego.ControllerComments{
            Method: "AprobacionAsignacionInicial",
            Router: `/AprobacionAsignacionInicial/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"],
        beego.ControllerComments{
            Method: "Aprobado",
            Router: `/Aprobado`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionController"],
        beego.ControllerComments{
            Method: "InformacionAsignacionInicial",
            Router: `/InformacionAsignacionInicial/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"],
        beego.ControllerComments{
            Method: "ValorMovimientoFuente",
            Router: `/ValorMovimientoFuente`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"],
        beego.ControllerComments{
            Method: "ValorMovimientoFuenteLista",
            Router: `/ValorMovimientoFuenteLista`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"],
        beego.ControllerComments{
            Method: "ValorMovimientoFuenteListaFunc",
            Router: `/ValorMovimientoFuenteListaFunc`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:AprobacionFuenteController"],
        beego.ControllerComments{
            Method: "ValorMovimientoFuenteTraslado",
            Router: `/ValorMovimientoFuenteTraslado`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id/:valor/:vigencia`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"],
        beego.ControllerComments{
            Method: "ArbolApropiaciones",
            Router: `/ArbolApropiaciones/:unidadEjecutora/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:ApropiacionController"],
        beego.ControllerComments{
            Method: "SaldoApropiacion",
            Router: `/SaldoApropiacion/:rubro/:unidadEjecutora/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "AprobarAnulacionDisponibilidad",
            Router: `/AprobarAnulacion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "AprobarAnulacion",
            Router: `/AprobarAnulacion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "ExpedirDisponibilidad",
            Router: `/ExpedirDisponibilidad`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "SaldoCdp",
            Router: `/SaldoCdp/:idPsql/:rubro`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "ValorDisponibilidadesFuenteRubroDependencia",
            Router: `/ValorDisponibilidadesFuenteRubroDependencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "DisponibilidadByNecesidad",
            Router: `DisponibilidadByNecesidad/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "ListaDisponibilidades",
            Router: `ListaDisponibilidades/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "InfoSolicitudDisponibilidadByID",
            Router: `SolicitudById/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:DisponibilidadController"],
        beego.ControllerComments{
            Method: "InfoSolicitudDisponibilidad",
            Router: `Solicitudes/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:FuenteFinanciamientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:FuenteFinanciamientoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/RegistrarFuente`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:FuenteFinanciamientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:FuenteFinanciamientoController"],
        beego.ControllerComments{
            Method: "RegistrarModificacionFuente",
            Router: `/RegistrarModificacionFuente`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "CreateHomologacion",
            Router: `/CreateHomologacion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "CreateRubroHomologado",
            Router: `/CreateRubroHomologado`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetAllRubrosHomologado",
            Router: `/GetAllRubrosHomologado/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetArbolRubrosHomologado",
            Router: `/GetArbolRubrosHomologado`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetHomologationNumberEntity",
            Router: `/GetHomologationNumberEntity`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:HomologacionRubroController"],
        beego.ControllerComments{
            Method: "GetHomologationNumberRubro",
            Router: `/GetHomologationNumberRubro/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:MovimientoApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:MovimientoApropiacionController"],
        beego.ControllerComments{
            Method: "AprobarMovimietnoApropiacion",
            Router: `/AprobarMovimietnoApropiacion/:unidadEjecutora/:vigencia`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:MovimientoApropiacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:MovimientoApropiacionController"],
        beego.ControllerComments{
            Method: "ComprobarMovimientoApropiacion",
            Router: `/ComprobarMovimientoApropiacion/:unidadEjecutora/:vigencia`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "AprobarAnulacion",
            Router: `/AprobarAnulacion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "CargueMasivoPr",
            Router: `/CargueMasivoPr`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "GetSolicitudesRp",
            Router: `/GetSolicitudesRp/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "GetSolicitudesRpByID",
            Router: `/GetSolicitudesRpById/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "ListaNecesidadesByRp",
            Router: `/ListaNecesidadesByRp/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "SaldoRp",
            Router: `/SaldoRp/:idPsql/:rubro`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "SolicitudesRpByDependencia",
            Router: `/SolicitudesRpByDependencia/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RegistroPresupuestalController"],
        beego.ControllerComments{
            Method: "ListaRp",
            Router: `ListaRp/:vigencia`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"],
        beego.ControllerComments{
            Method: "ArbolRubros",
            Router: `/ArbolRubros/:unidadEjecutora`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"],
        beego.ControllerComments{
            Method: "EliminarRubro",
            Router: `/EliminarRubro/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/presupuesto_mid/controllers:RubroController"],
        beego.ControllerComments{
            Method: "RegistrarRubro",
            Router: `/RegistrarRubro/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
