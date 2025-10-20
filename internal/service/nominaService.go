// Crear nómina calculando horas extras, bonificaciones y deducciones
// 📁 service/nomina_service.go - ACTUALIZADO
package service

import (
	"context"
	"fmt"
	"retoBack/internal/dao"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

type nominaService struct{}

var Nominas = &nominaService{}

// Create - VERSIÓN CORREGIDA para coincidir con el controller
func (s *nominaService) Create(ctx context.Context, nomina *entity.Nomina, horasExtras float64, diasFaltantes int) (*dto.NominaDTO, error) {
	fmt.Println("🎯 [NOMINA SERVICE] Create - Iniciando creación de nómina...")

	if nomina == nil {
		return nil, gerror.New("la entidad nómina no puede ser nula")
	}

	fmt.Printf("📝 [NOMINA SERVICE] Create - Datos recibidos: EmpleadoID=%d, HorasExtras=%.2f, DiasFaltantes=%d\n",
		nomina.EmpleadoId, horasExtras, diasFaltantes)

	// ✅ LOGS DETALLADOS PARA DEBUG
	fmt.Printf("🔍 [DEBUG SERVICE] Nomina entity recibida: %+v\n", nomina)

	// Obtener empleado con bolsa de empleo
	empleado, err := dao.Empleados().GetById(ctx, nomina.EmpleadoId)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] Create - Error al obtener empleado: %v\n", err)
		return nil, gerror.New("error al obtener empleado: " + err.Error())
	}
	if empleado == nil {
		fmt.Printf("❌ [NOMINA SERVICE] Create - Empleado no encontrado: ID=%d\n", nomina.EmpleadoId)
		return nil, gerror.New("empleado no encontrado")
	}

	// ✅ LOGS DETALLADOS DEL EMPLEADO
	fmt.Printf("🔍 [DEBUG SERVICE] Empleado obtenido: %s %s\n", empleado.Nombre, empleado.Apellido)
	fmt.Printf("🔍 [DEBUG SERVICE] BolsaEmpleoID del empleado: %d\n", empleado.BolsaEmpleoID)
	fmt.Printf("🔍 [DEBUG SERVICE] BolsaEmpleo es nil: %t\n", empleado.BolsaEmpleo == nil)

	if empleado.BolsaEmpleo != nil {
		fmt.Printf("🔍 [DEBUG SERVICE] BolsaEmpleo detalles - ID: %d, Puesto: %s, Salario: %.2f\n",
			empleado.BolsaEmpleo.ID, empleado.BolsaEmpleo.Puesto, empleado.BolsaEmpleo.Salario)
	} else {
		fmt.Printf("❌ [DEBUG SERVICE] BolsaEmpleo es NIL - No se cargó la relación\n")
		// Intentar cargar manualmente la bolsa de empleo
		if empleado.BolsaEmpleoID > 0 {
			fmt.Printf("🔄 [DEBUG SERVICE] Intentando cargar bolsa manualmente, ID: %d\n", empleado.BolsaEmpleoID)
			bolsa, err := dao.BolsaEmpleos().GetById(ctx, empleado.BolsaEmpleoID)
			if err != nil {
				fmt.Printf("❌ [DEBUG SERVICE] Error cargando bolsa manualmente: %v\n", err)
			} else if bolsa != nil {
				empleado.BolsaEmpleo = bolsa
				fmt.Printf("✅ [DEBUG SERVICE] Bolsa cargada manualmente - Puesto: %s, Salario: %.2f\n",
					bolsa.Puesto, bolsa.Salario)
			}
		}
	}

	if empleado.BolsaEmpleo == nil {
		fmt.Printf("❌ [NOMINA SERVICE] Create - Empleado sin bolsa empleo: ID=%d\n", nomina.EmpleadoId)
		return nil, gerror.New("empleado no tiene asignado un puesto con salario")
	}

	fmt.Printf("✅ [NOMINA SERVICE] Create - Empleado encontrado: %s %s, Puesto: %s, Salario: %.2f\n",
		empleado.Nombre, empleado.Apellido, empleado.BolsaEmpleo.Puesto, empleado.BolsaEmpleo.Salario)

	// Calcular valores de la nómina
	nomina.SalarioBase = empleado.BolsaEmpleo.Salario
	nomina.HorasExtras = horasExtras
	valorHora := nomina.SalarioBase / 240
	nomina.Bonificaciones = horasExtras * valorHora * 1.5
	diasFaltantes = nomina.DiasFaltantes
	nomina.Deducciones = float64(diasFaltantes) * (nomina.SalarioBase / 30)
	nomina.TotalPago = nomina.SalarioBase + nomina.Bonificaciones - nomina.Deducciones
	nomina.FechaPago = time.Now()

	fmt.Printf("🧮 [NOMINA SERVICE] Create - Cálculos: SalarioBase=%.2f, Bonificaciones=%.2f, Deducciones=%.2f, TotalPago=%.2f\n",
		nomina.SalarioBase, nomina.Bonificaciones, nomina.Deducciones, nomina.TotalPago)

	// Guardar
	if err := dao.Nominas.Create(ctx, nomina); err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] Create - Error al guardar en BD: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [NOMINA SERVICE] Create - Nómina guardada en BD: ID=%d\n", nomina.Id)

	// Retornar DTO con información del empleado
	dtoResult := &dto.NominaDTO{
		ID:             nomina.Id,
		EmpleadoID:     nomina.EmpleadoId,
		FechaPago:      nomina.FechaPago,
		SalarioBase:    nomina.SalarioBase,
		HorasExtras:    nomina.HorasExtras,
		Bonificaciones: nomina.Bonificaciones,
		Deducciones:    nomina.Deducciones,
		TotalPago:      nomina.TotalPago,
		EmpleadoNombre: fmt.Sprintf("%s %s", empleado.Nombre, empleado.Apellido),
		EmpleadoPuesto: empleado.BolsaEmpleo.Puesto,
	}

	fmt.Printf("🎉 [NOMINA SERVICE] Create - Nómina creada exitosamente: ID=%d, SalarioBase=%.2f\n", dtoResult.ID, dtoResult.SalarioBase)
	return dtoResult, nil
}

// Update - VERSIÓN CORREGIDA para coincidir con el controller
func (s *nominaService) Update(ctx context.Context, nomina *entity.Nomina, horasExtras float64, diasFaltantes int) error {
	fmt.Println("🎯 [NOMINA SERVICE] Update - Iniciando actualización de nómina...")

	if nomina == nil {
		return gerror.New("la entidad nómina no puede ser nula")
	}

	if nomina.Id == 0 {
		return gerror.New("ID de nómina requerido")
	}

	fmt.Printf("📝 [NOMINA SERVICE] Update - Datos: ID=%d, EmpleadoID=%d, HorasExtras=%.2f, DiasFaltantes=%d\n",
		nomina.Id, nomina.EmpleadoId, horasExtras, diasFaltantes)

	// Obtener empleado con bolsa de empleo
	empleado, err := dao.Empleados().GetById(ctx, nomina.EmpleadoId)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] Update - Error al obtener empleado: %v\n", err)
		return gerror.New("error al obtener empleado: " + err.Error())
	}
	if empleado == nil {
		fmt.Printf("❌ [NOMINA SERVICE] Update - Empleado no encontrado: ID=%d\n", nomina.EmpleadoId)
		return gerror.New("empleado no encontrado")
	}

	if empleado.BolsaEmpleo == nil {
		fmt.Printf("❌ [NOMINA SERVICE] Update - Empleado sin bolsa empleo: ID=%d\n", nomina.EmpleadoId)
		return gerror.New("empleado no tiene asignado un puesto con salario")
	}

	fmt.Printf("✅ [NOMINA SERVICE] Update - Empleado encontrado: %s %s, Salario: %.2f\n",
		empleado.Nombre, empleado.Apellido, empleado.BolsaEmpleo.Salario)

	nomina.SalarioBase = empleado.BolsaEmpleo.Salario
	nomina.HorasExtras = horasExtras
	valorHora := nomina.SalarioBase / 240
	nomina.Bonificaciones = horasExtras * valorHora * 1.5
	nomina.Deducciones = float64(diasFaltantes) * (nomina.SalarioBase / 30)
	nomina.TotalPago = nomina.SalarioBase + nomina.Bonificaciones - nomina.Deducciones

	fmt.Printf("🧮 [NOMINA SERVICE] Update - Cálculos: SalarioBase=%.2f, Bonificaciones=%.2f, Deducciones=%.2f, TotalPago=%.2f\n",
		nomina.SalarioBase, nomina.Bonificaciones, nomina.Deducciones, nomina.TotalPago)

	// Actualizar en BD usando GoFrame DAO
	err = dao.Nominas.Update(ctx, nomina)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] Update - Error al actualizar en BD: %v\n", err)
		return err
	}

	fmt.Printf("✅ [NOMINA SERVICE] Update - Nómina actualizada exitosamente: ID=%d\n", nomina.Id)
	return nil
}

// Los demás métodos (GetAll, GetById, Delete, GetByEmpleadoID) se mantienen igual...

// Eliminar nómina (MISMA LÓGICA)
func (s *nominaService) Delete(ctx context.Context, id int) error {
	fmt.Printf("🎯 [NOMINA SERVICE] Delete - Iniciando eliminación de nómina: ID=%d\n", id)

	if id == 0 {
		return gerror.New("ID de nómina requerido")
	}

	// Eliminar usando GoFrame DAO (SOLO CAMBIO DE dao.Nominas.Delete)
	err := dao.Nominas.Delete(ctx, id)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] Delete - Error al eliminar: %v\n", err)
		return err
	}

	fmt.Printf("✅ [NOMINA SERVICE] Delete - Nómina eliminada exitosamente: ID=%d\n", id)
	return nil
}

// Listar todas las nóminas (MISMA LÓGICA)
func (s *nominaService) GetAll(ctx context.Context) ([]*dto.NominaDTO, error) {
	fmt.Println("🎯 [NOMINA SERVICE] GetAll - Iniciando...")

	// 1. Obtener nóminas del DAO de GoFrame (SOLO CAMBIO DE dao.Nominas.GetAll)
	nominas, err := dao.Nominas.GetAll(ctx)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] GetAll - Error en dao.GetAll: %v\n", err)
		return nil, gerror.New("error al obtener nóminas de la base de datos: " + err.Error())
	}

	fmt.Printf("✅ [NOMINA SERVICE] GetAll - DAO retornó %d nóminas\n", len(nominas))

	// 2. Si no hay nóminas, retornar array vacío (MISMA LÓGICA)
	if len(nominas) == 0 {
		fmt.Println("ℹ️ [NOMINA SERVICE] GetAll - No hay nóminas en la base de datos")
		return []*dto.NominaDTO{}, nil
	}

	// 3. Procesar cada nómina y enriquecer con datos del empleado (MISMA LÓGICA)
	var result []*dto.NominaDTO
	for i, n := range nominas {
		if n == nil {
			fmt.Printf("⚠️ [NOMINA SERVICE] GetAll - Nómina en posición %d es NIL, saltando...\n", i)
			continue
		}

		fmt.Printf("📋 [NOMINA SERVICE] GetAll - Procesando nómina %d: ID=%d, EmpleadoID=%d\n", i, n.Id, n.EmpleadoId)

		var empleadoNombre, empleadoPuesto string
		var salarioBase float64 = n.SalarioBase

		// El DAO de GoFrame ya carga la relación Empleado automáticamente
		if n.Empleado != nil {
			empleadoNombre = fmt.Sprintf("%s %s", n.Empleado.Nombre, n.Empleado.Apellido)
			if n.Empleado.BolsaEmpleo != nil {
				empleadoPuesto = n.Empleado.BolsaEmpleo.Puesto
				// Si no hay salario en la nómina, usar el de la bolsa de empleo
				if salarioBase == 0 {
					salarioBase = n.Empleado.BolsaEmpleo.Salario
				}
			}
			fmt.Printf("👤 [NOMINA SERVICE] GetAll - Empleado encontrado: %s, Puesto: %s\n", empleadoNombre, empleadoPuesto)
		} else {
			fmt.Printf("⚠️ [NOMINA SERVICE] GetAll - Empleado es NIL para nómina ID=%d\n", n.Id)
			// Si no viene la relación, intentar cargar empleado manualmente
			empleado, err := dao.Empleados().GetById(ctx, n.EmpleadoId)
			if err == nil && empleado != nil {
				empleadoNombre = fmt.Sprintf("%s %s", empleado.Nombre, empleado.Apellido)
				if empleado.BolsaEmpleo != nil {
					empleadoPuesto = empleado.BolsaEmpleo.Puesto
					if salarioBase == 0 {
						salarioBase = empleado.BolsaEmpleo.Salario
					}
				}
				fmt.Printf("🔄 [NOMINA SERVICE] GetAll - Empleado cargado manualmente: %s\n", empleadoNombre)
			} else {
				empleadoNombre = fmt.Sprintf("Empleado %d", n.EmpleadoId)
			}
		}

		// 4. Crear DTO con información completa (MISMA LÓGICA)
		nominaDTO := &dto.NominaDTO{
			ID:             n.Id,
			EmpleadoID:     n.EmpleadoId,
			FechaPago:      n.FechaPago,
			SalarioBase:    salarioBase,
			HorasExtras:    n.HorasExtras,
			Bonificaciones: n.Bonificaciones,
			Deducciones:    n.Deducciones,
			TotalPago:      n.TotalPago,
			EmpleadoNombre: empleadoNombre,
			EmpleadoPuesto: empleadoPuesto,
		}

		result = append(result, nominaDTO)
		fmt.Printf("✅ [NOMINA SERVICE] GetAll - Nómina %d procesada: %s - %.2f\n", i, empleadoNombre, nominaDTO.TotalPago)
	}

	fmt.Printf("🎉 [NOMINA SERVICE] GetAll - Proceso completado: %d/%d nóminas procesadas\n", len(result), len(nominas))
	return result, nil
}

// Obtener nómina por ID (MISMA LÓGICA)
func (s *nominaService) GetById(ctx context.Context, id int) (*dto.NominaDTO, error) {
	fmt.Printf("🎯 [NOMINA SERVICE] GetById - Buscando nómina: ID=%d\n", id)

	if id == 0 {
		return nil, gerror.New("ID de nómina requerido")
	}

	// Obtener usando GoFrame DAO (SOLO CAMBIO DE dao.Nominas.GetById)
	nomina, err := dao.Nominas.GetById(ctx, id)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] GetById - Error en dao.GetById: %v\n", err)
		return nil, err
	}

	if nomina == nil {
		fmt.Printf("❌ [NOMINA SERVICE] GetById - Nómina no encontrada: ID=%d\n", id)
		return nil, gerror.New("nómina no encontrada")
	}

	fmt.Printf("✅ [NOMINA SERVICE] GetById - Nómina encontrada: ID=%d, EmpleadoID=%d\n", nomina.Id, nomina.EmpleadoId)

	var empleadoNombre, empleadoPuesto string
	salario := nomina.SalarioBase

	if nomina.Empleado != nil {
		empleadoNombre = fmt.Sprintf("%s %s", nomina.Empleado.Nombre, nomina.Empleado.Apellido)
		if nomina.Empleado.BolsaEmpleo != nil {
			empleadoPuesto = nomina.Empleado.BolsaEmpleo.Puesto
			if salario == 0 {
				salario = nomina.Empleado.BolsaEmpleo.Salario
			}
		}
		fmt.Printf("💰 [NOMINA SERVICE] GetById - Empleado: %s, Salario: %.2f\n", empleadoNombre, salario)
	} else {
		fmt.Printf("💰 [NOMINA SERVICE] GetById - Usando salario base de nómina: %.2f\n", salario)
	}

	dtoResult := &dto.NominaDTO{
		ID:             nomina.Id,
		EmpleadoID:     nomina.EmpleadoId,
		FechaPago:      nomina.FechaPago,
		SalarioBase:    salario,
		HorasExtras:    nomina.HorasExtras,
		Bonificaciones: nomina.Bonificaciones,
		Deducciones:    nomina.Deducciones,
		TotalPago:      nomina.TotalPago,
		EmpleadoNombre: empleadoNombre,
		EmpleadoPuesto: empleadoPuesto,
	}

	fmt.Printf("✅ [NOMINA SERVICE] GetById - Nómina retornada exitosamente: ID=%d\n", dtoResult.ID)
	return dtoResult, nil
}

// Obtener nóminas por ID de empleado (MISMA LÓGICA)
func (s *nominaService) GetByEmpleadoID(ctx context.Context, empleadoID int) ([]*dto.NominaDTO, error) {
	fmt.Printf("🎯 [NOMINA SERVICE] GetByEmpleadoID - Buscando nóminas del empleado: ID=%d\n", empleadoID)

	if empleadoID == 0 {
		return nil, gerror.New("ID de empleado requerido")
	}

	// Obtener usando GoFrame DAO (SOLO CAMBIO DE dao.Nominas.GetByEmpleadoID)
	nominas, err := dao.Nominas.GetByEmpleadoID(ctx, empleadoID)
	if err != nil {
		fmt.Printf("❌ [NOMINA SERVICE] GetByEmpleadoID - Error: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [NOMINA SERVICE] GetByEmpleadoID - Encontradas %d nóminas\n", len(nominas))

	var result []*dto.NominaDTO
	for _, nomina := range nominas {
		var empleadoNombre string

		if nomina.Empleado != nil {
			empleadoNombre = fmt.Sprintf("%s %s", nomina.Empleado.Nombre, nomina.Empleado.Apellido)
		} else {
			empleadoNombre = fmt.Sprintf("Empleado %d", nomina.EmpleadoId)
		}

		dto := &dto.NominaDTO{
			ID:             nomina.Id,
			EmpleadoID:     nomina.EmpleadoId,
			EmpleadoNombre: empleadoNombre,
			FechaPago:      nomina.FechaPago,
			SalarioBase:    nomina.SalarioBase,
			HorasExtras:    nomina.HorasExtras,
			Bonificaciones: nomina.Bonificaciones,
			Deducciones:    nomina.Deducciones,
			TotalPago:      nomina.TotalPago,
		}
		result = append(result, dto)
	}

	return result, nil
}
