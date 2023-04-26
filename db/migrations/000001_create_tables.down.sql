drop table user_roles;

drop table role_permissions;

drop table permissions;

drop table roles;

drop table user_img;

drop table persons;

drop table jobs;

drop table orgs;

drop table users;

drop function public.defaultrole();

drop function public.get_user_roles(uuid);

drop function public.delete_user(uuid);

drop function createuser(varchar, varchar, varchar, varchar, varchar, varchar, uuid, uuid, uuid);

drop extension "pgcrypto";

drop extension "uuid-ossp";