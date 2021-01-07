hello:
	@echo "Wassup!!!"
run:
	go run main.go
build:
	@echo "Building binary file for Linux..."
	GOOS=linux go build -o bin/kc
testrs:
	scp -P 1006 trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/checker.txt .
	scp -P 1006 trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/historical_delta.txt .
clean:
	rm *.log