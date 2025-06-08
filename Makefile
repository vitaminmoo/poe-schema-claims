
all: api/server.gen.go

api.json: api.jsonnet
	jsonnet api.jsonnet > api.json

api/server.gen.go: api.json cfg.yaml
	go tool oapi-codegen -config cfg.yaml api.json
