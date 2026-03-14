package cute

import (
	"context"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: HTTPTestMaker Level Handlers - Verify handlers are copied from options to HTTPTestMaker
func TestHTTPTestMaker_StoresRequestInformation(t *testing.T) {
	t.Parallel()

	handler := func(req *http.Request) ([]*allure.Parameter, error) {
		return allure.NewParameters("maker_key", "maker_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithRequestInformation(handler),
	)

	// Verify handlers are stored
	require.Len(t, maker.requestInformation, 1)
}

func TestHTTPTestMaker_StoresRequestInformationT(t *testing.T) {
	t.Parallel()

	handler := func(tProvider T, req *http.Request) ([]*allure.Parameter, error) {
		return allure.NewParameters("maker_t_key", "maker_t_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithRequestInformationT(handler),
	)

	// Verify handlers are stored
	require.Len(t, maker.requestInformationT, 1)
}

func TestHTTPTestMaker_StoresResponseInformation(t *testing.T) {
	t.Parallel()

	handler := func(resp *http.Response) ([]*allure.Parameter, error) {
		return allure.NewParameters("maker_resp_key", "maker_resp_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithResponseInformation(handler),
	)

	// Verify handlers are stored
	require.Len(t, maker.responseInformation, 1)
}

func TestHTTPTestMaker_StoresResponseInformationT(t *testing.T) {
	t.Parallel()

	handler := func(tProvider T, resp *http.Response) ([]*allure.Parameter, error) {
		return allure.NewParameters("maker_resp_t_key", "maker_resp_t_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithResponseInformationT(handler),
	)

	// Verify handlers are stored
	require.Len(t, maker.responseInformationT, 1)
}

// Test 2: Builder Level Handlers
func TestBuilder_RequestInformation(t *testing.T) {
	t.Parallel()

	maker := NewHTTPTestMaker()

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	maker.NewTestBuilder().
		Create().
		RequestInformation(func(req *http.Request) ([]*allure.Parameter, error) {
			return allure.NewParameters("builder_key", "builder_value"), nil
		}).
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress),
		).
		ExecuteTest(context.Background(), capturedT)

	// Verify: builder_key parameter appears
	assert.NotZero(t, capturedT.CapturedParamsCount())
	assert.Equal(t, "builder_value", capturedT.findParamByName("builder_key"))
}

func TestBuilder_ResponseInformation(t *testing.T) {
	t.Parallel()

	maker := NewHTTPTestMaker()

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	maker.NewTestBuilder().
		Create().
		ResponseInformation(func(resp *http.Response) ([]*allure.Parameter, error) {
			return allure.NewParameters("builder_resp_key", "builder_resp_value"), nil
		}).
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress),
		).
		ExecuteTest(context.Background(), capturedT)

	// Verify: builder_resp_key parameter appears
	assert.NotZero(t, capturedT.CapturedParamsCount())
	assert.Equal(t, "builder_resp_value", capturedT.findParamByName("builder_resp_key"))
}

// Test 3: Verify handlers are propagated from maker to test during builder creation
func TestHTTPTestMaker_PropagatesToTest(t *testing.T) {
	t.Parallel()

	globalHandler := func(req *http.Request) ([]*allure.Parameter, error) {
		return allure.NewParameters("global_key", "global_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithRequestInformation(globalHandler),
	)

	// Get the test created by NewTestBuilder
	cuteBuilder := maker.NewTestBuilder().Create().RequestBuilder(WithURI("")).(*cute)
	test := cuteBuilder.tests[0]

	// Verify: maker handlers are propagated to the test
	require.Len(t, test.RequestInformation, 1)
	require.NotNil(t, test.RequestInformation[0])
}

// Test 4: Table Test with fillProps()
func TestTableTest_WithMakerHandlers(t *testing.T) {
	t.Parallel()

	makerHandler := func(resp *http.Response) ([]*allure.Parameter, error) {
		return allure.NewParameters("maker_param", "maker_value"), nil
	}

	test1Handler := func(resp *http.Response) ([]*allure.Parameter, error) {
		return allure.NewParameters("test1_param", "test1_value"), nil
	}

	maker := NewHTTPTestMaker(
		WithResponseInformation(makerHandler),
	)

	tests := []*Test{
		{
			Name: "test1",
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress + "/test"),
				},
			},
			ResponseInformation: []ResponseInformation{test1Handler},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	maker.NewTestBuilder().
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), capturedT)

	// Verify: Both makerHandler and test1Handler execute
	assert.Equal(t, 2, capturedT.CapturedParamsCount())
	assert.Equal(t, "maker_value", capturedT.findParamByName("maker_param"))
	assert.Equal(t, "test1_value", capturedT.findParamByName("test1_param"))
}

// Test 5: Per-Table-Test Handlers
func TestTableTest_WithTableSpecificHandlers(t *testing.T) {
	t.Parallel()

	table1Tests := []*Test{
		{
			Name: "table1_test1",
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress + "/test"),
				},
			},
		},
		{
			Name: "table1_test2",
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress + "/test"),
				},
			},
		},
	}

	table2Tests := []*Test{
		{
			Name: "table2_test1",
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress + "/test"),
				},
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	NewHTTPTestMaker(
		WithRequestInformation(
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("global_key", "global_value"), nil
			})).
		NewTestBuilder().
		CreateTableTest().
		RequestInformationTable(
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("table1_key", "table1_value"), nil
			}).
		PutTests(table1Tests...).
		NextTest().
		CreateTableTest().
		RequestInformationTable(
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("table2_key", "table2_value"), nil
			}).
		PutTests(table2Tests...).
		ExecuteTest(context.Background(), capturedT)

	// For table1 tests: should have global + table1 handlers
	// For table2 tests: should have global + table2 handlers (not table1)
	// Total parameters: 2 (table1_test1) + 2 (table1_test2) + 2 (table2_test1) = 6
	// Each test should have: global_key + either table1_key or table2_key
	assert.Equal(t, 6, capturedT.CapturedParamsCount())

	table1Test1Params := capturedT.findAllParamsByTestName("table1_test1")
	table1Test2Params := capturedT.findAllParamsByTestName("table1_test2")
	table2Test1Params := capturedT.findAllParamsByTestName("table2_test1")

	assert.Len(t, table1Test1Params, 2)
	assert.Equal(t, "global_value", findParameterValue(table1Test1Params, "global_key"))
	assert.Equal(t, "table1_value", findParameterValue(table1Test1Params, "table1_key"))

	assert.Len(t, table1Test2Params, 2)
	assert.Equal(t, "global_value", findParameterValue(table1Test2Params, "global_key"))
	assert.Equal(t, "table1_value", findParameterValue(table1Test2Params, "table1_key"))

	assert.Len(t, table2Test1Params, 2)
	assert.Equal(t, "global_value", findParameterValue(table2Test1Params, "global_key"))
	assert.Equal(t, "table2_value", findParameterValue(table2Test1Params, "table2_key"))
}

// Test 6: Integration with existing test methods
func TestBuilder_ChainableWithOtherMethods(t *testing.T) {
	t.Parallel()

	maker := NewHTTPTestMaker()

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	maker.NewTestBuilder().
		Create().
		RequestInformation(func(req *http.Request) ([]*allure.Parameter, error) {
			return allure.NewParameters("chained_key", "chained_value"), nil
		}).
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress+"/test"),
		).
		ExpectStatus(200).
		ExecuteTest(context.Background(), capturedT)

	// Verify: chained_key parameter appears
	assert.NotZero(t, capturedT.CapturedParamsCount())
	assert.Equal(t, "chained_value", capturedT.findParamByName("chained_key"))
}

// Test 7: Multiple handlers at different levels
func TestMultipleLevels_AllHandlersExecute(t *testing.T) {
	t.Parallel()

	maker := NewHTTPTestMaker(
		WithRequestInformation(func(req *http.Request) ([]*allure.Parameter, error) {
			return allure.NewParameters("maker_key", "maker_value"), nil
		}),
	)

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	results := maker.NewTestBuilder().
		Create().
		RequestInformation(func(req *http.Request) ([]*allure.Parameter, error) {
			return allure.NewParameters("builder_key", "builder_value"), nil
		}).
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress+"/test"),
		).
		ExecuteTest(context.Background(), capturedT)

	require.NotNil(t, results)
	require.NotEmpty(t, results)
	// Verify: both maker and builder handlers execute
	assert.Equal(t, 2, capturedT.CapturedParamsCount())
	assert.Equal(t, "maker_value", capturedT.findParamByName("maker_key"))
	assert.Equal(t, "builder_value", capturedT.findParamByName("builder_key"))
}
