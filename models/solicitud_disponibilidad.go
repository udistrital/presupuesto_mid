package models

import (
	"time"
)

// SolicitudDisponibilidad ...
type SolicitudDisponibilidad struct {
	Id             int        `orm:"column(id);pk"`
	Numero         int        `orm:"column(numero)"`
	Vigencia       float64    `orm:"column(vigencia)"`
	FechaSolicitud time.Time  `orm:"column(fecha_solicitud);type(date)"`
	Necesidad      *Necesidad `orm:"column(necesidad);rel(fk)"`
	Expedida       bool       `orm:"column(expedida)"`
}
