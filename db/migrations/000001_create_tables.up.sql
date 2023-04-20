CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    id uuid default uuid_generate_v4() not null primary key ,
    login      varchar(100)                        not null,
    password   varchar(100)                        not null,
    created_at timestamp default now()             not null,
    updated_at timestamp default now(),
    deleted_at timestamp,
    email      varchar                             not null
);

create table user_img (
                          id uuid default uuid_generate_v4() not null primary key,
                          user_id uuid references users(id),
                          mime_type varchar(256) not null,
                          bucket_name varchar(512) not null,
                          file_name varchar(1024) not null,
                          created_at timestamp default now()             not null,
                          updated_at timestamp default now(),
                          deleted_at timestamp
);

create table jobs (
                      id uuid default uuid_generate_v4() not null primary key,
                      name varchar(512) not null
);

create table orgs (
                      id uuid default uuid_generate_v4() not null primary key,
                      name varchar(512) not null
);

create table persons (
                         id uuid default uuid_generate_v4() not null primary key,
                         user_id uuid not null references users(id),
                         last_name varchar(256) not null,
                         first_name varchar(256) not null,
                         sur_name varchar(256),
                         job_id uuid references jobs(id),
                         org_id uuid references orgs(id),
                         created_at timestamp default now()             not null,
                         updated_at timestamp default now(),
                         deleted_at timestamp
);

create table roles
(
    id uuid default uuid_generate_v4() not null primary key,
    name varchar(256) not null
);

create function defaultrole() returns uuid
    language plpgsql
as
$$
declare
roleId uuid;
begin
select id
into roleId
from roles r
where r.name = 'user';
return roleId;
end;
$$;

create table user_roles
(
    user_id uuid references users(id),
    role_id uuid default defaultrole() references roles(id)
);

create function get_user_roles(input_id uuid)
    returns TABLE(rolenames character varying)
    language plpgsql
as
$$
begin
return query
select
    r.name
from users u
         join user_roles ur on u.id = ur.user_id
         join roles r on ur.role_id = r.id
where u.id = input_id;
end;
$$;

create function delete_user(input_id uuid) returns boolean
    language plpgsql
as
$$
begin
update users set deleted_at = now() where id = input_id;
update persons set deleted_at = now() where user_id = input_id;
update user_img set deleted_at = now() where user_id = input_id;
return true;
end;
$$;

create function createUser(
    inp_login character varying,
    inp_password character varying,
    inp_email character varying,
    inp_lastname character varying,
    inp_firstname character varying,
    inp_surname character varying,
    inp_job_id uuid,
    inp_org_id uuid,
    inp_role_id uuid
) returns uuid
    language plpgsql
as
$$
declare
created_user_id uuid;
begin
insert into users (login, password, email)
values (inp_login, crypt(inp_password, gen_salt('md5')), inp_email)
    returning id into created_user_id;

case when inp_role_id is null then insert into user_roles (user_id) values (created_user_id);
else insert into user_roles(user_id, role_id) values (created_user_id, inp_role_id);
end case;

insert into persons (user_id, last_name, first_name, sur_name, job_id, org_id)
values (created_user_id, inp_lastname, inp_firstname, inp_surname, inp_job_id, inp_org_id);

return created_user_id;
end;
$$;

insert into roles (name)
values ('user'), ('executor'), ('admin');

insert into jobs (name)
values ('admin');

insert into orgs (name)
values ('Админская');

select createUser('admin',
                  'admin12345',
                  'admin@test.ru',
                  'testov',
                  'test',
                  'testovich',
                  (select id from jobs where name = 'admin'),
                  (select id from orgs where name = 'Админская'),
                  (select id from roles where name = 'admin')
           );