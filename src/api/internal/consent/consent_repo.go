package consent

import (
	"api/internal/db"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

// ConsentRepo represents the repository for consent operations
type ConsentRepo struct {
	DB *db.DB
}

// NewConsentRepo creates a new instance of ConsentRepo
func NewConsentRepo() *ConsentRepo {
	db := db.NewDB()
	return &ConsentRepo{DB: db}
}

// GetConsents fetches consents with pagination support
func (cr *ConsentRepo) GetConsents(params db.PaginationParams) (*db.PaginationResponse, error) {
	baseQuery := "SELECT * FROM consents"
	countQuery := "SELECT COUNT(*) FROM consents"

	// Get total count
	var total int64
	row, err := cr.DB.QueryRow(countQuery)
	if err != nil {
		return nil, fmt.Errorf("error counting consents: %w", err)
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error scanning total count: %w", err)
	}

	// Build pagination query
	paginationQuery, args, err := db.BuildPaginationQuery(params)
	if err != nil {
		return nil, fmt.Errorf("error building pagination query: %w", err)
	}

	// Execute paginated query
	rows, err := cr.DB.Query(baseQuery+paginationQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying consents: %w", err)
	}
	defer rows.Close()

	var consents []*Consent
	var lastID int

	for rows.Next() {
		var consent Consent
		err := rows.Scan(
			&consent.ConsentID,
			&consent.PatientID,
			&consent.SourceHospital,
			&consent.TargetHospital,
			&consent.Purpose,
			&consent.DataCategories,
			&consent.StartDate,
			&consent.ExpiryDate,
			&consent.Status,
			&consent.Version,
			&consent.Signature,
			&consent.CreatedAt,
			&consent.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning consent: %w", err)
		}
		consents = append(consents, &consent)
		lastID = consent.ConsentID
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating consents: %w", err)
	}

	// Calculate pagination metadata
	hasMore := false
	totalPages := 0
	if params.Limit > 0 {
		hasMore = len(consents) == params.Limit
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}

	response := &db.PaginationResponse{
		Data:       consents,
		Total:      total,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}

	// Add cursor/page specific fields
	switch params.Type {
	case db.CursorPagination:
		if hasMore {
			response.NextCursor = strconv.Itoa(lastID)
		}
	case db.PagePagination:
		response.Page = params.Page
	}

	return response, nil
}

// GetConsentByID retrieves a consent by its ID from the database
func (cr *ConsentRepo) GetConsentByID(id string) (*Consent, error) {
	var consent Consent
	row, err := cr.DB.QueryRow("SELECT * FROM consents WHERE consent_id = ?", id)

	if err != nil {
		return nil, err
	}

	row.Scan(
		&consent.ConsentID,
		&consent.PatientID,
		&consent.SourceHospital,
		&consent.TargetHospital,
		&consent.Purpose,
		&consent.DataCategories,
		&consent.StartDate,
		&consent.ExpiryDate,
		&consent.Status,
		&consent.Version,
		&consent.Signature,
		&consent.CreatedAt,
		&consent.UpdatedAt,
	)

	return &consent, nil
}

// InsertConsent inserts a new consent into the database
func (cr *ConsentRepo) InsertConsent(consent *Consent) (int, error) {
	_, err := cr.DB.Insert("INSERT INTO consents (patient_id, source_hospital, target_hospital, purpose, data_categories, start_date, expiry_date, status, version, signature, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		consent.PatientID, consent.SourceHospital, consent.TargetHospital, consent.Purpose, consent.DataCategories, consent.StartDate, consent.ExpiryDate, consent.Status, consent.Version, consent.Signature, consent.CreatedAt, consent.UpdatedAt)
	if err != nil {
		fmt.Printf("Error inserting consent: %v\n", err)
		return 0, err
	}
	return consent.ConsentID, nil
}

// UpdateConsent updates an existing consent in the database
func (cr *ConsentRepo) UpdateConsent(consent *Consent) (int, error) {
	_, err := cr.DB.Update("UPDATE consents SET patient_id=?, source_hospital=?, target_hospital=?, purpose=?, data_categories=?, start_date=?, expiry_date=?, status=?, version=?, signature=?, updated_at=? WHERE id=?",
		consent.PatientID, consent.SourceHospital, consent.TargetHospital, consent.Purpose, consent.DataCategories, consent.StartDate, consent.ExpiryDate, consent.Status, consent.Version, consent.Signature, consent.UpdatedAt, consent.ID)
	if err != nil {
		return 0, err
	}
	return consent.ConsentID, nil
}

// DeleteConsent deletes a consent from the database
func (cr *ConsentRepo) DeleteConsent(id int) (int, error) {
	_, err := cr.DB.Delete("DELETE FROM consents WHERE id=?", id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
