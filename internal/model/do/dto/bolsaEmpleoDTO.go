package dto

import "time"

type BolsaEmpleoDTO struct {
	ID               int        `json:"id"`
	Puesto           string     `json:"puesto"`
	Descripcion      string     `json:"descripcion"`
	Salario          float64    `json:"salario"`           // 💰 Sueldo ofrecido
	Estado           string     `json:"estado"`            // 📌 "DISPONIBLE" | "OCUPADO" | "CERRADO"
	DepartamentoID   int        `json:"departamento_id"`   // 🔗 FK al departamento
	EmpleadoID       int        `json:"empleado_id"`       // 🔗 FK al empleado (nullable)
	FechaPublicacion time.Time  `json:"fecha_publicacion"` // 🗓️ Fecha en que se publicó la oferta
	FechaCierre      *time.Time `json:"fecha_cierre"`      // 🕓 Fecha en que se cerró la oferta
	FechaOcupacion   *time.Time `json:"fecha_ocupacion"`   // 🧾 Fecha en que se ocupó el puesto
}
