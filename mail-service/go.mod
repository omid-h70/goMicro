module mail-service

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
	github.com/xhit/go-simple-mail/v2 v2.13.0
)

require (
	github.com/go-test/deep v1.1.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/toorop/go-dkim v0.0.0-20201103131630-e1cd1a0a5208 // indirect
)

replace helpers v0.0.0 => ../helpers
