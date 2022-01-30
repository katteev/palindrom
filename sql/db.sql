create schema if not exists test;

connect test;

create table test.users(ID int(10) primary key not null, fornavn varchar(50), etternavn varchar(50));

insert into test.users (ID, fornavn, etternavn) values (1, "Anna", "Hannah");

insert into test.users (ID, fornavn, etternavn) values (2, "AnnaTest", "Hannah2");
