-- Status and PaymentStatus enums will be stored as TEXT
-- Evidence table to store supporting documents
/*select *
 from evidence
 select *
 from loan_applications
 select *
 from payment_periods -- drop table evidence
 -- drop table loan_applications
 -- drop table payment_periods
 */
CREATE TABLE evidence (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    -- e.g., "INCOME_STATEMENT", "BANK_STATEMENT"
    description TEXT,
    url TEXT,
    uploaded_at TIMESTAMP NOT NULL,
    loan_application_id TEXT,
    FOREIGN KEY (loan_application_id) REFERENCES loan_applications(id)
);
-- Main loan applications table
CREATE TABLE loan_applications (
    id TEXT PRIMARY KEY,
    applicant_id TEXT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    term INTEGER NOT NULL,
    -- In months
    purpose TEXT,
    status TEXT NOT NULL,
    -- PENDING, REVIEWING, APPROVED, etc.
    credit_score INTEGER,
    interest_rate DECIMAL(5, 2),
    applied_at TIMESTAMP NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    approved_at TIMESTAMP,
    disbursed_at TIMESTAMP,
    CHECK (amount > 0),
    CHECK (term > 0),
    CHECK (interest_rate >= 0),
    CHECK (
        status IN (
            'PENDING',
            'REVIEWING',
            'APPROVED',
            'REJECTED',
            'DISBURSED',
            'COMPLETED',
            'DEFAULTED'
        )
    )
);
-- Payment periods table
CREATE TABLE payment_periods (
    id TEXT PRIMARY KEY,
    loan_id TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    interest_amount DECIMAL(15, 2) NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    paid_amount DECIMAL(15, 2) DEFAULT 0,
    fine_amount DECIMAL(15, 2) DEFAULT 0,
    status TEXT NOT NULL,
    -- PENDING, PAID, OVERDUE, INCOMPLETE
    paid_at TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loan_applications(id),
    CHECK (amount >= 0),
    CHECK (interest_amount >= 0),
    CHECK (principal_amount >= 0),
    CHECK (paid_amount >= 0),
    CHECK (fine_amount >= 0),
    CHECK (
        status IN ('PENDING', 'PAID', 'OVERDUE', 'INCOMPLETE')
    )
);
-- Indexes for better query performance
CREATE INDEX idx_loan_applications_status ON loan_applications(status);
CREATE INDEX idx_loan_applications_applicant ON loan_applications(applicant_id);
CREATE INDEX idx_payment_periods_loan ON payment_periods(loan_id);
CREATE INDEX idx_payment_periods_status ON payment_periods(status);
CREATE INDEX idx_payment_periods_due_date ON payment_periods(due_date);
CREATE INDEX idx_evidence_loan ON evidence(loan_application_id);
-- select * from evidence
-- select * from loan_applications
-- select * from payment_periods
/*

 delete from payment_periods
 */
-- Add credit scores table
CREATE TABLE credit_scores (
    id TEXT PRIMARY KEY,
    applicant_id TEXT NOT NULL,
    credit_score INTEGER NOT NULL,
    checked_at TIMESTAMP NOT NULL,
    source TEXT NOT NULL,
    -- e.g., "EXPERIAN", "EQUIFAX", "TRANSUNION"
    CHECK (
        credit_score >= 300
        AND credit_score <= 850
    )
);
-- Index for quick lookups
CREATE INDEX idx_credit_scores_applicant ON credit_scores(applicant_id, checked_at);
CREATE TABLE loan_payments (
    id VARCHAR(36) PRIMARY KEY,
    loan_id VARCHAR(36) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loan_applications(id),
    CHECK (status IN ('PENDING', 'PAID', 'OVERDUE'))
);
-- Indexes for better query performance
CREATE INDEX idx_loan_payments_loan_id ON loan_payments(loan_id);
CREATE INDEX idx_loan_payments_status ON loan_payments(status);
CREATE INDEX idx_loan_payments_due_date ON loan_payments(due_date);


 delete  from evidence
 delete from loan_applications
select * from loan_applications


UPDATE loan_applications 
		SET credit_score = 750 , interest_rate = 0.05
		WHERE id = 'APP-0010'

        UPDATE loan_applications 
                SET credit_score = 750 , interest_rate = 0.05
                WHERE id = 'APP-0010'


                SELECT id, applicant_id, amount, term, purpose, status,
			   credit_score, interest_rate, applied_at, last_updated_at,
			   approved_at, disbursed_at
		FROM loan_applications WHERE id = 'APP-0010'