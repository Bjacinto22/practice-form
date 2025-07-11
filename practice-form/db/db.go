package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"practice-form/models"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Conexión a la base de datos SQLite (usa un archivo .db local)
	DB, err = gorm.Open(sqlite.Open("applicants.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error al conectar con la base de datos:", err)
	}

	// Migración automática de modelos
	err = DB.AutoMigrate(
		&models.Applicant{},
		&models.Family{},
		&models.Child{},
		&models.Education{},
		&models.CertificationSkill{},
		&models.Experience{},
		&models.Reference{},
		&models.MedicalInfo{},
		&models.Emergency{},
		&models.JobApplication{},
	)
	if err != nil {
		log.Fatal("❌ Error al migrar los modelos:", err)
	}

	log.Println("✅ Base de datos conectada y migrada correctamente")
}
