package beego

import (
	"mime"
	"path/filepath"
	"strconv"

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
	if _, ok := ErrorMaps["401"]; !ok {
		Errorhandler("401", unauthorized)
	}

	if _, ok := ErrorMaps["402"]; !ok {
		Errorhandler("402", paymentRequired)
	}

	if _, ok := ErrorMaps["403"]; !ok {
		Errorhandler("403", forbidden)
	}

	if _, ok := ErrorMaps["404"]; !ok {
		Errorhandler("404", notFound)
	}

	if _, ok := ErrorMaps["405"]; !ok {
		Errorhandler("405", methodNotAllowed)
	}

	if _, ok := ErrorMaps["500"]; !ok {
		Errorhandler("500", internalServerError)
	}
	if _, ok := ErrorMaps["501"]; !ok {
		Errorhandler("501", notImplemented)
	}
	if _, ok := ErrorMaps["502"]; !ok {
		Errorhandler("502", badGateway)
	}

	if _, ok := ErrorMaps["503"]; !ok {
		Errorhandler("503", serviceUnavailable)
	}

	if _, ok := ErrorMaps["504"]; !ok {
		Errorhandler("504", gatewayTimeout)
	}
	return nil
}

func registerSession() error {
	if SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		if sessionConfig == "" {
			sessionConfig = `{"cookieName":"` + SessionName + `",` +
				`"gclifetime":` + strconv.FormatInt(SessionGCMaxLifetime, 10) + `,` +
				`"providerConfig":"` + filepath.ToSlash(SessionSavePath) + `",` +
				`"secure":` + strconv.FormatBool(EnableHttpTLS) + `,` +
				`"enableSetCookie":` + strconv.FormatBool(SessionAutoSetCookie) + `,` +
				`"domain":"` + SessionDomain + `",` +
				`"cookieLifeTime":` + strconv.Itoa(SessionCookieLifeTime) + `}`
		}
		GlobalSessions, err = session.NewManager(SessionProvider, sessionConfig)
		if err != nil {
			return err
		}
		go GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	if AutoRender {
		err := BuildTemplate(ViewsPath)
		if err != nil && RunMode == "dev" {
			Warn(err)
		}
	}
	return nil
}

func registerDocs() error {
	if EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
	return nil
}

func registerAdmin() error {
	if EnableAdmin {
		go beeAdminApp.Run()
	}
	return nil
}
