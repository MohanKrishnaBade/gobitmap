all:  run bad

hello:
	echo "Hello"
bad:
	 ./bin/app
run:
	go build -o bin/app