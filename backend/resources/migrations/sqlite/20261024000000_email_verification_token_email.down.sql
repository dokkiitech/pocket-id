PRAGMA foreign_keys= OFF;
BEGIN;

ALTER TABLE email_verification_tokens
    DROP COLUMN email;

COMMIT;
PRAGMA foreign_keys= ON;
