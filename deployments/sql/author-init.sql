CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS authors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    biography VARCHAR(255) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Membuat fungsi trigger untuk memperbarui kolom updated_at dan menaikkan version
CREATE OR REPLACE FUNCTION update_updated_at_and_version_authors()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   NEW.version = OLD.version + 1;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Membuat trigger yang memanggil fungsi di atas sebelum pembaruan baris
CREATE TRIGGER set_updated_at_and_version_authors
BEFORE UPDATE ON authors
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_and_version_authors();
