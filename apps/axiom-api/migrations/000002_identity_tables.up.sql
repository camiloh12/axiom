-- RLS helper: function to get current firm ID from session variable
CREATE OR REPLACE FUNCTION current_firm_id() RETURNS uuid AS $$
  SELECT current_setting('app.current_firm_id', true)::uuid;
$$ LANGUAGE sql STABLE;

-- Firms
CREATE TABLE firms (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name           text NOT NULL,
  slug           text NOT NULL UNIQUE,
  logo_url       text,
  timezone       text NOT NULL DEFAULT 'America/New_York',
  billing_contact_email text NOT NULL,
  subscription_tier text NOT NULL DEFAULT 'Growth'
    CHECK (subscription_tier IN ('Growth', 'Scale', 'Enterprise')),
  country        text NOT NULL DEFAULT 'US'
    CHECK (country IN ('US', 'CA')),
  staff_count_range text,
  primary_audit_types jsonb DEFAULT '[]',
  settings       jsonb NOT NULL DEFAULT '{}',
  created_at     timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE firms ENABLE ROW LEVEL SECURITY;
CREATE POLICY firms_isolation ON firms
  USING (id = current_firm_id());

-- Users
CREATE TABLE users (
  id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id                uuid REFERENCES firms(id),
  client_id              uuid,
  email                  text NOT NULL UNIQUE,
  display_name           text NOT NULL,
  role                   user_role NOT NULL,
  auth_method            auth_method NOT NULL DEFAULT 'Password',
  password_hash          text,
  notification_frequency notification_frequency NOT NULL DEFAULT 'Daily',
  tour_completed         boolean NOT NULL DEFAULT false,
  is_active              boolean NOT NULL DEFAULT true,
  created_at             timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT users_firm_xor_client CHECK (
    (firm_id IS NOT NULL AND client_id IS NULL) OR
    (firm_id IS NULL AND client_id IS NOT NULL)
  ),
  CONSTRAINT users_client_role CHECK (
    (role IN ('ClientAdmin', 'ClientUser') AND client_id IS NOT NULL) OR
    (role NOT IN ('ClientAdmin', 'ClientUser') AND firm_id IS NOT NULL)
  )
);

CREATE INDEX idx_users_firm_id ON users(firm_id);
CREATE INDEX idx_users_client_id ON users(client_id);
CREATE INDEX idx_users_email ON users(email);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_isolation ON users
  USING (firm_id = current_firm_id() OR firm_id IS NULL);

-- Clients
CREATE TABLE clients (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id               uuid NOT NULL REFERENCES firms(id),
  name                  text NOT NULL,
  industry              text,
  primary_contact_email text,
  created_at            timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_clients_firm_id ON clients(firm_id);

ALTER TABLE clients ENABLE ROW LEVEL SECURITY;
CREATE POLICY clients_isolation ON clients
  USING (firm_id = current_firm_id());

-- Add FK from users.client_id to clients
ALTER TABLE users ADD CONSTRAINT fk_users_client
  FOREIGN KEY (client_id) REFERENCES clients(id);

-- Invitations
CREATE TABLE invitations (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  firm_id         uuid NOT NULL REFERENCES firms(id),
  email           text NOT NULL,
  assigned_role   user_role NOT NULL,
  token_hash      text NOT NULL UNIQUE,
  status          invitation_status NOT NULL DEFAULT 'Sent',
  expires_at      timestamptz NOT NULL,
  reminder_sent_at timestamptz,
  invited_by_id   uuid NOT NULL REFERENCES users(id),
  accepted_at     timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_invitations_firm_id ON invitations(firm_id);
CREATE INDEX idx_invitations_token_hash ON invitations(token_hash);

ALTER TABLE invitations ENABLE ROW LEVEL SECURITY;
CREATE POLICY invitations_isolation ON invitations
  USING (firm_id = current_firm_id());
