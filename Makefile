build:
	docker build -t kalyuzhn_constanta .
start:
	docker run kalyuzhn_constanta -p 8080:8080 -d
clean:
	docker stop kalyuzhn_constanta
	docker rmi kalyuzhn_constanta
run: build start
test:
	curl -vvv -X POST http://localhost:8080/ \
		-d "{\"urls\":[\"https://ya.ru/\",\"https://google.com/\", \"wss://pornhub.com/path_that_dont_exists\"]}"
.PHONY: build start run clean test