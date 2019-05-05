package models

import (
	"time"
)

// MovimientoFuenteFinanciamientoApropiacion ...
type MovimientoFuenteFinanciamientoApropiacion struct {
	Id                                        int                              `orm:"column(id);pk;auto"`
	Fecha                                     time.Time                        `orm:"column(fecha_registro);type(date)"`
	Valor                                     float64                          `orm:"column(valor)"`
	TipoMovimiento                            *TipoMovimiento                  `orm:"column(tipo_movimiento);rel(fk)"`
	Descripcion                               string                           `orm:"column(descripcion);null"`
	FuenteFinanciamientoApropiacion           *FuenteFinanciamientoApropiacion `orm:"column(fuente_financiamiento_apropiacion);rel(fk)"`
	MovimientoFuenteFinanciamientoApropiacion int                              `orm:"column(movimiento_fuente_financiamiento_apropiacion);null"`
	TipoDocumento                             int         			   						 `orm:"column(tipo_documento)"`
	ValorDisponible                           float64
	ValorGastado															float64
}
