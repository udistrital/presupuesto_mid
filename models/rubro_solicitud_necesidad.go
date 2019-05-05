package models

// RubroSolicitudNecesidad ... 
type RubroSolicitudNecesidad struct {
	Id    int     `orm:"column(id);pk"`
	Rubro *Rubro  `orm:"column(rubro);rel(fk)"`
	Monto float64 `orm:"column(monto)"`
}
