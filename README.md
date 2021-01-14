# goWeather
Go Weather Server and Client

goWeatherServer принимает GET запросы на http://loclahost:8080 с пользовательским IP адресом заданным в X-FORWARDED-FOR header'е. 
goWeatherClient посылает запросы на goWeatherServer и возвращает температуру.

В терминале:
go run weatherServer.go

В новом терминале:
go run weatherClient.go
