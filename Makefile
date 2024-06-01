IMAGE_NAME = docker_image
DOCKERFILE = Dockerfile
CONTAINER_NAME = docker_container
build:
	docker build -t $(IMAGE_NAME) .
run:
	docker run -p 8080:8080 $(IMAGE_NAME)
clean:
	docker system prune -af