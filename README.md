# constanta-test


## Сборка и запуск
```
docker build -t constanta-test .
docker run -p 7000:7000 constanta-test
```
## Пример запроса
На вход в ручку `/v1/collect-data` ожидается POST-запрос с массивом URL
```
curl -d '{"urls": ["https://reqres.in/api/user/1", "https://reqres.in/api/user/2"]}' -H "Content-Type: application/json" -X POST http://localhost:7000/v1/collect-data
```
## Ответ
Ответ в случае успешного выполнения запроса:
```json
[
  {
    "resource": "https://reqres.in/api/user/1",
    "data": "{\"data\":{\"id\":1,\"name\":\"cerulean\",\"year\":2000,\"color\":\"#98B2D1\",\"pantone_value\":\"15-4020\"},\"support\":{\"url\":\"https://reqres.in/#support-heading\",\"text\":\"To keep ReqRes free, contributions towards server costs are appreciated!\"}}"
  },
  {
    "resource": "https://reqres.in/api/user/2",
    "data": "{\"data\":{\"id\":2,\"name\":\"fuchsia rose\",\"year\":2001,\"color\":\"#C74375\",\"pantone_value\":\"17-2031\"},\"support\":{\"url\":\"https://reqres.in/#support-heading\",\"text\":\"To keep ReqRes free, contributions towards server costs are appreciated!\"}}"
  }
]
```
### Ответ в случае невалидных входных данных или других ошибок.
Запрос:
```
curl -d '{"urls": ["google.com", "http//vk.com"]}' -H "Content-Type: application/json" -X POST http://localhost:7000/v1/collect-data
```
````json
{
    "error": "url google.com not valid: parse \"google.com\": invalid URI for request; url http//vk.com not valid: parse \"http//vk.com\": invalid URI for request"
}
````
## Выполненные ограничения
-  [x] для реализации задачи следует использовать Go 1.13 или выше
-  [x] использовать можно только компоненты стандартной библиотеки Go
-  [x] сервер не должен принимать запрос если количество url в в нем больше 20
-  [x] таймаут на запрос одного url - одна секунда
-  [x] в репозитории должен быть приложен Dockerfile, который позволяет собрать и запустить сервис
-  [x] пожалуйста, добавь README с описанием того, как сервис собирать и запускать (с использованием docker) и какого формата запрос он принимает, какой ответ нам ожидать
-  [x] будет супер, если ты добавишь минимальные комментарии в код, которые объяснят, что где происходит

Со звездочкой (хорошо бы сделать, но не обязательно):

- [x] для каждого входящего запроса должно быть не больше 4 **одновременных** исходящих
- [x] сервер не должен обслуживать больше чем 100 одновременных входящих http-запросов
- [ ] обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом 
- [x] сервис должен поддерживать 'graceful shutdown'
