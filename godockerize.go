package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	buildTemplate = `
FROM google/golang:stable
# Godep for vendoring
RUN go get github.com/tools/godep
# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

MAINTAINER dahernan@gmail.com
ENV APP_DIR $GOPATH{{.Appdir}}
 
# Set the entrypoint 
ENTRYPOINT ["/opt/app/{{.Entrypoint}}"]
ADD . $APP_DIR

# Compile the binary and statically link
RUN mkdir /opt/app
RUN cd $APP_DIR && godep restore
RUN cd $APP_DIR && CGO_ENABLED=0 go build -o /opt/app/{{.Entrypoint}} -ldflags '-d -w -s'

EXPOSE {{.Expose}}
`
)

type DockerInfo struct {
	Appdir     string
	Entrypoint string
	Expose     string
}

func main() {
	expose := flag.String("expose", "3000", "Port to expose in docker")

	flag.Parse()

	goPath := os.Getenv("GOPATH")
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	appdir := strings.Replace(dir, goPath, "", 1)

	_, entrypoint := path.Split(appdir)

	dockerInfo := DockerInfo{
		Appdir:     appdir,
		Entrypoint: entrypoint,
		Expose:     *expose,
	}

	generateDockerfile(dockerInfo)
}

func generateDockerfile(dockerInfo DockerInfo) {
	t := template.Must(template.New("buildTemplate").Parse(buildTemplate))

	f, err := os.Create("Dockerfile")
	if err != nil {
		fmt.Printf("Error wrinting Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, dockerInfo)

	fmt.Printf("Dockerfile generated, you can build the image with: \n")
	fmt.Printf("$ docker build -t %s .\n", dockerInfo.Entrypoint)
}
