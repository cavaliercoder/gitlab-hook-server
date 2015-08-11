APP = gitlab-hook-server

GO = go

ALL: $(APP)

$(APP): main.go config.go io.go rule.go systemhooks.go webhooks.go
	$(GO) build -x -o $(APP)

get-deps:
	$(GO) get -u code.google.com/p/go-uuid/uuid

test:
	$(GO) test -v

clean:
	$(GO) clean -x

.PHONY: all get-deps test clean
