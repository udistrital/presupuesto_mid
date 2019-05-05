package models

// Rubro ...
type Rubro struct {
	Id              int    `orm:"auto;column(id);pk"`
	Entidad         int    `orm:"column(entidad)"`
	Codigo          string `orm:"column(codigo)"`
	Descripcion     string `orm:"column(descripcion);null"`
	UnidadEjecutora int16  `orm:"column(unidad_ejecutora)"`
	Nombre          string `orm:"column(nombre);null"`
}
