package service

import (
	"context"
	"fmt"
	"retoBack/internal/dao"
	"retoBack/internal/model/entity"
	"strings"
	"unicode"

	"github.com/gogf/gf/v2/errors/gerror"
)

type departamentosService struct{}

var Departamentos = departamentosService{}

func (d *departamentosService) GetById(ctx context.Context, id int) (*entity.Departamento, error) {
	return dao.Departamentos().GetById(ctx, id)
}

func (d *departamentosService) GetAllBasic(ctx context.Context) ([]*entity.Departamento, error) {
	deps, err := dao.Departamentos().GetAllBasic(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("✅ Service - Departamentos básicos cargados: %d\n", len(deps))
	for i, dep := range deps {
		fmt.Printf("   %d. ID: %d, Nombre: %s, Código: %s\n", i+1, dep.ID, dep.Nombre, dep.Codigo)
	}
	return deps, nil
}

func (d *departamentosService) GetAllWithEmpleos(ctx context.Context) ([]*entity.Departamento, error) {
	deps, err := dao.Departamentos().GetAllWithEmpleos(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("✅ Service - Departamentos con empleos: %d\n", len(deps))
	for _, dep := range deps {
		fmt.Printf("   - %s: %d empleos en bolsa\n", dep.Nombre, len(dep.BolsaEmpleos))
	}
	return deps, nil
}

func (d *departamentosService) GetWithEmpleos(ctx context.Context, id int) (*entity.Departamento, error) {
	dep, err := dao.Departamentos().GetWithEmpleos(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("✅ Service - Departamento %s cargado con %d empleos\n", dep.Nombre, len(dep.BolsaEmpleos))
	return dep, nil
}

func (d *departamentosService) GetAll(ctx context.Context) ([]*entity.Departamento, error) {
	return d.GetAllBasic(ctx)
}

func (d *departamentosService) Create(ctx context.Context, departamento *entity.Departamento) error {
	fmt.Printf("=== CREAR DEPARTAMENTO ===\n")
	fmt.Printf("📥 Departamento a crear: %s\n", departamento.Nombre)

	fmt.Println("=== DEPARTAMENTOS EXISTENTES ===")
	todos, _ := dao.Departamentos().GetAllBasic(ctx)
	for i, dep := range todos {
		fmt.Printf("%d. ID: %d, Nombre: '%s', Código: '%s'\n",
			i+1, dep.ID, dep.Nombre, dep.Codigo)
	}

	existente, err := dao.Departamentos().FindByName(ctx, departamento.Nombre)
	if err != nil {
		fmt.Printf("❌ Error buscando por nombre: %v\n", err)
		return err
	}

	if existente != nil {
		fmt.Printf("❌ YA EXISTE un departamento con nombre '%s' (ID: %d)\n",
			departamento.Nombre, existente.ID)
		return gerror.New("El nombre ya está registrado")
	} else {
		fmt.Printf("✅ Nombre '%s' disponible\n", departamento.Nombre)
	}

	if departamento.Codigo == "" {
		codigo, err := d.generarCodigo(ctx, departamento.Nombre)
		if err != nil {
			return err
		}
		departamento.Codigo = codigo
		fmt.Printf("✅ Código generado: %s\n", departamento.Codigo)
	}

	existentePorCodigo, err := dao.Departamentos().FindByCodigo(ctx, departamento.Codigo)
	if err != nil {
		return err
	}
	if existentePorCodigo != nil {
		fmt.Printf("❌ Código '%s' ya existe, generando alternativo\n", departamento.Codigo)
		codigoAlternativo, err := d.generarCodigoAlternativo(ctx, departamento.Nombre)
		if err != nil {
			return err
		}
		departamento.Codigo = codigoAlternativo
		fmt.Printf("✅ Código alternativo: %s\n", departamento.Codigo)
	}

	fmt.Printf("🔄 Creando departamento: %s (%s)\n", departamento.Nombre, departamento.Codigo)
	return dao.Departamentos().Create(ctx, departamento)
}

func (d *departamentosService) Update(ctx context.Context, id int, departamento *entity.Departamento) error {
	// ✅ Forzar ID desde la URL para consistencia
	departamento.ID = uint(id)

	fmt.Printf("=== ACTUALIZAR DEPARTAMENTO ID: %d ===\n", id)
	fmt.Printf("📥 Nuevos datos: Nombre='%s', Código='%s'\n", departamento.Nombre, departamento.Codigo)

	// ✅ VALIDAR QUE EL CÓDIGO NO ESTÉ VACÍO
	if departamento.Codigo == "" {
		return gerror.New("El código del departamento es requerido")
	}

	// ✅ 1. Validar nombre único (excluyendo el actual)
	existentePorNombre, err := dao.Departamentos().FindByName(ctx, departamento.Nombre)
	if err != nil {
		return err
	}
	if existentePorNombre != nil && existentePorNombre.ID != departamento.ID {
		return gerror.Newf("Ya existe otro departamento con el nombre '%s'", departamento.Nombre)
	}

	// ✅ 2. Validar código único (excluyendo el actual)
	existentePorCodigo, err := dao.Departamentos().FindByCodigo(ctx, departamento.Codigo)
	if err != nil {
		return err
	}
	if existentePorCodigo != nil && existentePorCodigo.ID != departamento.ID {
		return gerror.Newf("El código '%s' ya está en uso por otro departamento", departamento.Codigo)
	}

	fmt.Printf("✅ Validaciones pasadas - Actualizando departamento\n")
	return dao.Departamentos().Update(ctx, departamento)
}

// ✅ MÉTODO DELETE ACTUALIZADO CON VALIDACIÓN PROFESIONAL
func (d *departamentosService) Delete(ctx context.Context, id int) error {
	fmt.Printf("=== VALIDANDO ELIMINACIÓN DEPARTAMENTO ID: %d ===\n", id)

	// ✅ USAR EL MÉTODO CORREGIDO
	departamento, err := dao.Departamentos().GetWithEmpleos(ctx, id)
	if err != nil {
		fmt.Printf("❌ Error obteniendo departamento: %v\n", err)
		return err
	}

	if departamento == nil {
		fmt.Printf("❌ Departamento no encontrado\n")
		return gerror.New("Departamento no encontrado")
	}

	fmt.Printf("🔍 Departamento '%s' tiene %d empleos en bolsa\n",
		departamento.Nombre, len(departamento.BolsaEmpleos))

	// ✅ Validar si tiene empleos en bolsa
	if len(departamento.BolsaEmpleos) > 0 {
		empleosActivos := d.contarEmpleosActivos(departamento.BolsaEmpleos)
		totalEmpleos := len(departamento.BolsaEmpleos)

		fmt.Printf("❌ NO SE PUEDE ELIMINAR - Departamento '%s' tiene %d empleos (%d activos)\n",
			departamento.Nombre, totalEmpleos, empleosActivos)

		var mensajeError string
		if empleosActivos > 0 {
			mensajeError = fmt.Sprintf(
				"No se puede eliminar el departamento '%s' porque tiene %d empleo(s) activo(s) en la bolsa de empleo",
				departamento.Nombre,
				empleosActivos,
			)
		} else {
			mensajeError = fmt.Sprintf(
				"No se puede eliminar el departamento '%s' porque tiene %d empleo(s) históricos en la bolsa",
				departamento.Nombre,
				totalEmpleos,
			)
		}

		return gerror.New(mensajeError)
	}

	fmt.Printf("✅ Departamento '%s' sin empleos en bolsa - Procediendo a eliminar\n", departamento.Nombre)
	return dao.Departamentos().Delete(ctx, id)
}

// ✅ MÉTODO AUXILIAR PARA CONTAR EMPLEOS ACTIVOS
func (d *departamentosService) contarEmpleosActivos(empleos []entity.BolsaEmpleo) int {
	count := 0
	for _, empleo := range empleos {
		if empleo.Estado == "DISPONIBLE" || empleo.Estado == "OCUPADO" {
			count++
		}
	}
	return count
}

// ✅ MÉTODO PARA VERIFICAR SI SE PUEDE ELIMINAR (PARA FRONTEND)
func (d *departamentosService) CanDelete(ctx context.Context, id int) (bool, string, error) {
	fmt.Printf("=== VERIFICANDO SI SE PUEDE ELIMINAR DEPARTAMENTO ID: %d ===\n", id)

	// 1. Primero obtener el departamento CON sus empleos en bolsa
	departamento, err := dao.Departamentos().GetWithEmpleos(ctx, id)
	if err != nil {
		fmt.Printf("❌ Error obteniendo departamento: %v\n", err)
		return false, "", err
	}

	if departamento == nil {
		fmt.Printf("❌ Departamento no encontrado\n")
		return false, "Departamento no encontrado", nil
	}

	fmt.Printf("🔍 Departamento '%s' tiene %d empleos en bolsa\n",
		departamento.Nombre, len(departamento.BolsaEmpleos))

	// 2. Validar si tiene empleos en bolsa
	if len(departamento.BolsaEmpleos) > 0 {
		empleosActivos := d.contarEmpleosActivos(departamento.BolsaEmpleos)
		totalEmpleos := len(departamento.BolsaEmpleos)

		fmt.Printf("❌ NO SE PUEDE ELIMINAR - Departamento '%s' tiene %d empleos (%d activos)\n",
			departamento.Nombre, totalEmpleos, empleosActivos)

		var mensaje string
		if empleosActivos > 0 {
			mensaje = fmt.Sprintf("No se puede eliminar porque tiene %d empleo(s) activo(s) en la bolsa", empleosActivos)
		} else {
			mensaje = fmt.Sprintf("No se puede eliminar porque tiene %d empleo(s) históricos en la bolsa", totalEmpleos)
		}

		return false, mensaje, nil
	}

	fmt.Printf("✅ Departamento '%s' sin empleos en bolsa - Se puede eliminar\n", departamento.Nombre)
	return true, "Departamento se puede eliminar", nil
}

// ✅ GENERAR CÓDIGO
func (d *departamentosService) generarCodigo(ctx context.Context, nombre string) (string, error) {
	var codigo strings.Builder

	for _, char := range nombre {
		if unicode.IsLetter(char) {
			codigo.WriteRune(unicode.ToUpper(char))
			if codigo.Len() >= 3 {
				break
			}
		}
	}

	for codigo.Len() < 3 {
		codigo.WriteString("X")
	}

	return codigo.String(), nil
}

// ✅ GENERAR CÓDIGO ALTERNATIVO
func (d *departamentosService) generarCodigoAlternativo(ctx context.Context, nombre string) (string, error) {
	baseCodigo, _ := d.generarCodigo(ctx, nombre)

	for i := 1; i <= 9; i++ {
		codigoAlternativo := baseCodigo + string(rune('0'+i))

		existente, err := dao.Departamentos().FindByCodigo(ctx, codigoAlternativo)
		if err != nil {
			return "", err
		}
		if existente == nil {
			return codigoAlternativo, nil
		}
	}

	return baseCodigo + "A", nil
}
