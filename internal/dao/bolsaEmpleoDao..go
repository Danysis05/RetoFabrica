package dao

import (
	"context"
	"fmt"
	"retoBack/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type bolsaEmpleosDao struct{}

var bolsaEmpleos = bolsaEmpleosDao{}

func BolsaEmpleos() *bolsaEmpleosDao {
	return &bolsaEmpleos
}

// Nombre de la tabla en BD
func (d *bolsaEmpleosDao) table() string {
	return "bolsa_empleos"
}

// Obtener por ID
func (d *bolsaEmpleosDao) GetById(ctx context.Context, id int) (*entity.BolsaEmpleo, error) {
	var b entity.BolsaEmpleo
	err := g.Model(d.table()).Where("id", id).Scan(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ➕ Crear un nuevo registro
func (d *bolsaEmpleosDao) Create(ctx context.Context, b *entity.BolsaEmpleo) error {
	data := g.Map{
		"puesto":            b.Puesto,
		"descripcion":       b.Descripcion,
		"salario":           b.Salario,
		"estado":            b.Estado,
		"departamento_id":   b.DepartamentoID,
		"fecha_publicacion": b.FechaPublicacion,
	}

	// Solo agregar empleado_id si es mayor que 0
	if b.EmpleadoID > 0 {
		data["empleado_id"] = b.EmpleadoID
	}
	if b.FechaOcupacion != nil && !b.FechaOcupacion.IsZero() {
		data["fecha_ocupacion"] = b.FechaOcupacion
	}
	if b.FechaCierre != nil && !b.FechaCierre.IsZero() {
		data["fecha_cierre"] = b.FechaCierre
	}

	fmt.Printf("🎯 Data para INSERT: %+v\n", data)
	fmt.Printf("🔍 Tabla: %s\n", d.table())

	result, err := g.Model(d.table()).Data(data).Insert()
	if err != nil {
		fmt.Printf("❌ Error en INSERT: %v\n", err)
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		b.ID = int(id)
		fmt.Printf("✅ INSERT exitoso - ID generado: %d\n", id)
	}
	return err
}

// ✏️ Actualizar registro existente
func (d *bolsaEmpleosDao) Update(ctx context.Context, b *entity.BolsaEmpleo) error {
	data := g.Map{
		"puesto":          b.Puesto,
		"descripcion":     b.Descripcion,
		"salario":         b.Salario,
		"estado":          b.Estado,
		"departamento_id": b.DepartamentoID,
	}

	if b.EmpleadoID > 0 {
		data["empleado_id"] = b.EmpleadoID
	} else {
		data["empleado_id"] = nil // quitar asignación si no hay empleado
	}
	if b.FechaOcupacion != nil && !b.FechaOcupacion.IsZero() {
		data["fecha_ocupacion"] = b.FechaOcupacion
	}
	if b.FechaCierre != nil && !b.FechaCierre.IsZero() {
		data["fecha_cierre"] = b.FechaCierre
	}

	_, err := g.Model(d.table()).Data(data).Where("id", b.ID).Update()
	return err
}

// ❌ Eliminar registro
func (d *bolsaEmpleosDao) Delete(ctx context.Context, id int) error {
	_, err := g.Model(d.table()).Where("id", id).Delete()
	return err
}

// 📋 Obtener todos los registros
func (d *bolsaEmpleosDao) GetAll(ctx context.Context) ([]*entity.BolsaEmpleo, error) {
	var lista []*entity.BolsaEmpleo
	err := g.Model(d.table()).
		Order("fecha_publicacion DESC").
		Scan(&lista)
	return lista, err
}

// 🔎 Buscar empleos por departamento
func (d *bolsaEmpleosDao) GetByDepartamento(ctx context.Context, departamentoID int) ([]*entity.BolsaEmpleo, error) {
	var lista []*entity.BolsaEmpleo
	err := g.Model(d.table()).
		Where("departamento_id", departamentoID).
		Order("fecha_publicacion DESC").
		Scan(&lista)
	return lista, err
}

// 🔎 Buscar empleos disponibles
func (d *bolsaEmpleosDao) GetDisponibles(ctx context.Context) ([]*entity.BolsaEmpleo, error) {
	var lista []*entity.BolsaEmpleo
	err := g.Model(d.table()).
		Where("estado", "DISPONIBLE").
		Scan(&lista)
	return lista, err
}
