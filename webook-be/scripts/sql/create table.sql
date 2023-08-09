create table webook.users
(
    id       bigint auto_increment
        primary key,
    email    varchar(191) null,
    password longtext     null,
    ctime    bigint       null,
    utime    bigint       null,
    constraint email
        unique (email)
);

