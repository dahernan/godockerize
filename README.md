# godockerize

I was tired to write the same Dockerfile again and again, so I did a
simple tool to generate, a Dockerfile for your Go microservice.

The Dockerfile is based in the golang image, it follows these steps.
* Based in golang:stable
* Use [godep](https://github.com/tools/godep) for vendoring the dependencies. (if you are not using godep, It will break).
* Uses the project name and the root directory as a ENTRYPOINT
* Restores the dependencies via 'godep restore'
* Compiles the project
* Exposes the port

The Dockerfile is inspired by this two blog post.
* http://blog.charmes.net/2014/11/release-go-code-and-others-via-docker.html?spref=tw
* https://blog.golang.org/docker

## Install
```
$ go get github.com/dahernan/godockerize

```

## Usage Example
```
# Go to the root directory of your project, for example
$ cd $GOPATH/src/github.com/dahernan/gopherscraper


# Run godockerize exposing the port 3001
$ godockerize -expose 3001

Dockerfile generated, you can build the image with:
$ docker build -t gopherscraper .

```

## Build from the scratch image instead of using golang image

!Warning Only tested in Mac, you have to enable cross compilation

If you use the flag 'scratch', the script is going to:

1) Crosscompile the project for linux
2) Generate a Dockerfile ready to build, from the scratch

```
# Go to the root directory of your project, for example
$ cd $GOPATH/src/github.com/dahernan/gopherscraper

# Run godockerize using -scratch
$ godockerize -expose 3001 -scratch

Dockerfile generated, you can build the image with:
$ docker build -t gopherscraper .

```


## How to enable cross compile on Mac to build Linux executable

Follow this steps

```
$ cd /usr/local/go/src

$ sudo GOOS=linux CGO_ENABLED=0 ./make.bash --no-clean

```

You can try it doing for example:

```
$ GOOS=linux CGO_ENABLED=0 go build -o gopherscraper

$ file gopherscraper
gopherscraper: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped

```


