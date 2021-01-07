run:
	@echo "Running test server..."
	go run main.go
build:
	@echo "Building binary file for Linux..."
	GOOS=linux go build -o bin/kc
testrs:
	@echo "Downloading resources from production server..."
	scp -P 1006 trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/checker.txt .
	scp -P 1006 trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/historical_delta.txt .
clean:
	@echo "Removing log files..."
	rm *.log
deployfe:
	@echo "Uploading frontend files to production server..."
	scp -P 1006 -r frontend/ trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/kerrigan-chart/
deploybe:
	@echo "Deploying server files to production server..."
	scp -P 1006 kc trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/kerrigan-chart/
deployall:
	@echo "Deploying all to production server..."
	scp -P 1006 -r bin/kc frontend/ trym@dace.dev:/home/trym/devev/to_the_moon/kerrigan/kerrigan-chart/