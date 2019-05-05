package models

type EstadoOrdenPago struct {
	Id                int     `orm:"column(id);pk"`
	Nombre            string  `orm:"column(nombre)"`
	Activo            bool    `orm:"column(activo)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	CodigoAbreviacion string  `orm:"column(codigo_abreviacion);null"`
	NumeroOrden       float64 `orm:"column(numero_orden);null"`
}
