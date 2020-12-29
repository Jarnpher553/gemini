package repo

type Repository interface {
	Read(interface{}, ...interface{}) error
	ReadAll(interface{}, ...interface{}) error
	Insert(interface{}) error
	Remove(interface{}, ...interface{}) error
	Modify(val interface{}) error
}
