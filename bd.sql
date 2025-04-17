create schema if not exists schema_call

create table if not exists schema_call.phone_call
(
    call_id   text   not null
        constraint phone_call_pk
            primary key,
    user_id   text   not null,
    leader_id text   not null,
    title     text      not null,
    start_time timestamp not null,
    status    text      not null,
    feedback  text      not null
);
