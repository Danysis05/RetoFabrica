package boot

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ✅ Inicializa la conexión global con GORM (PostgreSQL)
func InitGorm() {
	// Obtiene variables de entorno (puedes reemplazar con valores fijos si lo prefieres)
	user := os.Getenv("postgres")
	pass := os.Getenv("12345678")
	host := os.Getenv("localhost")
	port := os.Getenv("5432")
	name := os.Getenv("retoBD")
	sslmode := os.Getenv("disable") // puedes usar "disable" si estás en local

	// 🧩 DSN (Data Source Name) para PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Bogota",
		host, user, pass, name, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos con GORM (PostgreSQL): %v", err)
	}

	DB = db
	log.Println("✅ Conexión GORM (PostgreSQL) inicializada correctamente")
}
