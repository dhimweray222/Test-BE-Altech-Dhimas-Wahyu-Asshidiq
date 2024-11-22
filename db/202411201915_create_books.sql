CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- UUID as primary key with auto-generation
    title VARCHAR(255) NOT NULL,                     -- Title of the book
    description TEXT,                                -- Description of the book
    publish_date DATE,                               -- Publication date of the book
    author_id UUID,                                  -- Foreign key referencing authors (UUID)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Automatically sets the current timestamp
    CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES authors(id)  -- Foreign key constraint
);