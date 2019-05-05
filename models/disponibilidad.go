package models

import (
	"time"
)

// Disponibilidad ...
type Disponibilidad struct {
	Id                           int                             `orm:"auto;column(id);pk"`
	Vigencia                     float64                         `orm:"column(vigencia)"`
	NumeroDisponibilidad         float64                         `orm:"column(numero_disponibilidad);null"`
	Responsable                  int                             `orm:"column(responsable);null"`
	FechaRegistro                time.Time                       `orm:"column(fecha_registro);type(date);null"`
	Estado                       *EstadoDisponibilidad           `orm:"column(estado);rel(fk)"`
	Solicitud                    int                             `orm:"column(solicitud)"`
	DisponibilidadApropiacion    []*DisponibilidadApropiacion    `orm:"reverse(many)"`
	DisponibilidadProcesoExterno []*DisponibilidadProcesoExterno `orm:"reverse(many)"`
	DatosNecesidad               *Necesidad
}

//  DisponibilidadProcesoExterno ...
type DisponibilidadProcesoExterno struct {
	Id                 int                 `orm:"auto;column(id);pk"`
	TipoDisponibilidad *TipoDisponibilidad `orm:"column(tipo_disponibilidad);rel(fk)"`
	ProcesoExterno     int                 `orm:"column(proceso_externo)"`
	Disponibilidad     *Disponibilidad     `orm:"column(disponibilidad);rel(fk)"`
}

//TipoDisponibilidad ...
type TipoDisponibilidad struct {
	Id          int     `orm:"auto;column(id);pk"`
	Nombre      string  `orm:"column(nombre)"`
	Descripcion string  `orm:"column(descripcion);null"`
	Activo      bool    `orm:"column(activo)"`
	NumeroOrden float64 `orm:"column(numero_orden);null"`
}

// InfoSolDisp ...
type InfoSolDisp struct {
	SolicitudDisponibilidad SolicitudDisponibilidad
	DependenciaSolicitante  Dependencia
	DependenciaDestino      Dependencia
	Responsable             int
}
