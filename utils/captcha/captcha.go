package captcha

// modifiy and integrated to Beego in one file from https://github.com/dchest/captcha

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
)

var (
	defaultChars = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

const (
	challengeNums    = 6
	expiration       = 600
	fieldIdName      = "captcha_id"
	fieldCaptchaName = "captcha"
	cachePrefix      = "captcha_"
	urlPrefix        = "/captcha/"
)

type Captcha struct {
	store            cache.Cache
	urlPrefix        string
	FieldIdName      string
	FieldCaptchaName string
	StdWidth         int
	StdHeight        int
	ChallengeNums    int
	Expiration       int64
	CachePrefix      string
}

func (c *Captcha) key(id string) string {
	return c.CachePrefix + id
}

func (c *Captcha) genRandChars() []byte {
	return utils.RandomCreateBytes(c.ChallengeNums, defaultChars...)
}

func (c *Captcha) Handler(ctx *context.Context) {
	var chars []byte

	id := path.Base(ctx.Request.RequestURI)
	if i := strings.Index(id, "."); i != -1 {
		id = id[:i]
	}

	key := c.key(id)

	if v, ok := c.store.Get(key).([]byte); ok {
		chars = v
	} else {
		ctx.Output.SetStatus(404)
		ctx.WriteString("captcha not found")
		return
	}

	// reload captcha
	if len(ctx.Input.Query("reload")) > 0 {
		chars = c.genRandChars()
		if err := c.store.Put(key, chars, c.Expiration); err != nil {
			ctx.Output.SetStatus(500)
			ctx.WriteString("captcha reload error")
			beego.Error("Reload Create Captcha Error:", err)
			return
		}
	}

	img := NewImage(chars, c.StdWidth, c.StdHeight)
	if _, err := img.WriteTo(ctx.ResponseWriter); err != nil {
		beego.Error("Write Captcha Image Error:", err)
	}
}

func (c *Captcha) CreateCaptchaHtml() template.HTML {
	value, err := c.CreateCaptcha()
	if err != nil {
		beego.Error("Create Captcha Error:", err)
		return ""
	}

	// create html
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`+
		`<a class="captcha" href="javascript:">`+
		`<img onclick="this.src=('%s%s.png?reload='+(new Date()).getTime())" class="captcha-img" src="%s%s.png">`+
		`</a>`, c.FieldIdName, value, c.urlPrefix, value, c.urlPrefix, value))
}

func (c *Captcha) CreateCaptcha() (string, error) {
	// generate captcha id
	id := string(utils.RandomCreateBytes(15))

	// get the captcha chars
	chars := c.genRandChars()

	// save to store
	if err := c.store.Put(c.key(id), chars, c.Expiration); err != nil {
		return "", err
	}

	return id, nil
}

func (c *Captcha) VerifyReq(req *http.Request) bool {
	req.ParseForm()
	return c.Verify(req.Form.Get(c.FieldIdName), req.Form.Get(c.FieldCaptchaName))
}

func (c *Captcha) Verify(id string, challenge string) (success bool) {
	if len(challenge) == 0 || len(id) == 0 {
		return
	}

	var chars []byte

	key := c.key(id)

	if v, ok := c.store.Get(key).([]byte); ok && len(v) == len(challenge) {
		chars = v
	} else {
		return
	}

	defer func() {
		// finally remove it
		c.store.Delete(key)
	}()

	// verify challenge
	for i, c := range chars {
		if c != challenge[i]-48 {
			return
		}
	}

	return true
}

func NewCaptcha(urlPrefix string, store cache.Cache) *Captcha {
	cpt := &Captcha{}
	cpt.store = store
	cpt.FieldIdName = fieldIdName
	cpt.FieldCaptchaName = fieldCaptchaName
	cpt.ChallengeNums = challengeNums
	cpt.Expiration = expiration
	cpt.CachePrefix = cachePrefix
	cpt.StdWidth = stdWidth
	cpt.StdHeight = stdHeight

	if len(urlPrefix) == 0 {
		urlPrefix = urlPrefix
	}

	if urlPrefix[len(urlPrefix)-1] != '/' {
		urlPrefix += "/"
	}

	cpt.urlPrefix = urlPrefix

	return cpt
}

func NewWithFilter(urlPrefix string, store cache.Cache) *Captcha {
	cpt := NewCaptcha(urlPrefix, store)

	// create filter for serve captcha image
	beego.AddFilter(urlPrefix+":", "BeforeRouter", cpt.Handler)

	// add to template func map
	beego.AddFuncMap("create_captcha", cpt.CreateCaptchaHtml)

	return cpt
}
