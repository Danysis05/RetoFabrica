package mapper

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
)

// Entity -> DTO
func ToDepartamentoDTO(dep *entity.Departamento) *dto.DepartamentoDTO {
	if dep == nil {
		return nil
	}

	fmt.Printf("🔧 Mapper - Convirtiendo Entity a DTO: %s (ID: %d)\n", dep.Nombre, dep.ID)

	return &dto.DepartamentoDTO{
		ID:                dep.ID, // ✅ CORREGIDO: ya es uint, no convertir
		Nombre:            dep.Nombre,
		Codigo:            dep.Codigo,
		Descripcion:       dep.Descripcion,
		FechaCreacion:     dep.FechaCreacion,
		FechaModificacion: dep.FechaModificacion,
	}
}

// DTO -> Entity
func ToDepartamentoEntity(depDTO *dto.DepartamentoDTO) *entity.Departamento {
	if depDTO == nil {
		return nil
	}

	fmt.Printf("🔧 Mapper - Convirtiendo DTO a Entity: %s (ID: %d)\n", depDTO.Nombre, depDTO.ID)

	entity := &entity.Departamento{
		Nombre:      depDTO.Nombre,
		Codigo:      depDTO.Codigo,
		Descripcion: depDTO.Descripcion,
	}

	// Solo asignar ID si es mayor que 0
	if depDTO.ID > 0 {
		entity.ID = depDTO.ID // ✅ CORREGIDO: ya es uint, no convertir
		fmt.Printf("   - Asignando ID: %d (modo edición)\n", depDTO.ID)
	} else {
		fmt.Printf("   - Sin ID (modo creación)\n")
	}

	return entity
}

// Entity -> DTO (Lista)
func ToDepartamentoDTOList(deps []*entity.Departamento) []*dto.DepartamentoDTO {
	if deps == nil {
		return nil
	}

	fmt.Printf("🔧 Mapper - Convirtiendo lista de %d departamentos\n", len(deps))

	dtos := make([]*dto.DepartamentoDTO, len(deps))
	for i, dep := range deps {
		dtos[i] = ToDepartamentoDTO(dep)
	}
	return dtos
}
