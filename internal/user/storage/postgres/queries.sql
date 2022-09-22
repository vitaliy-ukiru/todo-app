--name: CreateUser :one
INSERT INTO users(email, username, password)
VALUES (pggen.arg('Email'),
        pggen.arg('Username'),
        pggen.arg('Password'))
RETURNING id, created_at;

--name: FindUserByID :one
SELECT id, email, username, created_at
FROM users
WHERE id = pggen.arg('ID')::uuid;

--name: FindUserByEmail :one
SELECT id, email, username, password, created_at
FROM users
WHERE email = pggen.arg('Email');

--name: UpdateUserPassword :exec
UPDATE users
SET password = pggen.arg('Password')
WHERE id = pggen.arg('ID')::uuid;


--name: DeleteUser :exec
DELETE
FROM users
WHERE id = pggen.arg('ID')::uuid;
