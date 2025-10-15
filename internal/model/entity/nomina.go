// 📁 entity/nomina.go
package entity

import "time"

type Nomina struct {
	Id             int        `json:"id" orm:"id,primary" description:"ID"`
	EmpleadoId     int        `json:"empleado_id" orm:"empleado_id" description:"ID del empleado"`
	Empleado       *Empleados `json:"empleado,omitempty" orm:"with:id=empleado_id" description:"Empleado relacionado"`
	FechaPago      time.Time  `json:"fecha_pago" orm:"fecha_pago" description:"Fecha de pago"`
	SalarioBase    float64    `json:"salario_base" orm:"salario_base" description:"Salario base"`
	HorasExtras    float64    `json:"horas_extras" orm:"horas_extras" description:"Horas extras"`
	Bonificaciones float64    `json:"bonificaciones" orm:"bonificaciones" description:"Bonificaciones"`
	Deducciones    float64    `json:"deducciones" orm:"deducciones" description:"Deducciones"`
	TotalPago      float64    `json:"total_pago" orm:"total_pago" description:"Total a pagar"`
}

func (n *Nomina) TableName() string {
	return "nomina"
}
