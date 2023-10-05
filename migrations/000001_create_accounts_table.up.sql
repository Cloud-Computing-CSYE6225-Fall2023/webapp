CREATE TABLE IF NOT EXISTS accounts (
  id VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255),
  account_created TIMESTAMP DEFAULT current_timestamp,
  account_updated TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_password ON accounts (password);
