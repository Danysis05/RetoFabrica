package controller

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Obtener todos los empleados
func EmpleadoShowAll(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Listando todos los empleados\n")

	empleados, err := service.Empleados.GetAllDTO(r.Context())
	if err != nil {
		fmt.Printf("❌ Error listando empleados: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando %d empleados\n", len(empleados))
	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    empleados,
		"count":   len(empleados),
	})
}

// Crear empleado
func EmpleadoCreate(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Creando nuevo empleado\n")
	fmt.Printf("📦 Body recibido: %s\n", r.GetBodyString())

	var req dto.EmpleadoDTO
	if err := r.Parse(&req); err != nil {
		fmt.Printf("❌ Error parseando request: %v\n", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida: " + err.Error(),
		})
		return
	}

	// Ignorar ID en creación
	req.ID = 0

	fmt.Printf("📥 DTO recibido del frontend:\n")
	fmt.Printf("   Nombre: %s\n", req.Nombre)
	fmt.Printf("   Apellido: %s\n", req.Apellido)
	fmt.Printf("   DocumentoTipo: %s\n", req.DocumentoTipo)
	fmt.Printf("   DocumentoNumero: %s\n", req.DocumentoNumero)
	fmt.Printf("   CorreoElectronico: %s\n", req.CorreoElectronico)
	fmt.Printf("   Ciudad: %s\n", req.Ciudad)
	fmt.Printf("   Direccion: %s\n", req.Direccion)
	fmt.Printf("   Telefono: %s\n", req.Telefono)
	fmt.Printf("   BolsaEmpleoID: %d\n", req.BolsaEmpleoID)

	// Validaciones
	if req.Nombre == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El nombre es requerido"})
		return
	}
	if req.Apellido == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El apellido es requerido"})
		return
	}
	if req.DocumentoNumero == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El número de documento es requerido"})
		return
	}
	if req.CorreoElectronico == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El correo electrónico es requerido"})
		return
	}

	// Crear usando el servicio
	if err := service.Empleados.Crear(r.Context(), &req); err != nil {
		fmt.Printf("❌ Error en servicio Crear: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al crear el empleado: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleado creado con éxito",
		"data":    req,
	})
}

// Crear empleado desde datos simples (alternativa)
func EmpleadoCreateSimple(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Creando empleado desde datos simples\n")

	// Obtener datos de los parámetros
	nombre := r.Get("nombre").String()
	apellido := r.Get("apellido").String()
	documentoTipo := r.Get("documentoTipo").String()
	documentoNumero := r.Get("documentoNumero").String()
	correoElectronico := r.Get("correoElectronico").String()
	ciudad := r.Get("ciudad").String()
	direccion := r.Get("direccion").String()
	telefono := r.Get("telefono").String()
	bolsaEmpleoID := r.Get("bolsaEmpleoID").Int()

	// Validaciones básicas
	if nombre == "" || apellido == "" || documentoNumero == "" || correoElectronico == "" {
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Nombre, apellido, documento y correo son requeridos",
		})
		return
	}

	if err := service.Empleados.CrearDesdeDatos(r.Context(),
		nombre, apellido, documentoTipo, documentoNumero,
		correoElectronico, ciudad, direccion, telefono,
		bolsaEmpleoID); err != nil {
		fmt.Printf("❌ Error creando empleado: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al crear el empleado: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Empleado creado exitosamente desde datos simples\n")
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleado creado con éxito",
	})
}

// Actualizar empleado
func EmpleadoUpdate(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Actualizando empleado\n")

	var req dto.EmpleadoDTO
	if err := r.Parse(&req); err != nil {
		fmt.Printf("❌ Error parseando request: %v\n", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		fmt.Printf("❌ ERROR: ID requerido\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID requerido para actualizar",
		})
		return
	}

	// ✅ LOG DETALLADO
	fmt.Printf("📝 Actualizando empleado ID: %d\n", req.ID)
	fmt.Printf("   Nuevo nombre: %s\n", req.Nombre)
	fmt.Printf("   Nuevo apellido: %s\n", req.Apellido)
	fmt.Printf("   Nuevo bolsaEmpleoID: %d\n", req.BolsaEmpleoID)

	// ✅ VALIDACIONES
	if req.Nombre == "" {
		fmt.Printf("❌ ERROR: Nombre está vacío\n")
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "El nombre es requerido",
		})
		return
	}

	if req.Apellido == "" {
		fmt.Printf("❌ ERROR: Apellido está vacío\n")
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "El apellido es requerido",
		})
		return
	}

	// ✅ CORREGIDO: Usar el método del service que acepta DTO
	if err := service.Empleados.Update(r.Context(), &req); err != nil {
		fmt.Printf("❌ Error actualizando empleado %d: %v\n", req.ID, err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al actualizar el empleado: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Empleado %d actualizado exitosamente\n", req.ID)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleado actualizado con éxito",
	})
}

// Eliminar empleado
func EmpleadoDelete(r *ghttp.Request) {
	// ✅ Obtener ID de query parameter: ?id=3
	idStr := r.Get("id").String()
	fmt.Printf("🗑️ Controller - Eliminando empleado ID: '%s'\n", idStr)

	if idStr == "" {
		fmt.Printf("❌ ID no proporcionado. Se debe usar: /empleados/delete?id=3\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID no proporcionado. Use: /empleados/delete?id=3",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Printf("❌ ID inválido: '%s'\n", idStr)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID debe ser un número válido mayor a 0",
		})
		return
	}

	fmt.Printf("🔍 Verificando existencia del empleado ID: %d\n", id)

	if err := service.Empleados.Delete(r.Context(), id); err != nil {
		fmt.Printf("❌ Error eliminando empleado %d: %v\n", id, err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al eliminar el empleado: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Empleado %d eliminado exitosamente\n", id)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleado eliminado con éxito",
	})
}

// Obtener empleado por ID
func EmpleadoGetByID(r *ghttp.Request) {
	idStr := r.Get("id").String()
	fmt.Printf("🔍 Controller - Obteniendo empleado ID: '%s'\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Printf("❌ ID inválido: '%s'\n", idStr)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	empleado, err := service.Empleados.GetByIdDTO(r.Context(), id)
	if err != nil {
		fmt.Printf("❌ Error obteniendo empleado %d: %v\n", id, err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al obtener el empleado: " + err.Error(),
		})
		return
	}

	if empleado == nil {
		fmt.Printf("❌ Empleado %d no encontrado\n", id)
		r.Response.WriteStatus(404, g.Map{
			"success": false,
			"error":   "Empleado no encontrado",
		})
		return
	}

	fmt.Printf("✅ Empleado %d encontrado: %s %s\n", id, empleado.Nombre, empleado.Apellido)
	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    empleado,
	})
}

// Buscar empleados con filtros
func EmpleadoSearch(r *ghttp.Request) {
	fmt.Printf("🔍 Controller - Buscando empleados con filtros\n")

	filters := make(map[string]interface{})

	// Obtener filtros de los parámetros
	if nombre := r.Get("nombre").String(); nombre != "" {
		filters["nombre"] = nombre
	}
	if apellido := r.Get("apellido").String(); apellido != "" {
		filters["apellido"] = apellido
	}
	if ciudad := r.Get("ciudad").String(); ciudad != "" {
		filters["ciudad"] = ciudad
	}
	if documentoTipo := r.Get("documentoTipo").String(); documentoTipo != "" {
		filters["documento_tipo"] = documentoTipo
	}

	fmt.Printf("   Filtros aplicados: %+v\n", filters)

	empleados, err := service.Empleados.FindWithFilters(r.Context(), filters)
	if err != nil {
		fmt.Printf("❌ Error buscando empleados: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al buscar empleados: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Encontrados %d empleados con los filtros aplicados\n", len(empleados))
	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    empleados,
		"count":   len(empleados),
		"filters": filters,
	})
}

// Verificar existencia de empleado
func EmpleadoExists(r *ghttp.Request) {
	documentoNumero := r.Get("documento").String()
	email := r.Get("email").String()

	fmt.Printf("🔍 Controller - Verificando existencia: documento=%s, email=%s\n", documentoNumero, email)

	if documentoNumero == "" && email == "" {
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Se requiere documento o email para verificar",
		})
		return
	}

	exists, err := service.Empleados.Exists(r.Context(), documentoNumero, email)
	if err != nil {
		fmt.Printf("❌ Error verificando existencia: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al verificar existencia: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"exists":  exists,
		"message": map[bool]string{
			true:  "El empleado ya existe",
			false: "El empleado no existe",
		}[exists],
	})
}

// Obtener estadísticas de empleados
func EmpleadoStats(r *ghttp.Request) {
	fmt.Printf("📊 Controller - Obteniendo estadísticas de empleados\n")

	stats, err := service.Empleados.GetStats(r.Context())
	if err != nil {
		fmt.Printf("❌ Error obteniendo estadísticas: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al obtener estadísticas: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Estadísticas obtenidas exitosamente\n")
	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    stats,
	})
}
