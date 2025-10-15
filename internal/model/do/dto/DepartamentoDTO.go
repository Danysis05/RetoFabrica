package dto

import "time"

type DepartamentoDTO struct {
	ID                uint      `json:"id"`
	Nombre            string    `json:"nombre"`
	Codigo            string    `json:"codigo"` // ← AGREGAR
	Descripcion       string    `json:"descripcion"`
	FechaCreacion     time.Time `json:"fecha_creacion"`     // ← AGREGAR
	FechaModificacion time.Time `json:"fecha_modificacion"` // ← AGREGAR
	// Eliminar BolsaEmpleoDTO - no es necesario
}
