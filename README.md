# goWeather
Go Weather Server and Client

goWeatherServer принимает GET запросы на http://loclahost:8080 с пользовательским IP адресом заданным в X-FORWARDED-FOR header'е. 
goWeatherClient посылает запросы на goWeatherServer и возвращает температуру.

В терминале:

export OPENWEATHER_API_KEY=\<Your OPENWEATHER API key\>

По умолчанию, для конвертации IP адреса в географические координаты используется https://ipapi.co/, что иногда заканчивается ошибкой RateLimited. 

Для использования https://geo.ipify.org/ для конвертации IP адреса в географические координаты, задайте IPIFY_API_KEY, например набрав в терминале export IPIFY_API_KEY=\<Your IPIFY API key\>

В терминале:

go run weatherServer.go

В новом терминале:

go run weatherClient.go
