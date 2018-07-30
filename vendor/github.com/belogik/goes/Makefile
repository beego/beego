help:
	@echo "Available targets:"
	@echo "- test: run tests"
	@echo "- installdependencies: installs dependencies declared in dependencies.txt"

installdependencies:
	cat dependencies.txt | xargs go get

test: installdependencies
	go test -i && go test
