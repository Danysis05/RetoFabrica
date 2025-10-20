package routes

import (
	"retoBack/internal/controller"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Register(s *ghttp.Server) {
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.POST("/login", controller.AdminLogin)
		group.POST("/register", controller.AdminRegister)
		group.PUT("/updateAdmin", controller.AdminUpdate)
		group.POST("/logout", controller.AdminLogout)
		group.GET("/all", controller.ShowAllAdmins)
	})

	s.Group("/empleados", func(group *ghttp.RouterGroup) {
		group.POST("/create", controller.EmpleadoCreate)
		group.GET("/list", controller.EmpleadoShowAll)
		group.PUT("/update", controller.EmpleadoUpdate)
		group.DELETE("/delete", controller.EmpleadoDelete)
		group.GET("/:id", controller.EmpleadoGetByID)
	})

	s.Group("/departamentos", func(group *ghttp.RouterGroup) {
		group.POST("/create", controller.DepartamentoCreate)
		group.GET("/list", controller.DepartamentoShowAll)
		group.PUT("/update/:id", controller.DepartamentoUpdate)
		group.DELETE("/delete/:id", controller.DepartamentoDelete)
		group.GET("/can-delete/:id", controller.DepartamentoCanDelete)
		group.GET("/with-empleos", controller.DepartamentoShowAllWithEmpleos)
		group.GET("/:id/with-empleos", controller.DepartamentoShowWithEmpleos)
	})
	s.Group("/bolsa", func(group *ghttp.RouterGroup) {
		group.GET("/", controller.BolsaController.Listar)
		group.GET("/departamento/:id", controller.BolsaController.ListarPorDepartamento)
		group.POST("/", controller.BolsaController.Crear)
		group.PUT("/:id", controller.BolsaController.Actualizar)
		group.DELETE("/:id", controller.BolsaController.Eliminar)
	})
	s.Group("/nomina", func(group *ghttp.RouterGroup) {
		group.GET("/all", controller.NominaShowAll)
		group.POST("/create", controller.NominaCreate)
		group.PUT("/:id", controller.NominaUpdate)
		group.DELETE("/:id", controller.NominaDelete)
		group.GET("/:id", controller.NominaGetByID)
		group.GET("/empleado-info", controller.NominaGetEmpleadoInfo)
	})
}
