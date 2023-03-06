CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table "user"
(
    id         uuid   default uuid_generate_v4()  not null
        primary key,
    login      varchar(100)                        not null,
    password   varchar(100)                        not null,
    created_at timestamp default now()             not null,
    updated_at timestamp,
    deleted_at timestamp,
    email      varchar                             not null
);

create table person
(
    id          uuid      default uuid_generate_v4() not null
        primary key,
    last_name   varchar(256)                        not null,
    first_name  varchar(256)                        not null,
    middle_name varchar(256),
    job_name    varchar(256),
    org_name    varchar(256),
    created_at  timestamp default now()             not null,
    updated_at  timestamp,
    deleted_at  timestamp,
    user_id     uuid                                not null
        constraint user_fk
            references "user"
);

create table roles
(
    id   uuid default uuid_generate_v4() not null
        primary key,
    name varchar(255)
);

create function defaultrole() returns uuid
    language plpgsql
as
$$
declare
roleid uuid;
begin
select id
into roleid
from roles r
where r.name = 'user';
return roleid;
end;
    $$;

create table user_roles
(
    user_id uuid
        references "user",
    role_id uuid default defaultrole()
        references roles
);

create function get_user(log character varying, pass character varying)
    returns TABLE(id uuid, rolename character varying)
    language plpgsql
as
$$
begin
return query
select
    u.id,
    r.name
from "user" u
         join user_roles ur on u.id = ur.user_id
         join roles r on ur.role_id = r.id
where u.login = log and u.password = pass;
end;
    $$;

create function get_user_roles(input_id uuid)
    returns TABLE(rolenames character varying)
    language plpgsql
as
$$
begin
return query
select
    r.name
from "user" u
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
        --Удалить связь с ролями
delete from user_roles where user_id = input_id;

--Удалить персону
delete from person where user_id = input_id;

--Удалить пользователя

delete from "user" where id = input_id;

return true;
end
    $$;


create function createuser(input_login character varying, input_password character varying, input_email character varying, input_lastname character varying, input_firstname character varying, input_middle_name character varying, input_job_name character varying, input_org_name character varying, input_role_name character varying) returns uuid
    language plpgsql
as
$$
declare
created_user_id uuid;
        input_role_id uuid;
begin

        --create user
insert into "user" (login, password, email) values (input_login, input_password, input_email) returning id into created_user_id;

--search role_id
case when input_role_name != '' then input_role_id = (select id from roles where "name" = input_role_name);
else input_role_id = null;
end case;

        --create user_roles
case
            when input_role_id is null
                then insert into user_roles (user_id) values (created_user_id);
else insert into user_roles (user_id, role_id) values (created_user_id, input_role_id);
end case;

        -- create person
insert into person (last_name, first_name, middle_name, job_name, org_name, user_id)
values (input_lastname, input_firstname, input_middle_name, input_job_name, input_org_name, created_user_id);

return created_user_id;

end;
    $$;


insert into roles (name)
values ('user'), ('executor'), ('admin');

select createUser('admin', '696f396768747265683565726865667364676577726767646663d033e22ae348aeb5660fc2140aec35850c4da997', 'itsm87@mail.ru', 'Admin', 'Admin', 'Admin', 'Admin', 'Admin', 'admin');