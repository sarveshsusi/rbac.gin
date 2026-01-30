-- =========================================================
-- 1. SETUP & EXTENSIONS
-- =========================================================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================================================
-- 2. USER MANAGEMENT & AUTH
-- =========================================================

-- USERS
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'support', 'customer')),
    is_active BOOLEAN DEFAULT TRUE,
    must_reset_password BOOLEAN DEFAULT FALSE,
    created_by UUID,
    two_fa_enabled BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);

-- REFRESH TOKENS
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- PASSWORD RESET TOKENS
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ
);

-- 2FA OTPS
CREATE TABLE IF NOT EXISTS two_fa_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    code TEXT,
    expires_at TIMESTAMPTZ,
    used BOOLEAN,
    created_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_two_fa_otps_user_id ON two_fa_otps(user_id);

-- REMEMBERED DEVICES
CREATE TABLE IF NOT EXISTS remembered_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ
);

-- =========================================================
-- 3. PROFILES
-- =========================================================

-- CUSTOMERS
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    company VARCHAR(150),
    phone VARCHAR(20),
    address TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- SUPPORT ENGINEERS
CREATE TABLE IF NOT EXISTS support_engineers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    designation VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- =========================================================
-- 4. PRODUCT CATALOG
-- =========================================================

-- BRANDS
CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

-- CATEGORIES
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

-- BRAND CATEGORIES (Many-to-Many)
CREATE TABLE IF NOT EXISTS brand_categories (
    brand_id UUID REFERENCES brands(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (brand_id, category_id)
);

-- MODELS
CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    brand_id UUID REFERENCES brands(id)
);
CREATE INDEX IF NOT EXISTS idx_models_brand_id ON models(brand_id);

-- PRODUCTS
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    brand_id UUID NOT NULL REFERENCES brands(id),
    model_id UUID NOT NULL REFERENCES models(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ
);

-- =========================================================
-- 5. RELATIONSHIPS & CONTRACTS
-- =========================================================

-- CUSTOMER PRODUCTS
CREATE TABLE IF NOT EXISTS customer_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id),
    product_id UUID REFERENCES products(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    UNIQUE(customer_id, product_id)
);
CREATE INDEX IF NOT EXISTS idx_customer_products_customer ON customer_products(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_products_product ON customer_products(product_id);

-- AMC CONTRACTS
CREATE TABLE IF NOT EXISTS amc_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_product_id UUID REFERENCES customer_products(id),
    sla_hours INT,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    status TEXT,
    created_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_amc_contracts_customer_product ON amc_contracts(customer_product_id);

-- =========================================================
-- 6. TICKETING SYSTEM
-- =========================================================

-- TICKETS
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id),
    product_id UUID REFERENCES products(id),
    amc_id UUID REFERENCES amc_contracts(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) CHECK (status IN ('Open', 'Assigned', 'In Progress', 'Closed')) DEFAULT 'Open',
    priority VARCHAR(30),
    support_mode VARCHAR(50), -- On-site, Remote, Phone
    service_call_type VARCHAR(50), -- Warranty, Service, AMC
    closure_proof_image TEXT,
    sla_hours INT,
    target_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tickets_customer_id ON tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

-- TICKET ASSIGNMENTS
CREATE TABLE IF NOT EXISTS ticket_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    engineer_id UUID REFERENCES users(id),
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_ticket_assignments_ticket ON ticket_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_assignments_engineer ON ticket_assignments(engineer_id);

-- TICKET STATUS HISTORY
CREATE TABLE IF NOT EXISTS ticket_status_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    old_status TEXT,
    new_status TEXT,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMPTZ
);

-- TICKET COMMENTS
CREATE TABLE IF NOT EXISTS ticket_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    user_id UUID REFERENCES users(id),
    comment TEXT,
    is_internal BOOLEAN,
    created_at TIMESTAMPTZ
);

-- TICKET ATTACHMENTS
CREATE TABLE IF NOT EXISTS ticket_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    file_url TEXT,
    file_type TEXT,
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ
);

-- TICKET FEEDBACK
CREATE TABLE IF NOT EXISTS ticket_feedbacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    engineer_id UUID REFERENCES users(id),
    rating INT,
    comment TEXT,
    created_at TIMESTAMPTZ
);

-- =========================================================
-- 7. OPERATIONS & LOGS
-- =========================================================

-- SERVICE VISITS
CREATE TABLE IF NOT EXISTS service_visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    engineer_id UUID REFERENCES users(id),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    notes TEXT
);

-- GPS LOGS
CREATE TABLE IF NOT EXISTS gps_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    engineer_id UUID REFERENCES users(id),
    latitude FLOAT,
    longitude FLOAT,
    logged_at TIMESTAMPTZ
);

-- DIGITAL SIGNATURES
CREATE TABLE IF NOT EXISTS digital_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    signed_by TEXT,
    file_url TEXT,
    signed_at TIMESTAMPTZ
);

-- AMC SCHEDULES
CREATE TABLE IF NOT EXISTS amc_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amc_id UUID REFERENCES amc_contracts(id),
    visit_date TIMESTAMPTZ,
    completed BOOLEAN,
    ticket_id UUID REFERENCES tickets(id)
);

-- AUDIT LOGS
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity TEXT,
    entity_id UUID,
    action TEXT,
    performed_by UUID NULL REFERENCES users(id),
    ip TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ
);

-- =========================================================
-- 8. ESCALATIONS
-- =========================================================

-- ESCALATION RULES
CREATE TABLE IF NOT EXISTS escalation_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    condition VARCHAR(50),
    after_mins INT,
    role TEXT,
    created_at TIMESTAMPTZ
);

-- TICKET ESCALATIONS
CREATE TABLE IF NOT EXISTS ticket_escalations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id),
    rule_id UUID REFERENCES escalation_rules(id),
    escalated_at TIMESTAMPTZ,
    resolved BOOLEAN
);

-- =========================================================
-- 9. SEED INITIAL DATA
-- =========================================================

-- Insert Admin User (Password: Admin@123)
INSERT INTO users (name, email, password, role, is_active, must_reset_password)
VALUES (
    'Admin User',
    'admin@gmail.com',
    '$2a$12$jPj7fs9fZFtH.8pxdH6yq.75.naM1nNl8X5Qe.XF8uzGYT01mP5sW', -- Bcrypt hash for Admin@123
    'admin',
    TRUE,
    FALSE
)
ON CONFLICT (email) DO NOTHING;

-- 10. DATA PATCHES & UPDATES
-- =========================================================

-- Ensure admin user has 2FA disabled
UPDATE users 
SET two_fa_enabled = false 
WHERE email = 'admin@gmail.com';

-- Update tickets with invalid status to 'Open'
DO $$
BEGIN
    -- Check if the constraint exists before attempting to update invalid data
    -- (This is a safety measure; usually, you'd update data BEFORE adding constraints)
    UPDATE tickets 
    SET status = 'Open' 
    WHERE status NOT IN ('Open', 'Assigned', 'In Progress', 'Closed');
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'Could not update ticket statuses: %', SQLERRM;
END $$;
