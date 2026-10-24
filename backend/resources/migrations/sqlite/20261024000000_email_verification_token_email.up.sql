PRAGMA foreign_keys= OFF;
BEGIN;

-- Existing tokens are not bound to an email address, so they are invalidated
DELETE FROM email_verification_tokens;

ALTER TABLE email_verification_tokens
    ADD COLUMN email TEXT NOT NULL DEFAULT '';

COMMIT;
PRAGMA foreign_keys= ON;
