package stablehash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastStableHasher_DoubleChild(t *testing.T) {
	hasher := NewFastStableHasher()

	root := Address{}.Root()
	hasher.Write(root.Child(1), nil)
	hasher.Write(root.Child(1), nil)

	out := hasher.Finish()
	assert.Equal(t, "261232071512772414229682083989926651266", out.String())
}

func TestFastStableHasher_SingleStructSingleField(t *testing.T) {
	assert.Equal(t, "102568403942768160221811810082933398928", FastStableHash(&One[U32]{One: U32(5)}).String())
}

func TestFastStableHasher_AddOptionalField(t *testing.T) {
	one := &One[U32]{One: U32(5)}
	two := &TwoOptional{One: U32(5), Two: None[U32]()}
	tuple := &Tuple2[*One[U32], *TwoOptional]{One: one, Two: two}

	assert.Equal(t, "210303380251691017811466509002544125279", FastStableHash(tuple).String())
}

func TestFastStableHasher_AddDefaultField(t *testing.T) {
	one := &One[String]{One: String("one")}
	two := &Two[String]{One: String("one"), Two: String("")}
	tuple := &Tuple2[*One[String], *Two[String]]{One: one, Two: two}

	assert.Equal(t, "337538645577122176555714212704832450090", FastStableHash(tuple).String())
}

type One[T StableHashable] struct {
	One T
}

func (o *One[T]) StableHash(addr FieldAddress, hasher StableHasher) {
	o.One.StableHash(addr.Child(0), hasher)
}

type Two[T StableHashable] struct {
	One T
	Two T
}

func (o *Two[T]) StableHash(addr FieldAddress, hasher StableHasher) {
	o.One.StableHash(addr.Child(0), hasher)
	o.Two.StableHash(addr.Child(1), hasher)
}

type TwoOptional struct {
	One U32
	Two Optional[U32]
}

func (o *TwoOptional) StableHash(addr FieldAddress, hasher StableHasher) {
	o.One.StableHash(addr.Child(0), hasher)
	o.Two.StableHash(addr.Child(1), hasher)
}

type Tuple2[T1 StableHashable, T2 StableHashable] struct {
	One T1
	Two T2
}

func (o *Tuple2[T1, T2]) StableHash(addr FieldAddress, hasher StableHasher) {
	o.One.StableHash(addr.Child(0), hasher)
	o.Two.StableHash(addr.Child(1), hasher)
}
