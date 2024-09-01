# novaro-server

## Login

http://127.0.0.1:8080/v1/auth/ to login with twitter directly

http://127.0.0.1:8080/v1/auth?icode=12345678 to login with invitation code firstly and then with twitter

## Invitation Codes

### Generate

```
GET http://127.0.0.1:8080/v1/api/invitation/codes/add

Cookie: {cookie}
```

Get the cookie after you login

### Check the exist invitation codes

invitation_codes table in the db


1.运行
```
go mod tidy
go run src/main.go | tee out.log
```
