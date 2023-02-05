# REST-API TodoApp

«TodoApp» is a convenient REST API for providing data from the server to the user of a web application or website.

## Installation
To build TodoApp from source, simply run `git clone https://github.com/pyuldashev912/TodoApp` and `cd` into the project source directory. Then run `make`. After this, you should have a binary called `todoapp` in the current directory.
```
$ make
go build -v ./cmd/todoapp
...
```
It is assumed that you are using PostgreSQL. Create a new database using `createdb todoapp`. Add the following application configurations to .env file into the project source directory:
```
BIND_ADDR = ":8080"
LOG_LEVEL = "info"
DATABASE_URL = "host=localhost dbname=todoapp sslmode=disable"
SESSION_KEY = "<generate session key>"
```
Launch the application
```
$ ./todoapp
INFO[0000] Listening...
```

# Endpoints
## Create a user
### Request
`POST /sign-up`
```
http POST localhost:8080/sign-up name=name email=email password=password
```
### Response
```
{
    "id": int,
    "name": string,
    "email": string,
}
```
## Login
### Request
`POST /sign-in`
```
http --session=user POST localhost:8080/sign-in email=email password=password
```
### Response
```
{
    "info": string
}
```
## Logout
### Request
`POST /users/logout`
```
http --session=user POST localhost:8080/users/logout
```
### Response
```
{
    "info": string
}
```
## Who am I
### Request
`GET /users/me`
```
http --session=user GET localhost:8080/users/me
```
### Response
```
{
    "id": int,
    "name": string,
    "email": string,
}
```
## Create a task
### Request
`POST /users/tasks`
```
http --session=user POST localhost:8080/users/tasks title="Some task" description="Some text"
```
### Response
```
{
    "id": int,
    "title": string,
    "description": string,
    "done": bool,
    "creation_date": string
}
```
## Get task
### Request
`GET /users/tasks/id`
```
http --session=user GET localhost:8080//users/tasks/id"
```
### Response
```
{
    "id": int,
    "title": string,
    "description": string,
    "done": bool,
    "creation_date": string
}
```
## Get done/underdone tasks
### Request
`GET /users/tasks?done=true/false`
```
http --session=user GET localhost:8080/users/tasks?done=true"
```
### Response
```
[
    {
        "id": int,
        "title": string,
        "description": string,
        "done": bool,
        "creation_date": string
    }
    ...
]
```
## Get all tasks
### Request
`GET /users/tasks`
```
http --session=user GET localhost:8080/users/tasks"
```
### Response
```
[
    {
        "id": int,
        "title": string,
        "description": string,
        "done": bool,
        "creation_date": string
    }
    ...
]
```
## Mark the task as completed
### Request
`PATCH /users/tasks/id`
```
http --session=user PATCH localhost:8080/users/tasks/id"
```
### Response
```
{
    "info": string
}
```
## Delete the task
### Request
`DELETE /users/tasks/id`
```
http --session=user delete localhost:8080/users/tasks/id"
```
### Response
```
{
    "info": string
}
```