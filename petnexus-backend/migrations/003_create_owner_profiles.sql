CREATE TABLE IF NOT EXISTS owner_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    gender VARCHAR(30),
    date_of_birth DATE,
    phone_number VARCHAR(30) NOT NULL,
    avatar_url TEXT,
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    province VARCHAR(100),
    district VARCHAR(100),
    subdistrict VARCHAR(100),
    postal_code VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_owner_profiles_user_id_unique
ON owner_profiles(user_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_owner_profiles_user'
          AND conrelid = 'owner_profiles'::regclass
    ) THEN
        ALTER TABLE owner_profiles
        ADD CONSTRAINT fk_owner_profiles_user
        FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END
$$;
