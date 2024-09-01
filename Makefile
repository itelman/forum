build:
	docker image build -t forum .
run:
	docker run -dp 8080:8080 --name forum forum 
stop:
	docker stop forum
remove:
	docker rm -f forum
clean: stop remove
	docker system prune -a
all: build run

