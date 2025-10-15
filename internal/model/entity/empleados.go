package entity

import "time"

type Empleados struct {
	ID                int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre            string `gorm:"size:100;not null" json:"nombre"`
	Apellido          string `gorm:"size:100;not null" json:"apellido"`
	DocumentoTipo     string `gorm:"size:50;not null" json:"documento_tipo"`
	DocumentoNumero   string `gorm:"size:50;not null;unique" json:"documento_numero"`
	CorreoElectronico string `gorm:"size:100;not null;unique" json:"correo_electronico"`
	Ciudad            string `gorm:"size:100;not null" json:"ciudad"`
	Direccion         string `gorm:"size:255;not null" json:"direccion"`
	Telefono          string `gorm:"size:20" json:"telefono"`

	BolsaEmpleoID int          `gorm:"index" json:"bolsa_empleo_id"`
	BolsaEmpleo   *BolsaEmpleo `gorm:"foreignKey:BolsaEmpleoID" json:"bolsa_empleo"`

	FechaCreacion     time.Time `gorm:"autoCreateTime" json:"fecha_creacion"`
	FechaModificacion time.Time `gorm:"autoUpdateTime" json:"fecha_modificacion"`

	Nominas []Nomina `gorm:"foreignKey:EmpleadoId" json:"nominas"`
}

func (e *Empleados) TableName() string {
	return "empleados"
}
