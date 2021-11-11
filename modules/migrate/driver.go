package migrate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source"
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/fs"
	nurl "net/url"
	"runtime"
	"strings"
	"sync"
)

type (
	LuaScriptDriverImpl struct {
		Schema    string
		safe      sync.RWMutex
		nfs       FileSystem
		curReader *luaSqlBuilderReaderCloser
	}

	FileSystem interface {
		fs.ReadFileFS
		fs.ReadDirFS
		io.Closer
	}

	luaSqlBuilderReaderCloser struct {
		fd         fs.File
		data       []byte
		identifier string
		method     string
		vm         *lua.LState
		safe       sync.RWMutex
		sqlBuffer  *bytes.Buffer
	}
)

const (
	scheme       = "lua"
	DriverScheme = "lua://"
	GBuffer = "__GBuffer"
)

func init() {
	source.Register(DriverScheme, &LuaScriptDriverImpl{})
}

func NewLuaSqlBuilderReaderCloser(fd fs.File) *luaSqlBuilderReaderCloser {
	var reader = new(luaSqlBuilderReaderCloser)
	reader.fd = fd
	reader.safe = sync.RWMutex{}
	reader.sqlBuffer = bytes.NewBufferString(``)
	return reader
}

func (reader *luaSqlBuilderReaderCloser) Close() error {
	if reader.fd != nil {
		if err := reader.fd.Close(); err != nil {
			return err
		}
		reader.clear()
	}
	return nil
}

func (reader *luaSqlBuilderReaderCloser) SetVm(vm *lua.LState) *luaSqlBuilderReaderCloser {
	reader.safe.Lock()
	defer reader.safe.Unlock()
	if reader.vm == nil && vm != nil {
		reader.vm = vm
	}
	return reader
}

func (reader *luaSqlBuilderReaderCloser) Read(method string) (sql []byte, identifier string, error error) {
	reader.safe.Lock()
	defer reader.safe.Unlock()
	if reader.method == "" {
		reader.method = method
	}
	if reader.method != method {
		reader.method = method
	}
	if reader.vm == nil {
		return nil, reader.getIdentifier(), errors.New("miss lua.LState(vm)")
	}
	reader.sqlBuffer.Reset()
	if err := reader.parse(); err != nil {
		return nil, reader.getIdentifier(), err
	}
	return reader.sqlBuffer.Bytes(), reader.getIdentifier(), error
}

func (reader *luaSqlBuilderReaderCloser) parse() error {
	// 是否已经加载过脚本
	if reader.data == nil || len(reader.data) <= 0 {
		var info, err = reader.fd.Stat()
		if err != nil {
			return err
		}
		// 1. 读取脚本
		reader.data = make([]byte, info.Size()+1)
		if _, err = reader.fd.Read(reader.data); err != nil {
			return err
		}
		// 2. 注入 sql buffer
		reader.vm.SetGlobal(GBuffer, &lua.LUserData{
			Value: reader.sqlBuffer,
			Env:   reader.vm.Env,
		})
		// 3. 载入脚本
		if err = reader.vm.DoString(string(reader.data)); err != nil {
			return err
		}
	}
	// 4. 执行脚本 相关方法
	if err := reader.vm.DoString(reader.getExecCode()); err != nil {
		return err
	}
	return nil
}

// 获取可执行脚本代码
func (reader *luaSqlBuilderReaderCloser) getExecCode() string {
	if reader.method == "" {
		reader.method = "up"
	}
	var (
		method = "safeUp"               // lua method 驼峰法
		mod    = reader.getIdentifier() // 模块名
	)
	switch strings.ToLower(reader.method) {
	case string(source.Up), "safeup":
		method = "safeUp"
	case string(source.Down), "safedown":
		method = "safeDown"
	default:
		// 自定义方法
		method = reader.method
	}
	return fmt.Sprintf(`migration=module("%s");migration.%s();'`, mod, method)
}

func (reader *luaSqlBuilderReaderCloser) getIdentifier() string {
	if reader.identifier == "" && reader.fd != nil {
		if info, err := reader.fd.Stat(); err == nil {
			reader.identifier = strings.TrimSuffix(info.Name(), ".lua")
		}
	}
	return reader.identifier
}

func (reader *luaSqlBuilderReaderCloser) Reset() {
	reader.safe.Lock()
	defer reader.safe.Unlock()
	if reader.data != nil {
		reader.data = nil
	}
	reader.fd = nil
	reader.method = ""
	reader.identifier = ""
	reader.sqlBuffer.Reset()
}

func (reader *luaSqlBuilderReaderCloser) clear() {
	reader.data = nil
	reader.fd = nil
	reader.method = ""
	reader.sqlBuffer = nil
	reader.identifier = ""
}

func NewDriver() *LuaScriptDriverImpl {
	var driver = new(LuaScriptDriverImpl)
	return driver.init()
}

func (driver *LuaScriptDriverImpl) init() *LuaScriptDriverImpl {
	driver.curReader = nil
	driver.Schema = DriverScheme
	driver.safe = sync.RWMutex{}
	return driver
}

func (driver *LuaScriptDriverImpl) InitFs() *LuaScriptDriverImpl {

	return driver
}

// Open 打开
func (driver *LuaScriptDriverImpl) Open(url string) (source.Driver, error) {
	// url eg: lua://luaVm/dir/?fs=local&root=/data
	if url == "" {
			return nil,errors.New("empty url")
	}
	var info, err = nurl.Parse(url)
	if err!=nil {
			return nil,err
	}
	if info.Scheme != DriverScheme && info.Scheme!= scheme {
			return nil,errors.New("unSupport scheme")
	}

	return nil, nil
}


func (driver *LuaScriptDriverImpl) newFs(fd fs.File) source.Driver {
	var d = driver.setFd(fd)
	if d != nil {
		runtime.SetFinalizer(d, (*LuaScriptDriverImpl).destroy)
	}
	return d
}

func (driver *LuaScriptDriverImpl) setFd(fd fs.File) *LuaScriptDriverImpl {
	driver.safe.Lock()
	defer driver.safe.Unlock()
	if driver.curReader == nil {
		driver.curReader = NewLuaSqlBuilderReaderCloser(fd)
		return driver
	}
	var d = NewDriver()
	d.nfs = driver.nfs
	d.curReader = NewLuaSqlBuilderReaderCloser(fd)
	return d
}

func (driver *LuaScriptDriverImpl) file(url string) string {
	if strings.HasPrefix(url, DriverScheme) {
		return strings.TrimPrefix(url, DriverScheme)
	}
	return url
}

func (driver *LuaScriptDriverImpl) Close() error {
	defer driver.reset()
	if driver.curReader != nil {
		if err := driver.curReader.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (driver *LuaScriptDriverImpl) reset() {
	driver.curReader = nil
}

func (driver *LuaScriptDriverImpl) First() (version uint, err error) {
	panic("implement me")
}

func (driver *LuaScriptDriverImpl) Prev(version uint) (prevVersion uint, err error) {
	panic("implement me")
}

func (driver *LuaScriptDriverImpl) Next(version uint) (nextVersion uint, err error) {
	panic("implement me")
}

func (driver *LuaScriptDriverImpl) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	panic("implement me")
}

func (driver *LuaScriptDriverImpl) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	panic("implement me")
}

func (driver *LuaScriptDriverImpl) destroy() {
	defer runtime.SetFinalizer(driver, nil)
	if driver.curReader == nil {
		return
	}
	_ = driver.Close()
}

// 加载目录中的脚本
func (driver *LuaScriptDriverImpl) load() error {

	return nil
}
