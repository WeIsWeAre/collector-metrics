module handlers

go 1.22.0

replace service => ../service

replace storage => ../storage

replace models => ../models

require service v0.0.0-00010101000000-000000000000

require storage v0.0.0-00010101000000-000000000000 //indirect

require models v0.0.0-00010101000000-000000000000

require github.com/stretchr/testify v1.9.0

require (
	github.com/davecgh/go-spew v1.1.1 //indirect
	github.com/gorilla/mux v1.8.1 
	github.com/pmezard/go-difflib v1.0.0 //indirect
	gopkg.in/yaml.v3 v3.0.1 //indirect
)
