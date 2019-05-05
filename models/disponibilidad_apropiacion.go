package models

// DisponibilidadApropiacion ...
type DisponibilidadApropiacion struct {
	Id                   int                 `orm:"auto;column(id);pk"`
	Disponibilidad       *Disponibilidad     `orm:"column(disponibilidad);rel(fk)"`
	Apropiacion          *Apropiacion        `orm:"column(apropiacion);rel(fk)"`
	Valor                float64             `orm:"column(valor);null"`
	FuenteFinanciamiento *FuenteFinanciacion `orm:"column(fuente_financiamiento);rel(fk);null"`
}
