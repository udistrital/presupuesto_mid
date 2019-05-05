package models

import "time"

// OrdenPago ...
type OrdenPago struct {
	Id               int               `orm:"column(id);pk;auto"`
	Vigencia         float64           `orm:"column(vigencia)"`
	ValorBase        float64           `orm:"column(valor_base)"`
	Convenio         int               `orm:"column(convenio);null"`
	SubTipoOrdenPago *SubTipoOrdenPago `orm:"column(sub_tipo_orden_pago);rel(fk)"`
	UnidadEjecutora  *UnidadEjecutora  `orm:"column(unidad_ejecutora);rel(fk)"`
	Liquidacion      int               `orm:"column(liquidacion);null"`
	EntradaAlmacen   int               `orm:"column(entrada_almacen);null"`
	Consecutivo      int               `orm:"column(consecutivo)"`
	Documento        int               `orm:"column(documento)"`
	FormaPago        *FormaPago        `orm:"column(forma_pago);rel(fk);null"`

	OrdenPagoEstadoOrdenPago      []*OrdenPagoEstadoOrdenPago      `orm:"reverse(many)"`
	OrdenPagoRegistroPresupuestal []*OrdenPagoRegistroPresupuestal `orm:"reverse(many)"`
	OrdenPagoCuentaEspecial       []*OrdenPagoCuentaEspecial       `orm:"reverse(many)"`
}

type FormaPago struct {
	Id                int     `orm:"column(id);pk;auto"`
	Nombre            string  `orm:"column(nombre)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	CodigoAbreviacion string  `orm:"column(codigo_abreviacion);null"`
	Activo            bool    `orm:"column(activo)"`
	NumeroOrden       float64 `orm:"column(numero_orden);null"`
}

type SubTipoOrdenPago struct {
	Id                int            `orm:"column(id);pk;auto"`
	Nombre            string         `orm:"column(nombre)"`
	Descripcion       string         `orm:"column(descripcion);null"`
	CodigoAbreviacion string         `orm:"column(codigo_abreviacion);null"`
	Activo            bool           `orm:"column(activo)"`
	NumeroOrden       float64        `orm:"column(numero_orden);null"`
	TipoOrdenPago     *TipoOrdenPago `orm:"column(tipo_orden_pago);rel(fk)"`
	GrupoSecuencia    string         `orm:"column(grupo_secuencia);null"`
}

type TipoOrdenPago struct {
	Id                int     `orm:"column(id);pk;auto"`
	Nombre            string  `orm:"column(nombre)"`
	Activo            bool    `orm:"column(activo)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	CodigoAbreviacion string  `orm:"column(codigo_abreviacion);null"`
	NumeroOrden       float64 `orm:"column(numero_orden);null"`
}

type OrdenPagoRegistroPresupuestal struct {
	Id                   int                   `orm:"column(id);pk;auto"`
	OrdenPago            *OrdenPago            `orm:"column(orden_pago);rel(fk)"`
	RegistroPresupuestal *RegistroPresupuestal `orm:"column(registro_presupuestal);rel(fk)"`
	FechaRegistro        time.Time             `orm:"column(fecha_registro);type(timestamp without time zone)"`
}

type OrdenPagoCuentaEspecial struct {
	Id             int        `orm:"column(id);pk;auto"`
	OrdenPago      *OrdenPago `orm:"column(orden_pago);rel(fk)"`
	CuentaEspecial int        `orm:"column(cuenta_especial);rel(fk)"`
	FormaPago      *FormaPago `orm:"column(forma_pago);rel(fk)"`
	ValorBase      float64    `orm:"column(valor_base)"`
	FechaRegistro  time.Time  `orm:"column(fecha_registro);type(timestamp without time zone)"`
	Usuario        int        `orm:"column(usuario)"`
}
