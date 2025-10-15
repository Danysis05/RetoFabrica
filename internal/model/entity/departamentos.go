package entity

import "time"

type Departamento struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Nombre      string `gorm:"size:100;not null;unique" json:"nombre"`
	Codigo      string `gorm:"size:10;not null;unique" json:"codigo"`
	Descripcion string `gorm:"size:255" json:"descripcion"`

	// RELACIONES CORREGIDAS:
	// Solo mantener la relación con BolsaEmpleo (que sí tiene DepartamentoID)
	BolsaEmpleos []BolsaEmpleo `gorm:"foreignKey:DepartamentoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// ELIMINAR esta relación - Empleados no tiene DepartamentoID
	// Empleados    []Empleados   `gorm:"foreignKey:DepartamentoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	FechaCreacion     time.Time `gorm:"autoCreateTime"`
	FechaModificacion time.Time `gorm:"autoUpdateTime"`
}

func (d *Departamento) TableName() string {
	return "departamentos"
}
func (d *Departamento) TieneEmpleosEnBolsa() bool {
	return len(d.BolsaEmpleos) > 0
}
func (d *Departamento) ContarEmpleosActivosEnBolsa() int {
	count := 0
	for _, empleo := range d.BolsaEmpleos {
		if empleo.Estado == "DISPONIBLE" || empleo.Estado == "OCUPADO" {
			count++
		}
	}
	return count
}
