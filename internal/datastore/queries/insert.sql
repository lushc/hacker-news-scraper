INSERT INTO top_stories (id, type, title, content, url, score, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO NOTHING