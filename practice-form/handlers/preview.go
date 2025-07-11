package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"practice-form/db"
	"practice-form/models"
	"strconv"
	"strings"
)

func HandlePreview(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID desde la URL: /applicant/:id/preview
	path := strings.TrimPrefix(r.URL.Path, "/applicant/")
	idStr := strings.TrimSuffix(path, "/preview")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var applicant models.Applicant
	err = db.DB.Preload("Families").
		Preload("Children").
		Preload("Education").
		Preload("Certifications").
		Preload("Experience").
		Preload("References").
		Preload("MedicalInfo").
		Preload("Emergencies").
		Preload("JobApplication").
		First(&applicant, id).Error

	if err != nil {
		http.Error(w, "Solicitante no encontrado", http.StatusNotFound)
		return
	}

	// Renderizar plantilla HTML
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "preview.html")))
	tmpl.Execute(w, applicant)
}
