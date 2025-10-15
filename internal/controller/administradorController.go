package controller

import (
	"retoBack/internal/model/do/dto"
	"retoBack/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func ShowAllAdmins(r *ghttp.Request) {
	admins, err := service.Administradores.GetAll(r.Context())
	if err != nil {
		r.Response.WriteJson(g.Map{"error": err.Error()})
		return
	}
	r.Response.WriteJson(g.Map{"admins": admins})
}

func AdminLogin(r *ghttp.Request) {
	var req dto.AdministradorDTO
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{"error": "Solicitud inválida"})
		return
	}

	adminDTO, err := service.Administradores.Login(r.Context(), req.Correo, req.Password)
	if err != nil {
		r.Response.WriteJson(g.Map{"error": err.Error()})
		return
	}
	if adminDTO == nil {
		r.Response.WriteJson(g.Map{"error": "Credenciales inválidas"})
		return
	}

	r.Session.Set("admin", adminDTO)
	r.Session.Set("correo", adminDTO.Correo)

	r.Response.WriteJson(g.Map{
		"message": "Login exitoso",
		"admin":   adminDTO,
	})
}

func AdminRegister(r *ghttp.Request) {
	var req dto.AdministradorDTO
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatus(400, g.Map{"error": "Solicitud inválida"})
		return
	}

	err := service.Administradores.Create(r.Context(), &dto.AdministradorDTO{
		Nombre:   req.Nombre,
		Correo:   req.Correo,
		Password: req.Password,
	})
	if err != nil {
		r.Response.WriteStatus(400, g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"message": "Registro exitoso"})
}

func AdminUpdate(r *ghttp.Request) {
	var req dto.AdministradorDTO
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatus(400, g.Map{"error": "Solicitud inválida"})
		return
	}

	// Obtener admin de la sesión - FORMA CORRECTA para v2
	adminSession := r.Session.MustGet("admin")
	if adminSession.IsEmpty() {
		r.Response.WriteStatus(401, g.Map{"error": "No autenticado"})
		return
	}

	// Convertir a DTO
	var admin *dto.AdministradorDTO
	if err := adminSession.Struct(&admin); err != nil {
		r.Response.WriteStatus(401, g.Map{"error": "Sesión inválida"})
		return
	}

	req.ID = admin.ID

	if err := service.Administradores.Update(r.Context(), &req); err != nil {
		r.Response.WriteStatus(400, g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{"message": "Actualización exitosa"})
}

func AdminLogout(r *ghttp.Request) {
	r.Session.Remove("admin")
	r.Session.Remove("correo")
	r.Response.WriteJson(g.Map{"message": "Logout exitoso"})
}
