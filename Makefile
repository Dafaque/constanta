build:
	docker build -t kalyuzhn_constanta .
start:
	docker run -p 8080:8080 -d kalyuzhn_constanta
run: build start
test:
	curl -vvv -X POST http://localhost:8080/ \
		-d "{\"urls\":[ \
			\"https://jsonplaceholder.typicode.com/photos\", \
			\"https://jsonplaceholder.typicode.com/comments\", \
			\"https://jsonplaceholder.typicode.com/posts\" \
		]}"
.PHONY: build start run test