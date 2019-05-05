package models

//Apropiacion ...
type Apropiacion struct {
	Id              int                `orm:"column(id);pk"`
	Vigencia        float64            `orm:"column(vigencia);null"`
	Rubro           *Rubro             `orm:"column(rubro);rel(fk)"`
	UnidadEjecutora int                `orm:"column(unidad_ejecutora);null"`
	Valor           float64            `orm:"column(valor);null"`
	Estado          *EstadoApropiacion `orm:"column(estado);rel(fk)"`
}
