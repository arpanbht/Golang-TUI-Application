build:
	@go build -o tuiapp.exe .

run: build
	@./tuiapp.exe

