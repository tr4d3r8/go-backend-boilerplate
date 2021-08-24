# backend-master-golang

This is a boilerplate backend written in Go that uses postgres as a datastore. 

## Features 
* database connection (postgres)
* database migration 
* crud code gen 
* crud unit testing 


## Dependencies 

* Go
* Docker
* Taskfile 
* golang-migrate
* sqlc

## Getting started

* clone the repo
* get dependencies ``` go mod tidy ```
* start run postgres in docker for dev ```run postgres```
* create database ```run createdb ```
* migrate database ```run migrate_up```
* run tests ``` run tests ```


## Generate CRUD functions 
* modify query in db/query/
* run ``` sqlc generate ``` 
* this will create or modify the code gen files in ../sqlc/ 
    * EXAMPLE:
        * Modifying ../query/account.sql
        * Running ``` sqlc generate ```
        * Will create/modify ../sqlc/account.sql.go
 


