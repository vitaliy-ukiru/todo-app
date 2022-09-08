create table if not exists "users"
(
    id         uuid        default gen_random_uuid() not null
        constraint user_pk
            primary key,
    email      text                                  not null,
    username   text                                  not null,
    password   text                                  not null,
    created_at timestamptz default now()             not null
);

create unique index if not exists users_email_uindex
    on "users" (email);

create table if not exists task_lists
(
    id         uuid default gen_random_uuid() not null
        constraint task_lists_pk
            primary key,
    creator_id uuid                           not null
        constraint task_lists_users_id_fk
            references users,
    title      text                           not null
);

create table if not exists tasks
(
    id         uuid        default gen_random_uuid() not null
        constraint tasks_pk
            primary key,
    creator_id uuid                                  not null
        constraint tasks_users_id_fk
            references users,
    list_id    uuid
        constraint tasks_task_list_id_fk
            references task_lists,
    title      text                                  not null,
    body       text                                  not null,
    done       bool        default false             not null,
    created_at timestamptz default now()             not null,
    updated_at timestamptz
);




