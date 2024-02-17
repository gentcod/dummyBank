A simple bank orchestrating how database transactions are carried out

## Database migrations:

There are two instances for database migrations:
  - One supports migrations using Goose[https://github.com/pressly/goose], the migration files for this can be found in the `sql/schemas` directory. Best case scenario is to use goose to run migrations locally. Reference the Makefile for more details. In order not to run migrations automatically after build, comment out the `runDBMigration` method/func in the main.go file
###
   - The other supports migrations using Golang-Migrate[https://github.com/golang-migrate/migrate]. It enables migrations to be run when the executaable file is built. The migration files are located in the `sql/migrations` directory. I utilized this based on issues with database migrations during containerization.