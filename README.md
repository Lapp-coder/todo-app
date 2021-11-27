# REST API for todo list application ```v2.1.0```

***

### To run the application in docker:
```
$ make build && make run
```

### Before launching the application for the first time:

* Create a .env file in the root of the application with the following contents:
  
```
POSTGRES_PASSWORD=<your-password>
SIGNING_KEY=<any-character-set>
SALT=<any-character-set>
TZ=<timezone>
```

* Apply migrations to the database:

```
$ export POSTGRES_PASSWORD=<your-password> 
```

```
$ make migrate-up
```

### Use the following to create documentation:
```
$ make swag
```
### Documentation can be found at: 
`http://<host>:<port>/swagger/index.html`

### You can also run the tests with the following command:
```
$ make test
```
