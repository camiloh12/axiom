-- name: CreateInvitation :one
INSERT INTO invitations (firm_id, email, assigned_role, token_hash, expires_at, invited_by_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInvitationByTokenHash :one
SELECT * FROM invitations WHERE token_hash = $1;

-- name: ListInvitationsByFirmID :many
SELECT * FROM invitations
WHERE firm_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: AcceptInvitation :one
UPDATE invitations
SET status = 'Accepted', accepted_at = now()
WHERE id = $1 AND status = 'Sent'
RETURNING *;

-- name: CancelInvitation :exec
UPDATE invitations SET status = 'Expired' WHERE id = $1 AND status = 'Sent';
