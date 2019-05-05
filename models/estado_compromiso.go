package models

//EstadoCompromiso ...
type EstadoCompromiso struct {
	Id          int    `orm:"column(id);pk;auto"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion);null"`
}
