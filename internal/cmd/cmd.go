package cmd

import (
	"context"
	"fmt"
	"log"
	"retoBack/internal/cmd/routes" // 👈 importa tu paquete de rutas
	"retoBack/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// Registrar todas las rutas aquí
			routes.Register(s)

			// Levantar el servidor
			s.Run()
			return nil
		},
	}
	Migrate = gcmd.Command{
		Name:  "migrate",
		Usage: "migrate",
		Brief: "Crea las tablas en la base de datos según las entidades definidas",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println("🚧 Ejecutando migraciones...")

			// Conexión PostgreSQL (ajusta tu info)
			dsn := "host=localhost user=postgres password=12345678 dbname=RetoBD port=5432 sslmode=disable TimeZone=America/Bogota"
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
			}

			fmt.Println("✅ Conectado a la base de datos PostgreSQL")

			// Migrar todas tus entidades
			err = db.AutoMigrate(
				&entity.Departamento{},
				&entity.BolsaEmpleo{},
				&entity.Empleados{},
				&entity.Nomina{},
			)
			if err != nil {
				log.Fatalf("❌ Error al migrar tablas: %v", err)
			}

			fmt.Println("🎉 Migración completada exitosamente")
			return nil
		},
	}
)
