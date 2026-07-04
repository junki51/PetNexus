CREATE TABLE IF NOT EXISTS breeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    species VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    name_th VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_breeds_species_name_unique
ON breeds(species, name);
CREATE INDEX IF NOT EXISTS idx_breeds_species ON breeds(species);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_breeds_species'
          AND conrelid = 'breeds'::regclass
    ) THEN
        ALTER TABLE breeds
        ADD CONSTRAINT chk_breeds_species CHECK (species IN ('dog', 'cat'));
    END IF;
END
$$;

INSERT INTO breeds (species, name, name_th) VALUES
    ('dog', 'Golden Retriever', NULL),
    ('dog', 'Labrador Retriever', NULL),
    ('dog', 'Poodle', NULL),
    ('dog', 'Shiba Inu', NULL),
    ('dog', 'Siberian Husky', NULL),
    ('dog', 'Chihuahua', NULL),
    ('dog', 'Pomeranian', NULL),
    ('dog', 'Thai Bangkaew', NULL),
    ('cat', 'Persian', NULL),
    ('cat', 'Scottish Fold', NULL),
    ('cat', 'British Shorthair', NULL),
    ('cat', 'Siamese', NULL),
    ('cat', 'Maine Coon', NULL),
    ('cat', 'Ragdoll', NULL),
    ('cat', 'Sphynx', NULL),
    ('cat', 'Domestic Shorthair', NULL)
ON CONFLICT (species, name) DO NOTHING;

CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_profile_id UUID NOT NULL,
    breed_id UUID,
    species VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(30),
    date_of_birth DATE,
    weight_kg NUMERIC(6,2),
    microchip_id VARCHAR(100),
    avatar_url TEXT,
    color VARCHAR(100),
    distinctive_marks TEXT,
    is_neutered BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pets_owner_profile_id ON pets(owner_profile_id);
CREATE INDEX IF NOT EXISTS idx_pets_breed_id ON pets(breed_id);
CREATE INDEX IF NOT EXISTS idx_pets_species ON pets(species);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_pets_owner_profile'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT fk_pets_owner_profile
        FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_pets_breed'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT fk_pets_breed
        FOREIGN KEY (breed_id) REFERENCES breeds(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_pets_species'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT chk_pets_species CHECK (species IN ('dog', 'cat'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_pets_gender'
          AND conrelid = 'pets'::regclass
    ) THEN
        ALTER TABLE pets
        ADD CONSTRAINT chk_pets_gender
        CHECK (gender IS NULL OR gender IN ('male', 'female', 'unknown'));
    END IF;
END
$$;
