package repo

import (
	"fmt"
	"github.com/Jarnpher553/gemini/log"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"reflect"
	"strings"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"

// repo 仓储类
type Repository struct {
	*gorm.DB
	Logger   *log.ZapLogger
	userName string
	password string
	addr     string
	host     string
	port     string
	dbName   string
	logMode  bool
}

// FieldName 字段类
type FieldName struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

// Option 配置函数
type Option func(*Repository)

// UserName 用户名配置
func UserName(userName string) Option {
	return func(repo *Repository) {
		repo.userName = userName
	}
}

func LogMode(mode bool) Option {
	return func(repo *Repository) {
		repo.logMode = mode
	}
}

// Pwd 密码配置
func Pwd(password string) Option {
	return func(repo *Repository) {
		repo.password = password
	}
}

func Addr(addr string) Option {
	return func(repo *Repository) {
		repo.addr = addr
	}
}

// Host 服务器配置
func Host(host string) Option {
	return func(repo *Repository) {
		repo.host = host
	}
}

// Port 端口配置
func Port(port string) Option {
	return func(repo *Repository) {
		repo.port = port
	}
}

// DbName 数据库名配置
func DbName(dbName string) Option {
	return func(repo *Repository) {
		repo.dbName = dbName
	}
}

// New 构造函数
func New(options ...Option) *Repository {
	repo := &Repository{
		Logger: log.Zap.Mark("repo"),
	}

	for i := range options {
		options[i](repo)
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", repo.userName, repo.password, repo.addr, repo.dbName))

	if err != nil {
		repo.Logger.Fatal(log.Message(err))
	}

	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	db.SetLogger(repo)
	db.LogMode(repo.logMode)
	repo.DB = db

	return repo
}

// Deprecated: NewFromConfigFile 通过配置文件实例化repo
/*func NewFromConfigFile(file *config.Config, fn *FieldName) *repo {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", file.GetString(fn.Username), file.GetString(fn.Password), file.GetString(fn.Host), file.GetString(fn.Port), file.GetString(fn.DbName)))

	if err != nil {
		entry.Fatalln(err)
	}

	return &repo{
		DB: db,
	}
}*/

// ReadAll 查询单条
func (repo *Repository) ReadAll(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Find(out, where...).Error
}

func (repo *Repository) ReadAllWithOrderBy(out interface{}, orderby interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Order(orderby).Find(out, where...).Error
}

// Read 查询单条记录
func (repo *Repository) Read(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.First(out, where...).Error
}

// Read 查询单条记录
func (repo *Repository) ReadWithOrderBy(out interface{}, orderby interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Order(orderby).First(out, where...).Error
}

func (repo *Repository) Exist(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	sql := repo.DB.First(out, where...)
	not := sql.RecordNotFound()
	if not || sql.Error != nil {
		return sql.Error
	} else {
		return nil
	}
}

// Remove 删除
func (repo *Repository) Remove(val interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	return repo.DB.Delete(val, where...).Error
}

// Remove 获取影响行数
func (repo *Repository) RemoveWithAffect(val interface{}, where ...interface{}) (affects int64, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	ret := repo.DB.Delete(val, where...)
	return ret.RowsAffected, ret.Error
}

// Insert 新增
func (repo *Repository) Insert(val interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Create(val).Error
}

func (repo *Repository) SoftRemove(value interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	columns := make(map[string]interface{})
	columns["is_active"] = false

	if where != nil {
		return repo.DB.Model(value).Where(where[0], where[1:]...).UpdateColumns(columns).Error
	}
	return repo.DB.Model(value).UpdateColumns(columns).Error
}

// Modify 修改
func (repo *Repository) Modify(val interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Save(val).Error
}

// 更改单个字段
func (repo *Repository) ModifyColumn(val interface{}, attr string, upValue interface{}, where ...interface{}) (affects int64, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	kind := reflect.TypeOf(val).Kind()

	if where != nil {
		if kind == reflect.String {
			db := repo.DB.Table(val.(string)).Where(where[0], where[1:]...).Update(attr, upValue)
			return db.RowsAffected, db.Error
		} else {
			db := repo.DB.Model(val).Where(where[0], where[1:]...).Update(attr, upValue)
			return db.RowsAffected, db.Error
		}
	}
	if kind == reflect.String {
		db := repo.DB.Table(val.(string)).Update(attr, upValue)
		return db.RowsAffected, db.Error
	} else {
		db := repo.DB.Model(val).Update(attr, upValue)
		return db.RowsAffected, db.Error
	}
}

// 更改多个字段
func (repo *Repository) ModifyColumns(val interface{}, columns interface{}, where ...interface{}) (affects int64, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	kind := reflect.TypeOf(val).Kind()

	if where != nil {
		if kind == reflect.String {
			db := repo.DB.Table(val.(string)).Where(where[0], where[1:]...).Updates(columns)
			return db.RowsAffected, db.Error
		} else {
			db := repo.DB.Model(val).Where(where[0], where[1:]...).Updates(columns)
			return db.RowsAffected, db.Error
		}
	}
	if kind == reflect.String {
		db := repo.DB.Table(val.(string)).Updates(columns)
		return db.RowsAffected, db.Error
	} else {
		db := repo.DB.Model(val).Updates(columns)
		return db.RowsAffected, db.Error
	}
}

// ModifyFunc 使用函数更新
func (repo *Repository) ModifyFunc(val interface{}, modifier func(interface{}), where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	valNew := reflect.New(t.Elem())
	valNew.Elem().Set(v.Elem())

	i := valNew.Interface()

	modifier(i)

	if where != nil {
		return repo.DB.Model(val).Where(where[0], where[1:]...).Updates(i).Error
	}
	return repo.DB.Model(val).Updates(i).Error
}

// Transaction 执行包装事务
func (repo *Repository) Transaction(f func(*Repository) error) (e error) {
	repoTx := repo.begin()

	defer func() {
		if r := recover(); r != nil {
			repoTx.Rollback()
			e = fmt.Errorf("%v", r)
		}
	}()

	if err := repoTx.Error; err != nil {
		return err
	}

	if err := f(repoTx); err != nil {
		repoTx.Rollback()
		return err
	}

	if err := repoTx.Commit().Error; err != nil {
		repoTx.Rollback()
		return err
	}
	return nil
}

// begin 开始一个事务
func (repo *Repository) begin() *Repository {
	//开始一个事务
	tx := repo.DB.Begin()

	//返回一个包含事务的repo
	return &Repository{
		userName: repo.userName,
		password: repo.password,
		host:     repo.host,
		addr:     repo.addr,
		port:     repo.port,
		dbName:   repo.dbName,
		DB:       tx,
		Logger:   repo.Logger,
	}
}

type Expression func(db *gorm.DB) *gorm.DB

func Page(pageNum int, perCount int) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((pageNum - 1) * perCount).Limit(perCount)
	}
}

func Model(value interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(value)
	}
}

func Table(name string) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(name)
	}
}

func Select(query interface{}, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	}
}

func Order(value interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

func Join(query string, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

func Where(query interface{}, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func Group(query string) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(query)
	}
}

// Query 查询列表（可分页）
func (repo *Repository) Query(out interface{}, count bool, exps ...Expression) (c int, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	db := repo.DB

	for _, exp := range exps {
		db = exp(db)
	}

	db = db.Scan(out)
	if err := db.Error; err != nil {
		e = err
		return
	}

	if count {
		if err := db.
			Offset(-1).
			Limit(-1).
			Count(&c).Error; err != nil {
			e = err
			return
		}
	}

	return
}

func Expr(expression string, args ...interface{}) interface{} {
	expr := gorm.Expr(expression, args...)
	return expr
}

func Escape(input string) string {
	output := strings.Replace(input, `\`, `\\`, -1)
	output = strings.Replace(output, `%`, `\%`, -1)
	output = strings.Replace(output, `'`, `\'`, -1)
	output = strings.Replace(output, `_`, `\_`, -1)
	output = strings.Replace(output, `"`, `\"`, -1)
	return output
}

func (repo *Repository) Print(args ...interface{}) {
	gorm.LogFormatter()
	formatter := gormLogFormatter(args...)
	source := strings.Split(formatter[0].(string), "/")
	l := repo.Logger.With(zap.String("source", strings.Join(source[len(source)-2:], "/")))
	if args[0] == "sql" {
		l.
			With(zap.String("cost", formatter[2].(string))).
			With(zap.String("sql", formatter[3].(string))).
			Info(formatter[4].(string))
	} else {
		l.
			Error(formatter[2].(*mysql.MySQLError).Message)
	}
}

func (repo *Repository) Migrate(initial func(*Repository), values ...interface{}) {
	repo.DB.AutoMigrate(values...)

	if initial != nil {
		initial(repo)
	}
}
