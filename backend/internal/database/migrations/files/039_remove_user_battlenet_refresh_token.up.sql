-- Migration pour supprimer la colonne encrypted_refresh_token
-- Battle.net ne fournit pas de refresh tokens

-- Supprimer la colonne encrypted_refresh_token
ALTER TABLE users DROP COLUMN IF EXISTS encrypted_refresh_token;
