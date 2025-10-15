package service

import (
	"context"
	"fmt"
	"retoBack/internal/dao"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/mapper"
)

type bolsaService struct{}

var BolsaService = bolsaService{}

// 📋 Listar todas las bolsas
func (s *bolsaService) Listar(ctx context.Context) ([]*dto.BolsaEmpleoDTO, error) {
	lista, err := dao.BolsaEmpleos().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToBolsaDTOList(lista), nil
}

// 🆕 Crear nueva bolsa
func (s *bolsaService) Crear(ctx context.Context, bDTO *dto.BolsaEmpleoDTO) error {
	fmt.Printf("🔄 Iniciando creación...\n")
	fmt.Printf("📥 DTO recibido: ID=%d, Puesto=%s\n", bDTO.ID, bDTO.Puesto)

	entity := mapper.ToBolsaEntity(bDTO)

	fmt.Printf("🎯 Entity creado - ID: %d, Puesto: %s\n", entity.ID, entity.Puesto)

	err := dao.BolsaEmpleos().Create(ctx, entity)
	if err != nil {
		fmt.Printf("❌ Error en DAO: %v\n", err)
		return err
	}

	fmt.Printf("✅ Creación completada exitosamente - Nuevo ID: %d\n", entity.ID)
	return nil
}

// ✏️ Actualizar bolsa existente
func (s *bolsaService) Actualizar(ctx context.Context, bDTO *dto.BolsaEmpleoDTO) error {
	actual, err := dao.BolsaEmpleos().GetById(ctx, bDTO.ID)
	if err != nil {
		return err
	}

	fmt.Printf("🔍 Validando cambios en bolsa ID %d\n", bDTO.ID)
	fmt.Printf("   Estado actual: %s → Estado nuevo: %s\n", actual.Estado, bDTO.Estado)
	fmt.Printf("   EmpleadoID actual: %v → EmpleadoID nuevo: %v\n", actual.EmpleadoID, bDTO.EmpleadoID)

	// ⚙️ VALIDACIÓN 1: No permitir liberar manualmente un puesto ocupado
	if actual.Estado == "OCUPADO" && bDTO.Estado == "DISPONIBLE" {
		return fmt.Errorf("no se puede liberar manualmente un puesto ocupado")
	}

	// ⚙️ VALIDACIÓN 2: No permitir desasignar empleado manualmente
	if actual.EmpleadoID != 0 && bDTO.EmpleadoID == 0 {
		return fmt.Errorf("no se puede desasignar manualmente un empleado")
	}

	// ⚙️ VALIDACIÓN 3: No permitir asignar empleado manualmente
	if actual.EmpleadoID == 0 && bDTO.EmpleadoID != 0 {
		return fmt.Errorf("no se puede asignar manualmente un empleado")
	}

	// ⚙️ VALIDACIÓN 4: No permitir cambiar estado de un puesto ocupado
	if actual.Estado == "OCUPADO" && bDTO.Estado != "OCUPADO" {
		return fmt.Errorf("no se puede cambiar el estado de un puesto ocupado")
	}

	entity := mapper.ToBolsaEntity(bDTO)
	return dao.BolsaEmpleos().Update(ctx, entity)
}

// ❌ Eliminar bolsa por ID
func (s *bolsaService) Eliminar(ctx context.Context, id int) error {
	bolsa, err := dao.BolsaEmpleos().GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el empleo: %v", err)
	}

	if bolsa.EmpleadoID != 0 {
		return fmt.Errorf("no se puede eliminar este empleo porque tiene un empleado asignado")
	}

	return dao.BolsaEmpleos().Delete(ctx, id)
}

// 🔎 Listar bolsas por departamento
func (s *bolsaService) ListarPorDepartamento(ctx context.Context, departamentoID int) ([]*dto.BolsaEmpleoDTO, error) {
	lista, err := dao.BolsaEmpleos().GetByDepartamento(ctx, departamentoID)
	if err != nil {
		return nil, err
	}

	return mapper.ToBolsaDTOList(lista), nil
}

// 🔍 Obtener bolsa por ID
func (s *bolsaService) ObtenerPorID(ctx context.Context, id int) (*dto.BolsaEmpleoDTO, error) {
	fmt.Printf("🔍 Buscando bolsa de empleo con ID %d\n", id)

	bolsa, err := dao.BolsaEmpleos().GetById(ctx, id)
	if err != nil {
		fmt.Printf("❌ Error al obtener bolsa ID %d: %v\n", id, err)
		return nil, fmt.Errorf("no se pudo obtener la bolsa con ID %d: %v", id, err)
	}

	if bolsa == nil {
		fmt.Printf("⚠️ No se encontró bolsa con ID %d\n", id)
		return nil, fmt.Errorf("no se encontró la bolsa con ID %d", id)
	}

	dtoBolsa := mapper.ToBolsaDTO(bolsa)
	fmt.Printf("✅ Bolsa encontrada: %s (Estado: %s)\n", dtoBolsa.Puesto, dtoBolsa.Estado)
	return dtoBolsa, nil
}
