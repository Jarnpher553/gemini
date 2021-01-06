package mongo

type DocElem struct {
	Name  string
	Value interface{}
}

type M map[string]interface{}

type E DocElem

type D []DocElem

type A []interface{}
