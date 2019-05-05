package models


//DestinoDisponibilidad ...
type DestinoDisponibilidad struct {
	Id          int    `orm:"column(id);pk"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion)"`
}
