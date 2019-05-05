package models

// RegistroPresupuestalDisponibilidadApropiacion ...
type RegistroPresupuestalDisponibilidadApropiacion struct {
	Id                        int                        `orm:"column(id);pk"`
	RegistroPresupuestal      *RegistroPresupuestal       `orm:"column(registro_presupuestal);rel(fk)"`
	DisponibilidadApropiacion *DisponibilidadApropiacion `orm:"column(disponibilidad_apropiacion);rel(fk)"`
	Valor                     float64                    `orm:"column(valor);null"`
}
