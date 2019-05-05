package models

import (
	"time"
)

//AnulacionRegistroPresupuestal ...
type AnulacionRegistroPresupuestal struct {
	Id                   int                        `orm:"auto;column(id);pk"`
	Consecutivo          int                        `orm:"column(consecutivo)"`
	Motivo               string                     `orm:"column(motivo)"`
	FechaRegistro        time.Time                  `orm:"column(fecha_registro);type(date)"`
	TipoAnulacion        map[string]interface{} `orm:"column(tipo_anulacion);rel(fk)"`
	EstadoAnulacion      map[string]interface{}          `orm:"column(estado_anulacion);rel(fk)"`
	JustificacionRechazo string                     `orm:"column(justificacion_rechazo);null"`
	Responsable          int                        `orm:"column(responsable)"`
	Solicitante          int                        `orm:"column(solicitante)"`
	Expidio              int                        `orm:"column(expidio)"`

	AnulacionRegistroPresupuestalDisponibilidadApropiacion []*AnulacionRegistroPresupuestalDisponibilidadApropiacion `orm:"reverse(many)"`
}
