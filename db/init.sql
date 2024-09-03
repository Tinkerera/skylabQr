CREATE TABLE url_mappings (
                              id SERIAL PRIMARY KEY,
                              short_url VARCHAR(255) NOT NULL UNIQUE,
                              original_url TEXT NOT NULL,
                              expiration_date TIMESTAMP
);
