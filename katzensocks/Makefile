warped?=true
CGO_CFLAGS_ALLOW="-DPARAMS=sphincs-shake-256f"
ldflags="-buildid= -X github.com/katzenpost/katzenpost/core/epochtime.WarpedEpoch=${warped} -X github.com/katzenpost/katzenpost/server/internal/pki.WarpedEpoch=${warped} -X github.com/katzenpost/katzenpost/minclient/pki.WarpedEpoch=${warped}"
uid=$(shell [ "$$SUDO_UID" != "" ] && echo "$$SUDO_UID" || id -u)
gid=$(shell [ "$$SUDO_GID" != "" ] && echo "$$SUDO_GID" || id -g)
docker_user?=$(shell if echo ${docker}|grep -q podman; then echo 0:0; else echo ${uid}:${gid}; fi)
docker=$(shell if which podman|grep -q .; then echo podman; else echo docker; fi)
distro=alpine
image=katzenpost-$(distro)_go_mod
docker_args=--user ${docker_user} -v $(shell readlink -f ..):/go/katzenpost -e GOCACHE=/tmp/gocache --network=host --rm

dockerdockertest:
	$(docker) run ${docker_args} $(image) \
		sh -c 'cd /go/katzenpost/katzensocks/client/; GORACE=history_size=7 CGO_CFLAGS_ALLOW=-DPARAMS=sphincs-shake-256f go test $(testargs) -ldflags ${ldflags} -tags=docker_test -race -v -timeout 10m -run Docker'
test:
	cd client/; GORACE=history_size=7 CGO_CFLAGS_ALLOW=-DPARAMS=sphincs-shake-256f go test $(testargs) -ldflags ${ldflags} -tags=docker_test -race -v -timeout 10m -run Docker

server/cmd/katzensocks/katzensocks:
	cd server/cmd/server && CGO_CFLAGS_ALLOW=-DPARAMS=sphincs-shake-256f go build -trimpath -ldflags ${ldflags}
client/cmd/client/client:
	cd client/cmd/client && CGO_CFLAGS_ALLOW=-DPARAMS=sphincs-shake-256f go build -trimpath -ldflags ${ldflags}
clean:
	rm -f client/cmd/client/client server/cmd/server/server
