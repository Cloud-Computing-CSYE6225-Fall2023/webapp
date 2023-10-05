CREATE TABLE IF NOT EXISTS assignments (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    name VARCHAR(50),
    points integer,
    num_of_attemps integer,
    deadline timestamp,
    assignment_created TIMESTAMP DEFAULT current_timestamp,
    assignment_updated TIMESTAMP DEFAULT current_timestamp
);
