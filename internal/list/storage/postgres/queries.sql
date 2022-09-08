--name: CreateList :one
INSERT INTO task_lists(creator_id, title)
VALUES (pggen.arg('CreatorID'), pggen.arg('Title'))
RETURNING id;


--name: FindListByID :one
SELECT id, creator_id, title
FROM task_lists
WHERE id = pggen.arg('ListID');

--name: FindUserLists :many
SELECT id, creator_id, title
FROM task_lists
WHERE creator_id = pggen.arg('UserID');


--name: UpdateListTitle :exec
UPDATE task_lists
SET title=pggen.arg('NewTitle')
WHERE id = pggen.arg('ListID');


--name: DeleteList :exec
DELETE FROM task_lists WHERE id=pggen.arg('ListID');
