package beego

import (
	"mime"
	"path/filepath"
	"strconv"

	"net/http"

	"github.com/astaxie/beego/session"
)

//
func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

// register default error http handlers, 404,401,403,500 and 503.
func registerDefaultErrorHandler() error {

	for e, h := range map[string]func(http.ResponseWriter, *http.Request){
		"401": unauthorized,
		"402": paymentRequired,
		"403": forbidden,
		"404": notFound,
		"405": methodNotAllowed,
		"500": internalServerError,
		"501": notImplemented,
		"502": badGateway,
		"503": serviceUnavailable,
		"504": gatewayTimeout,
	} {
		if _, ok := ErrorMaps[e]; !ok {
			ErrorHandler(e, h)
		}
	}
	return nil
}

func registerSession() error {
	if BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		if sessionConfig == "" {
			sessionConfig = `{"cookieName":"` + BConfig.WebConfig.Session.SessionName + `",` +
				`"gclifetime":` + strconv.FormatInt(BConfig.WebConfig.Session.SessionGCMaxLifetime, 10) + `,` +
				`"providerConfig":"` + filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig) + `",` +
				`"secure":` + strconv.FormatBool(BConfig.Listen.HTTPSEnable) + `,` +
				`"enableSetCookie":` + strconv.FormatBool(BConfig.WebConfig.Session.SessionAutoSetCookie) + `,` +
				`"domain":"` + BConfig.WebConfig.Session.SessionDomain + `",` +
				`"cookieLifeTime":` + strconv.Itoa(BConfig.WebConfig.Session.SessionCookieLifeTime) + `}`
		}
		GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, sessionConfig)
		if err != nil {
			return err
		}
		go GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	if BConfig.WebConfig.AutoRender {
		err := BuildTemplate(BConfig.WebConfig.ViewsPath)
		if err != nil && BConfig.RunMode == "dev" {
			Warn(err)
		}
	}
	return nil
}

func registerDocs() error {
	if BConfig.WebConfig.EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
	return nil
}

func registerAdmin() error {
	if BConfig.Listen.AdminEnable {
		go beeAdminApp.Run()
	}
	return nil
}
