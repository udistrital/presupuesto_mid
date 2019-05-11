package models

import "time"

type Producto struct {
	Id            int              `orm:"column(id);pk;auto"`
	Nombre        string           `orm:"column(nombre)"`
	Descripcion   string           `orm:"column(descripcion);null"`
	FechaRegistro time.Time        `orm:"column(fecha_registro);type(date)"`
	ProductoRubro []*ProductoRubro `orm:"reverse(many)"`
}
