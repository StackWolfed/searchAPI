echo "This DB is only for developing porpuses, it does not save data"
docker run --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=toor --name database -d mariadb
docker logs -f database