#!/bin/bash

echo -e "Building the project, please wait...."
go build

echo -e "Project started, the web application is running on port 8082: localhost:8082. Enjoy !"
./TwitterGo
