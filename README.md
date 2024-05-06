

# Golang sample Web service
This project is based on the on-line course for Go project "orders-api".
It is a training exercise. 

## Overview
The code uses the common orders api to help demonstrate the following in Go:
1. Basic Go modules
2. Class methods and inline functions
3. Structs and json annotations to help build REST responses
4. Routes and handler patterns for API endpoints
5. Channels, defer and go routines to handle thread logic
6. Redis API for persisting orders

## Usage
Various support scripts are added to make testing of the application easy.
A Makefile is provided that shows how to build, run and test the app and related services.
1. dft:			# The default make option - which lists these options (targets) from the Makefile
2. start: app	# Runs the main.go application
3. app:			# Calls the start
4. redis:		# Starts redis using local ./start-redis
5. build:		# Builds the main.go to generate the executable 
6. create:		# Uses runorder.sh to test create order
7. testdata:	# Uses runorder.sh to create 5 new orders for user rabby
8. list:		# Uses runorder.sh to list existing orders stored in redis
9. get:			# Uses runorder.sh to get a specific order
10. failget:	# Uses runorder.sh to test get order for invalid ID 1234567890
11. docker-redis:	# Stops any local redis and starts the docker version