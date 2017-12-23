#!/bin/bash

psql -q -U postgres < create_database.sql
echo -e "Database created."

psql -q -U postgres -d BJTUitter < create_tables.sql
echo -e "Tables created."

go get github.com/lib/pq && go get github.com/gorilla/mux >&2
echo -e "Go's packages installed"
echo -e "The setup of BJTUitter is done. Please use the launch.sh to start the web application."
