CREATE TABLE IF NOT EXISTS authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- UUID as primary key with auto-generation
    name VARCHAR(255) NOT NULL UNIQUE,                      -- Author's name
    bio TEXT,                                        -- Author's biography
    birth_date DATE,                                 -- Author's birth date
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Automatically sets the current timestamp
);