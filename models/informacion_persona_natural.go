package models

// InformacionPersonaNatural ...
type InformacionPersonaNatural struct {
	TipoDocumento          *ParametroEstandar `orm:"column(tipo_documento);rel(fk)"`
	Id                     string             `orm:"column(num_documento_persona);pk"`
	DigitoVerificacion     float64            `orm:"column(digito_verificacion)"`
	PrimerApellido         string             `orm:"column(primer_apellido)"`
	SegundoApellido        string             `orm:"column(segundo_apellido);null"`
	PrimerNombre           string             `orm:"column(primer_nombre)"`
	SegundoNombre          string             `orm:"column(segundo_nombre);null"`
	Cargo                  string             `orm:"column(cargo)"`
	Perfil                 *ParametroEstandar `orm:"column(perfil);rel(fk)"`
	Profesion              string             `orm:"column(profesion);null"`
	Especialidad           string             `orm:"column(especialidad);null"`
	MontoCapitalAutorizado float64            `orm:"column(monto_capital_autorizado);null"`
	Genero                 string             `orm:"column(genero);null"`
}
