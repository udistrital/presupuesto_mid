package models

import "time"

type ProductoRubro struct {
	Id                int       `orm:"column(id);pk;auto"`
	Rubro             *Rubro    `orm:"column(rubro);rel(fk)"`
	Producto          *Producto `orm:"column(producto);rel(fk)"`
	ValorDistribucion float64   `orm:"column(valor_distribucion)"`
	Activo            bool      `orm:"column(activo);default(true)"`
	FechaRegistro     time.Time `orm:"column(fecha_registro);auto_now_add"`
}
