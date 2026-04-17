-- name: CreateFirm :one
INSERT INTO firms (name, slug, billing_contact_email, country, staff_count_range, primary_audit_types)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFirmByID :one
SELECT * FROM firms WHERE id = $1;

-- name: UpdateFirm :one
UPDATE firms
SET name = COALESCE(sqlc.narg('name'), name),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    timezone = COALESCE(sqlc.narg('timezone'), timezone),
    billing_contact_email = COALESCE(sqlc.narg('billing_contact_email'), billing_contact_email),
    settings = COALESCE(sqlc.narg('settings'), settings)
WHERE id = $1
RETURNING *;
