# Имя Docker-образа
IMAGE_NAME = tg_bot
# Версия Docker-образа
IMAGE_VERSION = latest
# Путь к Dockerfile
DOCKERFILE_PATH = .
# Порт, который будет прокинут из контейнера на хост
HOST_PORT = 8080
# Построение и запуск контейнера
all: build run
hello:
	echo "hello"
# Построение Docker-образа
build:
	docker build -t $(IMAGE_NAME):$(IMAGE_VERSION) $(DOCKERFILE_PATH)
# Запуск Docker-контейнера
run:
	docker run -p $(HOST_PORT):8080 -d --name $(IMAGE_NAME) $(IMAGE_NAME):$(IMAGE_VERSION)
# Остановка и удаление Docker-контейнера
stop:
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true
# Очистка (удаление) Docker-образа
clean:
	docker rmi $(IMAGE_NAME):$(IMAGE_VERSION) || true
# Пересборка и перезапуск контейнера
rebuild: stop clean build run
.PHONY: build run stop clean rebuild