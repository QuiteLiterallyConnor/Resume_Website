# Resume_Website

 TO INSTALL:
- Install MySQL (sudo apt install mysql-server)
- Add user "users" with password "User1234":
        CREATE USER 'users1234'@'%' IDENTIFIED BY 'User1234';
        GRANT ALL PRIVILEGES ON user_logs.* TO 'users1234'@'%';
        FLUSH PRIVILEGES;

- Create Database:
    CREATE DATABASE IF NOT EXISTS user_logs;USE user_logs;

- Create table "accessed_parts":
    CREATE TABLE IF NOT EXISTS accessed_parts (
        id INT NOT NULL AUTO_INCREMENT,
        user_id INT NOT NULL,
        part VARCHAR(255) NOT NULL,
        time_accessed DATETIME NOT NULL,
        PRIMARY KEY (id),
        KEY user_id (user_id)
    );

- Create table "user_info":

CREATE TABLE IF NOT EXISTS user_info (
    id INT NOT NULL AUTO_INCREMENT,
    ip VARCHAR(45) NOT NULL,
    accessed_parts TEXT DEFAULT NULL,
    time_accessed DATETIME DEFAULT NULL,
    first_time_accessed DATETIME DEFAULT NULL,
    last_time_accessed DATETIME DEFAULT NULL,
    blacklisted TINYINT(1) DEFAULT 0,
    client_data TEXT DEFAULT NULL,
    country VARCHAR(100) DEFAULT NULL,
    city VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (id)
);

- Add ngrok key with:
    export NGROK_AUTH_TOKEN=""

- Build
    docker compose build

TO RUN:
- Docker 
    docker compose up -d
or 
- Standalone
    go run main.go
