package dao

import (
	"context"
	"fmt"
	"retoBack/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type departamentosDao struct{}

var departamentos = departamentosDao{}

func Departamentos() *departamentosDao {
	return &departamentos
}

func (d *departamentosDao) table() string {
	return "departamentos"
}

func (d *departamentosDao) GetById(ctx context.Context, id int) (*entity.Departamento, error) {
	var dep *entity.Departamento
	err := g.Model(d.table()).Where("id", id).Scan(&dep)
	return dep, err
}

// ✅ OPCIÓN 1: Solo datos básicos (PARA SELECTS/FORMULARIOS)
func (d *departamentosDao) GetAllBasic(ctx context.Context) ([]*entity.Departamento, error) {
	var deps []*entity.Departamento
	err := g.Model(d.table()).
		Fields("id, nombre, codigo, descripcion, fecha_creacion, fecha_modificacion").
		Order("nombre ASC").
		Scan(&deps)

	fmt.Printf("✅ DAO - Departamentos básicos cargados: %d\n", len(deps))
	return deps, err
}

// ✅ OPCIÓN 2: Con empleos (PARA REPORTES/ESTADÍSTICAS)
func (d *departamentosDao) GetAllWithEmpleos(ctx context.Context) ([]*entity.Departamento, error) {
	var deps []*entity.Departamento
	err := g.Model(d.table()).
		With("BolsaEmpleos"). // ✅ CARGA la relación con empleos
		Order("nombre ASC").
		Scan(&deps)

	fmt.Printf("✅ DAO - Departamentos con empleos: %d\n", len(deps))
	for _, dep := range deps {
		fmt.Printf("   - %s: %d empleos en bolsa\n", dep.Nombre, len(dep.BolsaEmpleos))
	}
	return deps, err
}

// ✅ OPCIÓN 3: Departamento específico con sus empleos - CORREGIDO
func (d *departamentosDao) GetWithEmpleos(ctx context.Context, id int) (*entity.Departamento, error) {
	var dep entity.Departamento

	// ✅ PRIMERO: Cargar el departamento básico
	err := g.Model(d.table()).
		Where("id", id).
		Scan(&dep)

	if err != nil {
		return nil, err
	}

	// ✅ SEGUNDO: Cargar los empleos manualmente para evitar error de mapeo
	var empleos []entity.BolsaEmpleo
	err = g.Model("bolsa_empleos").
		Where("departamento_id", id).
		Scan(&empleos)

	if err != nil {
		return nil, err
	}

	// ✅ ASIGNAR los empleos al departamento
	dep.BolsaEmpleos = empleos

	fmt.Printf("✅ DAO - Departamento %s cargado con %d empleos\n", dep.Nombre, len(dep.BolsaEmpleos))
	return &dep, nil
}

// ✅ OPCIÓN 4: Método alternativo sin relaciones (PARA ELIMINACIÓN)
// ✅ OPCIÓN 4: Método alternativo sin relaciones (PARA ELIMINACIÓN) - CORREGIDO
func (d *departamentosDao) GetWithEmpleosCount(ctx context.Context, id int) (int, error) {
	count, err := g.Model("bolsa_empleos"). // ✅ count y err
						Where("departamento_id", id).
						Count()

	if err != nil {
		return 0, err
	}

	fmt.Printf("✅ DAO - Departamento ID %d tiene %d empleos en bolsa\n", id, count)
	return count, nil
}
func (d *departamentosDao) Create(ctx context.Context, dep *entity.Departamento) error {
	fmt.Printf("🔄 DAO - Insertando departamento: Nombre='%s', Código='%s', ID actual=%d\n",
		dep.Nombre, dep.Codigo, dep.ID)

	// ✅ FORZAR ID A CERO para que la base de datos lo genere automáticamente
	dep.ID = 0

	// ✅ Crear data sin el ID explícito
	data := map[string]interface{}{
		"nombre":             dep.Nombre,
		"codigo":             dep.Codigo,
		"descripcion":        dep.Descripcion,
		"fecha_creacion":     gtime.Now(),
		"fecha_modificacion": gtime.Now(),
	}

	fmt.Printf("📤 DAO - Datos a insertar: %+v\n", data)

	result, err := g.Model(d.table()).Data(data).Insert()
	if err != nil {
		fmt.Printf("❌ DAO - Error en insert: %v\n", err)
		return err
	}

	// ✅ Obtener el ID generado por la base de datos
	id, err := result.LastInsertId()
	if err == nil {
		dep.ID = uint(id)
		fmt.Printf("✅ DAO - Insert exitoso - Nuevo ID: %d\n", dep.ID)
	} else {
		fmt.Printf("⚠️ DAO - Insert exitoso pero no se pudo obtener ID: %v\n", err)
	}

	return nil
}

func (d *departamentosDao) Update(ctx context.Context, dep *entity.Departamento) error {
	_, err := g.Model(d.table()).Data(dep).Where("id", dep.ID).Update()
	return err
}

func (d *departamentosDao) Delete(ctx context.Context, id int) error {
	_, err := g.Model(d.table()).Where("id", id).Delete()
	return err
}

// ✅ MÉTODO COMPATIBLE (por defecto sin empleos)
func (d *departamentosDao) GetAll(ctx context.Context) ([]*entity.Departamento, error) {
	return d.GetAllBasic(ctx)
}

func (d *departamentosDao) FindByName(ctx context.Context, nombre string) (*entity.Departamento, error) {
	var departamento entity.Departamento
	err := g.Model(d.table()).Where("nombre", nombre).Scan(&departamento)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if departamento.ID == 0 {
		return nil, nil
	}

	return &departamento, nil
}

func (d *departamentosDao) FindByCodigo(ctx context.Context, codigo string) (*entity.Departamento, error) {
	var departamento entity.Departamento
	err := g.Model(d.table()).Where("codigo", codigo).Scan(&departamento)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if departamento.ID == 0 {
		return nil, nil
	}

	return &departamento, nil
}
