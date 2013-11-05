package session

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	filepder      = &FileProvider{}
	gcmaxlifetime int64
)

type FileSessionStore struct {
	f      *os.File
	sid    string
	lock   sync.RWMutex
	values map[interface{}]interface{}
}

func (fs *FileSessionStore) Set(key, value interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.values[key] = value
	return nil
}

func (fs *FileSessionStore) Get(key interface{}) interface{} {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if v, ok := fs.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (fs *FileSessionStore) Delete(key interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	delete(fs.values, key)
	return nil
}

func (fs *FileSessionStore) Flush() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.values = make(map[interface{}]interface{})
	return nil
}

func (fs *FileSessionStore) SessionID() string {
	return fs.sid
}

func (fs *FileSessionStore) SessionRelease() {
	defer fs.f.Close()
	b, err := encodeGob(fs.values)
	if err != nil {
		return
	}
	fs.f.Truncate(0)
	fs.f.Seek(0, 0)
	fs.f.Write(b)
}

type FileProvider struct {
	maxlifetime int64
	savePath    string
}

func (fp *FileProvider) SessionInit(maxlifetime int64, savePath string) error {
	fp.maxlifetime = maxlifetime
	fp.savePath = savePath
	return nil
}

func (fp *FileProvider) SessionRead(sid string) (SessionStore, error) {
	err := os.MkdirAll(path.Join(fp.savePath, string(sid[0]), string(sid[1])), 0777)
	if err != nil {
		println(err.Error())
	}
	_, err = os.Stat(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	var f *os.File
	if err == nil {
		f, err = os.OpenFile(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid), os.O_RDWR, 0777)
	} else if os.IsNotExist(err) {
		f, err = os.Create(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	} else {
		return nil, err
	}
	os.Chtimes(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid), time.Now(), time.Now())
	var kv map[interface{}]interface{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(b)
		if err != nil {
			return nil, err
		}
	}
	f.Close()
	f, err = os.OpenFile(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid), os.O_WRONLY|os.O_CREATE, 0777)
	ss := &FileSessionStore{f: f, sid: sid, values: kv}
	return ss, nil
}

func (fp *FileProvider) SessionExist(sid string) bool {
	_, err := os.Stat(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	if err == nil {
		return true
	} else {
		return false
	}
}

func (fp *FileProvider) SessionDestroy(sid string) error {
	os.Remove(path.Join(fp.savePath))
	return nil
}

func (fp *FileProvider) SessionGC() {
	gcmaxlifetime = fp.maxlifetime
	filepath.Walk(fp.savePath, gcpath)
}

func (fp *FileProvider) SessionAll() int {
	a := &activeSession{}
	err := filepath.Walk(fp.savePath, func(path string, f os.FileInfo, err error) error {
		return a.visit(path, f, err)
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return 0
	}
	return a.total
}

func (fp *FileProvider) SessionRegenerate(oldsid, sid string) (SessionStore, error) {
	err := os.MkdirAll(path.Join(fp.savePath, string(oldsid[0]), string(oldsid[1])), 0777)
	if err != nil {
		println(err.Error())
	}
	err = os.MkdirAll(path.Join(fp.savePath, string(sid[0]), string(sid[1])), 0777)
	if err != nil {
		println(err.Error())
	}
	_, err = os.Stat(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	var newf *os.File
	if err == nil {
		return nil, errors.New("newsid exist")
	} else if os.IsNotExist(err) {
		newf, err = os.Create(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	}

	_, err = os.Stat(path.Join(fp.savePath, string(oldsid[0]), string(oldsid[1]), oldsid))
	var f *os.File
	if err == nil {
		f, err = os.OpenFile(path.Join(fp.savePath, string(oldsid[0]), string(oldsid[1]), oldsid), os.O_RDWR, 0777)
		io.Copy(newf, f)
	} else if os.IsNotExist(err) {
		newf, err = os.Create(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	} else {
		return nil, err
	}
	f.Close()
	os.Remove(path.Join(fp.savePath, string(oldsid[0]), string(oldsid[1])))
	os.Chtimes(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid), time.Now(), time.Now())
	var kv map[interface{}]interface{}
	b, err := ioutil.ReadAll(newf)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(b)
		if err != nil {
			return nil, err
		}
	}

	newf, err = os.OpenFile(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid), os.O_WRONLY|os.O_CREATE, 0777)
	ss := &FileSessionStore{f: newf, sid: sid, values: kv}
	return ss, nil
}

func gcpath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if (info.ModTime().Unix() + gcmaxlifetime) < time.Now().Unix() {
		os.Remove(path)
	}
	return nil
}

type activeSession struct {
	total int
}

func (self *activeSession) visit(paths string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	self.total = self.total + 1
	return nil
}

func init() {
	Register("file", filepder)
}
