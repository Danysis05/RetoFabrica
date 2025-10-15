package dao

import (
	"context"
	"fmt"
	"retoBack/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type nominaDao struct{}

var Nominas = &nominaDao{}

func (d *nominaDao) table() string {
	return "nomina"
}

// ----------------------
// CREAR NÓMINA
// ----------------------
func (d *nominaDao) Create(ctx context.Context, nomina *entity.Nomina) error {
	fmt.Printf("💾 [NOMINA DAO] Create - Guardando nómina: EmpleadoID=%d\n", nomina.EmpleadoId)

	// ✅ LOG DETALLADO DE LO QUE SE VA A GUARDAR
	fmt.Printf("🔍 [NOMINA DAO] Create - Datos de nómina a guardar:\n")
	fmt.Printf("   - EmpleadoID: %d\n", nomina.EmpleadoId)
	fmt.Printf("   - SalarioBase: %.2f\n", nomina.SalarioBase)
	fmt.Printf("   - HorasExtras: %.2f\n", nomina.HorasExtras)
	fmt.Printf("   - Bonificaciones: %.2f\n", nomina.Bonificaciones)
	fmt.Printf("   - Deducciones: %.2f\n", nomina.Deducciones)
	fmt.Printf("   - TotalPago: %.2f\n", nomina.TotalPago)

	result, err := g.DB().Model(d.table()).Data(g.Map{
		"empleado_id":    nomina.EmpleadoId,
		"fecha_pago":     nomina.FechaPago,
		"salario_base":   nomina.SalarioBase,
		"horas_extras":   nomina.HorasExtras,
		"bonificaciones": nomina.Bonificaciones,
		"deducciones":    nomina.Deducciones,
		"total_pago":     nomina.TotalPago,
	}).Insert()

	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] Create - Error: %v\n", err)
		return gerror.Wrap(err, "Error al crear nómina")
	}

	// Obtener ID generado
	if id, err := result.LastInsertId(); err == nil {
		nomina.Id = int(id)
	}

	fmt.Printf("✅ [NOMINA DAO] Create - Nómina guardada: ID=%d\n", nomina.Id)
	return nil
}

// ----------------------
// ACTUALIZAR NÓMINA
// ----------------------
func (d *nominaDao) Update(ctx context.Context, nomina *entity.Nomina) error {
	fmt.Printf("💾 [NOMINA DAO] Update - Actualizando nómina: ID=%d\n", nomina.Id)

	_, err := g.DB().Model(d.table()).Data(g.Map{
		"empleado_id":    nomina.EmpleadoId,
		"fecha_pago":     nomina.FechaPago,
		"salario_base":   nomina.SalarioBase,
		"horas_extras":   nomina.HorasExtras,
		"bonificaciones": nomina.Bonificaciones,
		"deducciones":    nomina.Deducciones,
		"total_pago":     nomina.TotalPago,
	}).Where("id", nomina.Id).Update()

	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] Update - Error: %v\n", err)
		return gerror.Wrap(err, "Error al actualizar nómina")
	}

	fmt.Printf("✅ [NOMINA DAO] Update - Nómina actualizada: ID=%d\n", nomina.Id)
	return nil
}

// ----------------------
// ELIMINAR NÓMINA
// ----------------------
func (d *nominaDao) Delete(ctx context.Context, id int) error {
	fmt.Printf("💾 [NOMINA DAO] Delete - Eliminando nómina: ID=%d\n", id)

	_, err := g.DB().Model(d.table()).Where("id", id).Delete()
	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] Delete - Error: %v\n", err)
		return gerror.Wrap(err, "Error al eliminar nómina")
	}

	fmt.Printf("✅ [NOMINA DAO] Delete - Nómina eliminada: ID=%d\n", id)
	return nil
}

// ----------------------
// OBTENER TODAS LAS NÓMINAS
// ----------------------
func (d *nominaDao) GetAll(ctx context.Context) ([]*entity.Nomina, error) {
	fmt.Println("🔍 [NOMINA DAO] GetAll - Obteniendo todas las nóminas...")

	var nominas []*entity.Nomina

	err := g.DB().Model(d.table()).
		OrderDesc("id").
		Scan(&nominas)

	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] GetAll - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener nóminas")
	}

	fmt.Printf("✅ [NOMINA DAO] GetAll - Cargadas %d nóminas\n", len(nominas))

	// Cargar relaciones manualmente
	for i, nomina := range nominas {
		fmt.Printf("🔄 [NOMINA DAO] GetAll - Cargando empleado para nómina %d (ID=%d)\n", i, nomina.Id)
		if err := d.cargarEmpleado(ctx, nomina); err != nil {
			return nil, err
		}
	}

	return nominas, nil
}

// ----------------------
// OBTENER NÓMINA POR ID
// ----------------------
func (d *nominaDao) GetById(ctx context.Context, id int) (*entity.Nomina, error) {
	fmt.Printf("🔍 [NOMINA DAO] GetById - Buscando nómina: ID=%d\n", id)

	var nomina entity.Nomina
	err := g.DB().Model(d.table()).
		Where("id", id).
		Scan(&nomina)

	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] GetById - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener nómina")
	}

	if nomina.Id == 0 {
		fmt.Printf("❌ [NOMINA DAO] GetById - Nómina no encontrada: ID=%d\n", id)
		return nil, nil
	}

	fmt.Printf("✅ [NOMINA DAO] GetById - Nómina encontrada: ID=%d, EmpleadoID=%d\n", nomina.Id, nomina.EmpleadoId)

	// Cargar empleado relacionado
	if err := d.cargarEmpleado(ctx, &nomina); err != nil {
		return nil, err
	}

	return &nomina, nil
}

// ----------------------
// OBTENER NÓMINAS POR EMPLEADO
// ----------------------
func (d *nominaDao) GetByEmpleadoID(ctx context.Context, empleadoID int) ([]*entity.Nomina, error) {
	fmt.Printf("🔍 [NOMINA DAO] GetByEmpleadoID - Buscando nóminas del empleado: ID=%d\n", empleadoID)

	var nominas []*entity.Nomina
	err := g.DB().Model(d.table()).
		Where("empleado_id", empleadoID).
		OrderDesc("id").
		Scan(&nominas)

	if err != nil {
		fmt.Printf("❌ [NOMINA DAO] GetByEmpleadoID - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener nóminas por empleado")
	}

	fmt.Printf("✅ [NOMINA DAO] GetByEmpleadoID - Encontradas %d nóminas\n", len(nominas))
	return nominas, nil
}

// ----------------------
// CARGAR EMPLEADO RELACIONADO - CON LOGS DETALLADOS
// ----------------------
func (d *nominaDao) cargarEmpleado(ctx context.Context, nomina *entity.Nomina) error {
	if nomina.EmpleadoId > 0 {
		fmt.Printf("🔄 [NOMINA DAO] cargarEmpleado - Cargando empleado ID: %d para nómina ID: %d\n",
			nomina.EmpleadoId, nomina.Id)

		empleado, err := Empleados().GetById(ctx, nomina.EmpleadoId)
		if err != nil {
			fmt.Printf("❌ [NOMINA DAO] cargarEmpleado - Error al cargar empleado: %v\n", err)
			return gerror.Wrap(err, "Error al cargar empleado para nómina")
		}

		if empleado == nil {
			fmt.Printf("❌ [NOMINA DAO] cargarEmpleado - Empleado no encontrado: ID=%d\n", nomina.EmpleadoId)
			return gerror.New("Empleado no encontrado")
		}

		// ✅ LOGS DETALLADOS DEL EMPLEADO CARGADO
		fmt.Printf("✅ [NOMINA DAO] cargarEmpleado - Empleado cargado: %s %s\n",
			empleado.Nombre, empleado.Apellido)
		fmt.Printf("🔍 [NOMINA DAO] cargarEmpleado - BolsaEmpleoID del empleado: %d\n", empleado.BolsaEmpleoID)
		fmt.Printf("🔍 [NOMINA DAO] cargarEmpleado - BolsaEmpleo es nil: %t\n", empleado.BolsaEmpleo == nil)

		if empleado.BolsaEmpleo != nil {
			fmt.Printf("💰 [NOMINA DAO] cargarEmpleado - BolsaEmpleo cargada: Puesto=%s, Salario=%.2f\n",
				empleado.BolsaEmpleo.Puesto, empleado.BolsaEmpleo.Salario)
		} else {
			fmt.Printf("⚠️ [NOMINA DAO] cargarEmpleado - BolsaEmpleo es NIL para empleado\n")

			// Intentar cargar manualmente si tiene BolsaEmpleoID
			if empleado.BolsaEmpleoID > 0 {
				fmt.Printf("🔄 [NOMINA DAO] cargarEmpleado - Intentando cargar bolsa manualmente, ID: %d\n",
					empleado.BolsaEmpleoID)
				bolsa, err := BolsaEmpleos().GetById(ctx, empleado.BolsaEmpleoID)
				if err != nil {
					fmt.Printf("❌ [NOMINA DAO] cargarEmpleado - Error cargando bolsa manualmente: %v\n", err)
				} else if bolsa != nil {
					empleado.BolsaEmpleo = bolsa
					fmt.Printf("✅ [NOMINA DAO] cargarEmpleado - Bolsa cargada manualmente: Puesto=%s, Salario=%.2f\n",
						bolsa.Puesto, bolsa.Salario)
				} else {
					fmt.Printf("❌ [NOMINA DAO] cargarEmpleado - Bolsa no encontrada: ID=%d\n", empleado.BolsaEmpleoID)
				}
			}
		}

		nomina.Empleado = empleado
		fmt.Printf("✅ [NOMINA DAO] cargarEmpleado - Empleado asignado a nómina\n")
	} else {
		fmt.Printf("⚠️ [NOMINA DAO] cargarEmpleado - Nómina sin EmpleadoID: ID=%d\n", nomina.Id)
	}
	return nil
}
