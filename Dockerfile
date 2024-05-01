FROM node:16 as react-build
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.16 as go-build
WORKDIR /go/src/app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/docusinger

FROM python:3.9-slim as python-env
WORKDIR /python
COPY python/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY python/ ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=react-build /app/build ./public
COPY --from=go-build /go/bin/docusinger .
COPY --from=python-env /python ./
EXPOSE 8080
CMD ["./docusinger"]
