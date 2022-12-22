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
protoc --pu_out=. testdata/pet.proto
```

then, you can see diagram as below:

![out](https://user-images.githubusercontent.com/35035802/209184957-62704129-7f6b-4738-98c7-e91617c7f9b9.png)
