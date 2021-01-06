package uuid

import guid "github.com/satori/go.uuid"

type GUID string

func New() GUID {
	return GUID(guid.NewV4().String())
}

func IsGUID(str string) error {
	_, err := guid.FromString(str)
	return err
}
