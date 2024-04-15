FROM golang:1.21-alpine AS stageOne
COPY . /
WORKDIR /
RUN ls
RUN go mod download
RUN go build -o /cmd/MyShoo /cmd

FROM scratch
COPY --from=stageOne /cmd/MyShoo /cmd/
COPY --from=stageOne /config/. /config/
COPY --from=stageOne /internal/templates/. /internal/templates/
CMD [ "/cmd/MyShoo" ]

# use command example (from parent dir of this project): 
#   docker build -t shoe-mart:1.22 . 