# User
POST /users/register

curl -X POST http://localhost:8080/users/register
-H "Content-Type: application/json"
-d '{"Name": "Hadi", "PhoneNumber": "09124"}'