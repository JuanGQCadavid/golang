1. go mod edit -replace example.com/greetings=../greetings -> This will replace the example.com with the relative path to the module

2. go mod tidy -> Sync modules dependencies with those ones that are requered but no yedt imported