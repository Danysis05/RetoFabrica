package dto

import "time"

type EmpleadoDTO struct {
	ID                int       `json:"id,omitempty"` // omitempty evita enviar ID en creación
	Nombre            string    `json:"nombre"`
	Apellido          string    `json:"apellido"`
	DocumentoTipo     string    `json:"documentoTipo"`
	DocumentoNumero   string    `json:"documentoNumero"`
	CorreoElectronico string    `json:"correoElectronico"`
	Ciudad            string    `json:"ciudad"`
	Direccion         string    `json:"direccion"`
	Telefono          string    `json:"telefono"`
	BolsaEmpleoID     int       `json:"bolsaEmpleoID"`
	BolsaPuesto       string    `json:"bolsaPuesto,omitempty"`
	FechaCreacion     time.Time `json:"fechaCreacion,omitempty"`
	FechaModificacion time.Time `json:"fechaModificacion,omitempty"`
}
