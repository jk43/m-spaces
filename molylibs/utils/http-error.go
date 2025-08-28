package utils

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ErrorCode int
type ErrorID string

const (
	ErrorCodeNone           ErrorCode = 0
	ErrorCodeEmailBeingUsed ErrorCode = iota + 1000
	ErrorCodeFailedMongo
	ErrorCodeFailedMysql
	ErrorCodeFailedJSON
	ErrorCodeFailedGRPC
	ErrorCodeFailedHashing
	ErrorCodePasswordNotMatch
	ErrorCodeFailedJWT
	ErrorCodeJWTTokenIsMissingInCookie
	ErrorCodeNoEmailVerificationTokenFound
	ErrorCodeEmailHasBeenVerifiedAlready
	ErrorCodeOrgID
	ErrorCodeEmailHasNotVerified
	ErrorCodeXHost
	ErrorCodeUserNotActive
	ErrorCodeUserNotFound
	ErrorCodeUserUpdateFailed
	ErrorCodeGetOrgSettings
	ErrorCodeCacheSet
	ErrorCodeEmptyPassword
	ErrorCodeFile
	ErrorCodeFileValidation
	ErrorCodeS3
	ErrorCodeFailedCasbin
	ErrorCodeMongoZeroMatchedCount
	ErrorCodeFailedParseRequest
	ErrorCodeUnauthorized
	ErrorCodeAccessDeniedFromBroker
	ErrorCodeNoRecordFound
	ErrorCodeServiceHandler
	ErrorCodeSecurityBreach
	ErrorCodeOAuth
	ErrorCodeMFA
	ErrorParsingLoginTokens
	ErrorCodeOAuthUserUpdate
	ErrorCodeCMSAccessViolation
	ErrorCodeVerificationCodeNotMatch
)

type ErrorDetails struct {
	StructField string    `json:"structField"`
	Field       string    `json:"field"`
	Code        ErrorCode `json:"code"`
	Error       string    `json:"error"`
	ID          ErrorID   `json:"id"`
}

var isPreDefinedHttpErrorsLoaded = false

type ErrorTemplate string

var logger *zerolog.Logger = getLogger()

const (
	UnableToProcessRequest        ErrorTemplate = "Unable to process request"
	UnauthorizedRequest           ErrorTemplate = "Unauthorized request"
	InternalServerErr             ErrorTemplate = "Internal server error"
	UnauthorizedRequestFromBroker ErrorTemplate = "Unauthorized request"
)

func loadPreDefinedHttpErrors() {
	errorTemplates.Templates[UnableToProcessRequest] = errorTemplateSet{
		error:    UnableToProcessRequest,
		httpCode: http.StatusInternalServerError,
	}
	errorTemplates.Templates[UnauthorizedRequest] = errorTemplateSet{
		error:    UnauthorizedRequest,
		httpCode: http.StatusUnauthorized,
	}
	errorTemplates.Templates[InternalServerErr] = errorTemplateSet{
		error:    InternalServerErr,
		httpCode: http.StatusInternalServerError,
	}
	errorTemplates.Templates[UnauthorizedRequestFromBroker] = errorTemplateSet{
		error:    UnauthorizedRequestFromBroker,
		httpCode: http.StatusUnauthorized,
	}
	isPreDefinedHttpErrorsLoaded = true
}

func getErrorTemplate(template ErrorTemplate) (ErrorDetails, ErrorID) {
	if !isPreDefinedHttpErrorsLoaded {
		loadPreDefinedHttpErrors()
	}
	tem := errorTemplates.Templates[template]
	return NewErrorDetails("", "", string(tem.error), 0)
}

type errorTemplateSet struct {
	httpCode int
	error    ErrorTemplate
}

type preDefinedHttpErrors struct {
	Templates map[ErrorTemplate]errorTemplateSet
}

var errorTemplates = &preDefinedHttpErrors{
	Templates: make(map[ErrorTemplate]errorTemplateSet),
}

func getPriority(code ErrorCode) ErrorPriority {
	priority := ErrorPriorityNotDefined
	switch code {
	case ErrorCodeNone,
		ErrorCodeEmailBeingUsed,
		ErrorCodePasswordNotMatch,
		ErrorCodeFailedJWT,
		ErrorCodeJWTTokenIsMissingInCookie,
		ErrorCodeNoEmailVerificationTokenFound,
		ErrorCodeEmailHasBeenVerifiedAlready,
		ErrorCodeEmailHasNotVerified,
		ErrorCodeUserNotActive,
		ErrorCodeUserNotFound,
		ErrorCodeEmptyPassword,
		ErrorCodeFileValidation:
		priority = ErrorPriorityLow
	case ErrorCodeFailedMongo,
		ErrorCodeFailedMysql,
		ErrorCodeCacheSet,
		ErrorCodeFile:
		priority = ErrorPriorityMedium
	case ErrorCodeFailedJSON,
		ErrorCodeFailedGRPC,
		ErrorCodeFailedHashing,
		ErrorCodeOrgID,
		ErrorCodeXHost,
		ErrorCodeUserUpdateFailed,
		ErrorCodeGetOrgSettings,
		ErrorCodeS3,
		ErrorCodeServiceHandler,
		ErrorCodeFailedCasbin:
		priority = ErrorPriorityHigh
	case ErrorCodeFailedParseRequest,
		ErrorCodeSecurityBreach,
		ErrorCodeCMSAccessViolation:
		priority = ErrorPrioritySecurity
	}
	return priority
}

// Warpper for NewPreDefinedHttpError
func NewPreDefinedHttpError(template ErrorTemplate, code ErrorCode, w http.ResponseWriter, err error, params ...any) {
	errorTemplate, _ := getErrorTemplate(template)
	errorTemplate.Code = code
	httpError := &HttpError{
		errorCode:    code,
		errorDetails: errorTemplate,
		httpCode:     errorTemplates.Templates[template].httpCode,
		errorID:      errorTemplate.ID,
		error:        err,
	}
	//for the caller in the log
	httpError.callerPc, httpError.callerFile, httpError.callerLine, _ = runtime.Caller(1)
	if err != nil {
		httpError.Log(nil, getPriority(code), params...)
	}
	httpError.Write(w)
}

func NewHttpError(structField, field, msg string, code ErrorCode, HttpCode int, err error) *HttpError {
	e := &HttpError{}
	e.error = err
	e.httpCode = HttpCode
	e.errorDetails, e.errorID = NewErrorDetails(structField, field, msg, code)
	//for the caller in the log
	e.callerPc, e.callerFile, e.callerLine, _ = runtime.Caller(1)
	return e
}

type ErrorPriority int

const ErrorPrioritySecurity ErrorPriority = 0
const ErrorPriorityHigh ErrorPriority = 1
const ErrorPriorityMedium ErrorPriority = 2
const ErrorPriorityLow ErrorPriority = 3
const ErrorPriorityNotDefined ErrorPriority = 4

type HttpError struct {
	errorCode      ErrorCode
	errorID        ErrorID
	errorDetails   ErrorDetails
	error          error
	message        string
	httpCode       int
	callerPc       uintptr
	callerFile     string
	callerLine     int
	callerFuncName string
}

func (e *HttpError) Log(strs map[string]string, priority ErrorPriority, params ...any) *HttpError {
	output := logger.Error().Err(e.error)
	output.Int("httpCode", e.httpCode)
	output.Str("errorId", string(e.errorID))
	output.Int("errorCode", int(e.errorCode))
	output.Str("priority", strconv.Itoa(int(priority)))

	funcName := runtime.FuncForPC(e.callerPc).Name()
	callerParts := strings.Split(e.callerFile, "/")
	caller := callerParts[len(callerParts)-1]

	output.Str("caller", caller+":"+strconv.Itoa(e.callerLine))
	output.Str("callerFunc", funcName)
	for k, v := range strs {
		output.Str(k, v)
	}
	for i, v := range params {
		output.Str("params-"+strconv.Itoa(i), fmt.Sprintf("%+v", v))
	}
	output.Send()
	return e
}

func (e *HttpError) Write(w http.ResponseWriter) {
	//hide error code for security reasons
	e.errorDetails.Code = ErrorCodeNone
	WriteJSON(w, e.httpCode, Response{Result: ERROR, Data: e.errorDetails})
}

// logger for http error.
func getLogger() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := log.
		With().
		Str("service", os.Getenv("SERVICE_NAME")).
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if os.Getenv("ENV") == "dev-local" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return &logger
}
