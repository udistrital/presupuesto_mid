package models

// TipoCompromisoTesoral ...
type TipoCompromisoTesoral struct {
	Id                  int                  `orm:"column(id);pk;auto"`
	Nombre              string               `orm:"column(nombre)"`
	EstadoActivo        bool                 `orm:"column(estado_activo)"`
	CategoriaCompromiso *CategoriaCompromiso `orm:"column(categoria_compromiso);rel(fk)"`
	Descripcion         string               `orm:"column(descripcion);null"`
}
