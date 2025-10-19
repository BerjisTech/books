-- Link purchases to Core billing intents
ALTER TABLE purchases ADD COLUMN IF NOT EXISTS payment_intent_id BIGINT;

