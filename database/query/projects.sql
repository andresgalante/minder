-- name: CreateProject :one
INSERT INTO projects (
    name,
    parent_id,
    metadata
) VALUES (
    $1, $2, sqlc.arg(metadata)::jsonb
) RETURNING *;

-- name: GetRootProjects :many
SELECT * FROM projects
WHERE parent_id IS NULL;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1 AND is_organization = FALSE LIMIT 1;

-- name: GetProjectByName :one
SELECT * FROM projects
WHERE name = $1 AND is_organization = FALSE LIMIT 1;

-- name: GetParentProjects :many
WITH RECURSIVE get_parents AS (
    SELECT id, parent_id, created_at FROM projects 
    WHERE projects.id = $1

    UNION

    (
        SELECT p.id, p.parent_id, p.created_at FROM projects p
        INNER JOIN get_parents gp ON p.id = gp.parent_id
        ORDER BY p.created_at ASC
    )
)
SELECT id FROM get_parents;

-- name: GetParentProjectsUntil :many
WITH RECURSIVE get_parents_until AS (
    SELECT id, parent_id, created_at FROM projects 
    WHERE projects.id = $1

    UNION

    (
        SELECT p.id, p.parent_id, p.created_at FROM projects p
        INNER JOIN get_parents_until gpu ON p.id = gpu.parent_id
        WHERE p.id != $2
        ORDER BY p.created_at ASC
    )
)
SELECT id FROM get_parents_until;

-- name: GetChildrenProjects :many
WITH RECURSIVE get_children AS (
    SELECT projects.id, projects.name, projects.metadata, projects.parent_id, projects.created_at, projects.updated_at FROM projects 
    WHERE projects.id = $1

    UNION

    (
        SELECT p.id, p.name, p.metadata, p.parent_id, p.created_at, p.updated_at FROM projects p
        INNER JOIN get_children gc ON p.parent_id = gc.id
        ORDER BY p.created_at ASC
    )
)
SELECT * FROM get_children;


-- name: DeleteProject :many
WITH RECURSIVE get_children AS (
    SELECT id, parent_id FROM projects
    WHERE projects.id = $1 AND projects.parent_id IS NOT NULL

    UNION

    SELECT p.id, p.parent_id FROM projects p
    INNER JOIN get_children gc ON p.parent_id = gc.id
)
DELETE FROM projects
WHERE id IN (SELECT id FROM get_children)
RETURNING id, name, metadata, created_at, updated_at, parent_id;