DROP TABLE IF EXISTS clientinfo;


create table clientinfo
(
id SERIAL PRIMARY KEY,
name VARCHAR(30),
phone VARCHAR(30),
city VARCHAR(30),
address VARCHAR(30),
region VARCHAR(30),
email VARCHAR(30)
);