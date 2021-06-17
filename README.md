# REST API for todo list application ```v1.1.0```

***

### To run the application:
```make build && make run```

### Before launching the application for the first time:

* Apply migrations to the database:

```make migrate```
  
* Create an .env file in the root of the application with the following contents:
```
POSTGRES_PASSWORD=YourPostgresPassword
SIGNING_KEY=YourSigningKey
```

### Use the following to create documentation:
```make swag```
### Documentation can be found at: http://localhost:8080/swagger/index.html

### You can also run the tests with the following command:
```make test```

