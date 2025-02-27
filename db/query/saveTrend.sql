-- name: SaveTrend :exec
INSERT INTO trends (trends_name,trends_location,trends_rank) VALUES (?, ?, ?);

-- name: DeleteTrend :exec
DELETE FROM trends;