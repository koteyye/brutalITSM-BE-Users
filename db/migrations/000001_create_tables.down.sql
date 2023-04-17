drop table persons;

drop table user_roles;

drop table roles;

drop table user_img;

drop table users;

drop table jobs;

drop table orgs;

drop function public.defaultrole();

drop function public.get_user_roles(uuid);

drop function public.delete_user(uuid);

drop function createuser(varchar, varchar, varchar, varchar, varchar, varchar, uuid, uuid, uuid);