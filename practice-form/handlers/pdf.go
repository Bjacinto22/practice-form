package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"practice-form/db"
	"practice-form/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateApplicantPDF(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del path /applicant/{id}/pdf
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "applicant" || parts[2] != "pdf" {
		http.Error(w, "Ruta no válida", http.StatusBadRequest)
		return
	}

	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID no válido", http.StatusBadRequest)
		return
	}

	// Buscar en DB
	var applicant models.Applicant
	err = db.DB.Preload("JobApplication").First(&applicant, id).Error
	if err != nil {
		http.Error(w, "Solicitante no encontrado", http.StatusNotFound)
		return
	}

	// Generar PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Resumen del Solicitante")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Nombre: "+applicant.FirstName+" "+applicant.LastName)
	pdf.Ln(8)
	pdf.Cell(0, 10, "DNI: "+applicant.DNI)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Email: "+applicant.Email)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Puesto Solicitado: "+applicant.JobApplication.AppliedPosition)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=applicant.pdf")
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Error al generar PDF", http.StatusInternalServerError)
		return
	}
}
