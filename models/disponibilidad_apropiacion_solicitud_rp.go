package models

// DisponibilidadApropiacionSolicitud_rp ...
type DisponibilidadApropiacionSolicitudRp struct {
	Id                        int
	DisponibilidadApropiacion int `orm:"column(disponibilidad_apropiacion)"`
	SolicitudRp               int `orm:"column(solicitud_rp)"`
	Monto                     int `orm:"column(monto)"`
}
