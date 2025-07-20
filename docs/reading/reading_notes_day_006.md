# Docker for Containerization

## Docker Images for Go

### Steps to Containerize Go Projects using Docker
1. Create a Dockerfile
        - start with optional parser directive to set the grammar rules
        - add the base image that the application will use
        ```
        # syntax=docker/dockerfile:1

        FROM golang:1.19
        ```
    2. Create a working directory inside the image
        - this tells docker the default destination for all commands after this
        - all relative paths will be with respoect to this dir
        ```
        WORKDIR /app
        ```
    3. Copy go.mod and go.sum to the conatainer
        - before we can do go mod download to get project from repo
        - we need to add the go.mod and go.sum so it knows what to install
        ```
        COPY go.mod go.sum ./
        ```
    4. Run the go mod download command
        - this will download the dependencies from the directory specified in the go.mod
        ```
        RUN go mod download
        ```
    5. Copy the source code to the image
        ```
        COPY *.go ./
        ```
    6. Compile the application
        ```
        RUN CGO_ENABLED=0 GOOS=linux go build -o /your-app
        ```
    7. Run the application
        - runs the binary using the go runtime
        ```
        CMD ["/your-app"]
2. Build the image
    - use the docker build --tag name-of-image .
    - reads the dockerfile from directory mentioned, here it is "."
    - tags it with a human readabble name
3. View image
    - docker image ls

4. Tag images
    - by default the tag is "latest"
    - to make a specific version of the image that can be called using image-name:v1.1
    - docker iamge tag your-docker-image:latest docker-image-name:v1.1

5. Remove images
    - docker image rm docker-image-name:v1.1
        - if another versionof this image exists with the same image id 
          it will only remove the tag not the image

6. Multi stage builds
    - currently the entire go toolchain will be present in the app even after building
    - to solve this we use multi stage bulids
    - Dockerfile.multistage
    ```
    FROM goland:1.19 AS build-stage

    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download

    COPY *.go/
    
    RUN CGO_ENABLED=0 GOOS=linux go build -o /your-app

    FROM build-stage AS run-test-stage
    RUN go test -v ./...

    FROM gcr.io/distroless/base-debian11 AS build-release-stage

    WORKDIR /
    
    COPY --from=build-stage /your-app /your-app

    EXPOSE 8080

    USER nonroot:nonroot

    ENTRYPOINT ["/your-app"]
    ```

    - this way the final image will only contain the files that were built using 
      the build image and then we can use the prod lightweight image and just copy
      over the contents from the previous stage


## Docker Compose
- tool for defining and running multi-container applications
- manage services, networks, volumes using YAML
compose.YAML
```
services:
    web:
        build: .
        ports:
            - "8080:8080"
    redis:
        image: "redis-alpine"
```
- the web service uses an image built from a dockerfile
- the redis service uses the image base from the public redis image on dockerhub

- docker compose up
- pulls redis image and builds the image from code

```
develop:
    watch:
       - action: sync
         path: .
         target: /code
```
- this can be used to watch for code changes
- when a file is changed compose syncs the file to the corresponding location
  under /code
- this way we dont have to restart the application to update it
- this works only for interactive shell run programs (python, perl, bash)

- Split up services
    - we can add other yaml files using
    ```
    include:
        - infra.yaml
    ```
    where include is a yaml list
