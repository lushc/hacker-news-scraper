CREATE TABLE IF NOT EXISTS top_stories (
    id INT PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    url TEXT NOT NULL,
    score INT DEFAULT 0,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);