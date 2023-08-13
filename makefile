APP_NAME = bank
TESTS_DIR = ./...
COVERAGE_DIR = ./coverage
COVERAGE_FILE = cover.out

GOCMD = go
GOTEST = $(GOCMD) test

tests:
	mkdir $(COVERAGE_DIR) -p && go clean -testcache && $(GOTEST) -v $(TESTS_DIR) -coverprofile=$(COVERAGE_DIR)/$(COVERAGE_FILE) && $(GOCMD) tool cover -html=$(COVERAGE_DIR)/$(COVERAGE_FILE)

run:
	sudo docker compose up --build
