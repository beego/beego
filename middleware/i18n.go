package middleware

//import (
//	"github.com/astaxie/beego/config"
//	"os"
//	"path"
//)

//type Translation struct {
//	filetype     string
//	CurrentLocal string
//	Locales      map[string]map[string]string
//}

//func NewLocale(filetype string) *Translation {
//	return &Translation{
//		filetype:     filetype,
//		CurrentLocal: "zh",
//		Locales:      make(map[string]map[string]string),
//	}
//}

//func (t *Translation) loadTranslations(dirPath string) error {
//	dir, err := os.Open(dirPath)
//	if err != nil {
//		return err
//	}
//	defer dir.Close()

//	names, err := dir.Readdirnames(-1)
//	if err != nil {
//		return err
//	}

//	for _, name := range names {
//		fullPath := path.Join(dirPath, name)

//		fi, err := os.Stat(fullPath)
//		if err != nil {
//			return err
//		}

//		if fi.IsDir() {
//			continue
//		} else {
//			if err := t.loadTranslation(fullPath, name); err != nil {
//				return err
//			}
//		}
//	}

//	return nil
//}

//func (t *Translation) loadTranslation(fullPath, locale string) error {

//	sourceKey2Trans, ok := t.Locales[locale]
//	if !ok {
//		sourceKey2Trans = make(map[string]string)

//		t.Locales[locale] = sourceKey2Trans
//	}

//	for _, m := range trf.Messages {
//		if m.Translation != "" {
//			sourceKey2Trans[sourceKey(m.Source, m.Context)] = m.Translation
//		}
//	}

//	return nil
//}

//func (t *Translation) SetLocale(local string) {
//	t.CurrentLocal = local
//}

//func (t *Translation) Translate(key string) string {
//	if ct, ok := t.Locales[t.CurrentLocal]; ok {
//		if v, o := ct[key]; o {
//			return v
//		}
//	}
//	return key
//}
