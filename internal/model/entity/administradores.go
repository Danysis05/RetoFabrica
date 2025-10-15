package entity

import "time"

type Administrador struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Nombre   string `gorm:"size:100;not null"`
	Correo   string `gorm:"size:100;unique;not null"`
	Password string `gorm:"not null"`

	FechaCreacion     time.Time `gorm:"autoCreateTime"`
	FechaModificacion time.Time `gorm:"autoUpdateTime"`
}

func (d *Administrador) TableName() string {
	return "administradores" // plural, coincidiendo con la tabla real
}
