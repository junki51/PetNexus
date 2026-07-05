CREATE TABLE IF NOT EXISTS clinic_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    clinic_name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(30),
    email VARCHAR(255),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_clinic_profiles_user_id_unique
ON clinic_profiles(user_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_clinic_profiles_user'
          AND conrelid = 'clinic_profiles'::regclass
    ) THEN
        ALTER TABLE clinic_profiles
        ADD CONSTRAINT fk_clinic_profiles_user
        FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END
$$;
