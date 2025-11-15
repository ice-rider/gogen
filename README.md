# GoGen - Генератор кода для Go проектов

Мощный CLI-инструмент для генерации компонентов Clean Architecture в Go проектах.

## 🚀 Возможности

- ✅ Генерация сущностей (Entities)
- ✅ Генерация репозиториев (Repository pattern)
- ✅ Генерация use cases (Business logic)
- ✅ Автоматическая генерация моков для тестирования
- ✅ Генерация unit-тестов
- ✅ Интерактивный режим для детальной настройки
- ✅ Поддержка PostgreSQL, MySQL, SQLite, MongoDB
- ✅ Каскадная система конфигурации
- ✅ Поддержка кастомных шаблонов

## 📦 Установка

### Из исходников
```shell
git clone https://github.com/ice-rider/gogen.git
cd gogen
make install
```

### Из бинарника
С помощью Go: 
`go install github.com/ice-rider/gogen/cmd/gogen@latest`

# 🎯 Быстрый старт
## Инициализация проекта
```shell
cd your-go-project
gogen init
```
  
  
## Простая генерация

```shell
Создать сущность User
gogen -d User

# Создать сущность и репозиторий
gogen -d User -r User

# Создать сущность, репозиторий и use case
gogen -d User -r User -uc CreateUser

# С тестами и моками
gogen -d User -r User -uc CreateUser -t -m
```

## Пакетная генерация

### Создать несколько компонентов за раз
```shell
gogen -d User -d Product -d Order \
      -r User -r Product -r Order \
      -uc CreateUser -uc CreateOrder \
      -t -m
```
## Интерактивный режим
### Базовый интерактивный режим
```shell
gogen -d User --interactive
```

### Полностью интерактивный режим
```shell
gogen interactive
```

# 📖 Примеры использования  

Создание CRUD для сущности:  

```shell
gogen -d Product \
      -r Product \
      -uc CreateProduct \
      -uc GetProduct \
      -uc UpdateProduct \
      -uc DeleteProduct \
      -t -m
```
С полями сущности:  
```shell
gogen -d User:Name:string,Email:string:required,Age:int
```
Dry-run (предпросмотр)
```shell
gogen -d User -r User -uc CreateUser --dry-run
```


# ⚙️ Конфигурация
Создайте `gogen.yaml` в корне проекта:
```yaml
version: "1.0"

paths:
  domain: "internal/domain"
  repository: "internal/adapters/repository"
  usecase: "internal/core/usecases"

naming:
  style: "snake_case"
  suffixes:
    repository: "Repo"

templates:
  entity: "my_templates/entity.tmpl"
```

Также вы можете изменить глобальный конфиг global.yaml, а также глобальные .templ файлы. Они будут находится около exe файла.

# 📁 Структура проекта
your-project/  
├── internal/  
│   ├── domain/           # Сущности и интерфейсы  
│   ├── repository/       # Реализации репозиториев  
│   ├── usecase/          # Бизнес-логика  
│   └── mocks/            # Моки для тестов  
├── gogen.yaml            # Конфигурация (опционально)  
└── go.mod  

# 🛠️ Разработка
### Клонировать репозиторий
```shell
git clone https://github.com/yourname/gogen.git
cd gogen
```

### Установить зависимости
```shell
go mod download
```

### Запустить тесты
```shell
make test
```

### Собрать
```shell
make build
```

### Запустить линтер
```shell
make lint
```

# 📝 Лицензия  
MIT License - see [LICENSE](https://github.com/ice-rider/gogen/blob/main/LICENSE) file for details. 

# 🤝 Вклад
Contributions are welcome! Please feel free to submit a Pull Request.