CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS loans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    loan_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    return_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL DEFAULT 'BORROWED', -- Status of the loan (e.g., BORROWED, RETURNED, LOST)
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_user_id ON loans (user_id);
CREATE INDEX idx_loan_book_id ON loans (book_id);
CREATE INDEX idx_loan_status ON loans (status);

-- Membuat fungsi trigger untuk memperbarui kolom updated_at dan menaikkan version
CREATE OR REPLACE FUNCTION update_updated_at_and_version_loans()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   NEW.version = OLD.version + 1;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Membuat trigger yang memanggil fungsi di atas sebelum pembaruan baris
CREATE TRIGGER set_updated_at_and_version_loans
BEFORE UPDATE ON loans
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_and_version_loans();
