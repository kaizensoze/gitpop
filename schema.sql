drop table if exists users;
create table users (
  id integer,
  username text,
  access_token text,
  primary key (id)
);

drop table if exists ignores;
create table ignores (
  user_id integer,
  id integer,
  starred integer,
  primary key (user_id, id)
);
