CREATE TYPE user_role AS ENUM (
  'FirmAdmin', 'Partner', 'Manager', 'Staff', 'EQReviewer',
  'ClientAdmin', 'ClientUser', 'ViewOnly'
);

CREATE TYPE auth_method AS ENUM ('Password', 'OAuth', 'SAML');

CREATE TYPE notification_frequency AS ENUM ('RealTime', 'Daily', 'Weekly');

CREATE TYPE invitation_status AS ENUM ('Sent', 'Accepted', 'Expired');
