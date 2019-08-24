
SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

docker:
	docker build -t opf --no-cache .

sources:
	@echo $(SRCS)