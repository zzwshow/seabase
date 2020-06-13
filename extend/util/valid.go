package util

import (
	"bytes"
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

var Vd *validator.Validate
var translators = make(map[string]ut.Translator, 2)

// 初始化验证器
func init() {
	Vd = validator.New()

	// 使用json标签
	Vd.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	// 注册英文错误提示器
	_en := en.New()
	translators["en"], _ = ut.New(_en, _en).GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(Vd, translators["en"])

	// 注册中文错误提示器
	_zh := zh.New()
	translators["zh"], _ = ut.New(_zh, _zh).GetTranslator("zh")
	_ = zh_translations.RegisterDefaultTranslations(Vd, translators["zh"])

	// 注册一个国家码验证器
	_ = Vd.RegisterValidation("cc", func(fl validator.FieldLevel) bool {
		ok, _ := regexp.MatchString(`^[1-9][0-9]{1,2}$`, fl.Field().String())
		return ok
	})
	RegisterTagTranslation("cc", map[string]string{
		"en": "{0} is a invalid cc.",
		"zh": "{0}不是一个可用的国家码",
	})

	// 注册一个手机号码验证器
	_ = Vd.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		ok, _ := regexp.MatchString(`^1[3-9][0-9]{9}$`, fl.Field().String())
		return ok
	})
	RegisterTagTranslation("phone", map[string]string{
		"en": "{0} is a invalid phone.",
		"zh": "{0}不是一个可用的手机号码",
	})

	// // 给required_with添加翻译
	// RegisterTagTranslation("required", map[string]string{
	// 	"en": "{0} is required field.",
	// 	"zh": "{0}是一个必填字段",
	// })
}

// 自定义翻译
func RegisterTagTranslation(tag string, messages map[string]string) {
	for lang, message := range messages {
		_ = Vd.RegisterTranslation(tag, translators[lang], func(ut ut.Translator) error {
			return ut.Add(tag, message, false)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})
	}
}

func AnalysisErrorByValidate(error error) (err error) {
	var (
		buffer   bytes.Buffer
		errorStr string
	)
	if error != nil {
		for _, err := range error.(validator.ValidationErrors) {
			errorStr = "、" + err.Translate(translators["zh"]) //自行选择使用哪种语言提示
			//errorStr = err.Translate(translators["en"]) + ", "
			buffer.WriteString(errorStr)
		}
		err = errors.New(strings.TrimRight(buffer.String(), ", "))
		return err
	}
	return nil
}
