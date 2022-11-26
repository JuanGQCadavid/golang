go test  -v -> Run test
go test  -v -cover -> Show coverage percentage
go test -coverprofile=coverage.out -> Print a report or the test coverage
go tool cover -html=coverage.out -> Transform from .out to html