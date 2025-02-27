# API Documentation Template

## Overview
A clear, concise description of your API that answers the fundamental questions using the 5W1H framework.

### What
- **API Name**: [Your API Name]
- **Version**: [e.g., v1.0.0]
- **Description**: A brief overview of what your API does and its primary functions
- **Key Features**:
  - Feature 1
  - Feature 2
  - Feature 3

### Who
- **Target Users**: Who is this API designed for?
  - Primary audience (e.g., developers, system administrators)
  - Technical expertise level required
- **Stakeholders**: Who is involved in the API lifecycle?
  - Maintainers
  - Contributors
  - Support team

### Where
- **Hosting Environment**:
  - Production endpoint: `https://api.example.com/v1`
  - Staging endpoint: `https://staging-api.example.com/v1`
  - Development endpoint: `https://dev-api.example.com/v1`
- **Source Code**: Repository location
- **Documentation**: Where to find additional resources

### When
- **Release Information**:
  - Initial release date
  - Current version release date
  - Update frequency
- **Availability**:
  - Service level agreements (SLA)
  - Maintenance windows
  - Rate limiting details

### Why
- **Purpose**: The problem this API solves
- **Benefits**:
  - Key advantages
  - Business value
  - Technical benefits
- **Use Cases**: Common scenarios where this API is valuable

### How
- **Authentication**:
  ```
  Authorization: Bearer <your_api_key>
  ```
  - How to obtain API keys
  - Security best practices

- **Request Format**:
  ```json
  {
    "required_field": "value",
    "optional_field": "value"
  }
  ```

- **Response Format**:
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

## Endpoints

### [Endpoint Name]

**GET /resource**

What:
- Purpose of this endpoint
- Data returned

Who:
- Required permissions
- Rate limits specific to this endpoint

Where:
- Complete endpoint URL
- Related endpoints

When:
- Best time to use this endpoint
- Processing time expectations
- Caching details

Why:
- Use cases for this endpoint
- Benefits of using this endpoint

How:
- Request parameters
- Headers required
- Example request:
  ```bash
  curl -X GET "https://api.example.com/v1/resource" \
       -H "Authorization: Bearer your_api_key"
  ```
- Example response:
  ```json
  {
    "status": "success",
    "data": {
      // Response data structure
    }
  }
  ```

## Error Handling

### Error Response Format
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      // Additional error context
    }
  }
}
```

### Common Error Codes
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `429`: Too Many Requests
- `500`: Internal Server Error

## Best Practices

### Security
- Use HTTPS for all requests
- Implement rate limiting
- Follow OAuth 2.0 standards
- Regular security audits
- Input validation

### Performance
- Implement caching
- Pagination for large datasets
- Compression for large payloads
- Asynchronous operations for long-running tasks

### Versioning
- Semantic versioning (MAJOR.MINOR.PATCH)
- Version in URL path
- Deprecation policy
- Migration guides

## Support
- Contact information
- Issue reporting process
- SLA details
- Community resources
- Changelog

## SDK & Tools
- Official SDKs
- Community libraries
- Testing tools
- Development environments

## Getting Started

### Prerequisites
- Required software
- Account setup
- API key acquisition

### Quick Start Guide
1. Install dependencies
2. Configure authentication
3. Make your first API call
4. Handle responses

### Examples
- Code samples in multiple languages
- Common use case implementations
- Best practice demonstrations

## Updates & Maintenance
- Release schedule
- Deprecation notices
- Breaking changes policy
- Migration guides

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


# Loan Process Test Scenarios

## 1. Application Submission Scenarios

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