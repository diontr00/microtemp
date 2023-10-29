build:   ## Build  
	go build  -o bin/main cmd/server/main.go

run:  build ## Run 
	./bin/main

template: 
	ruby ./initscript.rb $(NAME) 

help: 
	@bash ./script/print_help.sh
