create schema if not exists schema_call;

create table if not exists schema_call.phone_call
(
    call_id    text not null
        constraint phone_call_pk
            primary key,
    user_id    text,
    leader_id  text,
    title      text,
    status     text,
    feedback   text,
    start_time text
);

insert into schema_call.phone_call (call_id, user_id, leader_id, title, status, feedback, start_time)
values ('1111','user1','leader1', 'title1','status1', 'feedback1','01:01:01'),
       ('2222','user2','leader2', 'title2','status2', 'feedback2','02:02:02'),
       ('3333','user3','leader3', 'title3','status3', 'feedback3','03:03:03'),
       ('4444','user4','leader3', 'title4','status4', 'feedback4','04:04:04');



