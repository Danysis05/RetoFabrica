package dto

import "time"

type NominaDTO struct {
	ID             int       `json:"id"`
	EmpleadoID     int       `json:"empleado_id"`
	EmpleadoNombre string    `json:"empleado_nombre,omitempty"` // ← AGREGAR
	EmpleadoPuesto string    `json:"empleado_puesto,omitempty"` // ← AGREGAR
	FechaPago      time.Time `json:"fecha_pago"`
	SalarioBase    float64   `json:"salario_base"`
	HorasExtras    float64   `json:"horas_extras"`
	Bonificaciones float64   `json:"bonificaciones"`
	Deducciones    float64   `json:"deducciones"`
	TotalPago      float64   `json:"total_pago"`
}
