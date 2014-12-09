package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

	buildScratchTemplate = `
FROM scratch
ENTRYPOINT ["/{{.Entrypoint}}"]

# Add the binary
ADD {{.Entrypoint}} /
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
	fromTheScratch := flag.Bool("scratch", false, "Build the from the base image scratch")

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

	if *fromTheScratch {
		buildForLinux(dockerInfo)
		generateDockerfileFromScratch(dockerInfo)
	} else {
		generateDockerfile(dockerInfo)
	}

}

func generateDockerfile(dockerInfo DockerInfo) {
	t := template.Must(template.New("buildTemplate").Parse(buildTemplate))

	f, err := os.Create("Dockerfile")
	if err != nil {
		log.Fatal("Error wrinting Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, dockerInfo)

	fmt.Printf("Dockerfile generated, you can build the image with: \n")
	fmt.Printf("$ docker build -t %s .\n", dockerInfo.Entrypoint)
}

func generateDockerfileFromScratch(dockerInfo DockerInfo) {
	t := template.Must(template.New("buildScratchTemplate").Parse(buildScratchTemplate))

	f, err := os.Create("Dockerfile")
	if err != nil {
		log.Fatal("Error writing Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, dockerInfo)

	fmt.Printf("Dockerfile from the scratch generated, you can build the image with: \n")
	fmt.Printf("$ docker build -t %s .\n", dockerInfo.Entrypoint)

}

func buildForLinux(dockerInfo DockerInfo) {
	os.Setenv("GOOS", "linux")
	os.Setenv("CGO_ENABLED", "0")
	cmd := exec.Command("go", "build", "-o", dockerInfo.Entrypoint)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("%s", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("%s", err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%s", err)
	}
	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("%s", errBuf)
	}

}
