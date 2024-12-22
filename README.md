# Сервис подсчёта арифметических выражений

## Установка и запуск
1. Клонировать репозиторий:
    ```bash
    git clone https://github.com/karrless/yaly-1
    ```
2. Переместиться в клонированные репозиторий:
    ```bash
    cd yaly-1
    ```
3. Запустить программу:
   1. На стандартном порту `8080`:
        ```bash
        go run cmd/main/main.go
        ```
    2. На кастомном порту (незабудьте изменить его при curl запросах):
        ```bash
        go run cmd/main/main.go [ваш порт]
        ```

## Возможности

## Примеры
1. Примеры с нормальной работой:

    Для Windows:
    ```bash
    curl "localhost:8080/api/v1/calculate" `
    --header "Content-Type: application/json" `
    --data '{"expression":"2+2*(2+3)"}'
    ```

    Для Linux:
    ```bash
    curl "localhost:8080/api/v1/calculate" \
    --header "Content-Type: application/json" \
    --data '{"expression":"2.12 * -2"}'
    ```

2. Пример ошибочных запросов, выдающих `422 Unprocessable Entity`:
    
    Для Windows:
    ```bash
    curl "localhost:8080/api/v1/calculate" `
    --header "Content-Type: application/json" `
    --data '{"expression":"((2+2**2)"}'
    ```

    Для Linux:
    ```bash
    curl "localhost:8080/api/v1/calculate" \
    --header "Content-Type: application/json" \
    --data '{"expression":"2,2*3,3 * - 2"}'
    ```

3. `500 Internal Server Error` возращается, например, при числах, превышающих лимиты float64:

    Для Windows:
    ```bash
    curl "localhost:8080/api/v1/calculate" `
    --header "Content-Type: application/json" `
    --data '{"expression":"1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035375169860499105765512820762454900903893289440758685084551339423045832369032229481658085593321233482747978262041447231687381771809192998812504040261841248583600+1"}'
    ```

    Для Linux:
    ```bash
    curl "localhost:8080/api/v1/calculate" \
    --header "Content-Type: application/json" \
    --data '{"expression":"1797693134862315708145274237317043567980705675258449965989174768031572607800285387605895586327668781715404589535143824642343213268894641827684675467035375169860499105765512820762454900903893289440758685084551339423045832369032229481658085593321233482747978262041447231687381771809192998812504040261841248583600+1"}'
    ```

## Заметки
- Помимо кодов, о которых говорилось в задании: 200, 422, 500, - также присутсвуют два необходимых кода: `400 Bad Request` и `405 StatusMethodNotAllowed`.
- Так как в задании были указаны только 2 возможные ошибки: `Expression is not valid` и `Internal server error`, - то распространнёная ошибка `Divided by zero`, т.е. деление на 0, интерпретируется как `Expression is not valid`