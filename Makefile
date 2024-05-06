
dft:		# This lists the targets in this Makefile
	cat Makefile | grep "^[a-z-]*:"
	@echo Choose an explicit option from these

start: app	# Runs the main.go application
	go run main.go

app:		# Calls the start
	@echo Starting the orders using go run main.go

redis:		# Starts redis using local ./start-redis
	./start-redis

docker-redis:	# Stops any local redis and starts the docker version
	./stop-redis
	docker run -p 6379:6379 redis:latest

build:		# Builds the main.go to generate the executable 
	go build main.go

create:		# Uses runorder.sh to test create order
	./runorder.sh create

testdata:	# Uses runorder.sh to create 5 new orders for user rabby
	for i in 1 2 3 4 5; do .runorder.sh -name rabby create; done

list:		# Uses runorder.sh to list existing orders stored in redis
	./runorder.sh list | jq

get:		# Uses runorder.sh to get a specific order
	./runorder.sh get ${ID} | jq

failget:	# Uses runorder.sh to test get order for invalid ID 1234567890
	./runorder.sh -v get 1234567890


