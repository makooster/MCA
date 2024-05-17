create table if not exists actors (
    id serial primary key,
    full_name varchar(255),
    dorama_id integer,
    foreign key (dorama_id) references doramas(dorama_id)
 );