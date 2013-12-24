/**
 * package: file
 * User: gouki
 * Date: 2013-10-22 - 14:22
 */
package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

func init() {
	Register("file", NewFileCache())
}

// FileCacheItem is basic unit of file cache adapter.
// it contains data and expire time.
type FileCacheItem struct {
	Data       interface{}
	Lastaccess int64
	Expired    int64
}

var (
	FileCachePath           string = "cache" // cache directory
	FileCacheFileSuffix     string = ".bin"  // cache file suffix
	FileCacheDirectoryLevel int    = 2       // cache file deep level if auto generated cache files.
	FileCacheEmbedExpiry    int64  = 0       // cache expire time, default is no expire forever.
)

// FileCache is cache adapter for file storage.
type FileCache struct {
	CachePath      string
	FileSuffix     string
	DirectoryLevel int
	EmbedExpiry    int
}

// Create new file cache with no config.
// the level and expiry need set in method StartAndGC as config string.
func NewFileCache() *FileCache {
	//    return &FileCache{CachePath:FileCachePath, FileSuffix:FileCacheFileSuffix}
	return &FileCache{}
}

// Start and begin gc for file cache.
// the config need to be like {CachePath:"/cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}
func (this *FileCache) StartAndGC(config string) error {

	var cfg map[string]string
	json.Unmarshal([]byte(config), &cfg)
	//fmt.Println(cfg)
	if _, ok := cfg["CachePath"]; !ok {
		cfg["CachePath"] = FileCachePath
	}
	if _, ok := cfg["FileSuffix"]; !ok {
		cfg["FileSuffix"] = FileCacheFileSuffix
	}
	if _, ok := cfg["DirectoryLevel"]; !ok {
		cfg["DirectoryLevel"] = strconv.Itoa(FileCacheDirectoryLevel)
	}
	if _, ok := cfg["EmbedExpiry"]; !ok {
		cfg["EmbedExpiry"] = strconv.FormatInt(FileCacheEmbedExpiry, 10)
	}
	this.CachePath = cfg["CachePath"]
	this.FileSuffix = cfg["FileSuffix"]
	this.DirectoryLevel, _ = strconv.Atoi(cfg["DirectoryLevel"])
	this.EmbedExpiry, _ = strconv.Atoi(cfg["EmbedExpiry"])

	this.Init()
	return nil
}

// Init will make new dir for file cache if not exist.
func (this *FileCache) Init() {
	app := filepath.Dir(os.Args[0])
	this.CachePath = filepath.Join(app, this.CachePath)
	ok, err := exists(this.CachePath)
	if err != nil { // print error
		//fmt.Println(err)
	}
	if !ok {
		if err = os.Mkdir(this.CachePath, os.ModePerm); err != nil {
			//fmt.Println(err);
		}
	}
	//fmt.Println(this.getCacheFileName("123456"));
}

// get cached file name. it's md5 encoded.
func (this *FileCache) getCacheFileName(key string) string {
	m := md5.New()
	io.WriteString(m, key)
	keyMd5 := hex.EncodeToString(m.Sum(nil))
	cachePath := this.CachePath
	//fmt.Println("cachepath : " , cachePath)
	//fmt.Println("md5" , keyMd5);
	switch this.DirectoryLevel {
	case 2:
		cachePath = filepath.Join(cachePath, keyMd5[0:2], keyMd5[2:4])
	case 1:
		cachePath = filepath.Join(cachePath, keyMd5[0:2])
	}

	ok, err := exists(cachePath)
	if err != nil {
		//fmt.Println(err)
	}
	if !ok {
		if err = os.MkdirAll(cachePath, os.ModePerm); err != nil {
			//fmt.Println(err);
		}
	}
	return filepath.Join(cachePath, fmt.Sprintf("%s%s", keyMd5, this.FileSuffix))
}

// Get value from file cache.
// if non-exist or expired, return empty string.
func (this *FileCache) Get(key string) interface{} {
	filename := this.getCacheFileName(key)
	filedata, err := File_get_contents(filename)
	//fmt.Println("get length:" , len(filedata));
	if err != nil {
		return ""
	}
	var to FileCacheItem
	Gob_decode([]byte(filedata), &to)
	if to.Expired < time.Now().Unix() {
		return ""
	}
	return to.Data
}

// Put value into file cache.
// timeout means how long to keep this file, unit of ms.
// if timeout equals FileCacheEmbedExpiry(default is 0), cache this item forever.
func (this *FileCache) Put(key string, val interface{}, timeout int64) error {
	filename := this.getCacheFileName(key)
	var item FileCacheItem
	item.Data = val
	if timeout == FileCacheEmbedExpiry {
		item.Expired = time.Now().Unix() + (86400 * 365 * 10) // ten years
	} else {
		item.Expired = time.Now().Unix() + timeout
	}
	item.Lastaccess = time.Now().Unix()
	data, err := Gob_encode(item)
	if err != nil {
		return err
	}
	err = File_put_contents(filename, data)
	return err
}

// Delete file cache value.
func (this *FileCache) Delete(key string) error {
	filename := this.getCacheFileName(key)
	if ok, _ := exists(filename); ok {
		return os.Remove(filename)
	}
	return nil
}

// Increase cached int value.
// this value is saving forever unless Delete.
func (this *FileCache) Incr(key string) error {
	data := this.Get(key)
	var incr int
	fmt.Println(reflect.TypeOf(data).Name())
	if reflect.TypeOf(data).Name() != "int" {
		incr = 0
	} else {
		incr = data.(int) + 1
	}
	this.Put(key, incr, FileCacheEmbedExpiry)
	return nil
}

// Decrease cached int value.
func (this *FileCache) Decr(key string) error {
	data := this.Get(key)
	var decr int
	if reflect.TypeOf(data).Name() != "int" || data.(int)-1 <= 0 {
		decr = 0
	} else {
		decr = data.(int) - 1
	}
	this.Put(key, decr, FileCacheEmbedExpiry)
	return nil
}

// Check value is exist.
func (this *FileCache) IsExist(key string) bool {
	filename := this.getCacheFileName(key)
	ret, _ := exists(filename)
	return ret
}

// Clean cached files.
// not implemented.
func (this *FileCache) ClearAll() error {
	//this.CachePath .递归删除

	return nil
}

// check file exist.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Get bytes to file.
// if non-exist, create this file.
func File_get_contents(filename string) ([]byte, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return []byte(""), err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return []byte(""), err
	}
	data := make([]byte, stat.Size())
	result, err := f.Read(data)
	if int64(result) == stat.Size() {
		return data, err
	}
	return []byte(""), err
}

// Put bytes to file.
// if non-exist, create this file.
func File_put_contents(filename string, content []byte) error {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(content)
	return err
}

// Gob encodes file cache item.
func Gob_encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Gob decodes file cache item.
func Gob_decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&to)
}
