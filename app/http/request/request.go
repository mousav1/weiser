package request

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/mousav1/weiser/app/http/validation"
)

type Request struct {
	ctx      *fiber.Ctx
	validate *validation.Validation // افزودن فیلد validate به struct Request
}

func New(ctx *fiber.Ctx) (*Request, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	validate := validation.New()

	// تعریف تنظیمات خاص برای validator
	// validate.RegisterValidation("my_validation", func(fl validator.FieldLevel) bool {
	// 	// اعتبارسنجی خاص
	// 	return true
	// })

	// additional error checks can be performed here
	return &Request{ctx: ctx, validate: validate}, nil
}

func (r *Request) Bind(data interface{}) error {
	// Check content type
	contentType := r.ctx.Get("Content-Type")
	switch contentType {
	case "application/json":
		if err := r.ctx.BodyParser(data); err != nil {
			return err
		}
	case "application/x-www-form-urlencoded":
		// Parse form data
		formData := make(map[string]interface{})
		if err := r.ctx.BodyParser(&formData); err != nil {
			return err
		}

		// Map form data
		if err := mapstructure.Decode(formData, &data); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	// اعتبارسنجی داده‌های دریافتی
	if err := r.validate.Validate(data); err != nil {
		return err
	}

	return nil
}

func (r *Request) Input(key string, def ...string) string {
	if val := r.ctx.FormValue(key); val != "" {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (r *Request) All() (map[string][]string, error) {
	form, err := r.ctx.MultipartForm()
	if err != nil {
		return nil, err
	}
	return form.Value, nil
}

func (r *Request) Only(keys []string) map[string]interface{} {
	values := make(map[string]interface{})
	for _, key := range keys {
		values[key] = r.ctx.FormValue(key)
	}
	return values
}

func (r *Request) Except(keys []string) (map[string][]string, error) {
	values, err := r.All()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		delete(values, key)
	}
	return values, nil
}

func (r *Request) Method() string {
	return r.ctx.Method()
}

func (r *Request) Scheme() string {
	if r.ctx.Protocol() == "https" {
		return "https"
	}
	return "http"
}

func (r *Request) Host() string {
	return r.ctx.Hostname()
}

func (r *Request) Path() string {
	return r.ctx.Path()
}

func (r *Request) URL() string {
	return r.Scheme() + "://" + r.Host() + r.Path()
}

func (r *Request) IsAjax() bool {
	return strings.ToLower(r.ctx.Get("X-Requested-With")) == "xmlhttprequest"
}

func (r *Request) IsMethod(method string) bool {
	return r.Method() == strings.ToUpper(method)
}

func (r *Request) IsSecure() bool {
	return r.Scheme() == "https"
}

func (r *Request) Is(p string) bool {
	switch p {
	case "json":
		return strings.Contains(r.ctx.Get("Content-Type"), "application/json")
	case "html":
		return strings.Contains(r.ctx.Get("Content-Type"), "text/html")
	case "xml":
		return strings.Contains(r.ctx.Get("Content-Type"), "application/xml, text/xml")
	case "plain":
		return strings.Contains(r.ctx.Get("Content-Type"), "text/plain")
	default:
		return false
	}
}

func (r *Request) IP() string {
	return r.ctx.IP()
}

func (r *Request) UserAgent() string {
	return r.ctx.Get("User-Agent")
}

func (r *Request) Referer() string {
	return r.ctx.Get("Referer")
}

func (r *Request) Ajax() bool {
	return strings.ToLower(r.ctx.Get("X-Requested-With")) == "xmlhttprequest"
}

func (r *Request) Header(key string) string {
	return r.ctx.Get(key)
}

func (r *Request) Value(key string, def ...string) string {
	if val := r.ctx.Params(key); val != "" {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (r *Request) Int(key string, def ...int) (int, error) {
	val := r.Value(key)
	if val == "" && len(def) > 0 {
		return def[0], nil
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *Request) Float(key string, def ...float64) (float64, error) {
	val := r.Value(key)
	if val == "" && len(def) > 0 {
		return def[0], nil
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *Request) Bool(key string, def ...bool) (bool, error) {
	val := r.Value(key)
	if val == "" && len(def) > 0 {
		return def[0], nil
	}
	v, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return v, nil
}

func (r *Request) Root() string {
	proto := "http"
	if r.ctx.Protocol() == "https" {
		proto = "https"
	}
	return proto + "://" + r.ctx.Hostname()
}

func (r *Request) FullURL() string {
	return r.Root() + r.ctx.OriginalURL()
}

func (r *Request) IsMethodSafe() bool {
	method := r.ctx.Method()
	return method == http.MethodGet || method == http.MethodHead
}

func (r *Request) IsJson() bool {
	contentType := r.ctx.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}

func (r *Request) IsXml() bool {
	contentType := r.ctx.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/xml") || strings.HasPrefix(contentType, "text/xml")
}

func (r *Request) IsHtml() bool {
	contentType := r.ctx.Get("Content-Type")
	return strings.HasPrefix(contentType, "text/html")
}

// func (r *Request) Cookies() map[string]string {
// 	headers := r.ctx.Request().Header
// 	result := make(map[string]string)
// 	for _, cookie := range headers.PeekMulti("Cookie") {
// 		parts := bytes.SplitN(cookie, []byte("="), 2)
// 		if len(parts) == 2 {
// 			result[string(parts[0])] = string(parts[1])
// 		}
// 	}
// 	return result
// }

func (r *Request) Cookie(key string) (string, error) {
	cookie := r.ctx.Cookies(key)
	if cookie == "" {
		return "", errors.New("cookie not found")
	}
	return cookie, nil
}
func (r *Request) Has(key string) bool {
	return r.ctx.FormValue(key) != ""
}

func (r *Request) HasFile(key string) bool {
	_, err := r.ctx.FormFile(key)
	return err == nil
}

func (r *Request) File(key string) (*multipart.FileHeader, error) {
	return r.ctx.FormFile(key)
}

func (r *Request) AllFiles() map[string][]*multipart.FileHeader {
	form, err := r.ctx.MultipartForm()
	if err != nil {
		return nil
	}
	return form.File
}
