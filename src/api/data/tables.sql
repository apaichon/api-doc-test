select count(*) from api_logs;

select * from api_logs;

CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    status_id INTEGER NOT NULL
);

select * from users;

CREATE TABLE payments (
    payment_id TEXT PRIMARY KEY,
    amount REAL NOT NULL,
    payment_method TEXT NOT NULL,
    payment_date TEXT NOT NULL,
    pay_to TEXT NOT NULL,
    note TEXT,
    status TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);


-- Create Role table
DROP TABLE roles;
CREATE TABLE roles (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name TEXT NOT NULL,
    role_desc TEXT,
    is_super_admin BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    status_id INTEGER NOT NULL
);



-- Create UserRolesModel table

CREATE TABLE user_roles (
    user_role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    status_id INTEGER NOT NULL
);

-- Create RolePermissions table
DROP TABLE role_permissions;
CREATE TABLE role_permissions (
    role_permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    resource_type_id INTEGER NOT NULL,
    resource_name TEXT NOT NULL,
    can_execute BOOLEAN NOT NULL,
    can_read BOOLEAN NOT NULL,
    can_write BOOLEAN NOT NULL,
    can_delete BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    status_id INTEGER NOT NULL
);

CREATE TABLE api_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL,
    request_id TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    status_text TEXT NOT NULL,
    duration REAL NOT NULL,
    request_body TEXT,
    client_ip TEXT,
    client_browser TEXT,
    client_browser_version TEXT,
    client_os TEXT,
    client_os_version TEXT,
    client_device TEXT,
    user_id INTEGER,
    error TEXT
);

drop table api_logs;

CREATE TABLE api_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    log TEXT NOT NULL,
    created_at DATETIME default current_timestamp
);

drop table consents;
drop table consent_logs;
-- Consent table
CREATE TABLE IF NOT EXISTS consents (
    consent_id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_id TEXT NOT NULL,
    source_hospital TEXT NOT NULL,
    target_hospital TEXT NOT NULL,
    purpose TEXT NOT NULL,
    data_categories TEXT NOT NULL, -- Stored as JSON array
    start_date DATETIME NOT NULL,
    expiry_date DATETIME NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('PENDING', 'ACTIVE', 'REVOKED', 'EXPIRED')),
    version INTEGER NOT NULL DEFAULT 1,
    signature TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ConsentLog table
CREATE TABLE IF NOT EXISTS consent_logs (
    consent_log_id INTEGER PRIMARY KEY AUTOINCREMENT,
    consent_id int NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('CREATED', 'UPDATED', 'REVOKED', 'ACCESSED')),
    actor_id TEXT NOT NULL,
    actor_type TEXT NOT NULL CHECK (actor_type IN ('PATIENT', 'DOCTOR', 'SYSTEM')),
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    FOREIGN KEY (consent_id) REFERENCES consents(id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_consents_patient_id ON consents(patient_id);
CREATE INDEX IF NOT EXISTS idx_consent_logs_consent_id ON consent_logs(consent_id);
CREATE INDEX IF NOT EXISTS idx_consents_status ON consents(status);


select * from api_logs

select * from roles

select * from payments

SELECT COUNT(*) FROM payments 
                WHERE payment_method LIKE '%john%' OR pay_to LIKE '%john%' OR note LIKE '%john%'

INSERT INTO user_roles (role_id, user_id, created_at, created_by, updated_at, updated_by, status_id) VALUES
(1, 6, datetime('now'), 'system', datetime('now'), 'system', 1)


-- Sample data for roles
INSERT INTO roles (role_id, role_name, role_desc, is_super_admin, created_at, created_by, updated_at, updated_by, status_id) VALUES
(1, 'admin', 'System Administrator', true, datetime('now'), 'system', datetime('now'), 'system', 1),
(2, 'manager', 'Department Manager', false, datetime('now'), 'system', datetime('now'), 'system', 1),
(3, 'user', 'Regular User', false, datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 'guest', 'Guest User', false, datetime('now'), 'system', datetime('now'), 'system', 1);

-- Sample data for user_roles
INSERT INTO user_roles (role_id, user_id, created_at, created_by, updated_at, updated_by, status_id) VALUES
(3, 3, datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 4, datetime('now'), 'system', datetime('now'), 'system', 1),

-- Sample data for role_permissions
INSERT INTO role_permissions (role_id, resource_type_id, resource_name, can_execute, can_read, can_write, can_delete, created_at, created_by, updated_at, updated_by, status_id) VALUES
(1, 1, 'users', true, true, true, true, datetime('now'), 'system', datetime('now'), 'system', 1),
(2, 2, 'reports', false, true, false, false, datetime('now'), 'system', datetime('now'), 'system', 1),
(3, 3, 'products', true, true, true, false, datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 4, 'public', false, true, false, false, datetime('now'), 'system', datetime('now'), 'system', 1);

select * from  vw_user_permissions
drop view vw_user_permissions;
select * from payments
-- Create a view that combines roles, user_roles, and role_permissions
CREATE VIEW vw_user_permissions AS
SELECT 
    ur.user_id,
    r.role_id,
    r.role_name,
    r.is_super_admin,
    rp.role_permission_id,
    rp.resource_type_id,
    rp.resource_name,
    rp.can_execute,
    rp.can_read,
    rp.can_write,
    rp.can_delete
    /*CASE 
        WHEN r.is_super_admin = true THEN true 
        ELSE rp.can_execute 
    END as effective_can_execute,
    CASE 
        WHEN r.is_super_admin = true THEN true 
        ELSE rp.can_read 
    END as effective_can_read,
    CASE 
        WHEN r.is_super_admin = true THEN true 
        ELSE rp.can_write 
    END as effective_can_write,
    CASE 
        WHEN r.is_super_admin = true THEN true 
        ELSE rp.can_delete 
    END as effective_can_delete
    */
FROM 
    user_roles ur
    INNER JOIN roles r ON ur.role_id = r.role_id
    LEFT JOIN role_permissions rp ON r.role_id = rp.role_id
WHERE 
    ur.status_id = 1 
    AND r.status_id = 1 
    AND rp.status_id = 1;

-- Check if a user can perform an action on a resource

    SELECT 
    effective_can_read,
    effective_can_write,
    effective_can_execute,
    effective_can_delete
FROM vw_user_permissions 
WHERE user_id = 'admin123' 
AND resource_name = 'users';


select * from user_roles;
 SELECT 
   * from vw_user_permissions 
   where user_id = 5
   select * from role_permissions;

   select 


   INSERT INTO user_roles (user_role_id, role_id, user_id, created_at, created_by, updated_at, updated_by, status_id) VALUES
(1, 1, 'admin123', datetime('now'), 'system', datetime('now'), 'system', 1),
(2, 2, 'manager456', datetime('now'), 'system', datetime('now'), 'system', 1),
(3, 3, 'user789', datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 4, 'guest101', datetime('now'), 'system', datetime('now'), 'system', 1);

select * from roles;
select * from user_roles;

Select * from users;

SELECT * FROM users
             Where user_name like '%john%' OR password like '%john%' OR salt like '%john%'
            LIMIT 10 OFFSET 0


SELECT name FROM sqlite_master WHERE type='table';

select * from role_permissions;

INSERT INTO role_permissions (role_id, resource_type_id, resource_name, can_execute, can_read, can_write, can_delete, created_at, created_by, updated_at, updated_by, status_id) VALUES
(3, 1, '/api/users', true, true, false, false, datetime('now'), 'system', datetime('now'), 'system', 1),
(3, 1, '/api/products', false, true, false, false, datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 1, '/api/users', false, false, false, false, datetime('now'), 'system', datetime('now'), 'system', 1),
(4, 1, '/api/products', true, true, false, false, datetime('now'), 'system', datetime('now'), 'system', 1);


 SELECT sum(can_execute) as can_execute, sum(can_read) as can_read, sum(can_write) as can_write, sum(can_delete) as can_delete FROM vw_user_permissions
             Where user_id = 4 AND resource_type_id = 1 AND resource_name = '/api/users'

             select * from users
             select * from user_roles
             select count(*) as total_roles from user_roles
             ur left join roles r on ur.role_id = r.role_id
              where ur.user_id = 4 and r.is_super_admin =1
             select * from roles
             select * from role_permissions
             select * from vw_user_permissions


select * from user_roles where user_id = 3

select * from role_permissions where role_id = 3

select * from users

 INSERT INTO user_roles ( role_id, user_id, created_at, created_by, updated_at, updated_by, status_id) VALUES
(3, 5, datetime('now'), 'system', datetime('now'), 'system', 1)



SELECT sum(can_execute) as can_execute, sum(can_read) as can_read, sum(can_write) as can_write, sum(can_delete) as can_delete FROM vw_user_permissions
             Where user_id = 5 AND resource_type_id = 1 AND resource_name = '/api/users'


              select count(*) as total_roles from user_roles
             ur left join roles r on ur.role_id = r.role_id
              where ur.user_id = 5 and r.is_super_admin =1

              select count(*) as is_super_admin from user_roles ur left join roles r 
                         on ur.role_id = r.role_id
              where ur.user_id = 5 and r.is_super_admin =1

              select * from vw_user_permissions where resource_name = '/api/users'

              select * from role_permissions where resource_name = '/api/users'

              select * from user_roles where user_id = 5

               select * from vw_user_permissions where user_id = 4


               select * from users 

-- Consent table
CREATE TABLE IF NOT EXISTS consents (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    source_hospital TEXT NOT NULL,
    target_hospital TEXT NOT NULL,
    purpose TEXT NOT NULL,
    data_categories TEXT NOT NULL, -- Stored as JSON array
    start_date DATETIME NOT NULL,
    expiry_date DATETIME NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('PENDING', 'ACTIVE', 'REVOKED', 'EXPIRED')),
    version INTEGER NOT NULL DEFAULT 1,
    signature TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ConsentLog table
CREATE TABLE IF NOT EXISTS consent_logs (
    id TEXT PRIMARY KEY,
    consent_id TEXT NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('CREATED', 'UPDATED', 'REVOKED', 'ACCESSED')),
    actor_id TEXT NOT NULL,
    actor_type TEXT NOT NULL CHECK (actor_type IN ('PATIENT', 'DOCTOR', 'SYSTEM')),
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    FOREIGN KEY (consent_id) REFERENCES consents(id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_consents_patient_id ON consents(patient_id);
CREATE INDEX IF NOT EXISTS idx_consent_logs_consent_id ON consent_logs(consent_id);
CREATE INDEX IF NOT EXISTS idx_consents_status ON consents(status); 

-- PRAGMA table_info(evidence);
-- select * from evidence

-- Status and PaymentStatus enums will be stored as TEXT
-- Evidence table to store supporting documents
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

delete from loan_applications
select * from loan_applications

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