clean:
	rm -f /home/aidan/code/go/rt-flow/plugins/example.so

build:
	go build  -buildmode=plugin -o /home/aidan/code/go/rt-flow/plugins/example.so ./test/test.go

all:
	rm -f ./plugins/example.so
	go build  -buildmode=plugin -o /home/aidan/code/go/rt-flow/plugins/example.so nodes.go


run:
	go run main.go
