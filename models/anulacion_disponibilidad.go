package models

import (
	"time"
)

//AnulacionDisponibilidad ...
type AnulacionDisponibilidad struct {
	Id                                 int                                   `orm:"column(id);pk;auto"`
	Consecutivo                        int                                   `orm:"column(consecutivo)"`
	Motivo                             string                                `orm:"column(motivo)"`
	FechaRegistro                      time.Time                             `orm:"column(fecha_registro);type(date)"`
	TipoAnulacion                      map[string]interface{}            `orm:"column(tipo_anulacion);rel(fk)"`
	EstadoAnulacion                    map[string]interface{}                      `orm:"column(estado_anulacion);rel(fk)"`
	JustificacionRechazo               string                                `orm:"column(justificacion_rechazo);null"`
	Responsable                        int                                   `orm:"column(responsable)"`
	Solicitante                        int                                   `orm:"column(solicitante)"`
	Expidio                            int                                   `orm:"column(expidio)"`
	AnulacionDisponibilidadApropiacion []*AnulacionDisponibilidadApropiacion `orm:"reverse(many)"`
}
