package dao

import (
	"context"
	"fmt"
	"retoBack/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type empleadosDao struct{}

var empleados = empleadosDao{}

func Empleados() *empleadosDao {
	return &empleados
}

func (e *empleadosDao) table() string {
	return "empleados"
}

// ----------------------
// CREAR EMPLEADO
// ----------------------
func (e *empleadosDao) Create(ctx context.Context, empleado *entity.Empleados) error {
	fmt.Printf("💾 [EMPLEADO DAO] Create - Creando empleado: %s %s\n", empleado.Nombre, empleado.Apellido)

	result, err := g.DB().Model("empleados").Data(g.Map{
		"nombre":             empleado.Nombre,
		"apellido":           empleado.Apellido,
		"documento_tipo":     empleado.DocumentoTipo,
		"documento_numero":   empleado.DocumentoNumero,
		"correo_electronico": empleado.CorreoElectronico,
		"ciudad":             empleado.Ciudad,
		"direccion":          empleado.Direccion,
		"telefono":           empleado.Telefono,
		"bolsa_empleo_id":    empleado.BolsaEmpleoID,
		"fecha_creacion":     time.Now(),
	}).Insert()

	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] Create - Error: %v\n", err)
		return gerror.Wrap(err, "Error al crear empleado en la base de datos")
	}

	// Obtener el ID generado
	if id, err := result.LastInsertId(); err == nil {
		empleado.ID = int(id)
	}

	fmt.Printf("✅ [EMPLEADO DAO] Create - Empleado creado: ID=%d\n", empleado.ID)
	return nil
}

// ----------------------
// ACTUALIZAR EMPLEADO
// ----------------------
func (e *empleadosDao) Update(ctx context.Context, empleado *entity.Empleados) error {
	fmt.Printf("💾 [EMPLEADO DAO] Update - Actualizando empleado: ID=%d\n", empleado.ID)

	_, err := g.DB().Model("empleados").Data(g.Map{
		"nombre":             empleado.Nombre,
		"apellido":           empleado.Apellido,
		"documento_tipo":     empleado.DocumentoTipo,
		"documento_numero":   empleado.DocumentoNumero,
		"correo_electronico": empleado.CorreoElectronico,
		"ciudad":             empleado.Ciudad,
		"direccion":          empleado.Direccion,
		"telefono":           empleado.Telefono,
		"bolsa_empleo_id":    empleado.BolsaEmpleoID,
		"fecha_modificacion": time.Now(),
	}).Where("id", empleado.ID).Update()

	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] Update - Error: %v\n", err)
		return gerror.Wrap(err, "Error al actualizar empleado")
	}

	fmt.Printf("✅ [EMPLEADO DAO] Update - Empleado actualizado: ID=%d\n", empleado.ID)
	return nil
}

// ----------------------
// ELIMINAR EMPLEADO
// ----------------------
func (e *empleadosDao) Delete(ctx context.Context, id int) error {
	fmt.Printf("💾 [EMPLEADO DAO] Delete - Eliminando empleado: ID=%d\n", id)

	_, err := g.DB().Exec(ctx, "DELETE FROM empleados WHERE id=$1", id)
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] Delete - Error: %v\n", err)
		return gerror.Wrap(err, "Error al eliminar empleado")
	}

	fmt.Printf("✅ [EMPLEADO DAO] Delete - Empleado eliminado: ID=%d\n", id)
	return nil
}

// ----------------------
// OBTENER TODOS LOS EMPLEADOS
// ----------------------
func (e *empleadosDao) GetAll(ctx context.Context) ([]*entity.Empleados, error) {
	fmt.Println("🔍 [EMPLEADO DAO] GetAll - Obteniendo todos los empleados...")

	var empleadosList []*entity.Empleados

	err := g.Model(e.table()).
		OrderDesc("id").
		Scan(&empleadosList)
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetAll - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener empleados")
	}

	fmt.Printf("✅ [EMPLEADO DAO] GetAll - Encontrados %d empleados\n", len(empleadosList))

	for i, emp := range empleadosList {
		fmt.Printf("🔄 [EMPLEADO DAO] GetAll - Cargando bolsa para empleado %d: %s %s\n", i, emp.Nombre, emp.Apellido)
		if err := e.cargarBolsaEmpleo(ctx, emp); err != nil {
			return nil, err
		}
	}

	return empleadosList, nil
}

// ----------------------
// BUSCAR POR EMAIL
// ----------------------
func (e *empleadosDao) FindByEmail(ctx context.Context, email string) (*entity.Empleados, error) {
	fmt.Printf("🔍 [EMPLEADO DAO] FindByEmail - Buscando por email: %s\n", email)

	var emp entity.Empleados

	err := g.Model(e.table()).
		Where("correo_electronico = ?", email).
		Scan(&emp)

	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] FindByEmail - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al buscar empleado por email")
	}
	if emp.ID == 0 {
		fmt.Printf("❌ [EMPLEADO DAO] FindByEmail - Empleado no encontrado\n")
		return nil, nil
	}

	fmt.Printf("✅ [EMPLEADO DAO] FindByEmail - Empleado encontrado: %s %s\n", emp.Nombre, emp.Apellido)

	if err := e.cargarBolsaEmpleo(ctx, &emp); err != nil {
		return nil, err
	}

	return &emp, nil
}

// ----------------------
// BUSCAR POR DOCUMENTO
// ----------------------
func (e *empleadosDao) FindByDocumento(ctx context.Context, documento string) (*entity.Empleados, error) {
	fmt.Printf("🔍 [EMPLEADO DAO] FindByDocumento - Buscando por documento: %s\n", documento)

	var emp entity.Empleados

	err := g.Model(e.table()).
		Where("documento_numero = ?", documento).
		Scan(&emp)

	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] FindByDocumento - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al buscar empleado por documento")
	}
	if emp.ID == 0 {
		fmt.Printf("❌ [EMPLEADO DAO] FindByDocumento - Empleado no encontrado\n")
		return nil, nil
	}

	fmt.Printf("✅ [EMPLEADO DAO] FindByDocumento - Empleado encontrado: %s %s\n", emp.Nombre, emp.Apellido)

	if err := e.cargarBolsaEmpleo(ctx, &emp); err != nil {
		return nil, err
	}

	return &emp, nil
}

// ----------------------
// OBTENER POR ID - CON LOGS DETALLADOS
// ----------------------
func (e *empleadosDao) GetById(ctx context.Context, id int) (*entity.Empleados, error) {
	fmt.Printf("🔍 [EMPLEADO DAO] GetById - Buscando empleado ID: %d\n", id)

	var emp entity.Empleados
	err := g.Model(e.table()).
		Where("id = ?", id).
		Scan(&emp)

	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetById - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener empleado por ID")
	}
	if emp.ID == 0 {
		fmt.Printf("❌ [EMPLEADO DAO] GetById - Empleado no encontrado: ID=%d\n", id)
		return nil, gerror.New("Empleado no encontrado")
	}

	fmt.Printf("✅ [EMPLEADO DAO] GetById - Empleado encontrado: %s %s, BolsaEmpleoID: %d\n",
		emp.Nombre, emp.Apellido, emp.BolsaEmpleoID)

	// ✅ CARGAR BOLSA EMPLEO CON LOGS DETALLADOS
	if err := e.cargarBolsaEmpleo(ctx, &emp); err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetById - Error cargando bolsa: %v\n", err)
		return nil, err
	}

	fmt.Printf("🔍 [EMPLEADO DAO] GetById - Después de cargarBolsaEmpleo: BolsaEmpleo es nil: %t\n",
		emp.BolsaEmpleo == nil)

	if emp.BolsaEmpleo != nil {
		fmt.Printf("✅ [EMPLEADO DAO] GetById - BolsaEmpleo cargada: Puesto=%s, Salario=%.2f\n",
			emp.BolsaEmpleo.Puesto, emp.BolsaEmpleo.Salario)
	} else {
		fmt.Printf("❌ [EMPLEADO DAO] GetById - BolsaEmpleo NO cargada para empleado ID: %d\n", id)
	}

	return &emp, nil
}

// ----------------------
// BUSCAR CON FILTROS
// ----------------------
func (e *empleadosDao) FindWithFilters(ctx context.Context, filters g.Map) ([]*entity.Empleados, error) {
	fmt.Println("🔍 [EMPLEADO DAO] FindWithFilters - Buscando con filtros...")

	var empleadosList []*entity.Empleados
	model := g.Model(e.table())

	// Usar índices de parámetros correctos
	paramIndex := 1
	if nombre, ok := filters["nombre"]; ok && nombre != "" {
		model = model.Where("nombre ILIKE ?", "%"+nombre.(string)+"%")
		paramIndex++
	}
	if apellido, ok := filters["apellido"]; ok && apellido != "" {
		model = model.Where("apellido ILIKE ?", "%"+apellido.(string)+"%")
		paramIndex++
	}
	if ciudad, ok := filters["ciudad"]; ok && ciudad != "" {
		model = model.Where("ciudad = ?", ciudad)
		paramIndex++
	}
	if documentoTipo, ok := filters["documento_tipo"]; ok && documentoTipo != "" {
		model = model.Where("documento_tipo = ?", documentoTipo)
	}

	err := model.OrderDesc("id").Scan(&empleadosList)
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] FindWithFilters - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al buscar empleados con filtros")
	}

	fmt.Printf("✅ [EMPLEADO DAO] FindWithFilters - Encontrados %d empleados\n", len(empleadosList))

	for _, emp := range empleadosList {
		if err := e.cargarBolsaEmpleo(ctx, emp); err != nil {
			return nil, err
		}
	}

	return empleadosList, nil
}

// ----------------------
// VERIFICAR EXISTENCIA
// ----------------------
func (e *empleadosDao) Exists(ctx context.Context, documentoNumero, email string) (bool, error) {
	fmt.Printf("🔍 [EMPLEADO DAO] Exists - Verificando existencia: documento=%s, email=%s\n", documentoNumero, email)

	count, err := g.Model(e.table()).
		Where("(documento_numero=$1 OR correo_electronico=$2)", documentoNumero, email).
		Count()
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] Exists - Error: %v\n", err)
		return false, gerror.Wrap(err, "Error al verificar existencia de empleado")
	}

	exists := count > 0
	fmt.Printf("✅ [EMPLEADO DAO] Exists - Resultado: %t\n", exists)
	return exists, nil
}

// ----------------------
// CARGAR BOLSA DE EMPLEO RELACIONADA - CON LOGS DETALLADOS
// ----------------------
func (e *empleadosDao) cargarBolsaEmpleo(ctx context.Context, empleado *entity.Empleados) error {
	if empleado.BolsaEmpleoID > 0 {
		fmt.Printf("🔄 [EMPLEADO DAO] cargarBolsaEmpleo - Cargando bolsa ID: %d para empleado: %s %s\n",
			empleado.BolsaEmpleoID, empleado.Nombre, empleado.Apellido)

		var bolsa entity.BolsaEmpleo
		err := g.Model("bolsa_empleos").
			Where("id = ?", empleado.BolsaEmpleoID).
			Scan(&bolsa)

		if err != nil {
			fmt.Printf("❌ [EMPLEADO DAO] cargarBolsaEmpleo - Error: %v\n", err)
			return gerror.Wrap(err, "Error al cargar bolsa de empleo")
		}

		if bolsa.ID == 0 {
			fmt.Printf("⚠️ [EMPLEADO DAO] cargarBolsaEmpleo - Bolsa no encontrada: ID=%d\n", empleado.BolsaEmpleoID)
			return nil
		}

		fmt.Printf("✅ [EMPLEADO DAO] cargarBolsaEmpleo - Bolsa cargada: ID=%d, Puesto=%s, Salario=%.2f\n",
			bolsa.ID, bolsa.Puesto, bolsa.Salario)

		empleado.BolsaEmpleo = &bolsa
	} else {
		fmt.Printf("⚠️ [EMPLEADO DAO] cargarBolsaEmpleo - Empleado sin BolsaEmpleoID: %s %s\n",
			empleado.Nombre, empleado.Apellido)
	}
	return nil
}

// ----------------------
// ESTADÍSTICAS
// ----------------------
func (e *empleadosDao) GetStats(ctx context.Context) (g.Map, error) {
	fmt.Println("📊 [EMPLEADO DAO] GetStats - Obteniendo estadísticas...")

	total, err := g.Model(e.table()).Count()
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetStats - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener estadísticas")
	}

	var ciudadStats []g.Map
	err = g.Model(e.table()).
		Fields("ciudad, COUNT(*) as count").
		Group("ciudad").
		Scan(&ciudadStats)
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetStats - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener estadísticas por ciudad")
	}

	conBolsa, err := g.Model(e.table()).
		Where("bolsa_empleo_id IS NOT NULL").
		Count()
	if err != nil {
		fmt.Printf("❌ [EMPLEADO DAO] GetStats - Error: %v\n", err)
		return nil, gerror.Wrap(err, "Error al obtener estadísticas de bolsa")
	}

	stats := g.Map{
		"total_empleados":      total,
		"empleados_por_ciudad": ciudadStats,
		"con_bolsa_empleo":     conBolsa,
		"sin_bolsa_empleo":     total - conBolsa,
	}

	fmt.Printf("✅ [EMPLEADO DAO] GetStats - Estadísticas obtenidas: Total=%d, ConBolsa=%d\n", total, conBolsa)
	return stats, nil
}
