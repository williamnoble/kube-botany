-- name: GetPlant :one
SELECT *
FROM plants
WHERE id = $1 LIMIT 1;

-- name: Listplants :many
SELECT *
FROM plants
ORDER BY name;

-- name: CreatePlant :one
INSERT INTO plants (name,
                    can_die,
                    water_consumption_rate,
                    minimum_water_level,
                    water_level,
                    last_watered,
                    growth_rate,
                    growth,
                    growth_stage,
                    last_updated,
                    backdrop,
                    mascot)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

-- name: UpdatePlant :exec
UPDATE plants
set name = $2
WHERE id = $1;

-- name: DeletePlant :exec
DELETE
FROM plants
WHERE id = $1;