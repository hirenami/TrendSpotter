-- name: saveTrend :exec
INSERT INTO trends (trends_name,trends_location,trends_rank) VALUES (?, ?, ?);

-- name: deleteTrend :exec
DELETE FROM trends;