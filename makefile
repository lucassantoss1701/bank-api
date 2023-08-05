#APPLICATION PARAMETERS
APP_NAME = bank
TESTS_DIR = ./internal...
COVERAGE_DIR = ./coverage
COVERAGE_FILE=cover.out

#GO parameters
GOCMD=go
GOTEST=$(GOCMD) test

tests:
	mkdir $(COVERAGE_DIR) -p && $(GOTEST) -v $(TESTS_DIR) -coverprofile=$(COVERAGE_DIR)/$(COVERAGE_FILE) && $(GOCMD) tool cover -html=$(COVERAGE_DIR)/$(COVERAGE_FILE)