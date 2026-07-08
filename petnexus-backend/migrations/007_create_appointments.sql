CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_profile_id UUID NOT NULL,
    clinic_profile_id UUID NOT NULL,
    pet_id UUID NOT NULL,
    title VARCHAR(150),
    appointment_type VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    note TEXT,
    created_by_user_id UUID,
    created_by_role VARCHAR(20) NOT NULL,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_appointments_owner_profile_id ON appointments(owner_profile_id);
CREATE INDEX IF NOT EXISTS idx_appointments_clinic_profile_id ON appointments(clinic_profile_id);
CREATE INDEX IF NOT EXISTS idx_appointments_pet_id ON appointments(pet_id);
CREATE INDEX IF NOT EXISTS idx_appointments_scheduled_at ON appointments(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_clinic_scheduled_at ON appointments(clinic_profile_id, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_appointments_owner_scheduled_at ON appointments(owner_profile_id, scheduled_at);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_appointments_owner_profile' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT fk_appointments_owner_profile
        FOREIGN KEY (owner_profile_id) REFERENCES owner_profiles(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_appointments_clinic_profile' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT fk_appointments_clinic_profile
        FOREIGN KEY (clinic_profile_id) REFERENCES clinic_profiles(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_appointments_pet' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT fk_appointments_pet
        FOREIGN KEY (pet_id) REFERENCES pets(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_appointments_created_by_user' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT fk_appointments_created_by_user
        FOREIGN KEY (created_by_user_id) REFERENCES users(id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_appointments_type' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT chk_appointments_type CHECK (
            appointment_type IN ('checkup', 'vaccination', 'consultation', 'follow_up', 'grooming', 'emergency', 'other')
        );
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_appointments_status' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT chk_appointments_status CHECK (
            status IN ('requested', 'scheduled', 'checked_in', 'completed', 'cancelled')
        );
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_appointments_created_by_role' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT chk_appointments_created_by_role
        CHECK (created_by_role IN ('owner', 'clinic'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_appointments_duration' AND conrelid = 'appointments'::regclass) THEN
        ALTER TABLE appointments ADD CONSTRAINT chk_appointments_duration
        CHECK (duration_minutes BETWEEN 5 AND 480);
    END IF;
END
$$;
