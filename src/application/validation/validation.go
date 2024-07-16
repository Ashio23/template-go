package validation

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"template-go/src/application/utils"

	"github.com/go-playground/validator"
)

type ApiError struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error"`
}

func (a ApiError) Error() string {
	return a.ErrorMessage
}

func HandleValidationError(w http.ResponseWriter, err error) {
	errorData := ApiError{
		StatusCode:   http.StatusBadRequest,
		ErrorMessage: err.Error(),
	}
	utils.JSONResponse(w, errorData, errorData.StatusCode)
}

func Validate(r *http.Request, structType interface{}) error {
	var err error
	validate := validator.New()
	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "application/xml"), strings.Contains(contentType, "text/xml"):
		err = xml.NewDecoder(r.Body).Decode(structType)
	default:
		err = json.NewDecoder(r.Body).Decode(structType)
	}

	if err != nil {
		return fmt.Errorf("error decoding input: %v", err)
	}

	v := reflect.ValueOf(structType).Elem()
	t := reflect.TypeOf(structType).Elem()
	fields := reflect.VisibleFields(t)

	for _, field := range fields {
		headerTag := field.Tag.Get("header")
		queryTag := field.Tag.Get("query")
		uriTag := field.Tag.Get("uri")

		if headerTag == "" && queryTag == "" && uriTag == "" {
			continue
		}

		if headerTag != "" {
			err = validateInput(headerTag, "header", field, v, r)
		}

		if queryTag != "" {
			err = validateInput(queryTag, "query", field, v, r)
		}

		if uriTag != "" {
			err = validateInput(uriTag, "uri", field, v, r)
		}

		if err != nil {
			return fmt.Errorf("error parsing input: %v", err.Error())
		}
	}

	err = validate.Struct(structType)

	if err != nil {
		return err
	}
	return nil
}

func validateInput(tag string, tagName string, field reflect.StructField, v reflect.Value, r *http.Request) error {
	var input string
	fieldName := field.Name

	switch tagName {
	case "header":
		input = r.Header.Get(tag)
	case "query":
		input = r.URL.Query().Get(tag)
	case "uri":
		input = r.PathValue(tag)
	default:
		return fmt.Errorf("unsupported tag name: %v", tagName)
	}

	if input == "" {
		validate := field.Tag.Get("validate")
		if validate != "-" {
			return fmt.Errorf("couldn't get input data for tag %v: %v", tag, tagName)
		}
	}
	if err := mapInputToStruct(input, fieldName, v); err != nil {
		return err
	}
	return nil
}

func mapInputToStruct(input string, fieldName string, v reflect.Value) error {
	ok := canSetValue(v)
	if !ok {
		return nil
	}

	structField := v.FieldByName(fieldName)

	err := setValue(structField, input, fieldName)
	if err != nil {
		return err
	}
	return nil
}

func canSetValue(v reflect.Value) bool {
	return v.IsValid() && v.CanAddr() && v.CanSet()
}

func setValue(v reflect.Value, input string, fieldName string) error {
	k := v.Kind()
	t := v.Type()

	switch k {
	case reflect.String:
		numMethod := v.NumMethod()
		if numMethod == 0 {
			v.SetString(input)
		}

		_, ok := t.MethodByName("IsValid")
		if ok {
			convertedValue := reflect.ValueOf(input).Convert(t)
			isValidReflect := convertedValue.MethodByName("IsValid").Call(nil)
			isValid := isValidReflect[0].Interface().(bool)
			if !isValid {
				return fmt.Errorf("cannot parse value for field %v: %v is not a valid %v", fieldName, input, t.Name())
			}
			v.SetString(input)

		}
	case reflect.Bool:
		b, err := strconv.ParseBool(input)
		if err != nil {
			return fmt.Errorf("cannot parse bool value for field %v: %w", fieldName, err)
		}

		v.SetBool(b)
	case reflect.Int:
		i, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("cannot parse int value for field %v: %w", fieldName, err)
		}

		v.SetInt(int64(i))
	case reflect.Float64:
		fl, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("cannot parse float64 value for field %v: %w", fieldName, err)
		}

		v.SetFloat(fl)
	default:
		return setSliceValue(t, v, input, fieldName)
	}
	return nil
}

func setSliceValue(t reflect.Type, v reflect.Value, input string, fieldName string) error {
	switch t {
	case reflect.SliceOf(reflect.TypeOf("string")):
		v.Set(reflect.ValueOf([]string{input}))
	case reflect.SliceOf(reflect.TypeOf("bool")):
		b, err := strconv.ParseBool(input)
		if err != nil {
			return fmt.Errorf("cannot parse bool value for field %v: %w", fieldName, err)
		}

		v.Set(reflect.ValueOf([]bool{b}))
	case reflect.SliceOf(reflect.TypeOf("int")):
		i, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("cannot parse int value for field %v: %w", fieldName, err)
		}

		v.Set(reflect.ValueOf([]int{i}))
	case reflect.SliceOf(reflect.TypeOf("float64")):
		fl, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("cannot parse float64 value for field %v: %w", fieldName, err)
		}

		v.Set(reflect.ValueOf([]float64{fl}))
	default:
		sliceInput := strings.Split(input, ",")
		completeSlice := reflect.MakeSlice(t, 0, 0)

		for _, in := range sliceInput {
			convertedValue := reflect.ValueOf(in).Convert(t.Elem())
			isValidReflect := convertedValue.MethodByName("IsValid").Call(nil)
			isValid := isValidReflect[0].Interface().(bool)
			if !isValid {
				customType := t.Elem().Name()
				return fmt.Errorf("cannot parse value for field %v: %v is not a valid %v", fieldName, in, customType)
			}
			completeSlice = reflect.Append(completeSlice, convertedValue)
		}

		v.Set(completeSlice)
	}
	return nil
}
