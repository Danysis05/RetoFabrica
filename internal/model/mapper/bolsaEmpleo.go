package mapper

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
)

// 🔄 Entity → DTO
func ToBolsaDTO(b *entity.BolsaEmpleo) *dto.BolsaEmpleoDTO {
	if b == nil {
		return nil
	}

	dto := &dto.BolsaEmpleoDTO{
		ID:               b.ID,
		Puesto:           b.Puesto,
		Descripcion:      b.Descripcion,
		Salario:          b.Salario,
		Estado:           b.Estado,
		DepartamentoID:   b.DepartamentoID,
		EmpleadoID:       b.EmpleadoID,
		FechaPublicacion: b.FechaPublicacion,
		FechaOcupacion:   b.FechaOcupacion,
		FechaCierre:      b.FechaCierre,
	}

	fmt.Printf("🔧 Mapper Bolsa → DTO: %s (ID: %d, Estado: %s)\n", b.Puesto, b.ID, b.Estado)
	return dto
}

// 🔁 DTO → Entity para CREACIÓN (sin ID)
func ToBolsaEntityForCreate(bDTO *dto.BolsaEmpleoDTO) *entity.BolsaEmpleo {
	if bDTO == nil {
		return nil
	}

	entity := &entity.BolsaEmpleo{
		Puesto:           bDTO.Puesto,
		Descripcion:      bDTO.Descripcion,
		Salario:          bDTO.Salario,
		Estado:           bDTO.Estado,
		DepartamentoID:   bDTO.DepartamentoID,
		EmpleadoID:       bDTO.EmpleadoID,
		FechaPublicacion: bDTO.FechaPublicacion,
		FechaOcupacion:   bDTO.FechaOcupacion,
		FechaCierre:      bDTO.FechaCierre,
	}

	fmt.Printf("🆕 Mapper Bolsa - Creación: %s (Salario: %.2f)\n", bDTO.Puesto, bDTO.Salario)
	return entity
}

// 🔁 DTO → Entity para ACTUALIZACIÓN (con ID)
func ToBolsaEntityForUpdate(bDTO *dto.BolsaEmpleoDTO) *entity.BolsaEmpleo {
	if bDTO == nil {
		return nil
	}

	entity := &entity.BolsaEmpleo{
		ID:               bDTO.ID,
		Puesto:           bDTO.Puesto,
		Descripcion:      bDTO.Descripcion,
		Salario:          bDTO.Salario,
		Estado:           bDTO.Estado,
		DepartamentoID:   bDTO.DepartamentoID,
		EmpleadoID:       bDTO.EmpleadoID,
		FechaPublicacion: bDTO.FechaPublicacion,
		FechaOcupacion:   bDTO.FechaOcupacion,
		FechaCierre:      bDTO.FechaCierre,
	}

	fmt.Printf("✏️ Mapper Bolsa - Actualización: %s (ID: %d)\n", bDTO.Puesto, bDTO.ID)
	return entity
}

// 🔁 DTO → Entity (versión genérica - usa creación por defecto)
func ToBolsaEntity(bDTO *dto.BolsaEmpleoDTO) *entity.BolsaEmpleo {
	return ToBolsaEntityForCreate(bDTO)
}

// 🔄 Entity → DTO (para listas)
func ToBolsaDTOList(entities []*entity.BolsaEmpleo) []*dto.BolsaEmpleoDTO {
	if entities == nil {
		return nil
	}

	dtos := make([]*dto.BolsaEmpleoDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = ToBolsaDTO(entity)
	}

	fmt.Printf("📋 Mapper Bolsa - Lista convertida: %d registros\n", len(dtos))
	return dtos
}

// 🔄 Entity → DTO con relaciones completas
func ToBolsaDTOWithRelations(b *entity.BolsaEmpleo) *dto.BolsaEmpleoDTO {
	if b == nil {
		return nil
	}

	dto := ToBolsaDTO(b)

	if b.Departamento != nil {
		fmt.Printf("🏢 Mapper Bolsa - Con departamento: %s\n", b.Departamento.Nombre)
	}

	return dto
}

// ⚙️ Validación de estado permitido
func IsEstadoValido(estado string) bool {
	estadosValidos := map[string]bool{
		"DISPONIBLE": true,
		"OCUPADO":    true,
		"CERRADO":    true,
	}
	return estadosValidos[estado]
}

// ✅ Validar DTO antes de usarlo
func ValidateBolsaDTO(bDTO *dto.BolsaEmpleoDTO) error {
	if bDTO == nil {
		return fmt.Errorf("el DTO de bolsa es nulo")
	}

	if bDTO.Puesto == "" {
		return fmt.Errorf("el puesto es requerido")
	}

	if bDTO.Salario <= 0 {
		return fmt.Errorf("el salario debe ser mayor a 0")
	}

	if !IsEstadoValido(bDTO.Estado) {
		return fmt.Errorf("estado inválido: %s", bDTO.Estado)
	}

	if bDTO.DepartamentoID <= 0 {
		return fmt.Errorf("departamento inválido o no asignado")
	}

	return nil
}
