package convert

import (
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestSuccessMapInterfaceToStringMap(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		testData := make(map[string]interface{})
		testData["first"] = 1
		testData["second"] = true
		testData["third"] = "third"

		out, err := InterfaceToStringMap(testData)

		test.AssertThat(t, err, nil)
		test.AssertThat(t, out["first"], "1")
		test.AssertThat(t, out["second"], "true")
		test.AssertThat(t, out["third"], "third")
	})
}

func TestSuccessStructInterfaceToStringMap(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		type testData struct {
			first  float32
			second int32
			third  uint32
			fourth string
		}
		data := testData{1.5, -5, 5, "string"}
		out, err := InterfaceToStringMap(data)
		test.AssertThat(t, err, nil)
		test.AssertThat(t, out["first"], "1.5000000000")
		test.AssertThat(t, out["second"], "-5")
		test.AssertThat(t, out["third"], "5")
		test.AssertThat(t, out["fourth"], "string")
	})
}

func TestFailureStructInterfaceToStringMap(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		type testData struct {
			boolean bool
		}
		data := testData{true}
		_, err := InterfaceToStringMap(data)
		test.AssertThat(t, err, nil, "not")
	})
}

func TestFailureInterfaceToStringMap(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		testData := "testData"
		_, err := InterfaceToStringMap(testData)
		test.AssertThat(t, err, nil, "not")
	})
}
