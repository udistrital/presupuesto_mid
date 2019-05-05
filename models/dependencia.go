package models

//Dependencia ...
type Dependencia struct {
	Id                  int    `orm:"column(id);pk"`
	Nombre              string `orm:"column(nombre)"`
	JefeDependencia     int    `orm:"column(jefe_dependencia);rel(fk)"`
	OrdenadorGasto      int    `orm:"column(ordenador_gasto);rel(fk)"`
	InfoJefeDependencia *InformacionPersonaNatural
	InfoOrdenador       *InformacionPersonaNatural
}
