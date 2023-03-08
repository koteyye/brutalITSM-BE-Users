drop table person;

drop table user_roles;

drop table roles;

drop table user_img;

drop table "user";

drop function public.defaultrole();

drop function public.get_user(varchar, varchar);

drop function public.get_user_roles(uuid);

drop function public.delete_user(uuid);

drop function public.createuser(varchar, varchar, varchar, varchar, varchar, varchar, varchar, varchar, varchar);