---
theme: default
background: https://source.unsplash.com/collection/94734566/1920x1080
class: 'text-center'
highlighter: shiki
lineNumbers: false
drawings:
  persist: false
transition: slide-left
title: API Documentation and Testing with Go
css: /styles/style.css
---

# API Documentation and Testing with Go
Building Better APIs with Go

---
layout: two-cols
---

# Table of Contents

## Part 1 : Write Document Practical
  - 5W 1H
  - Write Comment Practical

## Part 2: Documentation Tools
 - Go Doc
 - Swagger
 - Api Doc
 - Docusarus

::right::

## Part 3: API Testing
 - Pyramid Test Conceptual
 - Unit Test
 - Integration Test

## Part 4: Testing Tools
 - Playwright 
 - Write documentation and test with Agent AI

---
layout: default
---

# Part 1: Core Information (5W1H)

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### What
- API Name & Version
- Core functionality
- Key features

### Who
- Target users
- Required expertise
- Stakeholders

### Where
- Hosting environments
- Source code location
- Documentation links

</v-clicks>

</div>
<div>

<v-clicks>

### When
- Release dates
- Update frequency
- Availability

### Why
- Problem solved
- Key benefits
- Use cases

### How
- Authentication
- Request formats
- Response structures

</v-clicks>

</div>
</div>

---
layout: default
---

# Best Practices for Writing Code Comments
1. Be Consistent with Comment Style
2. Document at the Right Level
3. Follow Documentation Tool Conventions
4. Include Essential Elements
5. Document Types Thoroughly
6. Use Clear and Concise Language
7. Keep Comments Updated
8. Use Code Examples
9. Include Best Practices
10. Use Markdown Formatting

<a href="./write-comment-practical.html"> Best Practices for Writing Code Comments</a>

---
layout: default
---

# API Documentation Structure

```mermaid
graph TD
    A[API Documentation] --> B[Core Information]
    A --> C[Technical Details]
    A --> D[Usage Guide]
    A --> E[Support]

    B --> B1[API Overview]
    B --> B2[Authentication]
    B --> B3[Base URLs]
    B --> B4[Versioning]

    C --> C1[Endpoints]
    C --> C2[Request/Response]
    C --> C3[Error Handling]
    C --> C4[Rate Limits]

    C1 --> C1a[Methods]
    C1 --> C1b[Parameters]
    C1 --> C1c[Headers]

    C2 --> C2a[Data Types]
    C2 --> C2b[Examples]
    C2 --> C2c[Schemas]

    D --> D1[Quick Start]
    D --> D2[Code Examples]
    D --> D3[Best Practices]
    D --> D4[Use Cases]

    E --> E1[Contact Info]
    E --> E2[Issue Reporting]
    E --> E3[Change Log]
    E --> E4[FAQ]
```

<a href="./api-docs-template.html"> API Document Template</a>

---
layout: two-cols
---

# Authentication

<v-clicks>

### Security First
- API key management
- Bearer token format
- OAuth 2.0 compliance
- Best practices

</v-clicks>

::right::

<v-click>

```bash
# Authentication Header
Authorization: Bearer <your_api_key>

# Example Request
curl -X GET "https://api.example.com/v1/resource" \
     -H "Authorization: Bearer your_api_key"
```

</v-click>

---
layout: default
---

# Request & Response Formats

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Request Format
```json
{
  "required_field": "value",
  "optional_field": "value"
}
```

</v-clicks>

</div>
<div>

<v-clicks>

### Response Format
```json
{
  "status": "success",
  "data": {
    "field1": "value1",
    "field2": "value2"
  },
  "metadata": {
    "timestamp": "2025-02-18T12:00:00Z"
  }
}
```

</v-clicks>

</div>
</div>

---
layout: default
---

# Error Handling

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Error Response Structure
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {}
  }
}
```

</v-clicks>

</div>
<div>

<v-clicks>

### Common Error Codes
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `429`: Too Many Requests
- `500`: Internal Server Error

</v-clicks>

</div>
</div>

---
layout: two-cols
---

# Best Practices

<v-clicks>

### Security
- HTTPS enforcement
- Rate limiting
- Input validation
- Regular audits

### Performance
- Caching strategies
- Pagination
- Compression
- Async operations

</v-clicks>

::right::

<v-clicks>

### Versioning
- Semantic versioning
- URL path versions
- Deprecation policy
- Migration guides

### Documentation
- Keep it updated
- Include examples
- Version control
- Changelog

</v-clicks>

---
layout: default
---

# Loan System API Documentation
Best Practices & Examples

<img src="/images/api-docs/best-practice.svg" style="width: 600px;">

---
layout: two-cols
---

# API Overview

<v-clicks>

### Base Information
- **API Version**: v1.0.0
- **Base URL**: `https://api.loan-system.com/v1`
- **Format**: JSON
- **Authentication**: Bearer Token

### Key Features
- Loan Application Processing
- Credit Score Checking
- Payment Management
- Document Handling

</v-clicks>

::right::

<br/>
<br/>

<v-click>

```bash
# Authentication Header
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

# Base Request Format
Content-Type: application/json
Accept: application/json
```

</v-click>

---
layout: default
---

# Endpoints Overview

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Loan Applications
```http
POST   /loans/apply
GET    /loans/{id}
PUT    /loans/{id}/approve
PUT    /loans/{id}/reject
```

### Credit Operations
```http
GET    /credit-score/{applicantId}
POST   /credit-check
PUT    /credit-update
```

</v-clicks>

</div>
<div>

<v-clicks>

### Payment Management
```http
POST   /payments/{loanId}
GET    /payments/{loanId}/schedule
PUT    /payments/{periodId}
```

### Document Management
```http
POST   /documents/upload
GET    /documents/{loanId}
PUT    /documents/{id}/verify
```

</v-clicks>

</div>
</div>

---
layout: default
---

# Detailed Endpoint: Loan Application

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### POST /loans/apply
Create a new loan application

**Request Body**:
```json
{
  "applicantId": "string",
  "amount": "number",
  "term": "number",
  "purpose": "string",
  "evidence": [{
    "type": "string",
    "description": "string",
    "url": "string"
  }]
}
```

</v-clicks>

</div>
<div>

<v-clicks>
<div style="overflow-y: auto;height: 500px;">
**Success Response (201)**:
```json
{
  "id": "LOAN-001",
  "status": "PENDING",
  "appliedAt": "2024-01-20T15:00:00Z",
  "applicantId": "APP-001",
  "amount": 10000,
  "term": 12,
  "evidence": [{
    "id": "DOC-001",
    "type": "INCOME_STATEMENT",
    "status": "PENDING"
  }]
}
```



**Error Response (400)**:
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid loan amount",
    "details": {
      "field": "amount",
      "constraint": "must be positive"
    }
  }
}
```
</div>

</v-clicks>

</div>
</div>

---
layout: default
---

# Error Handling Standards

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Common Error Codes
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": {
      "field": "affected_field",
      "reason": "specific reason"
    }
  }
}
```

### HTTP Status Codes
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 409: Conflict
- 422: Unprocessable Entity
- 429: Too Many Requests
- 500: Internal Server Error

</v-clicks>

</div>
<div>

<v-clicks>
<div style="overflow-y: auto;height: 500px;">

### Business Error Codes
```typescript
// Loan Application Errors
LOAN_AMOUNT_EXCEEDED
INSUFFICIENT_CREDIT_SCORE
MISSING_DOCUMENTS
INVALID_EVIDENCE

// Payment Errors
PAYMENT_OVERDUE
INSUFFICIENT_AMOUNT
INVALID_PAYMENT_PERIOD

// Document Errors
INVALID_DOCUMENT_TYPE
UPLOAD_FAILED
VERIFICATION_FAILED
```

### Error Response Examples
```json
{
  "error": {
    "code": "INSUFFICIENT_CREDIT_SCORE",
    "message": "Credit score below minimum",
    "details": {
      "minimum": 650,
      "current": 600
    }
  }
}
```
</div>

</v-clicks>

</div>
</div>

---
layout: default
---

# Security & Rate Limiting

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Authentication
```http
Authorization: Bearer <token>
```

### Rate Limits
- 100 requests per minute per IP
- 1000 requests per hour per token
- Burst: 10 concurrent requests

### Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

</v-clicks>

</div>
<div>

<v-clicks>

### Security Requirements
- TLS 1.2 or higher
- Token expiration: 1 hour
- IP whitelisting for admin APIs
- Regular security audits

### Response Headers
```http
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

</v-clicks>

</div>
</div>

---
layout: default
---
# Summary Part 1

<v-clicks>

- 5W1H is a framework for writing documentation
- Write Comment Practical
- API Document Template

</v-clicks>


---
layout: default
---

# Part 2: Documentation Tools

 - Go Doc
 - Swagger
 - Api Doc
 - Docusarus

---
 layout: default
---

# Go Doc  

<v-clicks>

## Go Doc Features
- Built-in documentation generator
- Doc comments (`//` and `/**/`)
- Formatting conventions
- godoc.org integration

</v-clicks>

::right::

<v-click>

<br/>
<br/>

```go
// Package apiserver provides REST API functionality
package apiserver

// User represents the user model
type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}
```
```bash
go get golang.org/x/tools/cmd/godoc
godoc -http=:6060 -index
```
// Multiple packages
```bash
# For multiple packages
for pkg in $(go list ./...); do
    dir="docs/${pkg#*/}"
    mkdir -p "$dir"
    godoc -url="/pkg/$pkg" > "$dir/index.html"
done
```
</v-click>

---
layout: default
---

# Custom Go Doc

<div  style="height: 400px;overflow-y: auto;">
```bash
#!/bin/bash
# generate-docs.sh

# Set output directory
OUTPUT_DIR="docs"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate main index
echo "Generating main index..."
godoc -url="/" > "$OUTPUT_DIR/index.html"

# Generate package documentation
echo "Generating package documentation..."
for pkg in $(go list ./...); do
    # Create directory structure
    dir="$OUTPUT_DIR/${pkg#*/}"
    mkdir -p "$dir"
    
    # Generate documentation
    echo "Processing $pkg..."
    godoc -url="/pkg/$pkg" > "$dir/index.html"
done

echo "Documentation generated in $OUTPUT_DIR"

```
</div>

---
layout: default
---

# APIDoc

<v-clicks>

- APIDoc is a tool that generates API documentation from Any source code.
```javascript
/**
 * @api {get} /user/:id Request User information
 * @apiName GetUser
 * @apiGroup User
 *
 * @apiParam {Number} id Users unique ID.
 *
 * @apiSuccess {String} firstname Firstname of the User.
 * @apiSuccess {String} lastname  Lastname of the User.
 */
```

- Install APIDoc

```bash
 sudo npm install -g apidoc
```
</v-clicks>

---
layout: default
---

# Swagger/OpenAPI with Go

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

- go-swagger implementation
- Annotation-based docs
- Swagger UI generation
- Service integration

</v-clicks>

</div>
<div>

<v-click>
```bash
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/swag
swag init -g cmd/server/main.go -o cmd/server/docs --parseDependency
```

```go
// Import swagger docs in server.go
_ "api/cmd/server/docs" // Import swagger docs
```


```go
// @title User Management API
// @version 1.0
// @description User management server
// @host localhost:8080
// @BasePath /api/v1

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
func (h *Handler) CreateUser() {}
```

</v-click>

</div>
</div>

---

# Docusaurus

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

## Key Features
- Built on React
- Markdown support
- Versioning system
- Search functionality
- i18n support
- Plugin ecosystem

## Quick Start
```bash
# Create new project
npx create-docusaurus@latest my-api-docs classic

# Start development server
cd my-api-docs
npm start
```

</v-clicks>

</div>
<div>

<v-clicks>

## Project Structure
```
my-api-docs/
├── docs/
│   ├── intro.md
│   └── api/
│       ├── overview.md
│       └── endpoints.md
├── src/
│   └── pages/
├── docusaurus.config.js
└── sidebars.js
```

## Example Doc
<div  style="height: 200px;overflow-y: auto;">
```md
---
sidebar_position: 1
---

# API Overview

Welcome to our API documentation.

## Authentication
Bearer token authentication...

## Base URL
https://api.example.com/v1
```
</div>

</v-clicks>

</div>
</div>

---
layout: default
---

# Docusaurus Configuration

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

## docusaurus.config.js
```javascript
module.exports = {
  title: 'API Documentation',
  tagline: 'Clear, versioned API docs',
  url: 'https://your-domain.com',
  baseUrl: '/',
  organizationName: 'your-org',
  projectName: 'api-docs',
  
  presets: [
    ['@docusaurus/preset-classic', {
      docs: {
        sidebarPath: 'sidebars.js',
        editUrl: 'https://github.com/org/repo/edit/main/',
      },
    }],
  ],
}
```

</v-clicks>

</div>
<div>

<v-clicks>

## Features Setup
- Versioning
```bash
npm run docusaurus docs:version 1.0.0
```

- Search
```javascript
themeConfig: {
  algolia: {
    appId: 'YOUR_APP_ID',
    apiKey: 'YOUR_API_KEY',
    indexName: 'YOUR_INDEX_NAME',
  },
}
```

- i18n
```javascript
i18n: {
  defaultLocale: 'en',
  locales: ['en', 'fr', 'es'],
}
```

</v-clicks>

</div>
</div>

---
layout: default
---

# Summary Part 2

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

## Documentation Tools Overview
- **Go Doc**: Native Go package docs, minimal setup
- **Swagger**: Interactive API docs, code generation
- **APIDoc**: Simple annotation-based docs
- **Docusaurus**: Full documentation websites

## Key Selection Factors
- Project size and complexity
- Target audience
- Maintenance requirements
- Team expertise

</v-clicks>

</div>
<div>

<v-clicks>

## Quick Selection Guide
| Tool | Best For | Key Feature |
|------|----------|-------------|
| Go Doc | Go packages | Code-first docs |
| Swagger | REST APIs | Interactive testing |
| APIDoc | Simple APIs | Quick setup |
| Docusaurus | Large projects | Full features |

## Additional Considerations
- Integration requirements
- Customization needs
- Versioning support
- Deployment options

</v-clicks>

</div>
</div>

---
layout: section
---

# Part 3: API Testing

---
layout: two-cols
---

# Unit Testing


### Table-Driven Tests:

1. Test Structure
   - Use table-driven tests to test multiple scenarios
   - Name tests clearly and descriptively
   - Group related tests together
   - Follow the AAA pattern (Arrange, Act, Assert)

2. Mocking
   - Create interfaces for dependencies
   - Implement mock versions for testing
   - Use dependency injection for better testability

::right::
<br/>
<br/>
<br/>

3. Test Coverage
   - Test happy paths and error cases
   - Test edge cases and boundary conditions
   - Include negative testing scenarios

4. Test Types
   - Unit tests (`TestCreateUser`)
   - Parallel tests (`TestUserService_CreateUser_Parallel`)
   - Benchmarks (`BenchmarkCreateUser`)
   - Examples (`ExampleCreateUser`)

---
layout: two-cols
---

# Unit Testing
### Table-Driven Tests:

5. Best Practices
   - Use subtests for better organization (`t.Run`)
   - Run tests in parallel when possible
   - Use descriptive error messages
   - Isolate tests (each test should be independent)
   - Keep tests maintainable and readable

6. Additional Testing Features
   - Benchmarking with `testing.B`
   - Example tests with output verification
   - Usage of test helpers and utilities
   - Comments explaining test purposes

::right::
<br/>
<br/>
<br/>
7. Error Handling
   - Test both success and error scenarios
   - Verify error messages
   - Use meaningful assertions

---
layout: default
---

# Understand the workflow
## Example: Loan Application Workflow

<div  style="height: 400px;overflow-y: auto;">
```mermaid
stateDiagram-v2
    [*] --> ApplicationSubmitted: Submit Application
    
    state "Application Process" as AP {
        ApplicationSubmitted --> DocumentVerification: Upload Evidence
        DocumentVerification --> CreditCheck: Verify Documents
        
        state CreditCheck {
            Check --> Score: Get Credit Score
            Score --> Risk: Calculate Risk
            Risk --> Interest: Determine Rate
        }
        
        CreditCheck --> ApplicationReview: Complete Check
    }
    
    ApplicationReview --> Rejected: Fail Criteria
    ApplicationReview --> Approved: Meet Criteria
    
    state "Loan Disbursement" as LD {
        Approved --> FundsTransfer: Initiate Transfer
        FundsTransfer --> Disbursed: Transfer Complete
        Disbursed --> PaymentSchedule: Generate Schedule
    }
    
    state "Payment Cycle" as PC {
        PaymentSchedule --> InvoiceGenerated: Generate Invoice
        
        state "Payment Processing" as PP {
            InvoiceGenerated --> PaymentDue: Wait for Due Date
            PaymentDue --> PaymentReceived: Payment Made
            PaymentDue --> Overdue: Miss Payment
            Overdue --> FineCalculation: Calculate Fine
            FineCalculation --> UpdatedInvoice: Add Fine
            UpdatedInvoice --> PaymentDue: New Due Date
        }
        
        PaymentReceived --> PaymentValidation: Validate Amount
        
        state PaymentValidation {
            state "Amount Check" as AC
            state "Complete Payment" as CP
            state "Partial Payment" as PP2
            
            AC --> CP: Full Amount
            AC --> PP2: Partial Amount
            PP2 --> UpdatedInvoice: Remaining Balance
        }
    }
    
    PaymentValidation --> PeriodClosed: Payment Complete
    PeriodClosed --> [*]: All Periods Paid
    
    note right of AP
        Documents required:
        - Income statement
        - Bank statements
        - Employment verification
    end note
    
    note right of CreditCheck
        Factors considered:
        - Credit score
        - Income level
        - Existing debt
        - Employment history
    end note
    
    note right of PC
        Payment cycle includes:
        - Monthly invoicing
        - Payment tracking
        - Late fee calculation
        - Status updates
    end note
```
</div>

---
layout: default
---

# Loan Process Test Scenarios


## 1. Application Submission Scenarios

<div  style="height: 400px;overflow-y: auto;">

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| APP-001 | Valid application submission | New | Amount: $10,000<br>Term: 12 months<br>Income: $5,000/month<br>All required docs | Status: PENDING<br>Application ID generated | - Verify all documents stored<br>- Check timestamp<br>- Verify notification sent |
| APP-002 | Missing required documents | New | Amount: $10,000<br>Term: 12 months<br>Missing bank statement | Error: Missing required documents<br>Status: INCOMPLETE | - Error message details<br>- Document checklist updated |
| APP-003 | Invalid loan amount | New | Amount: -$5,000<br>Term: 12 months<br>All docs | Error: Invalid amount<br>Status: REJECTED | - Validation error details |
| APP-004 | Invalid loan term | New | Amount: $10,000<br>Term: 0 months<br>All docs | Error: Invalid term<br>Status: REJECTED | - Term validation message |

## 2. Credit Check Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| CRD-001 | Excellent credit score | PENDING | Credit Score: 800<br>Income: $5,000/month | Status: REVIEWING<br>Risk: LOW | - Interest rate calculation<br>- Approval recommendation |
| CRD-002 | Poor credit score | PENDING | Credit Score: 550<br>Income: $5,000/month | Status: REJECTED<br>Risk: HIGH | - Rejection reason recorded |
| CRD-003 | Borderline credit case | PENDING | Credit Score: 650<br>Income: $5,000/month | Status: REVIEWING<br>Risk: MEDIUM | - Manual review flag<br>- Additional checks needed |
| CRD-004 | Income verification failed | PENDING | Credit Score: 750<br>Income docs invalid | Status: REJECTED<br>Error: Income verification failed | - Document validation errors |

## 3. Loan Approval Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| APR-001 | Standard approval | REVIEWING | Risk: LOW<br>Interest: 5% | Status: APPROVED<br>Payment schedule generated | - Schedule accuracy<br>- Interest calculations |
| APR-002 | High-risk approval | REVIEWING | Risk: HIGH<br>Interest: 12% | Status: APPROVED<br>Higher interest rate | - Risk factor documentation<br>- Additional terms |
| APR-003 | Conditional approval | REVIEWING | Risk: MEDIUM<br>Additional collateral | Status: APPROVED<br>With conditions | - Condition documentation<br>- Follow-up tasks |
| APR-004 | Manual rejection | REVIEWING | Risk: HIGH<br>Insufficient income | Status: REJECTED<br>Reason documented | - Rejection notification<br>- Appeal process |

## 4. Disbursement Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| DSB-001 | Successful disbursement | APPROVED | Valid bank details<br>All conditions met | Status: DISBURSED<br>Funds transferred | - Transfer confirmation<br>- Schedule activation |
| DSB-002 | Failed bank transfer | APPROVED | Invalid bank details | Error: Transfer failed<br>Status: APPROVED | - Error handling<br>- Retry mechanism |
| DSB-003 | Partial disbursement | APPROVED | Multiple tranches<br>Schedule defined | Status: PARTIALLY_DISBURSED | - Tranche schedule<br>- Partial activation |
| DSB-004 | Cancelled before disbursement | APPROVED | Cancellation request | Status: CANCELLED<br>No transfer | - Cancellation reason<br>- Cleanup actions |

## 5. Payment Processing Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| PAY-001 | On-time full payment | ACTIVE | Amount: Full due<br>Date: Before due | Status: PAID<br>Period closed | - Payment allocation<br>- Next period update |
| PAY-002 | Late payment with penalty | ACTIVE | Amount: Full + penalty<br>Date: After due | Status: PAID<br>Penalty applied | - Penalty calculation<br>- Payment allocation |
| PAY-003 | Partial payment | ACTIVE | Amount: 50% of due<br>Date: On due | Status: PARTIALLY_PAID | - Remaining balance<br>- Next due date |
| PAY-004 | Overpayment | ACTIVE | Amount: 120% of due | Status: PAID<br>Excess allocated | - Excess allocation<br>- Next period adjustment |

## 6. Late Payment Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| LTE-001 | Payment 1-30 days late | OVERDUE | Days late: 15<br>Amount: Full due | Status: PAID<br>Standard penalty | - Penalty calculation<br>- Credit report impact |
| LTE-002 | Payment 30-60 days late | OVERDUE | Days late: 45<br>Amount: Full due | Status: PAID<br>Higher penalty | - Escalated penalties<br>- Collection actions |
| LTE-003 | Default threshold reached | OVERDUE | Days late: 90+<br>No payment | Status: DEFAULT<br>Collection process | - Default procedures<br>- Legal actions |
| LTE-004 | Payment plan negotiation | OVERDUE | Restructure request | Status: RESTRUCTURED | - New payment schedule<br>- Terms modification |

## 7. Loan Closure Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| CLS-001 | Normal completion | ACTIVE | All payments made | Status: CLOSED<br>Completion letter | - Final calculations<br>- Document generation |
| CLS-002 | Early payoff | ACTIVE | Full remaining balance | Status: CLOSED<br>Early payoff | - Payoff calculation<br>- Fee adjustments |
| CLS-003 | Closure after restructure | RESTRUCTURED | Final payment made | Status: CLOSED<br>Modified terms met | - Modified terms check<br>- History documentation |
| CLS-004 | Write-off closure | DEFAULT | Approved write-off | Status: WRITTEN_OFF | - Write-off approvals<br>- Tax implications |

## 8. Special Case Scenarios

| ID | Scenario | Initial State | Input Data | Expected Result | Additional Checks |
|---|---|---|---|---|---|
| SPC-001 | Death of borrower | ACTIVE | Death certificate | Status: SPECIAL_HANDLING | - Insurance claims<br>- Estate process |
| SPC-002 | Bankruptcy filing | ACTIVE | Bankruptcy notice | Status: LEGAL_REVIEW | - Legal procedures<br>- Collection freeze |
| SPC-003 | Fraud detection | ACTIVE | Fraud indicators | Status: INVESTIGATION | - Investigation process<br>- Legal actions |
| SPC-004 | Natural disaster relief | ACTIVE | Disaster declaration | Status: PAYMENT_HOLIDAY | - Relief terms<br>- Documentation |

</div>

---
layout: default
---
# Example Unit Test

- Prepare data
```go
// Prepare data
loanApp := &loan.LoanApplication{
	ID:          "APP-0001",
	ApplicantID: generateUniqueID(), // Generate unique ID
	Amount:      10000,
	Term:        12,
}
```
- Call function
```go
// Call function
t.Run("Apply for loan", func(t *testing.T) {
	err := loanService.ApplyForLoan(loanApp, evidence)
})
``` 

- Assert result
```go
// Assert result
if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, loan.StatusPending, tt.application.Status)
				assert.False(t, tt.application.AppliedAt.IsZero())
			}
```

---
layout: default
---

# Integration Testing Best Practices


<img src="/images/testing/api-integration-testing.svg" style="height: 400px;overflow-y: auto;">
 
---
layout: two-cols
---

# Integration Testing

<v-clicks>

### HTTP Testing
- Test server setup
- Request handling
- Response validation
- Cleanup management

</v-clicks>

::right::

<v-click>
<br/>
<br/>
<div style="height: 450px;overflow-y: auto;">
```go
func TestLoanAPI_Integration(t *testing.T) {
	baseURL := "http://localhost:4000"

	t.Run("Complete loan application flow", func(t *testing.T) {
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
		applyBody := `{
			"id": "APP-0010",
			"applicant_id": "USR-001",
			"amount": 10000,
			"term": 12,
			"purpose": "Home Improvement",
			"evidence": [
				{
					"id": "EVI-008",
					"type": "INCOME_STATEMENT",
					"description": "Monthly Income Statement"
				}
			]
		}`

		applyResp, err := makeRequest(http.MethodPost,
			baseURL+"/loans/apply",
			applyBody,
			token,
		)

		fmt.Println(token)
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
  }
```
</div>
</v-click>

---
layout: default
---

# Best Practices

<div grid="~ cols-2 gap-4">
<div>

<v-clicks>

### Organization
- Table-driven tests
- Separate test types
- Helper utilities

### Data Management
- Test fixtures
- Cleanup routines
- Database handling

</v-clicks>

</div>
<div>

<v-clicks>

### Quality Assurance
- Coverage tracking
- Error verification
- Edge case testing

### CI/CD Integration
- Automated testing
- Coverage reports
- Pipeline integration

</v-clicks>

</div>
</div>

---
layout: two-cols
---

# Test Reporting

<v-clicks>

### Automated Reports
- Test results summary
- Coverage metrics
- Performance benchmarks
- HTML report generation

### Integration
- CI pipeline reports
- Team dashboards
- Trend analysis
- Alert configuration

</v-clicks>

::right::

<v-click>

```go {all|2-8|10-16}
// test_report.go
type TestReport struct {
    TotalTests   int       `json:"total"`
    PassedTests  int       `json:"passed"`
    FailedTests  int       `json:"failed"`
    Coverage     float64   `json:"coverage"`
    Duration     string    `json:"duration"`
}

func generateReport() {
    cmd := exec.Command("go", "test", 
        "-json", "./...", 
        "-coverprofile=coverage.out")
    output, _ := cmd.Output()
    // Process test output
    report := parseTestOutput(output)
}
```

```bash
# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Run tests with JSON output
go test -json ./... | tee test-report.json
```

</v-click>

---
layout: section
---

# Part 4: Testing Tools

---
layout: default
---

# Playwright

<v-clicks>

- Playwright is a browser automation tool
- It can be used to test web applications
- It can be used to test API
- It can be used to test mobile applications

</v-clicks>

## Install Playwright

```bash
npm install playwright
```



---
layout: default
---

# Playwright Test

```typescript
test.describe('Loan API Tests', () => {
  test.describe('Apply for Loan', () => {
    test('should successfully apply for a loan', async ({ request }) => {
      // Create a loan application
      const loanApplication = createSampleLoanApplication();
      
      const response = await request.post(`${API_BASE_URL}/loans/apply`, {
        data: loanApplication
      });
      
      expect(response.status()).toBe(201);
      
      const responseBody = await response.json();
      expect(responseBody.ID).toBe(loanApplication.ID);
      expect(responseBody.Status).toBe('PENDING');
    }); 
  });
});
```

---
layout: default
---

# Agentic AI with Testing

<v-clicks>

- Agentic AI Testing is a new way to write document and test APIs

</v-clicks>

## Install Clade Desktop

```bash
brew install clade-desktop
```

## Install Clade CLI

```bash
brew install clade-cli
```

## Login to Clade

```bash
clade login
```

---
layout: default
---

## Install Playwright MCP Server

```bash
npm install -g @executeautomation/playwright-mcp-server
```

## Install Filesystem MCP Server

```bash
npm install -g @modelcontextprotocol/server-filesystem
```


## Setup Filesystem and Playwright MCP Server Config

```json
{
    "mcpServers": {
        "playwright": {
            "command": "npx",
            "args": [
                "-y",
                "@executeautomation/playwright-mcp-server"
            ]
        },
        "filesystem": {
            "command": "npx",
            "args": [
                "-y",
                "@modelcontextprotocol/server-filesystem",
                "/Users/a677018/Projects/courses/api-design-v3/api-design-v3/src"
            ]
        }
    }
}
```

---
layout: default
---

# Reopen Clade Desktop
## Write prompt to  write document and test
```bash
I need help with API documentation and automated testing for my project. Please:

1. Write API documentation for the following endpoint:
[Insert your endpoint code/details here]
Include:
- Endpoint URL and method
- Request parameters/body format
- Response format and status codes
- Error handling
- Authentication requirements
- Example requests and responses

2. Create Playwright API tests for this endpoint covering:
- Happy path scenarios
- Error/edge cases
- Authentication handling
- Response validation

The API is currently running at http://localhost:4000. Create tests in the /tests directory using TypeScript
```

---
layout: center
class: "text-center"
---

# Thank You!

[Documentation](https://golang.org/doc) · [Testing](https://golang.org/pkg/testing/) · [Tools](https://github.com/golang/go/wiki/Projects)
