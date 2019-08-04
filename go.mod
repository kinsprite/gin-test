module github.com/kinsprite/gintest

go 1.12

require (
	github.com/99designs/gqlgen v0.9.1
	github.com/elastic/go-sysinfo v1.0.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/jinzhu/gorm v1.9.10
	github.com/json-iterator/go v1.1.6
	github.com/kinsprite/gqlgen-todos v0.0.0-20190730175632-9367229f2208
	github.com/kinsprite/producttest v0.0.6
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/prometheus/procfs v0.0.2 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	github.com/vektah/gqlparser v1.1.2
	go.elastic.co/apm v1.4.0
	go.elastic.co/apm/module/apmgin v1.4.0
	go.elastic.co/apm/module/apmgrpc v1.4.0
	go.elastic.co/apm/module/apmhttp v1.4.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/sys v0.0.0-20190710143415-6ec70d6a5542 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.22.0
)

replace google.golang.org/grpc v1.22.0 => github.com/grpc/grpc-go v1.22.0
