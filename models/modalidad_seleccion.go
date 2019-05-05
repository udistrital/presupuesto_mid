package models


// ModalidadSeleccion ...
type ModalidadSeleccion struct {
	Id          int    `orm:"column(id);pk"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion);null"`
}
