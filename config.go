package beego

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	bComment = []byte{'#'}
	bEmpty   = []byte{}
	bEqual   = []byte{'='}
	bDQuote  = []byte{'"'}
)

// A Config represents the configuration.
type Config struct {
	filename string
	comment  map[int][]string  // id: []{comment, key...}; id 1 is for main comment.
	data     map[string]string // key: value
	offset   map[string]int64  // key: offset; for editing.
	sync.RWMutex
}

// ParseFile creates a new Config and parses the file configuration from the
// named file.
func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		file.Name(),
		make(map[int][]string),
		make(map[string]string),
		make(map[string]int64),
		sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)

	for nComment, off := 0, int64(1); ; {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}

		off += int64(len(line))

		if bytes.HasPrefix(line, bComment) {
			line = bytes.TrimLeft(line, "#")
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}
		if comment.Len() != 0 {
			cfg.comment[nComment] = []string{comment.String()}
			comment.Reset()
			nComment++
		}

		val := bytes.SplitN(line, bEqual, 2)
		if bytes.HasPrefix([]byte(strings.TrimSpace(string(val[1]))), bDQuote) {
			val[1] = bytes.Trim([]byte(strings.TrimSpace(string(val[1]))), `"`)
		}

		key := strings.TrimSpace(string(val[0]))
		cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
		cfg.data[key] = strings.TrimSpace(string(val[1]))
		cfg.offset[key] = off
	}
	return cfg, nil
}

// Bool returns the boolean value for a given key.
func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key])
}

// Int returns the integer value for a given key.
func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

func (c *Config) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key], 10, 64)
}

// Float returns the float value for a given key.
func (c *Config) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key], 64)
}

// String returns the string value for a given key.
func (c *Config) String(key string) string {
	return c.data[key]
}

// WriteValue writes a new value for key.
func (c *Config) SetValue(key, value string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.data[key]; !found {
		return errors.New("key not found: " + key)
	}

	c.data[key] = value
	return nil
}

func ParseConfig() (err error) {
	AppConfig, err = LoadConfig(AppConfigPath)
	if err != nil {
		return err
	} else {
		HttpAddr = AppConfig.String("httpaddr")
		if v, err := AppConfig.Int("httpport"); err == nil {
			HttpPort = v
		}
		if maxmemory, err := AppConfig.Int64("maxmemory"); err == nil {
			MaxMemory = maxmemory
		}
		AppName = AppConfig.String("appname")
		if runmode := AppConfig.String("runmode"); runmode != "" {
			RunMode = runmode
		}
		if autorender, err := AppConfig.Bool("autorender"); err == nil {
			AutoRender = autorender
		}
		if autorecover, err := AppConfig.Bool("autorecover"); err == nil {
			RecoverPanic = autorecover
		}
		if pprofon, err := AppConfig.Bool("pprofon"); err == nil {
			PprofOn = pprofon
		}
		if views := AppConfig.String("viewspath"); views != "" {
			ViewsPath = views
		}
		if sessionon, err := AppConfig.Bool("sessionon"); err == nil {
			SessionOn = sessionon
		}
		if sessProvider := AppConfig.String("sessionprovider"); sessProvider != "" {
			SessionProvider = sessProvider
		}
		if sessName := AppConfig.String("sessionname"); sessName != "" {
			SessionName = sessName
		}
		if sesssavepath := AppConfig.String("sessionsavepath"); sesssavepath != "" {
			SessionSavePath = sesssavepath
		}
		if sessMaxLifeTime, err := AppConfig.Int("sessiongcmaxlifetime"); err == nil && sessMaxLifeTime != 0 {
			int64val, _ := strconv.ParseInt(strconv.Itoa(sessMaxLifeTime), 10, 64)
			SessionGCMaxLifetime = int64val
		}
		if usefcgi, err := AppConfig.Bool("usefcgi"); err == nil {
			UseFcgi = usefcgi
		}
		if enablegzip, err := AppConfig.Bool("enablegzip"); err == nil {
			EnableGzip = enablegzip
		}
		if directoryindex, err := AppConfig.Bool("directoryindex"); err == nil {
			DirectoryIndex = directoryindex
		}
		if hotupdate, err := AppConfig.Bool("hotupdate"); err == nil {
			EnableHotUpdate = hotupdate
		}
		if timeout, err := AppConfig.Int64("httpservertimeout"); err == nil {
			HttpServerTimeOut = timeout
		}
		if errorsshow, err := AppConfig.Bool("errorsshow"); err == nil {
			ErrorsShow = errorsshow
		}
		if copyrequestbody, err := AppConfig.Bool("copyrequestbody"); err == nil {
			CopyRequestBody = copyrequestbody
		}
		if xsrfkey := AppConfig.String("xsrfkey"); xsrfkey != "" {
			XSRFKEY = xsrfkey
		}
		if enablexsrf, err := AppConfig.Bool("enablexsrf"); err == nil {
			EnableXSRF = enablexsrf
		}
		if expire, err := AppConfig.Int("xsrfexpire"); err == nil {
			XSRFExpire = expire
		}
		if tplleft := AppConfig.String("templateleft"); tplleft != "" {
			TemplateLeft = tplleft
		}
		if tplright := AppConfig.String("templateright"); tplright != "" {
			TemplateRight = tplright
		}
	}
	return nil
}
