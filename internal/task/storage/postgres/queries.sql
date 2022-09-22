--name: CreateTask :one
INSERT INTO tasks(creator_id, list_id, title, body)
VALUES (pggen.arg('CreatorID'),
        pggen.arg('ListID'),
        pggen.arg('Title'),
        pggen.arg('Body'))
RETURNING id, created_at;

--name: CreateTaskInDefaultList :one
INSERT INTO tasks(creator_id, title, body)
VALUES (pggen.arg('CreatorID'),
        pggen.arg('Title'),
        pggen.arg('Body'))
RETURNING id, created_at;


--name: FindTaskByID :one
SELECT id,
       creator_id,
       list_id,
       title,
       body,
       done,
       created_at,
       updated_at
FROM tasks
WHERE id = pggen.arg('ListID')::uuid
LIMIT 1;



--name: FindTaskInMainList :many
SELECT id,
       title,
       body,
       done,
       created_at,
       updated_at
FROM tasks
WHERE creator_id = pggen.arg('ID')::uuid
  AND list_id IS NULL;


--name: FindBasicTaskInList :many
SELECT id,
       title,
       body,
       done,
       created_at,
       updated_at
FROM tasks
WHERE list_id = pggen.arg('ListID')::uuid;

--name: UpdateTaskTitle :exec
UPDATE tasks
SET title = pggen.arg('Title')
WHERE id = pggen.arg('ID')::uuid;

--name: UpdateTaskBody :exec
UPDATE tasks
SET body = pggen.arg('Body')
WHERE id = pggen.arg('ID')::uuid;

--name: UpdateTaskStatus :exec
UPDATE tasks
SET done = pggen.arg('Status')
WHERE id = pggen.arg('ID')::uuid;


--name: ChangeTaskStatus :one
UPDATE tasks
SET done = not done
WHERE id = pggen.arg('ID')::uuid
RETURNING done;

--name: DeleteTask :exec
DELETE
FROM tasks
WHERE id = pggen.arg('ID')::uuid;
