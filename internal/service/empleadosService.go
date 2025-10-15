package service

import (
	"context"
	"fmt"
	"retoBack/internal/dao"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/entity"
	"retoBack/internal/model/mapper"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

type empleadosService struct{}

var Empleados = empleadosService{}

// ----------------------
// Obtener empleado
// ----------------------

func (s *empleadosService) GetById(ctx context.Context, id int) (*entity.Empleados, error) {
	return dao.Empleados().GetById(ctx, id)
}

func (s *empleadosService) GetByIdDTO(ctx context.Context, id int) (*dto.EmpleadoDTO, error) {
	empleado, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToEmpleadoDTO(empleado), nil
}

func (s *empleadosService) GetAll(ctx context.Context) ([]*entity.Empleados, error) {
	empleados, err := dao.Empleados().GetAll(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "Error al obtener todos los empleados")
	}
	return empleados, nil
}

func (s *empleadosService) GetAllDTO(ctx context.Context) ([]*dto.EmpleadoDTO, error) {
	empleados, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToEmpleadoDTOList(empleados), nil
}

func (s *empleadosService) FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]*dto.EmpleadoDTO, error) {
	empleados, err := dao.Empleados().FindWithFilters(ctx, filters)
	if err != nil {
		return nil, gerror.Wrap(err, "Error al buscar empleados con filtros")
	}
	return mapper.ToEmpleadoDTOList(empleados), nil
}

func (s *empleadosService) Exists(ctx context.Context, documentoNumero, email string) (bool, error) {
	exists, err := dao.Empleados().Exists(ctx, documentoNumero, email)
	if err != nil {
		return false, gerror.Wrap(err, "Error al verificar existencia de empleado")
	}
	return exists, nil
}

func (s *empleadosService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats, err := dao.Empleados().GetStats(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "Error al obtener estadísticas")
	}
	return stats, nil
}

// ----------------------
// Crear empleado
// ----------------------

func (s *empleadosService) Crear(ctx context.Context, empleadoDTO *dto.EmpleadoDTO) error {
	fmt.Printf("📦 Body recibido: %+v\n", empleadoDTO)

	sanitizedDTO := mapper.SanitizeEmpleadoDTO(empleadoDTO)
	empleado := mapper.ToEmpleadoEntityFromDTO(sanitizedDTO)

	// Evitar insertar relación
	empleado.BolsaEmpleo = nil

	if err := s.validarEmpleado(ctx, empleado); err != nil {
		return err
	}

	if err := dao.Empleados().Create(ctx, empleado); err != nil {
		return gerror.Wrap(err, "Error al crear empleado")
	}

	if empleado.BolsaEmpleoID > 0 {
		if err := s.asignarBolsa(ctx, empleado.BolsaEmpleoID, empleado.ID); err != nil {
			_ = dao.Empleados().Delete(ctx, empleado.ID)
			return gerror.Wrap(err, "Error al actualizar bolsa de empleo")
		}
	}

	fmt.Printf("✅ Empleado creado exitosamente: %s %s (ID: %d)\n", empleado.Nombre, empleado.Apellido, empleado.ID)
	return nil
}

func (s *empleadosService) CrearDesdeDatos(ctx context.Context, nombre, apellido, docTipo, docNumero, correo, ciudad, direccion, telefono string, bolsaID int) error {
	empleado := mapper.ToEmpleadoEntityFromData(nombre, apellido, docTipo, docNumero, correo, ciudad, direccion, telefono, bolsaID)
	empleado.BolsaEmpleo = nil

	if err := s.validarEmpleado(ctx, empleado); err != nil {
		return err
	}

	if err := dao.Empleados().Create(ctx, empleado); err != nil {
		return gerror.Wrap(err, "Error al crear empleado")
	}

	if empleado.BolsaEmpleoID > 0 {
		if err := s.asignarBolsa(ctx, empleado.BolsaEmpleoID, empleado.ID); err != nil {
			_ = dao.Empleados().Delete(ctx, empleado.ID)
			return gerror.Wrap(err, "Error al actualizar bolsa de empleo")
		}
	}

	fmt.Printf("✅ Empleado creado exitosamente: %s %s (ID: %d)\n", empleado.Nombre, empleado.Apellido, empleado.ID)
	return nil
}

// ----------------------
// Actualizar empleado
// ----------------------

func (s *empleadosService) Update(ctx context.Context, empleadoDTO *dto.EmpleadoDTO) error {
	empleado := mapper.ToEmpleadoEntityFromDTO(empleadoDTO)
	empleado.BolsaEmpleo = nil

	if err := s.validarDuplicados(ctx, empleado); err != nil {
		return err
	}

	if err := dao.Empleados().Update(ctx, empleado); err != nil {
		return gerror.Wrap(err, "Error al actualizar empleado")
	}

	fmt.Printf("✅ Empleado actualizado exitosamente: %s %s (ID: %d)\n", empleado.Nombre, empleado.Apellido, empleado.ID)
	return nil
}

// ----------------------
// Eliminar empleado
// ----------------------

func (s *empleadosService) Delete(ctx context.Context, id int) error {
	empleado, err := s.GetById(ctx, id)
	if err != nil {
		return gerror.Wrap(err, "Error al obtener empleado para eliminar")
	}

	if empleado.BolsaEmpleoID > 0 {
		if err := s.liberarBolsa(ctx, empleado.BolsaEmpleoID); err != nil {
			return err
		}
	}

	if err := dao.Empleados().Delete(ctx, id); err != nil {
		return gerror.Wrap(err, "Error al eliminar empleado")
	}

	fmt.Printf("✅ Empleado eliminado exitosamente: %s %s (ID: %d)\n", empleado.Nombre, empleado.Apellido, empleado.ID)
	return nil
}

// ----------------------
// Helpers privados
// ----------------------

func (s *empleadosService) validarEmpleado(ctx context.Context, e *entity.Empleados) error {
	if e.BolsaEmpleoID > 0 {
		bolsa, err := dao.BolsaEmpleos().GetById(ctx, e.BolsaEmpleoID)
		if err != nil {
			return gerror.Wrap(err, "Error al validar bolsa de empleo")
		}
		if bolsa.Estado != "DISPONIBLE" || bolsa.EmpleadoID != 0 {
			return gerror.New("El empleo seleccionado no está disponible")
		}
	}
	return s.validarDuplicados(ctx, e)
}

func (s *empleadosService) validarDuplicados(ctx context.Context, e *entity.Empleados) error {
	if e.DocumentoNumero != "" {
		if existing, _ := dao.Empleados().FindByDocumento(ctx, e.DocumentoNumero); existing != nil && existing.ID != e.ID {
			return gerror.New("Ya existe un empleado con ese número de documento")
		}
	}
	if e.CorreoElectronico != "" {
		if existing, _ := dao.Empleados().FindByEmail(ctx, e.CorreoElectronico); existing != nil && existing.ID != e.ID {
			return gerror.New("Ya existe un empleado con ese correo electrónico")
		}
	}
	return nil
}

func (s *empleadosService) asignarBolsa(ctx context.Context, bolsaID, empleadoID int) error {
	bolsa, err := dao.BolsaEmpleos().GetById(ctx, bolsaID)
	if err != nil {
		return gerror.Wrap(err, "Error al obtener bolsa de empleo")
	}

	bolsa.Estado = "OCUPADO"
	bolsa.EmpleadoID = empleadoID
	now := time.Now()
	bolsa.FechaOcupacion = &now

	if err := dao.BolsaEmpleos().Update(ctx, bolsa); err != nil {
		return gerror.Wrap(err, "Error al actualizar bolsa de empleo")
	}

	fmt.Printf("🔧 Bolsa actualizada: ID %d → Estado: OCUPADO, EmpleadoID: %d\n", bolsaID, empleadoID)
	return nil
}

func (s *empleadosService) liberarBolsa(ctx context.Context, bolsaID int) error {
	bolsa, err := dao.BolsaEmpleos().GetById(ctx, bolsaID)
	if err != nil {
		return gerror.Wrap(err, "Error al obtener bolsa de empleo")
	}

	bolsa.Estado = "DISPONIBLE"
	bolsa.EmpleadoID = 0
	bolsa.FechaOcupacion = nil

	if err := dao.BolsaEmpleos().Update(ctx, bolsa); err != nil {
		return gerror.Wrap(err, "Error al actualizar bolsa de empleo")
	}

	fmt.Printf("🔧 Bolsa liberada: ID %d → Estado: DISPONIBLE\n", bolsaID)
	return nil
}
