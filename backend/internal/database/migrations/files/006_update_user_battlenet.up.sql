-- D'abord, supprimer la contrainte unique sur battle_net_id
ALTER TABLE users 
    DROP CONSTRAINT IF EXISTS users_battle_net_id_key;

-- Ajouter les nouvelles colonnes
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS encrypted_refresh_token BYTEA,
    ADD COLUMN IF NOT EXISTS last_token_refresh TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- Renommer encrypted_token en encrypted_access_token (si la colonne existe)
DO $$ 
BEGIN 
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'encrypted_token'
    ) THEN
        ALTER TABLE users RENAME COLUMN encrypted_token TO encrypted_access_token;
    END IF;
END $$;

-- Convertir battle_net_id en VARCHAR tout en préservant les données existantes
ALTER TABLE users 
    ALTER COLUMN battle_net_id TYPE VARCHAR(255) USING 
        CASE 
            WHEN battle_net_id IS NOT NULL THEN battle_net_id::VARCHAR(255)
            ELSE NULL
        END;

-- Créer le nouvel index unique qui permet les valeurs NULL
CREATE UNIQUE INDEX IF NOT EXISTS users_battle_net_id_key 
    ON users(battle_net_id) 
    WHERE battle_net_id IS NOT NULL;

-- Nettoyer les anciennes colonnes si elles existent
DO $$ 
BEGIN 
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'battle_net_refresh_token'
    ) THEN
        ALTER TABLE users DROP COLUMN battle_net_refresh_token;
    END IF;
END $$;