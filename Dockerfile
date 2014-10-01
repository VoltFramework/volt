FROM golang:onbuild
ENTRYPOINT ["go-wrapper", "run"]
EXPOSE 8080
