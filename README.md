# Golang

Instructions:
Download the package and compress. Move the unzipped file into your Go working directory(.../GoPath/src/).

Folder
1. byt            : Byte Buffer (source code).
2. ws             : WebSocket (client,server,session) (source code).
3. example        : Network Communication with Websocket and Buffer.
4. github.com.zip : websocket dependencies files.

Tips:
Open the example folder and modify the Host parameters of c/client.go file and s/server.go file.
The Host parameters is the IP4 address of your PC.
Run client.go and server.go to test respectively (Running these two files in the vscode editor or using GitBash).

Dependency:
github.com/gorilla/websocket (If you use websocket).
Using ByteBuffer does not require any dependencies.

Command: 
go run client.go / go run server.go.
