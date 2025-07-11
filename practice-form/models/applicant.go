package models

import "time"

type Applicant struct {
	ID            uint `gorm:"primaryKey"`
	FirstName     string
	LastName      string
	PlaceOfBirth  string
	DateOfBirth   time.Time
	Age           int
	Address       string
	MaritalStatus string
	Gender        string
	Nationality   string
	DNI           string
	HomePhone     *string
	MobilePhone   string
	Email         string
	Photo         string

	Families       []Family             `gorm:"foreignKey:ApplicantID"`
	Children       []Child              `gorm:"foreignKey:ApplicantID"`
	Education      []Education          `gorm:"foreignKey:ApplicantID"`
	Certifications []CertificationSkill `gorm:"foreignKey:ApplicantID"`
	Experience     []Experience         `gorm:"foreignKey:ApplicantID"`
	References     []Reference          `gorm:"foreignKey:ApplicantID"`
	MedicalInfo    MedicalInfo          `gorm:"foreignKey:ApplicantID"`
	Emergencies    []Emergency          `gorm:"foreignKey:ApplicantID"`
	JobApplication JobApplication       `gorm:"foreignKey:ApplicantID"`
}

type Family struct {
	ID           uint `gorm:"primaryKey"`
	ApplicantID  uint
	FullName     string
	Relationship string
	PhoneNumber  string
}

type Child struct {
	ID          uint `gorm:"primaryKey"`
	ApplicantID uint
	FullName    string
	DateOfBirth time.Time
}

type Education struct {
	ID          uint `gorm:"primaryKey"`
	ApplicantID uint
	Level       string
	Institution string
	Location    string
	DegreeTitle string
}

type CertificationSkill struct {
	ID             uint `gorm:"primaryKey"`
	ApplicantID    uint
	Certifications string
	HobbiesSkills  string
}

type Experience struct {
	ID               uint `gorm:"primaryKey"`
	ApplicantID      uint
	CompanyName      string
	StartDate        time.Time
	EndDate          time.Time
	PositionHeld     string
	LastSalary       float64
	ReasonForLeaving string
}

type Reference struct {
	ID          uint `gorm:"primaryKey"`
	ApplicantID uint
	FullName    string
	Address     string
	PhoneNumber string
	Position    string
}

type JobApplication struct {
	ID                uint `gorm:"primaryKey"`
	ApplicantID       uint
	AppliedPosition   string
	StartDate         time.Time
	InformationSource string
	ExpectedSalary    float64
	Signature         string
	Place             string
	Date              time.Time
}

type MedicalInfo struct {
	ID               uint `gorm:"primaryKey"`
	ApplicantID      uint
	IllnessAnswer    string
	IllnessDetail    string
	AllergyAnswer    string
	AllergyDetail    string
	MedicationAnswer string
	MedicationDetail string
	SurgeryAnswer    string
	SurgeryDetail    string
}

type Emergency struct {
	ID           uint `gorm:"primaryKey"`
	ApplicantID  uint
	FullName     string
	Address      string
	PhoneNumber  string
	Relationship string
}
