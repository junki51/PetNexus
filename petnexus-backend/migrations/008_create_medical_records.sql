CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_profile_id UUID NOT NULL,
    pet_id UUID NOT NULL,
    appointment_id UUID,
    created_by_user_id UUID NOT NULL,
    visit_at TIMESTAMPTZ NOT NULL,
    chief_complaint TEXT NOT NULL,
    clinical_findings TEXT,
    diagnosis TEXT,
    treatment_plan TEXT,
    medications TEXT,
    follow_up_instructions TEXT,
    next_follow_up_at TIMESTAMPTZ,
    weight_kg NUMERIC(6,2),
    temperature_c NUMERIC(5,2),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_profile_id
ON medical_records(clinic_profile_id);

CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id
ON medical_records(pet_id);

CREATE INDEX IF NOT EXISTS idx_medical_records_clinic_visit_at
ON medical_records(clinic_profile_id, visit_at DESC);

CREATE INDEX IF NOT EXISTS idx_medical_records_pet_visit_at
ON medical_records(pet_id, visit_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_medical_records_appointment_id_unique
ON medical_records(appointment_id)
WHERE appointment_id IS NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_clinic_profile'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_clinic_profile
        FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_pet'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_pet
        FOREIGN KEY (pet_id) REFERENCES pets(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_appointment'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_appointment
        FOREIGN KEY (appointment_id) REFERENCES appointments(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_medical_records_created_by_user'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT fk_medical_records_created_by_user
        FOREIGN KEY (created_by_user_id) REFERENCES users(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_weight_positive'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_weight_positive
        CHECK (weight_kg IS NULL OR weight_kg > 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_temperature_positive'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_temperature_positive
        CHECK (temperature_c IS NULL OR temperature_c > 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_medical_records_next_follow_up_after_visit'
          AND conrelid = 'medical_records'::regclass
    ) THEN
        ALTER TABLE medical_records
        ADD CONSTRAINT chk_medical_records_next_follow_up_after_visit
        CHECK (next_follow_up_at IS NULL OR next_follow_up_at >= visit_at);
    END IF;
END
$$;
