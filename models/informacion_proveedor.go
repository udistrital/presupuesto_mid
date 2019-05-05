package models

// InformacionProveedor ...
type InformacionProveedor struct {
	Id           int    `orm:"column(id_proveedor);pk"`
	Tipopersona  string `orm:"column(tipopersona)"`
	NumDocumento string `orm:"column(num_documento)"`
	NomProveedor string `orm:"column(nom_proveedor);null"`
}
