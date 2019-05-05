package models

//AnulacionRegistroPresupuestalDisponibilidadApropiacion ...
type AnulacionRegistroPresupuestalDisponibilidadApropiacion struct {
	Id                                            int                                            `orm:"auto;column(id);pk"`
	AnulacionRegistroPresupuestal                 *AnulacionRegistroPresupuestal                 `orm:"column(anulacion_registro_presupuestal);rel(fk)"`
	RegistroPresupuestalDisponibilidadApropiacion *RegistroPresupuestalDisponibilidadApropiacion `orm:"column(registro_presupuestal_disponibilidad_apropiacion);rel(fk)"`
	Valor                                         float64                                        `orm:"column(valor)"`
}
