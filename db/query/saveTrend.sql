-- name: SaveTrend :exec
INSERT INTO trends (trends_name,trends_location,trends_rank,trends_endtimestamp,trends_increase_percentage) VALUES (?, ?, ?, ?, ?);

-- name: DeleteTrend :exec
DELETE FROM trends;