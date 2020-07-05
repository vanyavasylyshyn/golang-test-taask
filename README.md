# golang-test-task
Authorization service for vacancy test task.
App is running on https://vasylyshyn-auth-service.herokuapp.com/

## Routes

### GET /api/credentials/:userID/new
Create credentials for user. Pass user id in place of userID

### POST /api/credentials/refresh
Refresh credentials. Pass refresh credential to Authorization Bearer Token.

### DELETE /api/credentials/:userID/destroy/all/refresh
Delete all refresh credentials for user.

### DELETE /api/credentials/destroy/refresh
Delete specific refresh credential. Pass refresh credential to Authorization Bearer Token.
