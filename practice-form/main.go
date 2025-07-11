package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"practice-form/db"
	"practice-form/handlers"'
	"practice-form/models"
	"strconv"
	"strings"
	"time"
)

func handleFormSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Validación de campos obligatorios
	requiredFields := []string{
		"first_name", "last_name", "place_of_birth", "date_of_birth", "age",
		"address", "marital_status", "gender", "nacionality", "dni",
		"mobile_phone", "email", "applied_position", "job_start_date",
		"expected_salary", "signature", "signature_place", "signature_datetime",
	}
	for _, field := range requiredFields {
		if r.FormValue(field) == "" {
			http.Error(w, fmt.Sprintf("El campo '%s' es obligatorio", field), http.StatusBadRequest)
			return
		}
	}

	// Datos personales
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	placeOfBirth := r.FormValue("place_of_birth")
	dateOfBirth, err := time.Parse("2006-01-02", r.FormValue("date_of_birth"))
	if err != nil {
		http.Error(w, "Fecha de nacimiento inválida", http.StatusBadRequest)
		return
	}
	age, _ := strconv.Atoi(r.FormValue("age"))
	address := r.FormValue("address")
	maritalStatus := r.FormValue("marital_status")
	gender := r.FormValue("gender")
	nationality := r.FormValue("nacionality")
	dni := r.FormValue("dni")
	homePhone := r.FormValue("home_phone") // opcional
	mobilePhone := r.FormValue("mobile_phone")
	email := r.FormValue("email")

	// Foto (obligatoria)
	var filename string
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "La fotografía es obligatoria", http.StatusBadRequest)
		return
	}
	defer file.Close()
	os.MkdirAll("uploads", os.ModePerm)
	filename = fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), handler.Filename)
	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Error al guardar la imagen", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Error al copiar la imagen", http.StatusInternalServerError)
		return
	}

	// Crear y guardar solicitante
	applicant := models.Applicant{
		FirstName:     firstName,
		LastName:      lastName,
		PlaceOfBirth:  placeOfBirth,
		DateOfBirth:   dateOfBirth,
		Age:           age,
		Address:       address,
		MaritalStatus: maritalStatus,
		Gender:        gender,
		Nationality:   nationality,
		DNI:           dni,
		MobilePhone:   mobilePhone,
		Email:         email,
		Photo:         filename,
	}
	if homePhone != "" {
		applicant.HomePhone = &homePhone
	}
	if err := db.DB.Create(&applicant).Error; err != nil {
		http.Error(w, "Error al guardar el solicitante", http.StatusInternalServerError)
		return
	}

	// Familiares (opcional)
	familiares := []models.Family{
		{ApplicantID: applicant.ID, FullName: r.FormValue("spouse_name"), Relationship: "Cónyuge", PhoneNumber: r.FormValue("spouse_phone")},
		{ApplicantID: applicant.ID, FullName: r.FormValue("father_name"), Relationship: "Padre", PhoneNumber: r.FormValue("father_phone")},
		{ApplicantID: applicant.ID, FullName: r.FormValue("mother_name"), Relationship: "Madre", PhoneNumber: r.FormValue("mother_phone")},
	}
	for _, fam := range familiares {
		if fam.FullName != "" {
			db.DB.Create(&fam)
		}
	}

	// Hijos (opcional)
	childrenNames := r.Form["children_name[]"]
	childrenBirthdates := r.Form["children_birthdate[]"]
	for i := range childrenNames {
		if childrenNames[i] != "" && childrenBirthdates[i] != "" {
			if dob, err := time.Parse("2006-01-02", childrenBirthdates[i]); err == nil {
				db.DB.Create(&models.Child{ApplicantID: applicant.ID, FullName: childrenNames[i], DateOfBirth: dob})
			}
		}
	}

	// Formación académica
	for i := 1; i <= 3; i++ {
		level := map[int]string{1: "Secundaria", 2: "Educación Superior", 3: "Otros"}[i]
		inst := r.FormValue("education_institution_" + strconv.Itoa(i))
		loc := r.FormValue("education_place_" + strconv.Itoa(i))
		title := r.FormValue("education_title_" + strconv.Itoa(i))
		if inst != "" || loc != "" || title != "" {
			db.DB.Create(&models.Education{
				ApplicantID: applicant.ID,
				Level:       level,
				Institution: inst,
				Location:    loc,
				DegreeTitle: title,
			})
		}
	}

	// Certificaciones y habilidades
	certs := r.Form["certifications[]"]
	skills := r.Form["hobbies_skills[]"]
	for i := 0; i < len(certs) || i < len(skills); i++ {
		var cert, skill string
		if i < len(certs) {
			cert = certs[i]
		}
		if i < len(skills) {
			skill = skills[i]
		}
		if cert != "" || skill != "" {
			db.DB.Create(&models.CertificationSkill{
				ApplicantID:    applicant.ID,
				Certifications: cert,
				HobbiesSkills:  skill,
			})
		}
	}

	// Experiencia laboral (opcional)
	empresas := r.Form["empresa[]"]
	ingresos := r.Form["ingreso[]"]
	salidas := r.Form["salida[]"]
	puestos := r.Form["puesto[]"]
	salarios := r.Form["salario[]"]
	motivos := r.Form["motivo_salida[]"]
	for i := range empresas {
		if empresas[i] == "" && puestos[i] == "" {
			continue
		}
		start, _ := time.Parse("2006-01-02", ingresos[i])
		end, _ := time.Parse("2006-01-02", salidas[i])
		sal, _ := strconv.ParseFloat(salarios[i], 64)
		db.DB.Create(&models.Experience{
			ApplicantID:      applicant.ID,
			CompanyName:      empresas[i],
			StartDate:        start,
			EndDate:          end,
			PositionHeld:     puestos[i],
			LastSalary:       sal,
			ReasonForLeaving: motivos[i],
		})
	}

	// Referencias personales (obligatorio al menos una)
	if len(r.Form["ref_full_name"]) == 0 || r.Form["ref_full_name"][0] == "" {
		http.Error(w, "Debe proporcionar al menos una referencia", http.StatusBadRequest)
		return
	}
	refNames := r.Form["ref_full_name"]
	refAddrs := r.Form["ref_address"]
	refPhones := r.Form["ref_phone"]
	refPositions := r.Form["ref_position"]
	for i := range refNames {
		if refNames[i] != "" {
			db.DB.Create(&models.Reference{
				ApplicantID: applicant.ID,
				FullName:    refNames[i],
				Address:     refAddrs[i],
				PhoneNumber: refPhones[i],
				Position:    refPositions[i],
			})
		}
	}

	// Información médica (obligatoria)
	medical := models.MedicalInfo{
		ApplicantID:      applicant.ID,
		IllnessAnswer:    r.FormValue("illness_answer"),
		IllnessDetail:    r.FormValue("illness_detail"),
		AllergyAnswer:    r.FormValue("allergy_answer"),
		AllergyDetail:    r.FormValue("allergy_detail"),
		MedicationAnswer: r.FormValue("medication_answer"),
		MedicationDetail: r.FormValue("medication_detail"),
		SurgeryAnswer:    r.FormValue("surgery_answer"),
		SurgeryDetail:    r.FormValue("surgery_detail"),
	}
	db.DB.Create(&medical)

	// Contactos de emergencia (obligatorio al menos uno)
	emergencyNames := r.Form["emergency_name[]"]
	emergencyAddrs := r.Form["emergency_address[]"]
	emergencyPhones := r.Form["emergency_phone[]"]
	emergencyRelations := r.Form["emergency_relation[]"]
	if len(emergencyNames) == 0 || emergencyNames[0] == "" {
		http.Error(w, "Debe agregar al menos un contacto de emergencia", http.StatusBadRequest)
		return
	}
	for i := range emergencyNames {
		if emergencyNames[i] != "" {
			db.DB.Create(&models.Emergency{
				ApplicantID:  applicant.ID,
				FullName:     emergencyNames[i],
				Address:      emergencyAddrs[i],
				PhoneNumber:  emergencyPhones[i],
				Relationship: emergencyRelations[i],
			})
		}
	}

	// Solicitud de empleo
	expectedSalary, _ := strconv.ParseFloat(r.FormValue("expected_salary"), 64)
	startDate, _ := time.Parse("2006-01-02", r.FormValue("job_start_date"))
	signatureDate, _ := time.Parse("2006-01-02 15:04:05", r.FormValue("signature_datetime"))
	db.DB.Create(&models.JobApplication{
		ApplicantID:       applicant.ID,
		AppliedPosition:   r.FormValue("applied_position"),
		StartDate:         startDate,
		InformationSource: r.FormValue("information_source"),
		ExpectedSalary:    expectedSalary,
		Signature:         r.FormValue("signature"),
		Place:             r.FormValue("signature_place"),
		Date:              signatureDate,
	})

	http.Redirect(w, r, fmt.Sprintf("/applicant/%d/preview", applicant.ID), http.StatusSeeOther)
}

func main() {
	db.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/form.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/applicant/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/preview") {
			handlers.HandlePreview(w, r)
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/submit", handleFormSubmit)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	http.HandleFunc("/applicant/pdf", handlers.GenerateApplicantPDF)

	fmt.Println("Server http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
