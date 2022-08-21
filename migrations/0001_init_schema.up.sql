BEGIN;

CREATE TABLE IF NOT EXISTS public.todo (
    id CHAR(36) NOT NULL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT(500)
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

COMMIT;
