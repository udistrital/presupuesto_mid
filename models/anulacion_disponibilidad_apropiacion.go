package models

//AnulacionDisponibilidadApropiacion ...
type AnulacionDisponibilidadApropiacion struct {
	Id                        int                        `orm:"column(id);pk"`
	DisponibilidadApropiacion *DisponibilidadApropiacion `orm:"column(disponibilidad_apropiacion);rel(fk)"`
	Valor                     float64                    `orm:"column(valor)"`
	Anulacion                 *AnulacionDisponibilidad   `orm:"column(anulacion);rel(fk)"`
}
