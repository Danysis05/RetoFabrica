package dto

import "time"

type AdministradorDTO struct {
	ID                uint      `json:"id"`
	Nombre            string    `json:"nombre"`
	Correo            string    `json:"correo"`
	Password          string    `json:"password"`
	FechaCreacion     time.Time `json:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion"`
}
