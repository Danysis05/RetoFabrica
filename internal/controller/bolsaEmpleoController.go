package controller

import (
	"fmt"
	"retoBack/internal/model/do/dto"
	"retoBack/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type bolsaController struct{}

var BolsaController = bolsaController{}

// Listar todas las bolsas
func (c *bolsaController) Listar(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Listando todas las bolsas de empleo\n")

	data, err := service.BolsaService.Listar(r.Context())
	if err != nil {
		fmt.Printf("❌ Error listando bolsas: %v\n", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Error al cargar los empleos",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando %d empleos\n", len(data))
	r.Response.WriteJson(data)
}

// Listar bolsas por departamento
func (c *bolsaController) ListarPorDepartamento(r *ghttp.Request) {
	depIDStr := r.Get("id").String()
	fmt.Printf("🔄 Controller - Listando bolsas por departamento: %s\n", depIDStr)

	depID, err := strconv.Atoi(depIDStr)
	if err != nil || depID <= 0 {
		fmt.Printf("❌ ID de departamento inválido: %s\n", depIDStr)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "ID de departamento inválido",
		})
		return
	}

	data, err := service.BolsaService.ListarPorDepartamento(r.Context(), depID)
	if err != nil {
		fmt.Printf("❌ Error listando bolsas por departamento %d: %v\n", depID, err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Error al cargar los empleos del departamento",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller - Enviando %d empleos del departamento %d\n", len(data), depID)
	r.Response.WriteJson(data)
}

// Crear bolsa
func (c *bolsaController) Crear(r *ghttp.Request) {
	fmt.Printf("🔄 Controller - Creando nueva bolsa de empleo\n")
	fmt.Printf("📦 Body recibido: %s\n", r.GetBodyString())

	var b dto.BolsaEmpleoDTO
	err := r.Parse(&b)
	if err != nil {
		fmt.Printf("❌ Error parseando request: %v\n", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Solicitud inválida: " + err.Error(),
		})
		return
	}

	fmt.Printf("📥 DTO recibido del frontend: %+v\n", b)

	// ✅ VALIDACIONES
	if b.Puesto == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El puesto es requerido"})
		return
	}
	if len(b.Puesto) < 3 {
		r.Response.WriteJson(g.Map{"success": false, "error": "El puesto debe tener al menos 3 caracteres"})
		return
	}
	if b.Descripcion == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "La descripción es requerida"})
		return
	}
	if len(b.Descripcion) < 10 {
		r.Response.WriteJson(g.Map{"success": false, "error": "La descripción debe tener al menos 10 caracteres"})
		return
	}
	if b.Salario <= 0 {
		r.Response.WriteJson(g.Map{"success": false, "error": "El salario debe ser mayor a 0"})
		return
	}
	if b.DepartamentoID <= 0 {
		r.Response.WriteJson(g.Map{"success": false, "error": "El departamento es requerido"})
		return
	}

	estadosValidos := map[string]bool{"DISPONIBLE": true, "OCUPADO": true, "CERRADO": true}
	if !estadosValidos[b.Estado] {
		r.Response.WriteJson(g.Map{"success": false, "error": "Estado inválido. Debe ser: DISPONIBLE, OCUPADO o CERRADO"})
		return
	}

	fmt.Printf("✅ Validaciones pasadas - Enviando al servicio...\n")

	if err := service.BolsaService.Crear(r.Context(), &b); err != nil {
		fmt.Printf("❌ Error en servicio: %v\n", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Error al crear el empleo: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Controller: Creación completada exitosamente\n")
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleo creado exitosamente",
	})
}

// Actualizar bolsa
func (c *bolsaController) Actualizar(r *ghttp.Request) {
	idStr := r.Get("id").String()
	fmt.Printf("🔄 Controller - Actualizando bolsa ID: %s\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		r.Response.WriteJson(g.Map{"success": false, "error": "ID inválido"})
		return
	}

	var b dto.BolsaEmpleoDTO
	if err := r.Parse(&b); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "error": "Solicitud inválida: " + err.Error()})
		return
	}

	b.ID = id

	fmt.Printf("📝 Actualizando empleo ID: %d\n", id)

	if b.Puesto == "" {
		r.Response.WriteJson(g.Map{"success": false, "error": "El puesto es requerido"})
		return
	}
	if b.Salario <= 0 {
		r.Response.WriteJson(g.Map{"success": false, "error": "El salario debe ser mayor a 0"})
		return
	}

	if err := service.BolsaService.Actualizar(r.Context(), &b); err != nil {
		fmt.Printf("❌ Error actualizando empleo %d: %v\n", id, err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   "Error al actualizar el empleo: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Empleo %d actualizado exitosamente\n", id)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleo actualizado exitosamente",
	})
}

// Eliminar bolsa
func (c *bolsaController) Eliminar(r *ghttp.Request) {
	idStr := r.Get("id").String()
	fmt.Printf("🗑️ Controller - Eliminando bolsa ID: %s\n", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		r.Response.WriteStatusExit(400, g.Map{"success": false, "error": "ID inválido"})
		return
	}

	fmt.Printf("🔍 Verificando si el empleo %d tiene empleado asignado...\n", id)

	// 🧩 Verificar si tiene un empleado asignado
	empleo, err := service.BolsaService.ObtenerPorID(r.Context(), id)
	if err != nil {
		fmt.Printf("❌ Error buscando empleo: %v\n", err)
		r.Response.WriteStatusExit(404, g.Map{"success": false, "error": "Empleo no encontrado"})
		return
	}

	// 💡 Ahora EmpleadoID es int, por lo que se usa 0 para indicar "sin asignar"
	if empleo.EmpleadoID != 0 {
		fmt.Printf("🚫 No se puede eliminar: tiene empleado asignado (EmpleadoID=%d)\n", empleo.EmpleadoID)
		r.Response.WriteStatusExit(409, g.Map{
			"success": false,
			"error":   "No se puede eliminar el empleo porque tiene un empleado asignado",
		})
		return
	}

	fmt.Printf("🧹 Eliminando empleo ID %d...\n", id)
	if err := service.BolsaService.Eliminar(r.Context(), id); err != nil {
		fmt.Printf("❌ Error eliminando empleo %d: %v\n", id, err)
		r.Response.WriteStatusExit(500, g.Map{
			"success": false,
			"error":   "Error al eliminar el empleo: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Empleo %d eliminado exitosamente\n", id)
	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Empleo eliminado exitosamente",
	})
}
