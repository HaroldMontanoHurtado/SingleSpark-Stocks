CREATE DATABASE IF NOT EXISTS singlespark;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
