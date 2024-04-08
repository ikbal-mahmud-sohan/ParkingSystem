-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS vehicle
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_no      VARCHAR(20) NOT NULL,
    entry           TIMESTAMP DEFAULT current_timestamp,
    exit            TIMESTAMP,
    parked_hours    INT NOT NULL DEFAULT 0,
    fee             INT NOT NULL DEFAULT 0,
    slot_id         UUID REFERENCES parking_slot(id) ON DELETE SET NULL,
    created_at      TIMESTAMP DEFAULT current_timestamp,
    updated_at      TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS vehicle;
