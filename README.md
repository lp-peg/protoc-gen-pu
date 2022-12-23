# protoc-gen-pu

protoc plugin to generate a [plant UML](https://plantuml.com/) class diagram.

## prerequisites
- go ~> v1.19
- protoc ~> 3.15


## install

```sh
go install github.com/lp-peg/protoc-gen-pu@latest
```

## usage

```sh
protoc --pu_out=. testdata/pet.proto
```

then, you can see diagram as below:

<p align="center">
  <img src="https://user-images.githubusercontent.com/35035802/209184957-62704129-7f6b-4738-98c7-e91617c7f9b9.png" />
</p>

## options

You can use these plugin options.  
Pass option via `--pu_opt` flag.

| option key | description      | default  |
| ---------- | ---------------- | -------- |
| out        | output file name | out.pu   |
| skinparams | [skin param](https://plantuml.com/en/skinparam). Required format is `<param>:<value>`. You can pass this opt multiple times. | `"linetype:ortho"` |
| circle     | `hide` or `show` circle | `hide` |

**example:**
```sh
❯ protoc -Itestdata \
  --pu_opt=out=pet.pu \
  --pu_opt=skinparams=linetype:ortho \
  --pu_opt=skinparams=classFontSize:10 \
  --pu_opt=circle=show \
  --pu_out=. \
  testdata/pet.proto
```

```
@startuml


skinparam linetype ortho
skinparam classFontSize 10

entity Pet {
  name (STRING)
  age (INT32)
  animal (Animal<FK>)
}

entity Animal {
  name (STRING)
  category (AnimalCategory<FK>)
}

entity AnimalCategory {
  UNSPECIFIED (ENUM: 0)
  DOG (ENUM: 1)
  CAT (ENUM: 2)
  COW (ENUM: 3)
  BEAR (ENUM: 4)
}

Animal <-- Pet
AnimalCategory <-- Animal

@enduml
```