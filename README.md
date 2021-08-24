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



