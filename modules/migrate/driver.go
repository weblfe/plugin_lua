package migrate

import (
	"github.com/golang-migrate/migrate/v4/source"
	"io"
	"io/fs"
	"runtime"
	"strings"
	"sync"
)

type (

	luaScriptDriverImpl struct {
		Schema   string
		FsDriver FileSystem
		curFd    fs.File
		safe     sync.RWMutex
	}

	FileSystem interface {
		fs.ReadFileFS
		fs.ReadDirFS
	}

)

const (
	schema       = "lua"
	DriverSchema = "lua://"
)

func init() {
	source.Register(schema, NewDriver())
}

func NewDriver() *luaScriptDriverImpl {
	var driver = new(luaScriptDriverImpl)
	return driver.init()
}

func (driver *luaScriptDriverImpl) init() *luaScriptDriverImpl {
	driver.curFd = nil
	driver.Schema = DriverSchema
	driver.safe = sync.RWMutex{}
	return driver
}

func (driver *luaScriptDriverImpl) RegisterDriver(fs FileSystem) *luaScriptDriverImpl {
	if driver.FsDriver == nil {
		driver.FsDriver = fs
	}
	return driver
}

func (driver *luaScriptDriverImpl) Open(url string) (source.Driver, error) {
	var (
		file    = driver.file(url)
		fd, err = driver.FsDriver.Open(file)
	)
	if err != nil {
		return driver, err
	}
	var d = driver.newFs(fd)
	// 加载文件
	if impl, ok := d.(*luaScriptDriverImpl); ok {
		if err = impl.load(); err != nil {
			return d, err
		}
	}
	return d, nil
}

func (driver *luaScriptDriverImpl) newFs(fd fs.File) source.Driver {
	var d = driver.setFd(fd)
	if d != nil {
		runtime.SetFinalizer(d, (*luaScriptDriverImpl).destroy)
	}
	return d
}

func (driver *luaScriptDriverImpl) setFd(fd fs.File) *luaScriptDriverImpl {
	driver.safe.Lock()
	defer driver.safe.Unlock()
	if driver.curFd == nil {
		driver.curFd = fd
		return driver
	}
	var d = NewDriver()
	d.curFd = fd
	d.FsDriver = driver.FsDriver
	return d
}

func (driver *luaScriptDriverImpl) file(url string) string {
	if strings.HasPrefix(url, DriverSchema) {
		return strings.TrimPrefix(url, DriverSchema)
	}
	return url
}

func (driver *luaScriptDriverImpl) Close() error {
	defer driver.reset()
	if driver.curFd != nil {
		if err := driver.curFd.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (driver *luaScriptDriverImpl) reset() {
	driver.curFd = nil
}

func (driver *luaScriptDriverImpl) First() (version uint, err error) {
	panic("implement me")
}

func (driver *luaScriptDriverImpl) Prev(version uint) (prevVersion uint, err error) {
	panic("implement me")
}

func (driver *luaScriptDriverImpl) Next(version uint) (nextVersion uint, err error) {
	panic("implement me")
}

func (driver *luaScriptDriverImpl) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	panic("implement me")
}

func (driver *luaScriptDriverImpl) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	panic("implement me")
}

func (driver *luaScriptDriverImpl) destroy() {
	defer runtime.SetFinalizer(driver, nil)
	if driver.curFd == nil {
		return
	}
	_ = driver.Close()
}

// 加载目录中的脚本
func (driver *luaScriptDriverImpl) load() error {

	return nil
}
