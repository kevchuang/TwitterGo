# BJTUitter

This document describe how to install and deploy our web application **BJTUitter**.

## Prerequisites

1. [PostgreSQL 9.6](https://www.postgresql.org/download/ "PostgreSQL download")
2. [Go 1.9.2](https://www.golang.org/dl "Go download")
3. Pgcrypto

## Installation

The installation is divide in 4 parts:
1. Installation, creation and deployment of the PostgreSQL database.
2. Installation of multiple go packages.
3. Configuration of the different components.
4. Deployment.

You can either follow the instructions line as follow, or use the bash scripts **deploy.sh** and **launch.sh** at the root of the repository (**$ROOT**).

### The database

Before anything else, you need to download the pgcrypto extension use by our project. On Ubuntu distribution you can add it like that:

* `sudo apt-get install postgresql-contrib-9.6`

Once you have setup PostgreSQL you can create the database that **BJTUitter** will use. Use this command line at the **$ROOT** to create the database:

* `psql -U postgres < create_database.sql`

And then use this command line (always at the **$ROOT**) in order to create the tables:

* `psql -U postgres -d BJTUitter < create_tables.sql`

### Go packages

Be sure that your **$GO_ROOT** and **$GO_PATH** are setup. See the [documentations](https://golang.org/doc/install "Goland doc") of Golang for more details
Download the go packages that the application use via this command:

* `go get github.com/lib/pq && go get github.com/gorilla/mux`

### Configuration

You can configure the connection of the web application to the PostgreSQL's database throught this file **$ROOT/main.go** at the line 17:

* `db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=BJTUitter sslmode=disable")`

### Deployment

Now you can try our web application **BJTUitter**, simply use this command line at the **$ROOT**:

* `go build && ./TwitterGo`

You can also launch our multiple test to insure our code integrity. Use this command line at the **$ROOT**:

* `go test`

### Contributors

* **Kevin Chuang**

* **Sylvain Birukoff**

* **Louis Chaumier**
