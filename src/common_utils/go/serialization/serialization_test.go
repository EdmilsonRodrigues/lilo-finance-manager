package serialization_test

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"
)

func TestBindArray(t *testing.T) {
	t.Run("should successfully bind array", func(t *testing.T) {
		assertion := func(value string, value2 string) bool {
			model := TestModel{Value: value, Value2: value2}
			response, err := serialization.BindArray[*TestSerializer]([]TestModel{model})
			if err != nil {
				t.Error(err)
				return false
			}
			if len(response) != 1 {
				t.Errorf("expected 1 response, got %d", len(response))
				return false
			}
			if response[0].Value != value {
				t.Errorf("expected %s, got %v", value, response)
				return false
			}
			if response[0].Value2 != value2 {
				t.Errorf("expected %s, got %v", value2, response)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should fail to bind array if model is not expected", func(t *testing.T) {
		assertion := func(value string) bool {
			_, err := serialization.BindArray[*TestSerializer]([]any{value})
			if err == nil {
				t.Error("expected error, got nil")
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})
}

func TestFilterSerializerFields(t *testing.T) {
	t.Run("should filter fields", func(t *testing.T) {
		assertion := func(value string, value2 string) bool {
			serializer := TestSerializer{Value: value, Value2: value2}
			fields := []string{"value"}
			filtered, err := serialization.FilterSerializerFields(&serializer, fields)
			if err != nil {
				t.Error(err)
				return false
			}
			if len(filtered) != 1 {
				t.Errorf("expected 1 field, got %d", len(filtered))
				return false
			}
			if filtered["value"] != value {
				t.Errorf("expected %s, got %v", value, filtered)
				return false
			}
			if _, exists := filtered["value2"]; exists {
				t.Errorf("expected empty string, got %v", filtered)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should return all if fields is empty", func(t *testing.T) {
		assertion := func(value string, value2 string) bool {
			serializer := TestSerializer{Value: value, Value2: value2}
			filtered, err := serialization.FilterSerializerFields(&serializer, []string{})
			if err != nil {
				t.Error(err)
				return false
			}
			if len(filtered) != 2 {
				t.Errorf("expected 2 fields, got %d", len(filtered))
				return false
			}
			if filtered["value"] != value {
				t.Errorf("expected %s, got %v", value, filtered)
				return false
			}
			if filtered["value2"] != value2 {
				t.Errorf("expected %s, got %v", value2, filtered)
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})
}

func TestCreatePaginatedResponse(t *testing.T) {
	t.Run("should create paginated response", func(t *testing.T) {
		assertion := func(value string, value2 string) bool {
			response, err := serialization.CreatePaginatedResponse(1, 10, 100, serialization.QueryConditions{}, []serialization.Serializer{&TestSerializer{Value: value, Value2: value2}}, []string{})
			if err != nil {
				t.Error(err)
				return false
			}
			if response.Status != "success" {
				t.Errorf("expected success, got %s", response.Status)
				return false
			}
			if response.Data.Page != 1 {
				t.Errorf("expected 1, got %d", response.Data.Page)
				return false
			}
			if response.Data.Total != 10 {
				t.Errorf("expected 10, got %d", response.Data.Total)
				return false
			}
			if response.Data.Size != 10 {
				t.Errorf("expected 10, got %d", response.Data.Size)
				return false
			}
			if response.Data.TotalItems != 100 {
				t.Errorf("expected 100, got %d", response.Data.TotalItems)
				return false
			}
			if len(response.Data.Items) != 1 {
				t.Errorf("expected 1, got %d", len(response.Data.Items))
				return false
			}
			if response.Data.Items[0]["value"] != value {
				t.Errorf("expected %s, got %v", value, response.Data.Items[0])
				return false
			}
			if response.Data.Items[0]["value2"] != value2 {
				t.Errorf("expected %s, got %v", value2, response.Data.Items[0])
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("should create paginated response with no size", func(t *testing.T) {
		assertion := func() bool {

			response, err := serialization.CreatePaginatedResponse(0,0,0, serialization.QueryConditions{}, []serialization.Serializer{}, []string{"Value", "Value2"})
			if err != nil {
				t.Error(err)
				return false
			}
			if response.Status != "success" {
				t.Errorf("expected success, got %s", response.Status)
				return false
			}
			if response.Data.Page != 0 {
				t.Errorf("expected 1, got %d", response.Data.Page)
				return false
			}
			if response.Data.Total != 0 {
				t.Errorf("expected 0, got %d", response.Data.Total)
				return false
			}
			if response.Data.Size != 0 {
				t.Errorf("expected 0, got %d", response.Data.Size)
				return false
			}
			if response.Data.TotalItems != 0 {
				t.Errorf("expected 0, got %d", response.Data.TotalItems)
				return false
			}
			if len(response.Data.Items) != 0 {
				t.Errorf("expected 0, got %d", len(response.Data.Items))
				return false
			}
			return true
		}
		if err := quick.Check(assertion, &quick.Config{
			MaxCount: 1000,
		}); err != nil {
			t.Error(err)
		}
	})
}

type TestSerializer struct {
	Value  string `json:"value"`
	Value2 string `json:"value2"`
}
type TestModel struct {
	Value  string
	Value2 string
}

func (s *TestSerializer) Marshal(filters serialization.QueryConditions) (serialization.JSONResponse, error) {
	return serialization.JSONResponse{
		Status: "success",
		Data: serialization.DataItem{
			"value":  s.Value,
			"value2": s.Value2,
		},
	}, nil
}

func (s *TestSerializer) BindModel(model interface{}) error {
	modelValue, ok := model.(TestModel)
	if !ok {
		return fmt.Errorf("model is not of type TestModel")
	}
	s.Value = modelValue.Value
	s.Value2 = modelValue.Value2
	return nil
}
