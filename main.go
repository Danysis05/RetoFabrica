package main

import (
	"fmt"
	"log"
	"retoBack/internal/cmd/routes"
	"retoBack/internal/model/entity"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Middleware CORS
func CORSMiddleware(r *ghttp.Request) {
	r.Response.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization, X-Requested-With")
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Response.Header().Set("Access-Control-Max-Age", "86400")

	if r.Method == "OPTIONS" {
		r.Response.WriteStatus(200)
		return
	}

	r.Middleware.Next()
}

func main() {
	fmt.Println("🚀 Iniciando servidor...")

	// =========================================================
	// 🔧 1. Conexión a la base de datos (PostgreSQL)
	// =========================================================
	dsn := "host=localhost user=postgres password=12345678 dbname=retoBD port=5432 sslmode=disable TimeZone=America/Bogota"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	fmt.Println("✅ Conectado a la base de datos PostgreSQL")

	// =========================================================
	// 🧱 2. Ejecutar migraciones automáticamente
	// =========================================================
	err = db.AutoMigrate(
		&entity.Departamento{},
		&entity.BolsaEmpleo{},
		&entity.Empleados{},
		&entity.Nomina{},
	)
	if err != nil {
		log.Fatalf("❌ Error al migrar tablas: %v", err)
	}
	fmt.Println("🎉 Migraciones ejecutadas correctamente")

	// =========================================================
	// 🌐 3. Configurar y levantar el servidor GoFrame
	// =========================================================
	s := g.Server()
	s.SetPort(8000)
	s.Use(CORSMiddleware)
	routes.Register(s)

	fmt.Println("🚀 Servidor corriendo en http://localhost:8000")
	fmt.Println("🌐 CORS habilitado para: http://localhost:4200")
	fmt.Println("✅ GoFrame configurado correctamente")
	fmt.Println("✅ DAOs y Services listos para usar")

	s.Run()
}
