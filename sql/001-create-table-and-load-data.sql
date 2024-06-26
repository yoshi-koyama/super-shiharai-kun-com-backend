CREATE TABLE company (
    id CHAR(26) PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    representative_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    address VARCHAR(255) NOT NULL
);

CREATE TABLE user (
    id CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company(id)
);

CREATE TABLE client (
    id CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    representative_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    address VARCHAR(255) NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company(id)
);

CREATE TABLE client_bank_account (
    id CHAR(26) PRIMARY KEY,
    client_id CHAR(26) NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    branch_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (client_id) REFERENCES client(id)
);

CREATE TABLE invoice (
    id CHAR(26) PRIMARY KEY,
    company_id CHAR(26) NOT NULL,
    client_id CHAR(26) NOT NULL,
    issue_date DATE NOT NULL,
    payment_amount DECIMAL(10, 2) NOT NULL,
    fee DECIMAL(10, 2),
    fee_rate DECIMAL(5, 2),
    tax DECIMAL(10, 2),
    tax_rate DECIMAL(5, 2),
    invoice_amount DECIMAL(10, 2) NOT NULL,
    payment_due_date DATE NOT NULL,
    status ENUM('unprocessed', 'processing', 'paid', 'error') NOT NULL,
    FOREIGN KEY (company_id) REFERENCES company(id),
    FOREIGN KEY (client_id) REFERENCES client(id),
    INDEX payment_due_date_index (payment_due_date)
);

-- 初期データを登録
INSERT INTO company (id, company_name, representative_name, phone_number, postal_code, address) VALUES
('01J13SNTHMBQ8A5039QM3QB3JC', 'ABC Corporation', 'John Doe', '0123456789', '1234567', '123 Main St, City, Country'),
('01J13SP0GWQQ88GB4NWYC9NK47', 'XYZ Inc.', 'Jane Smith', '0987654321', '7654321', '456 Market St, City, Country');

INSERT INTO user (id, company_id, full_name, email_address) VALUES
('01J13SQDGNE4GXWFP91FK61T12', '01J13SNTHMBQ8A5039QM3QB3JC', 'Alice Johnson', 'alice@example.com'),
('01J13SQP75MV04K2B9TXNNB1FP', '01J13SP0GWQQ88GB4NWYC9NK47', 'Bob Brown', 'bob@example.com');

INSERT INTO client (id, company_id, company_name, representative_name, phone_number, postal_code, address) VALUES
('01J13SR5J68KDGNQ8Y6HECCSA5', '01J13SNTHMBQ8A5039QM3QB3JC', 'Client A', 'Tom Lee', '1122334455', '1112222', '789 High St, City, Country'),
('01J13SRB54ANAZ7HE3YFRXWW5Q', '01J13SP0GWQQ88GB4NWYC9NK47', 'Client B', 'Sara Wong', '5544332211', '3334444', '101 Low St, City, Country');

INSERT INTO client_bank_account (id, client_id, bank_name, branch_name, account_number, account_name) VALUES
('01J13SRVSK3BP6HHBNSEMQ8D2H', '01J13SR5J68KDGNQ8Y6HECCSA5', 'Bank A', 'Branch A', '1234567890', 'Tom Lee'),
('01J13SS0REDNCE2HGGC581BNZA', '01J13SRB54ANAZ7HE3YFRXWW5Q', 'Bank B', 'Branch B', '0987654321', 'Sara Wong');

INSERT INTO invoice (id, company_id, client_id, issue_date, payment_amount, fee, fee_rate, tax, tax_rate, invoice_amount, payment_due_date, status) VALUES
('01J13SS623G1G3AMBD7YRMSC1D', '01J13SNTHMBQ8A5039QM3QB3JC', '01J13SR5J68KDGNQ8Y6HECCSA5', '2023-06-01', 1000.00, 50.00, 5.00, 80.00, 8.00, 1130.00, '2023-07-01', 'unprocessed'),
('01J13SSC6KDYG1FXW2TYTQ0GE7', '01J13SP0GWQQ88GB4NWYC9NK47', '01J13SRB54ANAZ7HE3YFRXWW5Q', '2023-06-15', 2000.00, 100.00, 5.00, 160.00, 8.00, 1130.00, '2023-07-01', 'processing');

