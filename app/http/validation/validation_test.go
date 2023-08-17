package validation

import (
	"testing"
)

func TestValidation_Validate(t *testing.T) {
	v := New()

	// Define a test struct
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	// Create an instance of the test struct
	testData := TestStruct{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	// Perform validation
	errors, err := v.Validate(testData)

	if err != nil {
		t.Fatalf("Validation error: %v", err)
	}

	if len(errors) != 0 {
		t.Errorf("Expected no validation errors, but got %d errors", len(errors))
	}
}

func TestValidation_Validate_WithErrors(t *testing.T) {
	v := New()

	// Define a test struct
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	// Create an instance of the test struct with invalid data
	testData := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	// Perform validation
	errors, err := v.Validate(testData)

	if err == nil {
		t.Fatal("Expected validation error, but got nil")
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 validation errors, but got %d errors", len(errors))
	}
}

func TestValidation_CreateErrorResponse(t *testing.T) {
	v := New()

	// Create a slice of validation errors
	validationErrors := []ValidationError{
		{Field: "Name", Error: "Field validation for Name failed on required"},
		{Field: "Email", Error: "Field validation for Email failed on email"},
	}

	// Create an error response
	errorResponse := v.CreateErrorResponse(validationErrors)

	// Verify the error response
	expectedMessage := "Validation error. (2 errors)"
	if errorResponse.Message != expectedMessage {
		t.Errorf("Expected error message %q, but got %q", expectedMessage, errorResponse.Message)
	}

	expectedErrors := map[string][]string{
		"Name":  {"Field validation for Name failed on required"},
		"Email": {"Field validation for Email failed on email"},
	}
	if !compareErrorMaps(errorResponse.Errors, expectedErrors) {
		t.Errorf("Expected error map %+v, but got %+v", expectedErrors, errorResponse.Errors)
	}
}

func compareErrorMaps(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, aErrors := range a {
		bErrors, ok := b[key]
		if !ok {
			return false
		}
		if len(aErrors) != len(bErrors) {
			return false
		}
		for i, aError := range aErrors {
			if aError != bErrors[i] {
				return false
			}
		}
	}
	return true
}
