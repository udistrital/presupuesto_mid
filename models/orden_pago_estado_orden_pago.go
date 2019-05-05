package models

import (
	"time"


)

type OrdenPagoEstadoOrdenPago struct {
	Id              int              `orm:"column(id);pk;auto"`
	OrdenPago       *OrdenPago       `orm:"column(orden_pago);rel(fk)"`
	EstadoOrdenPago *EstadoOrdenPago `orm:"column(estado_orden_pago);rel(fk)"`
	FechaRegistro   time.Time        `orm:"column(fecha_registro);type(date)"`
	Usuario         int              `orm:"column(usuario);null"`
}
