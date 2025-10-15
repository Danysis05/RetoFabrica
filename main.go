// 📁 main.go - VERSIÓN SOLO GOFRAIME
package main

import (
	"fmt"
	"retoBack/internal/cmd/routes"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Middleware CORS
func CORSMiddleware(r *ghttp.Request) {
	// Configurar headers CORS
	r.Response.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization, X-Requested-With")
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Response.Header().Set("Access-Control-Max-Age", "86400")

	// Manejar preflight OPTIONS request
	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(200)
		return
	}

	r.Middleware.Next()
}

func main() {
	// ----- Servidor GoFrame -----
	s := g.Server()
	s.SetPort(8000)

	// ----- AGREGAR MIDDLEWARE CORS -----
	s.Use(CORSMiddleware)

	// ----- Registrar todas las rutas -----
	routes.Register(s)

	// Mensaje de confirmación
	fmt.Println("🚀 Servidor corriendo en http://localhost:8000")
	fmt.Println("🌐 CORS habilitado para: http://localhost:4200")
	fmt.Println("✅ GoFrame configurado correctamente")
	fmt.Println("✅ DAOs y Services listos para usar")

	// ----- Levantar servidor -----
	s.Run()
}
