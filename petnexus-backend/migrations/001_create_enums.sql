CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('owner', 'clinic', 'clinic_staff', 'admin');
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            WHERE t.typname = 'user_role'
              AND e.enumlabel = 'owner'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'owner';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            WHERE t.typname = 'user_role'
              AND e.enumlabel = 'clinic'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'clinic';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            WHERE t.typname = 'user_role'
              AND e.enumlabel = 'clinic_staff'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'clinic_staff';
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_enum e
            JOIN pg_type t ON t.oid = e.enumtypid
            WHERE t.typname = 'user_role'
              AND e.enumlabel = 'admin'
        ) THEN
            ALTER TYPE user_role ADD VALUE 'admin';
        END IF;
    END IF;
END
$$;
