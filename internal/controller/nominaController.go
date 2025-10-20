package controller

import (
	"fmt"
	"retoBack/internal/model/entity"
	"retoBack/internal/service"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// NominaShowAll obtiene todas las nóminas
func NominaShowAll(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaShowAll - Iniciando...")

	nominas, err := service.Nominas.GetAll(r.Context())
	if err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaShowAll - Error: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"message": "Error interno del servidor al obtener nóminas",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("✅ [CONTROLLER] NominaShowAll - Éxito: %d nóminas encontradas\n", len(nominas))

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    nominas,
		"count":   len(nominas),
	})
}

// NominaCreate crea una nueva nómina
func NominaCreate(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaCreate - Iniciando...")

	var req struct {
		EmpleadoID    int     `json:"empleado_id"`
		HorasExtras   float64 `json:"horas_extras"`
		DiasFaltantes int     `json:"dias_faltantes"`
	}

	if err := r.Parse(&req); err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaCreate - Error parseando request: %v\n", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("📝 [CONTROLLER] NominaCreate - Datos recibidos: EmpleadoID=%d, HorasExtras=%.2f, DiasFaltantes=%d\n",
		req.EmpleadoID, req.HorasExtras, req.DiasFaltantes)

	// Crear entidad Nomina básica
	nominaEntity := &entity.Nomina{
		EmpleadoId: req.EmpleadoID,
	}

	// Llamar al service para crear nómina
	dtoNomina, err := service.Nominas.Create(r.Context(), nominaEntity, req.HorasExtras, req.DiasFaltantes)
	if err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaCreate - Error en service: %v\n", err)

		// ✅ MEJORAR manejo de errores
		errorMsg := err.Error()
		statusCode := 500
		userMessage := "Error al crear nómina"

		// ✅ DETECTAR ERRORES ESPECÍFICOS
		if strings.Contains(errorMsg, "no tiene asignado un puesto con salario") {
			statusCode = 400
			userMessage = "El empleado seleccionado no tiene un puesto asignado con salario"
		} else if strings.Contains(errorMsg, "empleado no encontrado") {
			statusCode = 400
			userMessage = "El empleado seleccionado no existe"
		} else if strings.Contains(errorMsg, "error al obtener empleado") {
			statusCode = 400
			userMessage = "Error al obtener información del empleado"
		} else if strings.Contains(errorMsg, "Solicitud inválida") {
			statusCode = 400
			userMessage = "Datos de solicitud inválidos"
		}

		r.Response.WriteStatus(statusCode, g.Map{
			"success": false,
			"error":   userMessage,
			"details": errorMsg,
		})
		return
	}

	fmt.Printf("✅ [CONTROLLER] NominaCreate - Éxito: Nómina ID=%d creada\n", dtoNomina.ID)

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Nómina creada con éxito",
		"data":    dtoNomina,
	})
}

// NominaUpdate actualiza una nómina existente
func NominaUpdate(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaUpdate - Iniciando...")

	// Obtener ID de la URL
	id := r.GetRouter("id").Int()
	if id == 0 {
		fmt.Printf("❌ [CONTROLLER] NominaUpdate - ID no proporcionado en URL\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID de nómina requerido en la URL",
		})
		return
	}

	fmt.Printf("📝 [CONTROLLER] NominaUpdate - ID de nómina: %d\n", id)

	var req struct {
		EmpleadoID    int     `json:"empleado_id"`
		HorasExtras   float64 `json:"horas_extras"`
		DiasFaltantes int     `json:"dias_faltantes"`
	}

	if err := r.Parse(&req); err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaUpdate - Error parseando request: %v\n", err)
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "Solicitud inválida",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("📝 [CONTROLLER] NominaUpdate - Datos: EmpleadoID=%d, HorasExtras=%.2f, DiasFaltantes=%d\n",
		req.EmpleadoID, req.HorasExtras, req.DiasFaltantes)

	// Crear entidad con el ID de la URL
	nominaEntity := &entity.Nomina{
		Id:         id,
		EmpleadoId: req.EmpleadoID,
	}

	if err := service.Nominas.Update(r.Context(), nominaEntity, req.HorasExtras, req.DiasFaltantes); err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaUpdate - Error en service: %v\n", err)

		// ✅ MEJORAR manejo de errores
		errorMsg := err.Error()
		statusCode := 500
		userMessage := "Error al actualizar nómina"

		if strings.Contains(errorMsg, "no tiene asignado un puesto con salario") {
			statusCode = 400
			userMessage = "El empleado seleccionado no tiene un puesto asignado con salario"
		} else if strings.Contains(errorMsg, "empleado no encontrado") {
			statusCode = 400
			userMessage = "El empleado seleccionado no existe"
		} else if strings.Contains(errorMsg, "nómina no encontrada") {
			statusCode = 404
			userMessage = "La nómina que intenta actualizar no existe"
		}

		r.Response.WriteStatus(statusCode, g.Map{
			"success": false,
			"error":   userMessage,
			"details": errorMsg,
		})
		return
	}

	fmt.Printf("✅ [CONTROLLER] NominaUpdate - Éxito: Nómina ID=%d actualizada\n", id)

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Nómina actualizada con éxito",
	})
}

// NominaDelete elimina una nómina
func NominaDelete(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaDelete - Iniciando...")

	// Obtener ID de la URL
	id := r.GetRouter("id").Int()
	if id == 0 {
		fmt.Printf("❌ [CONTROLLER] NominaDelete - ID no proporcionado en URL\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID de nómina requerido en la URL",
		})
		return
	}

	fmt.Printf("📝 [CONTROLLER] NominaDelete - ID de nómina a eliminar: %d\n", id)

	if err := service.Nominas.Delete(r.Context(), id); err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaDelete - Error en service: %v\n", err)

		// ✅ MEJORAR manejo de errores
		errorMsg := err.Error()
		statusCode := 500
		userMessage := "Error al eliminar nómina"

		if strings.Contains(errorMsg, "nómina no encontrada") {
			statusCode = 404
			userMessage = "La nómina que intenta eliminar no existe"
		}

		r.Response.WriteStatus(statusCode, g.Map{
			"success": false,
			"error":   userMessage,
			"details": errorMsg,
		})
		return
	}

	fmt.Printf("✅ [CONTROLLER] NominaDelete - Éxito: Nómina ID=%d eliminada\n", id)

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Nómina eliminada con éxito",
	})
}

// NominaGetByID obtiene una nómina por ID
func NominaGetByID(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaGetByID - Iniciando...")

	// Obtener ID de la URL
	id := r.GetRouter("id").Int()
	if id == 0 {
		fmt.Printf("❌ [CONTROLLER] NominaGetByID - ID no proporcionado en URL\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID de nómina requerido en la URL",
		})
		return
	}

	fmt.Printf("📝 [CONTROLLER] NominaGetByID - Buscando nómina ID: %d\n", id)

	nomina, err := service.Nominas.GetById(r.Context(), id)
	if err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaGetByID - Error en service: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al obtener nómina",
			"details": err.Error(),
		})
		return
	}

	if nomina == nil {
		fmt.Printf("❌ [CONTROLLER] NominaGetByID - Nómina no encontrada: %d\n", id)
		r.Response.WriteStatus(404, g.Map{
			"success": false,
			"error":   "Nómina no encontrada",
		})
		return
	}

	fmt.Printf("✅ [CONTROLLER] NominaGetByID - Éxito: Nómina ID=%d encontrada\n", id)

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    nomina,
	})
}

// Obtener información del empleado con salario para nómina
func NominaGetEmpleadoInfo(r *ghttp.Request) {
	fmt.Println("🎯 [CONTROLLER] NominaGetEmpleadoInfo - Iniciando...")

	empleadoID := r.Get("empleadoId").Int()
	if empleadoID == 0 {
		fmt.Printf("❌ [CONTROLLER] NominaGetEmpleadoInfo - ID de empleado no proporcionado\n")
		r.Response.WriteStatus(400, g.Map{
			"success": false,
			"error":   "ID de empleado requerido",
		})
		return
	}

	fmt.Printf("🔍 [CONTROLLER] NominaGetEmpleadoInfo - Buscando empleado: ID=%d\n", empleadoID)

	// Obtener empleado con bolsa de empleo
	empleado, err := service.Empleados.GetById(r.Context(), empleadoID)
	if err != nil {
		fmt.Printf("❌ [CONTROLLER] NominaGetEmpleadoInfo - Error: %v\n", err)
		r.Response.WriteStatus(500, g.Map{
			"success": false,
			"error":   "Error al obtener empleado",
		})
		return
	}

	if empleado == nil {
		fmt.Printf("❌ [CONTROLLER] NominaGetEmpleadoInfo - Empleado no encontrado: ID=%d\n", empleadoID)
		r.Response.WriteStatus(404, g.Map{
			"success": false,
			"error":   "Empleado no encontrado",
		})
		return
	}

	// Verificar si tiene bolsa de empleo con salario
	var salarioBase float64 = 0
	var puesto string = ""
	var tieneBolsaActiva bool = false

	if empleado.BolsaEmpleoID > 0 {
		bolsa, err := service.BolsaService.ObtenerPorID(r.Context(), empleado.BolsaEmpleoID)
		if err == nil && bolsa != nil && bolsa.Estado == "OCUPADO" {
			salarioBase = bolsa.Salario
			puesto = bolsa.Puesto
			tieneBolsaActiva = true
			fmt.Printf("💰 [CONTROLLER] NominaGetEmpleadoInfo - Salario encontrado: %.2f, Puesto: %s\n", salarioBase, puesto)
		} else {
			fmt.Printf("⚠️ [CONTROLLER] NominaGetEmpleadoInfo - Empleado tiene BolsaEmpleoID=%d pero no se pudo cargar\n", empleado.BolsaEmpleoID)
		}
	} else {
		fmt.Printf("⚠️ [CONTROLLER] NominaGetEmpleadoInfo - Empleado no tiene BolsaEmpleoID asignado\n")
	}

	response := g.Map{
		"success": true,
		"data": g.Map{
			"id":               empleado.ID,
			"nombre":           empleado.Nombre,
			"apellido":         empleado.Apellido,
			"documentoNumero":  empleado.DocumentoNumero,
			"puesto":           puesto,
			"salarioBase":      salarioBase,
			"tieneBolsaActiva": tieneBolsaActiva,
			"bolsaEmpleoId":    empleado.BolsaEmpleoID,
		},
	}

	fmt.Printf("✅ [CONTROLLER] NominaGetEmpleadoInfo - Éxito: %s %s, Salario: %.2f\n",
		empleado.Nombre, empleado.Apellido, salarioBase)

	r.Response.WriteJson(response)
}
