#!/bin/bash

docker exec port-service_db_1 bash -c 'mysql -uroot -proot -D ports-test -e "SET FOREIGN_KEY_CHECKS = 0;
truncate port;
SET FOREIGN_KEY_CHECKS = 1;"'