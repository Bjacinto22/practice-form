package handlers

import (
	"net/http"
	"practice-form/db"
	"practice-form/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func saveAssociations(tx *gorm.DB, applicant *models.Applicant) {
	for i := range applicant.Families {
		applicant.Families[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Families[i])
	}
	for i := range applicant.Children {
		applicant.Children[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Children[i])
	}
	for i := range applicant.Education {
		applicant.Education[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Education[i])
	}
	for i := range applicant.Certifications {
		applicant.Certifications[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Certifications[i])
	}
	for i := range applicant.Experience {
		applicant.Experience[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Experience[i])
	}
	for i := range applicant.References {
		applicant.References[i].ApplicantID = applicant.ID
		tx.Create(&applicant.References[i])
	}
	if applicant.JobApplication.ApplicantID == 0 {
		applicant.JobApplication.ApplicantID = applicant.ID
		tx.Create(&applicant.JobApplication)
	}
	for i := range applicant.Emergencies {
		applicant.Emergencies[i].ApplicantID = applicant.ID
		tx.Create(&applicant.Emergencies[i])
	}
}

func CreateApplicant(c *gin.Context) {
	var input models.Applicant

	// 1. Validar la estructura del JSON entrante
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en el formato del JSON: " + err.Error()})
		return
	}

	// 2. Guardar el Applicant principal
	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el solicitante: " + err.Error()})
		return
	}

	// 3. Guardar relaciones asociadas (si vienen anidadas en el JSON)
	saveAssociations(db.DB, &input)

	// 4. Confirmación
	c.JSON(http.StatusOK, gin.H{"message": "Formulario guardado exitosamente", "applicant_id": input.ID})
}
