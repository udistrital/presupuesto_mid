package models

// DependenciaDependencia ...
type DependenciaDependencia struct {
	Id               int `orm:"column(id);pk"`
	DependenciaPadre int `orm:"column(dependencia_padre);rel(fk)"`
	DependenciaHijo  int `orm:"column(dependencia_hijo);rel(fk)"`
}

//DependenciaNecesidad ...
type DependenciaNecesidad struct {
	Id                         int        `orm:"column(id);pk"`
	JefeDependenciaSolicitante int        `orm:"column(jefe_dependencia_solicitante)"`
	JefeDependenciaDestino     int        `orm:"column(jefe_dependencia_destino)"`
	Necesidad                  *Necesidad `orm:"column(necesidad)"`
	OrdenadorGasto             int        `orm:"column(ordenador_gasto)"`
}
