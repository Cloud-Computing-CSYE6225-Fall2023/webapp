CREATE TABLE IF NOT EXISTS submissions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    assignment_id VARCHAR(255),
    submission_url VARCHAR(255),
    submission_created TIMESTAMP DEFAULT current_timestamp,
    submission_updated TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_user_assignment_id ON submissions (user_id, assignment_id);