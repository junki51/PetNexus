ALTER TABLE pets ADD COLUMN IF NOT EXISTS public_pet_id VARCHAR(50);

UPDATE pets SET public_pet_id = NULL
WHERE public_pet_id IS NOT NULL AND BTRIM(public_pet_id) = '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_pets_public_pet_id_unique
ON pets(public_pet_id);

DO $$
DECLARE
    pet_record RECORD;
    candidate VARCHAR(50);
BEGIN
    FOR pet_record IN
        SELECT id FROM pets
        WHERE public_pet_id IS NULL
    LOOP
        LOOP
            candidate := 'PNX-PET-' || UPPER(SUBSTRING(ENCODE(gen_random_bytes(6), 'hex'), 1, 6));
            BEGIN
                UPDATE pets SET public_pet_id = candidate WHERE id = pet_record.id;
                EXIT;
            EXCEPTION WHEN unique_violation THEN
                -- Generate another candidate and retry this pet.
            END;
        END LOOP;
    END LOOP;
END
$$;

ALTER TABLE pets ALTER COLUMN public_pet_id SET NOT NULL;
