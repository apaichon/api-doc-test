package consent

import (
    "time"

)

// Domain Models
type Consent struct {
    ConsentID       int       `json:"consent_id"`
    PatientID       string    `json:"patient_id"`
    SourceHospital  string    `json:"source_hospital"`
    TargetHospital  string    `json:"target_hospital"`
    Purpose         string    `json:"purpose"`
    DataCategories  []string  `json:"data_categories"`
    StartDate       time.Time `json:"start_date"`
    ExpiryDate      time.Time `json:"expiry_date"`
    Status          string    `json:"status"` // PENDING, ACTIVE, REVOKED, EXPIRED
	Version         int       `json:"version"`
    Signature       string    `json:"signature"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type ConsentLog struct {
	ConsentLogID int       `json:"consent_log_id"`
    ConsentID   int       `json:"consent_id"`
    Action      string    `json:"action"` // CREATED, UPDATED, REVOKED, ACCESSED
    ActorID     string    `json:"actor_id"`
    ActorType   string    `json:"actor_type"` // PATIENT, DOCTOR, SYSTEM
    Timestamp   time.Time `json:"timestamp"`
    Description string    `json:"description"`
}