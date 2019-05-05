package models

// Necesidad ...
type Necesidad struct {
	Id                          int     `orm:"column(id);pk"`
	Numero                      int     `orm:"column(numero);null"`
	Vigencia                    float64 `orm:"column(vigencia)"`
	UnidadEjecutora             int     `orm:"column(unidad_ejecutora)"`
	Objeto                      string  `orm:"column(objeto)"`
	Valor                       float64 `orm:"column(valor)"`
	Justificacion               string  `orm:"column(justificacion)"`
	ModalidadSeleccion          *ModalidadSeleccion
	DatosDependenciaSolicitante *Dependencia
}
