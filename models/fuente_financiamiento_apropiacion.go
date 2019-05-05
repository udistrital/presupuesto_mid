package models

// FuenteFinanciamientoApropiacion ...
type FuenteFinanciamientoApropiacion struct {
	Id                   int                   `orm:"column(id);pk;auto"`
	Apropiacion          *Apropiacion          `orm:"column(apropiacion);rel(fk)"`
	FuenteFinanciamiento *FuenteFinanciamiento `orm:"column(fuente_financiamiento);rel(fk)"`
	Dependencia          int                   `orm:"column(dependencia)"`
}
