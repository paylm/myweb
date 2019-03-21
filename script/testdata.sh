#!/bin/bash

table="project_book"
database="ppl"
host="127.0.0.1"

# create test data for project_book
for i in {100..300} ; do
	mysql -u paylm -p123456 -h ${host} ${database} -e "insert into ${table} insert into project_book(project_id,by_user,stat) values($i,13,0);"
done

