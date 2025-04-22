create table schema_call.table_name
(
    call_id    text
        constraint table_name_pk
            primary key,
    user_id    text,
    leader_id  text,
    title      text,
    status     text,
    feedback   text,
    start_time time
);
