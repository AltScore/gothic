package xapi

import (
	"encoding/json"
	"testing"

	"github.com/AltScore/gothic/pkg/xerrors"
	"github.com/AltScore/gothic/pkg/xvalidator"
	"github.com/labstack/echo/v4"
	"github.com/nsf/jsondiff"
)

func TestErrorNormalizerMiddleware(t *testing.T) {
	f := ErrorNormalizerMiddleware()

	tests := []struct {
		name     string
		err      error
		wantJson string
	}{
		{
			name:     "nil",
			err:      nil,
			wantJson: `null`,
		},
		{
			name:     "echo.HTTPError",
			err:      echo.NewHTTPError(400, "Bad Request"),
			wantJson: `{"error":{"code":"invalid-argument","message":"Bad Request"}}`,
		},
		{
			name:     "xerrors.HttpError",
			err:      xerrors.NewDuplicateError("some-entity", "it already exists", "%s", "1234"),
			wantJson: `{"error":{"code":"duplicate","message":"duplicate: some-entity: 1234: it already exists"}}`,
		},
		{
			name: "fieldError",
			err:  makeFieldError(),
			wantJson: `{
	            "error": {
	                "code": "validation_error",
	                "details": [
	                    {
	                        "FailedTag": "required",
	                        "Field": "SomeField",
	                        "Namespace": "testStruct.SomeField",
	                        "Param": "",
	                        "StructNamespace": "testStruct.SomeField",
	                        "Tag": "required",
	                        "Value": ""
	                    }
	                ],
	                "message": "Validation failed"
	            }
	        }`,
		},
		{
			name: "fieldErrors",
			err:  makeFieldErrors(),
			wantJson: `{
                "error": {
                    "code": "validation_error",
                    "details": [        
                        {
                            "FailedTag": "required",    
                            "Field": "SomeField",
                            "Namespace": "testStruct.SomeField",    
                            "Param": "",
                            "StructNamespace": "testStruct.SomeField",  
                            "Tag": "required",  
                            "Value": "" 
                        },
                        {   
                            "FailedTag": "gt",      
                            "Field": "ANumber", 
                            "Namespace": "testStruct.ANumber",  
                            "Param": "5",   
                            "StructNamespace": "testStruct.ANumber",    
                            "Tag": "gt",    
                            "Value": 3  
                        }
                    ],
                    "message": "Validation failed"  
                }
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			next := func(c echo.Context) error {
				return tt.err
			}

			newErr := f(next)(nil)

			if newErr == nil {
				if tt.wantJson != "null" {
					t.Errorf("ErrorNormalizerMiddleware() = %v, want %v", newErr, tt.wantJson)
				}
				return
			}

			options := jsondiff.DefaultConsoleOptions()
			options.SkipMatches = true

			actualBytes, _ := json.Marshal(newErr.(*echo.HTTPError).Message)

			differences, explanation := jsondiff.Compare([]byte(tt.wantJson), actualBytes, &options)

			if differences != jsondiff.FullMatch {
				t.Errorf("ErrorNormalizerMiddleware() = %v, want %v, explanation: %v", string(actualBytes), tt.wantJson, explanation)
			}
		})
	}
}

type testStruct struct {
	SomeField string `json:"some_field" validate:"required"`
	ANumber   int    `json:"a_number" validate:"required,gt=5"`
}

func makeFieldError() error {
	return xvalidator.Struct(testStruct{ANumber: 7})
}

func makeFieldErrors() error {
	return xvalidator.Struct(testStruct{ANumber: 3})
}
