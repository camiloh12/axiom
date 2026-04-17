-- name: CreateClient :one
INSERT INTO clients (firm_id, name, industry, primary_contact_email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients WHERE id = $1;

-- name: ListClientsByFirmID :many
SELECT * FROM clients
WHERE firm_id = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateClient :one
UPDATE clients
SET name = COALESCE(sqlc.narg('name'), name),
    industry = COALESCE(sqlc.narg('industry'), industry),
    primary_contact_email = COALESCE(sqlc.narg('primary_contact_email'), primary_contact_email)
WHERE id = $1
RETURNING *;
