# Переменные
REGISTRY_ID := crpbhnv13e7950b9a7hh
IMAGE_NAME := scavenger
TAG := latest
IMAGE := cr.yandex/$(REGISTRY_ID)/$(IMAGE_NAME):$(TAG)
CONTAINER_NAME := scavenger-container

.PHONY: build push deploy

# Сборка образа
build:
	docker build -t $(IMAGE) .

# Загрузка в Container Registry
push: build
	docker push $(IMAGE)

# Деплой (создаёт новую ревизию)
deploy: push
	yc serverless container revision deploy \
		--container-name $(CONTAINER_NAME) \
		--image $(IMAGE) \
		--memory 128M \
		--execution-timeout 5s \
		--concurrency 1 \
		--environment PORT=8080

# Создание контейнера (если ещё не создан)
create-container:
	yc serverless container create --name $(CONTAINER_NAME)

# Первичный деплой (создаёт контейнер и ревизию)
first-deploy: push
	yc serverless container revision deploy \
		--container-name $(CONTAINER_NAME) \
		--image $(IMAGE) \
		--memory 128M \
		--execution-timeout 5s \
		--concurrency 1 \
		--environment PORT=8080 || \
	(make create-container && make deploy)
