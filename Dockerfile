FROM golang:1.21-alpine AS stageOne
ADD . /
WORKDIR /
RUN go mod download
RUN go build -o /cmd/MyShoo /cmd

FROM scratch
COPY --from=stageOne /cmd/MyShoo /cmd/
COPY --from=stageOne /config/. /config/
COPY --from=stageOne /internal/templates/. /internal/templates/
COPY --from=stageOne /.env /
CMD [ "/cmd/MyShoo" ]
