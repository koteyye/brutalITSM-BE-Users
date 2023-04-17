# Brutal-ITSM-Users
Микросервис для хранения пользователей и авторизации в brutal itsm
[![logo](https://github.com/koteyye/brutal-itsm-fe/blob/master/public/logoend.png?raw=true "logo")](https://github.com/koteyye/brutal-itsm-fe/blob/master/public/logoend.png?raw=true "logo")
## REST
### AUTH
- **POST /auth/sign-in**  - получить JWT;
- **GET /auth/me** - получить пользователя по JWT;

### USERS
- **POST /api/users/create**  - создать пользователя;
- **POST /api/users/avatar/upload/:id** - загрузить аватар пользователя;
- **GET /api/users/**   - получить список пользователей;
- **GET /api/users/:id**  - получить пользователя по ID;
- **GET /api/users/roles** - получить список ролей;

### SEARCH
- **GET /api/search/job/:jobName** - получить список должностей по совпадению имени;
- **GET /api/search/org/:orgName** - получить список организаций по совпадению имени;

## gRPC
- **GetUserByToken** - проверка JWT и получение данных о пользователе


# Связанные проекты
- **FrontEnd на Vue 3** - https://github.com/koteyye/brutalITSM (приостановлен)
- **FrontEnd на React JS** - https://github.com/koteyye/brutal-itsm-fe
- **Brutal-ITSM-Trabls** - https://github.com/koteyye/brutalITSM-BE-Trabls