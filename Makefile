build:
	docker image build -t forum:alpha .
run:
	docker run -d -p 8081:8081 --name myforum forum:alpha 
stop:
	docker stop $$(docker ps -aq)
remove:
	docker rm $$(docker ps -aq)
clean: stop remove
	docker system prune -a
all: build run

