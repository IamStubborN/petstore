-- +migrate Up
-- +migrate StatementBegin
set search_path=public;
create schema if not exists public;
alter schema public owner to petstore;

create table if not exists category
(
    id bigserial not null
        constraint category_pk
            primary key,
    name varchar(50) not null,
    price numeric(15,2)       not null
);

alter table category owner to petstore;

create table if not exists tag
(
    id bigserial not null
        constraint tag_pk
            primary key,
    name varchar(50) not null
);

alter table tag owner to petstore;

create table if not exists pet_status
(
    id bigserial not null
        constraint pet_status_pk
            primary key,
    name varchar(50) not null
);

alter table pet_status owner to petstore;

create table if not exists pet
(
    id bigserial not null
        constraint pet_pk
            primary key,
    category_id bigint not null
        constraint category___fk
            references category
            on update set default on delete set default,
    name varchar(50) not null,
    photo_urls text[],
    pet_status_id bigint not null
        constraint pet_status_id__fk
            references pet_status
);

alter table pet owner to petstore;

create table if not exists user_status
(
    id bigserial not null
        constraint user_status_pk
            primary key,
    name varchar(20) not null,
    allowed_methods text[] default '{GET}'::text[]
);

alter table user_status owner to petstore;

create table if not exists "user"
(
    id bigserial not null
        constraint user_pk
            primary key,
    user_name varchar(50) not null,
    first_name varchar(50),
    last_name varchar(50),
    email varchar(50) not null,
    password varchar(100) not null,
    phone varchar(20) not null,
    user_status_id bigint not null
        constraint user_status___fk
            references user_status
);

alter table "user" owner to petstore;

create unique index if not exists user_user_name_uindex
    on "user" (user_name);

create table if not exists order_status
(
    id bigserial not null
        constraint order_status_pk
            primary key,
    name varchar(50) not null
);

alter table order_status owner to petstore;

create table if not exists "order"
(
    id              bigserial             not null
        constraint order_pk
            primary key,
    pet_id          bigint                not null
        constraint pet_id___fk
            references pet
            on update cascade on delete cascade,
    quantity        integer               not null,
    ship_date       timestamp             not null,
    order_status_id integer               not null
        constraint order_status___fk
            references order_status
            on update cascade on delete cascade,
    complete        boolean default false not null,
    user_id         bigint                not null
        constraint user_id___fk
            references "user"
            on update cascade on delete cascade
);

alter table "order"
    owner to petstore;

create table if not exists pet_tag
(
    pet_id bigint not null
        constraint pet_id_fk
            references pet
            on update cascade on delete cascade,
    tag_id bigint not null
        constraint tag___fk
            references tag,
    constraint pet_tag_pk
        primary key (pet_id, tag_id)
);

alter table pet_tag owner to petstore;

create or replace view user_info(id, user_name, first_name, last_name, email, password, phone, user_status, allowed_methods) as
SELECT u.id,
       u.user_name,
       u.first_name,
       u.last_name,
       u.email,
       u.password,
       u.phone,
       us.name AS user_status,
       us.allowed_methods
FROM ("user" u
         JOIN user_status us ON ((u.user_status_id = us.id)))
group by u.id, u.user_name, u.first_name, u.last_name, u.email, u.password, u.phone, us.name, us.allowed_methods;

alter table user_info owner to petstore;

create or replace view pet_info(id, category_id, category_name, price, name, photo_urls, pet_status_id,
                     pet_status_name) as
SELECT p.id,
       p.category_id,
       ct.name           AS category_name,
       ct.price,
       p.name,
       p.photo_urls,
       ps.id             AS pet_status_id,
       ps.name           AS pet_status_name
FROM (pet p
         JOIN category ct ON ((p.category_id = ct.id))
         JOIN pet_status ps ON ((p.pet_status_id = ps.id)))
GROUP BY p.id, p.category_id, ct.name, ct.price, p.name, p.photo_urls, ps.id, ps.name;

alter table pet_info owner to petstore;

create or replace view order_info(id, user_id, pet_id, pet_name, quantity, pet_status, ship_date, order_status, complete) as
SELECT o.id,
       o.user_id,
       o.pet_id,
       p.name  AS pet_name,
       o.quantity,
       ps.name AS pet_status,
       o.ship_date,
       os.name AS order_status,
       o.complete
FROM (((("order" o
    JOIN pet p ON ((o.pet_id = p.id)))
    JOIN pet_status ps ON ((p.pet_status_id = ps.id)))
    JOIN "user" u ON ((o.user_id = u.id)))
    JOIN order_status os ON ((o.order_status_id = os.id)))
GROUP BY o.id, u.user_name, p.name, o.quantity, ps.name, os.name;

alter table order_info owner to petstore;

create view invoice_info(id, user_name, pet, category, ship_date, quantity, price) as
SELECT o.id,
       u.user_name,
       p.name AS pet,
       c.name AS category,
       o.ship_date,
       o.quantity,
       c.price
FROM ((("order" o
    JOIN "user" u ON ((o.user_id = u.id)))
    JOIN pet p ON ((o.pet_id = p.id)))
    JOIN category c ON ((p.category_id = c.id)))
GROUP BY o.id, u.user_name, p.name, c.name, o.ship_date, o.quantity, c.price;

alter table invoice_info owner to petstore;

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
drop view if exists invoice_info cascade;
drop view if exists pet_info cascade;
drop view if exists order_info cascade;
drop view if exists user_info cascade;
drop table IF EXISTS category cascade;
drop table IF EXISTS "order" cascade;
drop table IF EXISTS order_status cascade;
drop table IF EXISTS pet cascade;
drop table IF EXISTS pet_status cascade;
drop table IF EXISTS pet_tag cascade;
drop table IF EXISTS tag cascade;
drop table IF EXISTS "user" cascade;
drop table IF EXISTS user_status cascade;
-- +migrate StatementEnd
