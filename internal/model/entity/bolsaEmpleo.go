package entity

import "time"

type BolsaEmpleo struct {
	ID             int     `gorm:"primaryKey;autoIncrement" json:"id"` // Cambiado a int
	Puesto         string  `gorm:"size:100;not null" json:"puesto"`
	Descripcion    string  `gorm:"type:text;not null" json:"descripcion"`
	Salario        float64 `gorm:"type:decimal(10,2);not null" json:"salario"`
	Estado         string  `gorm:"size:20;not null;default:'DISPONIBLE'" json:"estado"`
	DepartamentoID int     `gorm:"not null" json:"departamento_id"`       // Cambiado a int
	EmpleadoID     int     `gorm:"column:empleado_id" json:"empleado_id"` // Cambiado a int

	FechaPublicacion time.Time  `gorm:"type:date;autoCreateTime" json:"fecha_publicacion"`
	FechaOcupacion   *time.Time `gorm:"type:date" json:"fecha_ocupacion"`
	FechaCierre      *time.Time `gorm:"type:date" json:"fecha_cierre"`

	// Relaciones
	Departamento *Departamento `gorm:"foreignKey:DepartamentoID" json:"departamento"`
}

func (b *BolsaEmpleo) TableName() string {
	return "bolsa_empleos"
}
