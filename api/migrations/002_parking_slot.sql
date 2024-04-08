-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS parking_slot
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lot_id          UUID REFERENCES parking_lot(id) ON DELETE CASCADE,
    slot_number     INT CHECK (slot_number >= 0),
    status          VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    created_at      TIMESTAMP DEFAULT current_timestamp,
    updated_at      TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS parking_slot;
