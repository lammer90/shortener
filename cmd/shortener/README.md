# cmd/shortener

postgresql://localhost:5432/plotnikov?user=postgres&password=1234

В данной директории будет содержаться код, который скомпилируется в бинарное приложение

Определить параметры сборки:
go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit=release_20_5 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/shortener/main.go