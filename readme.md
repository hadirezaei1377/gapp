# Game App


related to https://github.com/hadirezaei1377/goboot

match two persons and serve some questions for them, after finishing the match, indicate the winner.

** rest api
** time free
** questions are graded, hard, medium, easy
** questions have different categories
** leader word shows the winner, at each time who is the winner, at each category  who is the winner, at each grade who is the winner

** panel admin ---> can see list of users, games, scores , ...



db migration:
go install github.com/rubenv/sql-migrate/...@latest
sql-migrate up -env="production" -config=repository/mysql/dbconfig.yml
sql-migrate down -env="production" -config=repository/mysql/dbconfig.yml -limit=1
sql-migrate status -env="production" -config=repository/mysql/dbconfig.yml

