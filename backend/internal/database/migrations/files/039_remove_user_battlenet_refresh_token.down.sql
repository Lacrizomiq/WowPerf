-- Migration pour ajouter la colonne encrypted_refresh_token
-- Battle.net ne fournit pas de refresh tokens

-- Ajouter la colonne encrypted_refresh_token
ALTER TABLE users ADD COLUMN IF NOT EXISTS encrypted_refresh_token BYTEA;