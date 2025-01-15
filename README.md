# Resume_Website

 TO INSTALL:
- Install MySQL (sudo apt install mysql-server)
- Add user "users" with password "User1234":
        CREATE USER 'users1234'@'%' IDENTIFIED BY 'User1234';
        GRANT ALL PRIVILEGES ON user_logs.* TO 'users1234'@'%';
        FLUSH PRIVILEGES;

TO RUN:
-  
    sh run_docker.sh
or 
    go run main.go
