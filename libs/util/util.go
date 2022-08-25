package util

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"reporting/libs/logger"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/go-rel/rel"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rotisserie/eris"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

var (
	ErrFilterNil   = eris.New("filter is nil")
	ErrFilterEmpty = eris.New("filter is empty")
)

var (
	ErrInvalid = errors.New("invalid")
)

type CTXValue string

const (
	CTXTrackerID  = CTXValue("CTX.Tracker.ID")
	CTXJWTPayload = CTXValue("CTX.JWTPayload")
)

func GetTracker(ctx context.Context) string {
	v, _ := ctx.Value(CTXTrackerID).(string)
	return v
}

func GetJWTPayload(ctx context.Context) JWTPayload {
	v, _ := ctx.Value(CTXJWTPayload).(JWTPayload)
	return v
}

var (
	Render = render.New()
)

type JWTPayload struct {
	UserID     uint64 `json:"user_id"`
	MerchantID uint64 `json:"merchat_id"`
	jwt.RegisteredClaims
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
	Err     any    `json:"error"`
}

func HTTPResponse(rw http.ResponseWriter, code int, message string, payload any, err any) {
	Render.JSON(rw, code, Response{
		Code:    code,
		Message: message,
		Payload: payload,
		Err:     err,
	})
}

func ErrorHTTPResponse(ctx context.Context, rw http.ResponseWriter, err error) {
	var (
		code          = http.StatusInternalServerError
		IsErrNotFound = eris.Is(err, rel.ErrNotFound)
		IsErrInvalid  = eris.Is(err, ErrInvalid)
		unpack        = eris.Unpack(err)
		trackerId     = GetTracker(ctx)
	)

	if !IsErrInvalid && !IsErrNotFound {
		logger.Log.With(
			zap.String("tracker_id", trackerId),
			zap.Any("error", eris.ToJSON(err, true)),
		).Error(unpack.ErrRoot.Msg)

		logger.Console.With(
			zap.String("tracker_id", trackerId),
		).Error(unpack.ErrRoot.Msg)
	}

	Err := map[string]any{
		"tracker_id": GetTracker(ctx),
	}

	if IsErrNotFound {
		code = http.StatusNotFound
		Err = nil
	}

	if IsErrInvalid {
		code = http.StatusBadRequest
		Err = nil
	}

	HTTPResponse(rw, code, unpack.ErrRoot.Msg, nil, Err)
}

func ValidationHTTPResponse(rw http.ResponseWriter, err []ValidationError) {
	var (
		code = http.StatusUnprocessableEntity
		msg  = "unprocessable request body, an error occured"
	)
	HTTPResponse(rw, code, msg, nil, map[string]any{
		"validation": err,
	})
}

func RequestBodyValidation(rw http.ResponseWriter, body io.ReadCloser, data any) bool {
	if data == nil {
		panic("data decode and validation is nil")
	}

	if err := json.NewDecoder(body).Decode(data); err != nil {
		var (
			code = http.StatusBadRequest
			msg  = "parse request body, an error occured"
		)
		HTTPResponse(rw, code, msg, nil, nil)
		return false
	}

	if errs := Validation(data); len(errs) > 0 {
		ValidationHTTPResponse(rw, errs)
		return false
	}

	return true
}

var (
	// Setup Validator
	Validate = validator.New()

	// Setup Validation Message Translation
	enTrans  = en.New()
	uni      = ut.New(enTrans, enTrans)
	Trans, _ = uni.GetTranslator("en")
)

type ValidationError struct {
	Key     string `json:"key"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func NewValidator() {
	// use the names which have been specified for JSON representations of structs, rather than normal Go field names
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Setup Error Message Translation
	en_translations.RegisterDefaultTranslations(Validate, Trans)
}

func Validation(data any) []ValidationError {
	err := Validate.Struct(data)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		panic("error invalid validation error")
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var (
		length     = len(errs)
		lengthName = 0
	)

	if length > 0 {
		rType := reflect.TypeOf(data)
		if rType.Kind() == reflect.Ptr {
			rType = rType.Elem()
		}
		lengthName = len(rType.Name())
	}

	val := make([]ValidationError, length)
	for i := 0; i < len(errs); i++ {
		// cut string namespace
		// note: add a `+1` to cut dot (.)
		key := errs[i].Namespace()[lengthName+1:]
		val[i] = ValidationError{
			Key:     key,
			Rule:    errs[i].Tag(),
			Message: errs[i].Translate(Trans),
		}
	}

	return val
}
