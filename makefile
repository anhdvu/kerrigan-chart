run:
	@echo "Running test server..."
	go run main.go
build:
	@echo "Building binary file for Linux..."
	GOOS=linux go build -o bin/kc
test:
	@echo "Downloading resources from production server..."
	scp -P 1006 trym@dace.dev:/home/trym/devev/to_the_moon/sentry/data/*.json .
clean:
	@echo "Removing log files..."
	rm *.log
deployfe:
	@echo "Uploading frontend files to production server..."
	scp -P 1006 -r frontend/ trym@dace.dev:/home/trym/devev/to_the_moon/kchart/
deploybe:
	@echo "Deploying server files to production server..."
	scp -P 1006 bin/kc trym@dace.dev:/home/trym/devev/to_the_moon/kchart/
deployall:
	@echo "Deploying all to production server..."
	scp -P 1006 -r bin/kc frontend/ ui/ trym@dace.dev:/home/trym/devev/to_the_moon/kchart/