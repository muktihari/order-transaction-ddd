# Order Transaction in DDD (Domain-Driven Design)

The idea of this repository is to demonstrate my approach to implement DDD for **Order Transaction** service in Golang as well as structuring code, logging, instrumenting, unit testing, etc. This is an experimental approach inspired by [https://github.com/marcusolsson/goddd](https://github.com/marcusolsson/goddd) and shall not be taken for granted as best practice or guidance.

**transaction** package is the heart of the domain model, and **Order** struct is the central class of the domain model. Everything start from there.

## Running the Application
Application is running default on port 8080 using inmem repository. See **Makefile** for shortcut. Example:
```sh
make run
```
Run Local using mongo as repository:
```sh
make run-mongo-migrate
```
Run via **docker-compose**:
```sh
docker-compose up
```

*Note: mongodb should have replica to enable transaction.*

Unit tests and PostgreSQL repository has not been fully completed though.