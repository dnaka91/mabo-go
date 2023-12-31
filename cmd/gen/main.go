package main

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

const (
	arrayMax = 32
	tupleMax = 12
)

func main() {
	if err := array(); err != nil {
		panic(err)
	}

	if err := tuple(); err != nil {
		panic(err)
	}

	if err := nonZero(); err != nil {
		panic(err)
	}
}

func array() error {
	f := NewFile("buf")
	f.HeaderComment("Code generated by ./cmd/gen. DO NOT EDIT.")

	for i := 1; i <= arrayMax; i++ {
		f.Add(
			encodeArray(i), Line(),
			decodeArray(i), Line(),
		)
	}

	return f.Save("buf/array.go")
}

func encodeArray(size int) *Statement {
	return Func().Id(fmt.Sprintf("EncodeArray%v", size)).Types(
		Id("T").Any(),
	).Params(
		Id("w").Index().Byte(),
		Id("array").Index(Lit(size)).Id("T"),
		Id("encode").Func().Params(
			Index().Byte(),
			Id("T"),
		).Index().Byte(),
	).Index().Byte().Block(
		Id("w").Op("=").Id("EncodeU64").Call(
			Id("w"),
			Uint64().Call(Len(Id("array"))),
		),
		Line(),
		For(
			List(Id("_"), Id("value")).Op(":=").Range().Id("array"),
		).Block(
			Id("w").Op("=").Id("encode").Call(Id("w"), Id("value")),
		),
		Line(),
		Return(Id("w")),
	)
}

func decodeArray(size int) *Statement {
	return Func().Id(fmt.Sprintf("DecodeArray%v", size)).Types(
		Id("T").Any(),
	).Params(
		Id("r").Index().Byte(),
		Id("decode").Func().Params(
			Index().Byte(),
		).Params(
			Index().Byte(),
			Id("T"),
			Error(),
		),
	).Params(
		Index().Byte(),
		Index(Lit(size)).Id("T"),
		Error(),
	).Block(
		List(Id("r"), Id("size"), Err()).Op(":=").Id("DecodeU64").Call(Id("r")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Index(Lit(size)).Id("T").Values(), Err()),
		),
		Line(),
		Id("array").Op(":=").Index(Lit(size)).Id("T").Values(),
		For(
			Id("i").Op(":=").Lit(0),
			Id("i").Op("<").Id("int").Call(Id("size")),
			Id("i").Op("++"),
		).Block(
			List(Id("r2"), Id("value"), Err()).Op(":=").Id("decode").Call(Id("r")),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Index(Lit(size)).Id("T").Values(), Err()),
			),
			Line(),
			Id("array").Index(Id("i")).Op("=").Id("value"),
			Id("r").Op("=").Id("r2"),
		),
		Line(),
		For(
			Id("i").Op(":=").Lit(size),
			Id("i").Op("<").Id("int").Call(Id("size")),
			Id("i").Op("++"),
		).Block(
			List(Id("r2"), Id("_"), Err()).Op(":=").Id("decode").Call(Id("r")),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Index(Lit(size)).Id("T").Values(), Err()),
			),
			Id("r").Op("=").Id("r2"),
		),
		Line(),
		Return(Id("r"), Id("array"), Nil()),
	)
}

func tuple() error {
	f := NewFile("mabo")
	f.HeaderComment("Code generated by ./cmd/gen. DO NOT EDIT.")

	for i := 2; i <= tupleMax; i++ {
		f.Add(
			defineTuple(i), Line(),
		)
	}

	return f.Save("tuple.go")
}

func defineTuple(size int) *Statement {
	types := make([]Code, size)
	for i := 0; i < size; i++ {
		types[i] = Id(fmt.Sprintf("T%v", i+1)).Any()
	}

	fields := make([]Code, size)
	for i := 0; i < size; i++ {
		fields[i] = Id(fmt.Sprintf("F%v", i+1)).Id(fmt.Sprintf("T%v", i+1))
	}

	return Type().Id(fmt.Sprintf("Tuple%v", size)).Types(types...).Struct(fields...)
}

func nonZero() error {
	f := NewFile("mabo")
	f.HeaderComment("Code generated by ./cmd/gen. DO NOT EDIT.")

	f.Add(
		defineNonZeroInt("U8", Uint8(), false), Line(),
		defineNonZeroInt("U16", Uint16(), false), Line(),
		defineNonZeroInt("U32", Uint32(), false), Line(),
		defineNonZeroInt("U64", Uint64(), false), Line(),
		defineNonZeroBigInt("U128", false), Line(),
		defineNonZeroInt("I8", Int8(), true), Line(),
		defineNonZeroInt("I16", Int16(), true), Line(),
		defineNonZeroInt("I32", Int32(), true), Line(),
		defineNonZeroInt("I64", Int64(), true), Line(),
		defineNonZeroBigInt("I128", true), Line(),
		defineNonZeroCollection("String", String()), Line(),
		defineNonZeroCollection("Bytes", Index().Byte()), Line(),
		defineNonZeroCollectionGen("Vec",
			Id("T"),
			Id("T").Any(),
			Index().Id("T"),
		), Line(),
		defineNonZeroCollectionGen("HashMap",
			List(Id("K"), Id("V")),
			List(Id("K").Comparable(), Id("V").Any()),
			Map(Id("K")).Id("V"),
		), Line(),
	)

	return f.Save("non_zero.go")
}

func defineNonZeroInt(name string, ty *Statement, signed bool) *Statement {
	ident := fmt.Sprintf("NonZero%v", name)
	op := ">"
	if signed {
		op = "!="
	}

	return Type().Id(ident).Struct(
		Id("value").Add(ty),
	).
		Line().
		Func().Params(
		Id("v").Id(ident),
	).Id("Get").Params().Add(ty).Block(
		Return(
			Id("v").Dot("value"),
		),
	).
		Line().Line().
		Func().Id(fmt.Sprintf("New%v", ident)).Params(
		Id("value").Add(ty),
	).Params(Id(ident), Bool()).Block(
		If(Id("value").Op(op).Lit(0)).Block(
			Return(
				Id(ident).Values(Dict{
					Id("value"): Id("value"),
				}),
				True(),
			),
		),
		Return(Id(ident).Values(), False()),
	)
}

func defineNonZeroBigInt(name string, signed bool) *Statement {
	ident := fmt.Sprintf("NonZero%v", name)
	ty := Op("*").Qual("math/big", "Int")
	op := ">"
	if signed {
		op = "!="
	}

	return Type().Id(ident).Struct(
		Id("value").Add(ty),
	).
		Line().
		Func().Params(
		Id("v").Id(ident),
	).Id("Get").Params().Add(ty).Block(
		Return(
			Id("v").Dot("value"),
		),
	).
		Line().Line().
		Func().Id(fmt.Sprintf("New%v", ident)).Params(
		Id("value").Add(ty),
	).Params(Id(ident), Bool()).Block(
		If(Id("value").Dot("Cmp").Call(
			Qual("math/big", "NewInt").Call(Lit(0)),
		).Op(op).Lit(0)).Block(
			Return(
				Id(ident).Values(Dict{
					Id("value"): Id("value"),
				}),
				True(),
			),
		),
		Return(Id(ident).Values(), False()),
	)
}

func defineNonZeroCollection(name string, ty *Statement) *Statement {
	ident := fmt.Sprintf("NonZero%v", name)

	return Type().Id(ident).Struct(
		Id("value").Add(ty),
	).
		Line().
		Func().Params(
		Id("v").Id(ident),
	).Id("Get").Params().Add(ty).Block(
		Return(
			Id("v").Dot("value"),
		),
	).
		Line().Line().
		Func().Id(fmt.Sprintf("New%v", ident)).Params(
		Id("value").Add(ty),
	).Params(Id(ident), Bool()).Block(
		If(Len(Id("value")).Op(">").Lit(0)).Block(
			Return(
				Id(ident).Values(Dict{
					Id("value"): Id("value"),
				}),
				True(),
			),
		),
		Return(Id(ident).Values(), False()),
	)
}

func defineNonZeroCollectionGen(name string, genName, genType, ty *Statement) *Statement {
	ident := fmt.Sprintf("NonZero%v", name)

	return Type().Id(ident).Types(genType).Struct(
		Id("value").Add(ty),
	).
		Line().
		Func().Params(
		Id("v").Id(ident).Types(genName),
	).Id("Get").Params().Add(ty).Block(
		Return(
			Id("v").Dot("value"),
		),
	).
		Line().Line().
		Func().Id(fmt.Sprintf("New%v", ident)).Types(genType).Params(
		Id("value").Add(ty),
	).Params(Id(ident).Types(genName), Bool()).Block(
		If(Len(Id("value")).Op(">").Lit(0)).Block(
			Return(
				Id(ident).Types(genName).Values(Dict{
					Id("value"): Id("value"),
				}),
				True(),
			),
		),
		Return(Id(ident).Types(genName).Values(), False()),
	)
}
