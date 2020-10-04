all: main.wasm

run: main.wasm
	~/go/bin/goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))' &

main.wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm

clean:
	rm -f main.wasm
