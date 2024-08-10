package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/mousav1/weiser/app/cookies"
	"github.com/mousav1/weiser/app/http/validation"
	"github.com/mousav1/weiser/app/session"
)

type Request struct {
	ctx      *fiber.Ctx
	validate *validation.Validation // افزودن فیلد validate به struct Request
}

func New(ctx *fiber.Ctx) (*Request, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	validate := validation.New()

	return &Request{ctx: ctx, validate: validate}, nil
}

func (r *Request) Bind(data interface{}) error {
	contentType := r.ctx.Get("Content-Type")

	switch contentType {
	case "application/json":
		if err := r.ctx.BodyParser(data); err != nil {
			return fmt.Errorf("failed to parse JSON body: %w", err)
		}
	case "application/x-www-form-urlencoded":
		if err := r.parseFormData(data); err != nil {
			return fmt.Errorf("failed to parse form data: %w", err)
		}
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	if err := r.validateData(data); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseFormData(data interface{}) error {
	formData := make(map[string]interface{})
	if err := r.ctx.BodyParser(&formData); err != nil {
		return err
	}
	return mapstructure.Decode(formData, &data)
}

func (r *Request) validateData(data interface{}) error {
	validationErrors, err := r.validate.Validate(data)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if validationErrors != nil {
		response := r.validate.CreateErrorResponse(validationErrors)
		jsonResp, _ := json.Marshal(response)
		return fiber.NewError(http.StatusBadRequest, string(jsonResp))
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
	contentType := r.ctx.Get("Content-Type")
	switch p {
	case "json":
		return strings.Contains(contentType, "application/json")
	case "html":
		return strings.Contains(contentType, "text/html")
	case "xml":
		return strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml")
	case "plain":
		return strings.Contains(contentType, "text/plain")
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
	return r.IsAjax()
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
	return strconv.Atoi(val)
}

func (r *Request) Float(key string, def ...float64) (float64, error) {
	val := r.Value(key)
	if val == "" && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(val, 64)
}

func (r *Request) Bool(key string, def ...bool) (bool, error) {
	val := r.Value(key)
	if val == "" && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(val)
}

func (r *Request) Root() string {
	proto := "http"
	if r.IsSecure() {
		proto = "https"
	}
	return proto + "://" + r.Host()
}

func (r *Request) FullURL() string {
	return r.Root() + r.ctx.OriginalURL()
}

func (r *Request) IsMethodSafe() bool {
	method := r.Method()
	return method == http.MethodGet || method == http.MethodHead
}

func (r *Request) IsJson() bool {
	return r.Is("json")
}

func (r *Request) IsXml() bool {
	return r.Is("xml")
}

func (r *Request) IsHtml() bool {
	return r.Is("html")
}

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

func (r *Request) getSessionID() (string, error) {
	manager := session.GetSessionManager()
	sessionID, err := cookies.GetCookie(r.ctx, "weiser_session")
	if err != nil {
		session := manager.StartSession(r.ctx)
		return session.ID, nil
	}
	if err := manager.CheckExpiration(sessionID.(string)); err != nil {
		return "", fmt.Errorf("session ID is invalid or has expired: %w", err)
	}
	return sessionID.(string), nil
}

func (r *Request) Getsession(key string) interface{} {
	sessionID, err := r.getSessionID()
	if err != nil {
		log.Printf("Failed to get session ID: %v\n", err)
		return nil
	}
	return session.GetSessionManager().Get(key, sessionID)
}

func (r *Request) Setsession(key string, value interface{}) {
	sessionID, err := r.getSessionID()
	if err != nil {
		log.Printf("Failed to get session ID: %v\n", err)
		return
	}
	session.GetSessionManager().Set(key, value, sessionID)
}

func (r *Request) Deletesession(key string) {
	sessionID, err := r.getSessionID()
	if err != nil {
		log.Printf("Failed to get session ID: %v\n", err)
		return
	}
	session.GetSessionManager().Delete(key, sessionID)
}
