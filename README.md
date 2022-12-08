# homeHandler_backend
Backend application to get life smarter. Actually the main feature consist into an offer to the user a tool platform utils to mange daily task (for example track daily food plan, get an idea of the expense of weekly food plan etc.)

## First configuration
Those are first step to setup the project.

- Create at the root level of the project a text file called secrets.txt and type in a secure password that wille be used
as your postgres database password.

- Create an environment file and name it dev.env or prod.env based on the environment where the app is launched. In the env file you need to create the following keys:
 - DB_STRING=postgresConnectionString
 - DB_PASSWORD=databasePassword
 - JWT_SECRET=secretUsedToSignTheJWT
 - SERVER_PORT=serverPortNumber
 - MAIL_SERVER_HOST=mailServerHostname
 - MAIL_SERVER_PORT=mailServerPort

