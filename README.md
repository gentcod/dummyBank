A simple bank orchestrating how database transactions are carried out

## Database migrations:

There are two scenarios for database migrations:
  - One involves migrations using Goose[https://github.com/pressly/goose], the migration files for this can be found in the `sql/schemas` directory. Goose can be used to run migrations locally; reference the `Makefile` file for more details. Goose is also used to run DB migrations when executing workflow for test, reference the `test.yaml` file in `.github/workflows` directory for more information In order not to run migrations automatically after build, comment out the `runDBMigration` method/func in the main.go file.
###
   - The other involves migrations using Golang-Migrate[https://github.com/golang-migrate/migrate]. It enables migrations to be run when the executaable file is built. The migration files are located in the `sql/migrations` directory. I utilized this based on issues with database migrations during containerization.