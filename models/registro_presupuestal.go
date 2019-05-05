package models

import (
	"time"
)

// RegistroPresupuestal ...
type RegistroPresupuestal struct {
	Id                                            int                                              `orm:"column(id);pk;auto"`
	Vigencia                                      float64                                          `orm:"column(vigencia)"`
	FechaRegistro                                 time.Time                                        `orm:"column(fecha_registro);type(date);null"`
	Responsable                                   int                                              `orm:"column(responsable);null"`
	Estado                                        *EstadoRegistroPresupuestal                      `orm:"column(estado);rel(fk)"`
	NumeroRegistroPresupuestal                    int                                              `orm:"column(numero_registro_presupuestal)"`
	Beneficiario                                  int                                              `orm:"column(beneficiario);null"`
	TipoCompromiso                                *Compromiso                                      `orm:"column(tipo_compromiso);rel(fk)"`
	NumeroCompromiso                              int                                              `orm:"column(numero_compromiso)"`
	Solicitud                                     int                                              `orm:"column(solicitud)"`
	RegistroPresupuestalDisponibilidadApropiacion []*RegistroPresupuestalDisponibilidadApropiacion `orm:"reverse(many)"`
	DatosSolicitud                                *SolicitudRp
}

// DatosRubroRegistroPresupuestal ...
type DatosRubroRegistroPresupuestal struct {
	Id                 int
	Disponibilidad     *Disponibilidad
	Apropiacion        *Apropiacion
	FuenteFinanciacion *FuenteFinanciacion
	Valor              float64
	ValorAsignado      float64
	Saldo              float64
}

// DatosRegistroPresupuestal ...
type DatosRegistroPresupuestal struct {
	Rp     *RegistroPresupuestal
	Rubros []DatosRubroRegistroPresupuestal
}

// InfoSolRp ...
type InfoSolRp struct {
	Solicitud *SolicitudRp
	Rubros    []DisponibilidadApropiacionSolicitudRp
}
