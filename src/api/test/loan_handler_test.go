package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"api/internal/loan"
)

func generateTestData() (string, string, float64) {
	timestamp := time.Now().UnixNano()
	loanID := fmt.Sprintf("LOAN-%d", timestamp)
	applicantID := fmt.Sprintf("APP-%d", timestamp)
	amount := 10000.0 + float64(timestamp%1000)
	return loanID, applicantID, amount
}

func makeRequest(method, url, body, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	client := &http.Client{}
	return client.Do(req)
}

func TestLoanAPI_Integration(t *testing.T) {
	baseURL := "http://127.0.0.1:4000"

	t.Run("Complete loan application flow", func(t *testing.T) {
		// Generate unique test data
		loanID, applicantID, amount := generateTestData()

		// Step 1: Login
		loginResp, err := makeRequest(http.MethodPost,
			baseURL+"/api/login",
			`{"username": "admin", "password": "password1234"}`,
			"",
		)
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}
		defer loginResp.Body.Close()

		if loginResp.StatusCode != http.StatusOK {
			t.Fatalf("Login failed: got status %v", loginResp.StatusCode)
		}

		var loginResponse map[string]interface{}
		if err := json.NewDecoder(loginResp.Body).Decode(&loginResponse); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		token := loginResponse["data"].(map[string]interface{})["token"].(string)

		// Step 2: Apply for loan
		applyBody := fmt.Sprintf(`{
			"id": "%s",
			"applicant_id": "%s",
			"amount": %.2f,
			"term": 12,
			"purpose": "Home Improvement",
			"evidence": [
				{
					"id": "EVI-%s",
					"type": "INCOME_STATEMENT",
					"description": "Monthly Income Statement"
				}
			]
		}`, loanID, applicantID, amount, time.Now().Format("20060102150405"))

		applyResp, err := makeRequest(http.MethodPost,
			baseURL+"/loans/apply",
			applyBody,
			token,
		)
		if err != nil {
			t.Fatalf("Failed to make apply request: %v", err)
		}
		defer applyResp.Body.Close()

		if applyResp.StatusCode != http.StatusCreated {
			t.Fatalf("Apply loan failed: got status %v", applyResp.StatusCode)
		}

		var loanApp loan.LoanApplication
		if err := json.NewDecoder(applyResp.Body).Decode(&loanApp); err != nil {
			t.Fatalf("Failed to decode apply response: %v", err)
		}

		// Step 3: Update credit score
		creditScore := 700 + (time.Now().UnixNano() % 100)
		creditBody := fmt.Sprintf(`{
			"loan_id": "%s",
			"credit_score": %d,
			"interest_rate": 0.05
		}`, loanApp.ID, creditScore)

		creditResp, err := makeRequest(http.MethodPost,
			baseURL+"/loans/updateCreditScore?loanID="+loanApp.ID,
			creditBody,
			token,
		)
		if err != nil {
			t.Fatalf("Failed to make credit score update request: %v", err)
		}
		defer creditResp.Body.Close()

		if creditResp.StatusCode != http.StatusOK {
			t.Fatalf("Update credit score failed: got status %v", creditResp.StatusCode)
		}

		// Step 4: Review application
		reviewResp, err := makeRequest(http.MethodGet,
			fmt.Sprintf(baseURL+"/loans/review?loanID=%s", loanApp.ID),
			"",
			token,
		)
		if err != nil {
			t.Fatalf("Failed to make review request: %v", err)
		}
		defer reviewResp.Body.Close()

		if reviewResp.StatusCode != http.StatusOK {
			t.Fatalf("Review application failed: got status %v", reviewResp.StatusCode)
		}

		// Step 5: Approve application
		interestRate := 0.05 + (float64(time.Now().UnixNano()%10) / 100)
		approveBody := fmt.Sprintf(`{
			"loan_id": "%s",
			"interest_rate": %.3f
		}`, loanApp.ID, interestRate)

		approveResp, err := makeRequest(http.MethodPost,
			baseURL+"/loans/approve",
			approveBody,
			token,
		)
		if err != nil {
			t.Fatalf("Failed to make approve request: %v", err)
		}
		defer approveResp.Body.Close()

		if approveResp.StatusCode != http.StatusAccepted {
			t.Fatalf("Approve loan failed: got status %v", approveResp.StatusCode)
		}

		var approvedLoan loan.LoanApplication
		if err := json.NewDecoder(approveResp.Body).Decode(&approvedLoan); err != nil {
			t.Fatalf("Failed to decode approve response: %v", err)
		}

		if approvedLoan.Status != loan.StatusApproved {
			t.Errorf("Expected loan status %v, got %v", loan.StatusApproved, approvedLoan.Status)
		}
	})
}
