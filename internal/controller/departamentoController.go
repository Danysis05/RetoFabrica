package controller

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/model/mapper"
	"retoBack/internal/service"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ✅ PARA FORMULARIOS (sin empleos) - ENDPOINT PRINCIPAL
func DepartamentoShowAll(r *ghttp.Request) {
	departamentos, err := service.Departamentos.GetAllBasic(r.Context())
	if err != nil {
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando %d departamentos (sin empleos)\n", len(departamentos))

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    departamentos,
		"count":   len(departamentos),
	})
}

// ✅ PARA REPORTES (con empleos) - ENDPOINT ADICIONAL
func DepartamentoShowAllWithEmpleos(r *ghttp.Request) {
	departamentos, err := service.Departamentos.GetAllWithEmpleos(r.Context())
	if err != nil {
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando %d departamentos (con empleos)\n", len(departamentos))

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    departamentos,
		"count":   len(departamentos),
	})
}

// ✅ DEPARTAMENTO ESPECÍFICO CON EMPLEOS
func DepartamentoShowWithEmpleos(r *ghttp.Request) {
	idStr := r.GetRouter("id").String() // ✅ CORREGIDO: Usar GetRouter para parámetros URL
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		r.Response.WriteStatus(400, g.Map{"error": "ID inválido"})
		return
	}

	departamento, err := service.Departamentos.GetWithEmpleos(r.Context(), id)
	if err != nil {
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if departamento == nil {
		r.Response.WriteStatus(404, g.Map{
			"success": false,
			"message": "Departamento no encontrado",
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando departamento %s con %d empleos\n",
		departamento.Nombre, len(departamento.BolsaEmpleos))

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    departamento,
	})
}

// Crear departamento
func DepartamentoCreate(r *ghttp.Request) {
	fmt.Println("=== DEPARTAMENTO CREATE ===")
	fmt.Println("Body recibido:", r.GetBodyString())

	var req dto.DepartamentoDTO
	if err := r.Parse(&req); err != nil {
		fmt.Println("❌ Error parseando request:", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ DTO parseado: %+v\n", req)

	departamentoEntity := mapper.ToDepartamentoEntity(&req)
	fmt.Printf("✅ Entidad mapeada: %+v\n", departamentoEntity)

	if err := service.Departamentos.Create(r.Context(), departamentoEntity); err != nil {
		fmt.Println("❌ Error en servicio Create:", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	fmt.Println("✅ Departamento creado exitosamente")
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Departamento creado con éxito",
		"data":    departamentoEntity,
	})
}

// ✅ Actualizar departamento - CORREGIDO
func DepartamentoUpdate(r *ghttp.Request) {
	idStr := r.GetRouter("id").String() // ✅ CORREGIDO: Usar GetRouter para parámetros URL
	fmt.Printf("🔄 Controller - Actualizando departamento ID: '%s'\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Printf("❌ ID inválido: '%s'\n", idStr)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	var req dto.DepartamentoDTO
	if err := r.Parse(&req); err != nil {
		fmt.Printf("❌ Error parseando request: %v\n", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida: " + err.Error(),
		})
		return
	}

	// ✅ UPDATE: El código ES REQUERIDO (usar el existente)
	if strings.TrimSpace(req.Codigo) == "" {
		fmt.Printf("❌ UPDATE: Código vacío no permitido en actualización\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "El código del departamento es requerido para actualizar",
		})
		return
	}

	fmt.Printf("📥 UPDATE - ID: %d, Nombre: '%s', Código: '%s'\n",
		id, req.Nombre, req.Codigo)

	departamentoEntity := mapper.ToDepartamentoEntity(&req)

	if err := service.Departamentos.Update(r.Context(), id, departamentoEntity); err != nil {
		fmt.Printf("❌ Error en servicio Update: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("✅ Departamento %d actualizado exitosamente\n", id)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Departamento actualizado con éxito",
	})
}

// ✅ Eliminar departamento - CORREGIDO
func DepartamentoDelete(r *ghttp.Request) {
	idStr := r.GetRouter("id").String() // ✅ CORREGIDO: Usar GetRouter para parámetros URL
	fmt.Printf("🗑️ Controller - Eliminando departamento ID: '%s'\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Printf("❌ ID inválido: '%s'\n", idStr)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID debe ser un número válido mayor a 0",
		})
		return
	}

	if err := service.Departamentos.Delete(r.Context(), id); err != nil {
		fmt.Printf("❌ Error eliminando departamento %d: %v\n", id, err)

		// ✅ MEJOR MANEJO DE ERRORES
		errorMsg := err.Error()
		statusCode := 500

		if strings.Contains(errorMsg, "No se puede eliminar") {
			statusCode = 423 // Locked
		}

		r.Response.WriteStatus(statusCode, g.Map{
			"success":      false,
			"error":        errorMsg,
			"errorType":    "BLOCKED_BY_EMPLOYMENTS",
			"departmentId": id,
		})
		return
	}

	fmt.Printf("✅ Departamento %d eliminado exitosamente\n", id)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Departamento eliminado con éxito",
	})
}

// ✅ VERIFICAR SI SE PUEDE ELIMINAR - CORREGIDO
func DepartamentoCanDelete(r *ghttp.Request) {
	idStr := r.GetRouter("id").String() // ✅ CORREGIDO: Usar GetRouter para parámetros URL
	fmt.Printf("🔍 Controller - Verificando si se puede eliminar departamento ID: '%s'\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	canDelete, reason, err := service.Departamentos.CanDelete(r.Context(), id)
	if err != nil {
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":      true,
		"canDelete":    canDelete,
		"reason":       reason,
		"departmentId": id,
	})
}
