package mapper

import (
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
	"time"
)

// Entity -> DTO
func ToNominaDTO(n *entity.Nomina) *dto.NominaDTO {
	if n == nil {
		return nil
	}
	return &dto.NominaDTO{
		ID:             n.Id,
		EmpleadoID:     n.EmpleadoId,
		FechaPago:      n.FechaPago,
		SalarioBase:    n.SalarioBase,
		HorasExtras:    n.HorasExtras,
		Bonificaciones: n.Bonificaciones,
		Deducciones:    n.Deducciones,
		TotalPago:      n.TotalPago,
	}
}

// DTO -> Entity (general)
func ToNominaEntity(nDTO *dto.NominaDTO) *entity.Nomina {
	if nDTO == nil {
		return nil
	}
	return &entity.Nomina{
		Id:             nDTO.ID,
		EmpleadoId:     nDTO.EmpleadoID,
		FechaPago:      nDTO.FechaPago,
		SalarioBase:    nDTO.SalarioBase,
		HorasExtras:    nDTO.HorasExtras,
		Bonificaciones: nDTO.Bonificaciones,
		Deducciones:    nDTO.Deducciones,
		TotalPago:      nDTO.TotalPago,
	}
}

// DTO -> Entity para creación (ignora ID y FechaPago)
func ToNominaEntityForCreate(nDTO *dto.NominaDTO) *entity.Nomina {
	if nDTO == nil {
		return nil
	}
	return &entity.Nomina{
		EmpleadoId:  nDTO.EmpleadoID,
		HorasExtras: nDTO.HorasExtras,
		// Salario, Bonificaciones, Deducciones y TotalPago se calculan en el service
		FechaPago: time.Time{},
	}
}

// DTO -> Entity para actualización (usa ID obligatorio)
func ToNominaEntityForUpdate(nDTO *dto.NominaDTO) *entity.Nomina {
	if nDTO == nil {
		return nil
	}
	return &entity.Nomina{
		Id:          nDTO.ID,
		EmpleadoId:  nDTO.EmpleadoID,
		HorasExtras: nDTO.HorasExtras,
		// Salario, Bonificaciones, Deducciones y TotalPago se pueden recalcular en el service
		FechaPago: nDTO.FechaPago,
	}
}
