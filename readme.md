# Gapp

# Migrations
```bash
go install github.com/rubenv/sql-migrate/...@latest
sql-migrate up -env="production" -config=repository/mysql/dbconfig.yml
sql-migrate down -env="production" -config=repository/mysql/dbconfig.yml -limit=1
sql-migrate status -env="production" -config=repository/mysql/dbconfig.yml
```