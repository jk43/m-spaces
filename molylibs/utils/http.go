package utils

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SuccessResponse struct {
	Result any `json:"result"`
}

const (
	SUCCESS               = "success"
	ERROR                 = "error"
	VERIFICATION_REQUIRED = "verification_required"
)

type Response struct {
	Result string `json:"result"`
	Data   any    `json:"data"`
}

type Payload[T any] struct {
	Data T `json:"data"`
}

type FormBuilderField struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type FormBuilderFieldRequired struct {
	Name  string `json:"name"`
	Value any    `json:"value" validate:"required`
}

type FormBuilderFields = map[string]FormBuilderField

func (p *Payload[T]) GetData() *T {
	return &p.Data
}

func NewErrorDetails(structField, field, err string, code ErrorCode) (ErrorDetails, ErrorID) {
	eid := GetErrorID()
	return ErrorDetails{
		StructField: structField,
		Field:       field,
		Code:        code,
		Error:       err,
		ID:          eid,
	}, eid
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	body := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	defer body.Close()

	dec := json.NewDecoder(body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only a single JSON value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data Response, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func GetErrorID() ErrorID {
	sum := sha256.Sum256([]byte(fmt.Sprint(time.Now().UnixNano())))
	return ErrorID(fmt.Sprintf("%x", sum))
}

type SimpleHttpClient struct {
	Method   string
	URL      string
	Body     io.Reader
	Header   http.Header
	Testing  bool           //for testing
	Error    error          //for testing
	Response *http.Response //for testing
}

func (c *SimpleHttpClient) Send() (*http.Response, error) {
	if c.Testing {
		return c.SendTesting()
	}
	req, err := http.NewRequest(c.Method, c.URL, c.Body)
	if err != nil {
		return nil, err
	}
	req.Header = c.Header
	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}

func (c *SimpleHttpClient) SendTesting() (*http.Response, error) {
	return c.Response, c.Error
}

// Validate login over JWT
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.Split(authHeader, " ")[1]
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		host, _ := GetXHost(r)
		claims, err := GetJWTClaims(token, host)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
		return
	})
}

// Parse JWT Claims
func ParseJWTClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = r.WithContext(context.WithValue(r.Context(), "claims", nil))
			next.ServeHTTP(w, r)
			return
		}
		token := strings.Split(authHeader, " ")[1]
		if token == "" {
			r = r.WithContext(context.WithValue(r.Context(), "claims", nil))
			next.ServeHTTP(w, r)
			return
		}
		host, _ := GetXHost(r)
		claims, err := GetJWTClaims(token, host)
		if err != nil {
			NewPreDefinedHttpError(UnauthorizedRequest, ErrorCodeFailedJWT, w, err)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
		return
	})
}

// Please enable middleware for this function
// mux.Use(ParseJWTClaims)
func GetClaims(r *http.Request) *JWTClaims {
	value := r.Context().Value("claims")
	if value == nil {
		return nil
	}
	claims := r.Context().Value("claims").(*JWTClaims)
	return claims
}

// type HttpTest struct {
// 	Method  string
// 	Target  string
// 	Payload string
// }

// func (ht *HttpTest) SimpleHttpTest(method string)
func SimpleHttpTest(method string, hdr http.HandlerFunc, body string) (*httptest.ResponseRecorder, *http.Request) {
	var reader io.Reader
	reader = strings.NewReader(body)
	req := httptest.NewRequest(method, "/", reader)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hdr)
	handler.ServeHTTP(rr, req)
	return rr, req
}

type Host = string
type OrgID = string
type UserID = string

func GetXHost(r *http.Request) (Host, error) {
	if len(r.Header["X-Host"]) > 0 {
		return r.Header["X-Host"][0], nil
	}
	// for local development only
	if os.Getenv("ENV") == "local-dev" {
		return r.Host, nil
	}
	return "", fmt.Errorf("Empty host")
}

type Request[T any] struct {
	Payload *Payload[T]
	Claims  *JWTClaims
	UserID  UserID
	OrgID   OrgID
	Host    Host
	Role    UserRole
	IP      string
}

func Par(r *http.Request, p []string) map[string]string {
	params := make(map[string]string)
	for _, s := range p {
		params[s] = chi.URLParam(r, s)
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

type ChiParam string

func ParseGetRequest(w http.ResponseWriter, r *http.Request, p ...ChiParam) (*Request[MapStringSlice[[]string]], error) {
	var err error
	req := Request[MapStringSlice[[]string]]{}
	req.Host, _ = GetXHost(r)
	req.IP = r.Header.Get("X-Forwarded")
	claims := GetClaims(r)
	if claims == nil {
		//guest
		req.Role = RoleGuest
		// due to the import cycle, we cannot use service.GetOrgID here
	} else {
		req.Claims = GetClaims(r)
		req.UserID = req.Claims.Subject
		req.OrgID = req.Claims.OrganizationID
		req.Role = req.Claims.Role
		if len(req.Claims.Audience) > 0 && req.Host != req.Claims.Audience[0] {
			fmt.Println("Disabled for websocket debugging")
			//return nil, errors.New("Host and Audience do not match")
		}
	}
	if err != nil {
		return nil, err
	}
	req.Payload = &Payload[MapStringSlice[[]string]]{
		Data: make(map[string][]string),
	}

	var payload Payload[MapStringSlice[[]string]]
	payload.Data = make(map[string][]string)
	req.Payload = &payload

	for _, s := range p {
		req.Payload.Data[string(s)] = []string{chi.URLParam(r, string(s))}
	}

	for k, v := range r.URL.Query() {
		req.Payload.Data[k] = strings.Split(v[0], ",")
	}

	ReadJSON(w, r, req.Payload)
	return &req, nil
}

func ParseRequest[T any](w http.ResponseWriter, r *http.Request, p ...string) (*Request[T], error) {
	var err error
	if r.Method == "GET" {
		return nil, errors.New("method with GET should call ParseGetRequest")
	}
	req := Request[T]{}
	req.Host, _ = GetXHost(r)
	req.IP = r.Header.Get("X-Forwarded")
	claims := GetClaims(r)
	if claims == nil {
		req.Role = RoleGuest
	} else {
		req.Claims = GetClaims(r)
		req.UserID = req.Claims.Subject
		req.OrgID = req.Claims.OrganizationID
		req.Role = req.Claims.Role
		if req.Host != req.Claims.Audience[0] {
			return nil, errors.New("Host and Audience do not match")
		}
	}
	if err != nil {
		return nil, err
	}
	req.Payload = &Payload[T]{
		Data: *new(T),
	}
	ReadJSON(w, r, req.Payload)
	return &req, nil
}

type VerificationInstruction struct {
	URL     string `json:"url"`
	KeyName string `json:"keyName"`
	Method  string `json:"method"`
	Message string `json:"message"`
}

type Pagination struct {
	Page        int64  `schema:"page"`
	RowsPerPage int64  `schema:"rowsPerPage"`
	RowsNumber  int64  `schema:"rowsNumber"`
	SortBy      string `schema:"sortBy"`
	Descending  bool   `schema:"descending"`
}

func (p *Pagination) GetMongoOptions() *options.FindOptions {
	order := -1
	if !p.Descending {
		order = 1
	}
	opts := options.Find()
	opts.SetLimit(p.RowsPerPage)
	opts.SetSkip(p.RowsPerPage * (p.Page - 1))
	opts.SetSort(bson.D{{Key: p.SortBy, Value: order}})
	return opts
}

// QueryExtractor is used to extract query parameters from the request and store them in a struct
func QueryExtractor[T any](q *MapStringSlice[[]string], e *T) map[string]string {
	t := reflect.TypeOf(*e)
	totalFields := t.NumField()
	keyToRemove := make([]string, totalFields)
	for i := 0; i < totalFields; i++ {
		field := t.Field(i)
		keyToRemove[i] = field.Tag.Get("schema")
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.Decode(e, *q) // Convert *QueryValue[[]string] to map[string][]string
	for _, k := range keyToRemove {
		delete(*q, k)
	}
	return nil
}
