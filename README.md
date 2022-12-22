# protoc-gen-pu

protoc plugin to generate [plant UML](https://plantuml.com/) class diagram.

## prerequisites
- go ~> v1.19
- protoc ~> 3.15


## install

```sh
go install github.com/lp-peg/protoc-gen-pu
```

## usage

```sh
protoc --pu_out=. testdata/
```

then, you can see diagram as below:

