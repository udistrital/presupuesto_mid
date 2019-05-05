package models

// FuenteFinanciacion ...
type FuenteFinanciacion struct {
	Id                       int    `orm:"column(id);pk;auto"`
	Descripcion              string `orm:"column(descripcion);null"`
	Nombre                   string `orm:"column(nombre)"`
	Codigo                   string `orm:"column(codigo)"`
	TipoFuenteFinanciamiento map[string]interface{}
}
