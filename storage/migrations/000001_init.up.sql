create schema if not exists schema_call;

create table schema_call.phone_call
(
    call_id    text
        constraint phone_call_pk
            primary key,
    user_id    text,
    leader_id  text,
    title      text,
    status     text,
    feedback   text,
    start_time time
);


