# crypto-server

## external libraries

- Used goralla (mux and websockets) library to work with rest apis and websockets over tcp.
- Inclued dependencies in the go.mod file, project should build successfully. [I can share the build file if you dont have access to mux libraries]

## how to run?

- `make run` should run the code or use `go run main.go` for the project root directory.
- server would be running on the `localhost` unless modified in main.go. I have declated `host` and `port` as constants at the start of the file on purpose, just so it would easy to test with.
- Alternatively, one could use vscode and launch in debug mode.
- once, the server is up and running `curl http://localhost:8080/currency/all` and `curl http://localhost:8080/currency/BTCUSDT`
- Alternatively, we can also use a web browser client or postman client.

## anything that cloud be added give the time?

- I would add a few more test cases for sockets as well.
- I also had TODOs in the code, so they might also need to be addressed.
