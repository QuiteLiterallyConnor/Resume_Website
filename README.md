# Resume_Website

 TO INSTALL:
- Install MySQL (sudo apt install mysql-server)
- Add user "users" with password "User1234":
        CREATE USER 'users1234'@'%' IDENTIFIED BY 'User1234';
        GRANT ALL PRIVILEGES ON user_logs.* TO 'users1234'@'%';
        FLUSH PRIVILEGES;

- Add ngrok key with:
    export NGROK_AUTH_TOKEN=""

- Build
    docker compose build
    go build -o main main.go

TO RUN:
-  
    docker compose up -d
or 
    go run main.go
