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
                          user_id uuid,
                          mime_type varchar(256) not null,
                          bucket_name varchar(512) not null,
                          file_name varchar(1024) not null,
                          created_at timestamp default now()             not null,
                          updated_at timestamp default now(),
                          deleted_at timestamp,
                          foreign key (user_id) references users(id)
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
                         user_id uuid not null,
                         last_name varchar(256) not null,
                         first_name varchar(256) not null,
                         sur_name varchar(256),
                         job_id uuid,
                         org_id uuid,
                         created_at timestamp default now()             not null,
                         updated_at timestamp default now(),
                         deleted_at timestamp,
                         foreign key (user_id) references users(id),
                         foreign key (job_id) references jobs(id),
                         foreign key (org_id) references orgs(id)
);

create table roles
(
    id uuid default uuid_generate_v4() not null primary key,
    name varchar(256) not null
);

create table permissions
(
    id uuid default uuid_generate_v4() not null primary key,
    permissionCode varchar(256) not null,
    created_at timestamp default now()             not null,
    updated_at timestamp default now(),
    deleted_at timestamp
);

create table role_permissions
(
    roles_id uuid not null,
    permission_id uuid not null,
    foreign key (roles_id) references roles(id),
    foreign key (permission_id) references permissions(id)
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
    user_id uuid,
    role_id uuid default defaultrole(),
    foreign key (user_id) references users(id),
    foreign key (role_id) references roles(id)
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


insert into permissions (permissioncode)
values ('newsRead'), ('newsWrite'), ('adminPanelRead'), ('adminPanelWrite');

with users_cte as (
    select id "roleid" from roles where name = 'admin'
)
insert into role_permissions(roles_id, permission_id) select (select roleid from users_cte), id from permissions;

create function getUserById(userId uuid)
    returns table (id uuid, login varchar(100), last_name varchar(256), first_name varchar(256), sur_name varchar(256), job_name varchar(512), org_name varchar(512), role_list varchar[], permissions varchar[], avatar json)
    language plpgsql
as
$$
begin
    return query
        select u.id, u.login, p.last_name, p.first_name, p.sur_name, j.name job_name, o.name org_name,
               (select array_agg(r.name) from roles r join user_roles ur on r.id = ur.role_id where ur.user_id = u.id) role_list,
               (with userPermissions as (
                   select distinct p.permissioncode
                   from permissions p
                            join role_permissions rp on p.id = rp.permission_id
                            join user_roles ur on ur.role_id = rp.roles_id
                   where ur.user_id = userId
               )
                select array_agg(permissioncode) from userPermissions) permissions,
               json_build_object('mimeType', ui.mime_type, 'bucketName', ui.bucket_name, 'fileName', ui.file_name) avatar
        from "users" u
                 join persons p on u.id = p.user_id
                 join jobs j on j.id = p.job_id
                 join orgs o on o.id = p.org_id
                 left join user_img ui on u.id = ui.user_id
        where u.id = userId and u.deleted_at is null and ui.deleted_at is null;
end;
$$;