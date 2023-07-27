create function getUserList()
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
       )
        select array_agg(permissioncode) from userPermissions) permissions,
       json_build_object('mimeType', ui.mime_type, 'bucketName', ui.bucket_name, 'fileName', ui.file_name) avatar
from "users" u
         join persons p on u.id = p.user_id
         join jobs j on j.id = p.job_id
         join orgs o on o.id = p.org_id
         left join user_img ui on u.id = ui.user_id
where u.deleted_at is null and ui.deleted_at is null;
end;
$$;