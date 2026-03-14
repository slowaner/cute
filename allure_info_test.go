package cute

import (
	"context"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

// TableBuilder Tests - Only table tests validate Test struct field behavior
// Note: Test-level AllureInfo/Links/Labels fields are only applied in table tests.
// Plain builder tests would test builder methods, not Test struct field propagation.

func TestTableBuilder_AllureFieldsPerTest(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Create tests with different AllureInfo/Links/Labels
	tests := []*Test{
		{
			Name: "test1",
			AllureInfo: AllureInformation{
				Title:       "Test 1 Title",
				Description: "Test 1 Description",
			},
			AllureLabels: AllureLabels{
				Epic:    "Epic1",
				Feature: "Feature1",
			},
			AllureLinks: AllureLinks{
				TestCase: "TC-1",
			},
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress),
				},
			},
			Expect: &Expect{
				Code: 200,
			},
		},
		{
			Name: "test2",
			AllureInfo: AllureInformation{
				Title:       "Test 2 Title",
				Description: "Test 2 Description",
			},
			AllureLabels: AllureLabels{
				Epic:    "Epic2",
				Feature: "Feature2",
			},
			AllureLinks: AllureLinks{
				TestCase: "TC-2",
			},
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress),
				},
			},
			Expect: &Expect{
				Code: 200,
			},
		},
	}

	NewHTTPTestMaker().
		NewTestBuilder().
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), capturedT)

	// Verify test1 captured data
	captured1Info := capturedT.GetCapturedAllureInfo("test1")
	captured1Labels := capturedT.GetCapturedAllureLabels("test1")
	captured1Links := capturedT.GetCapturedAllureLinks("test1")

	assert.Equal(t, "Test 1 Title", captured1Info.Title)
	assert.Equal(t, "Epic1", captured1Labels.Epic)
	assert.Equal(t, "TC-1", captured1Links.TestCase)

	// Verify test2 captured data
	captured2Info := capturedT.GetCapturedAllureInfo("test2")
	captured2Labels := capturedT.GetCapturedAllureLabels("test2")
	captured2Links := capturedT.GetCapturedAllureLinks("test2")

	assert.Equal(t, "Test 2 Title", captured2Info.Title)
	assert.Equal(t, "Epic2", captured2Labels.Epic)
	assert.Equal(t, "TC-2", captured2Links.TestCase)
}

func TestTableBuilder_CombinedBuilderAndTestFields(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Test with both builder-level and test-level fields
	tests := []*Test{
		{
			Name: "test_with_own_fields",
			AllureLabels: AllureLabels{
				Story: "Test Story Override",
			},
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress),
				},
			},
			Expect: &Expect{
				Code: 200,
			},
		},
	}

	NewHTTPTestMaker().
		NewTestBuilder().
		Epic("Builder Epic").
		Feature("Builder Feature").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), capturedT)

	// Verify both builder and test fields are present
	capturedLabels := capturedT.GetCapturedAllureLabels("test_with_own_fields")

	assert.Equal(t, "Builder Epic", capturedLabels.Epic)
	assert.Equal(t, "Builder Feature", capturedLabels.Feature)
	assert.Equal(t, "Test Story Override", capturedLabels.Story)
}

func TestTableBuilder_TestFieldsOverrideBuilder(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Test with conflicting builder-level and test-level fields
	// Validates that test-level fields override builder-level fields
	tests := []*Test{
		{
			Name: "test_with_overrides",
			AllureInfo: AllureInformation{
				Title:       "Test Level Title",
				Description: "Test Level Description",
			},
			AllureLabels: AllureLabels{
				Epic:    "Test Level Epic",
				Feature: "Test Level Feature",
			},
			AllureLinks: AllureLinks{
				TestCase: "TEST-TC-123",
			},
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress),
				},
			},
			Expect: &Expect{
				Code: 200,
			},
		},
	}

	NewHTTPTestMaker().
		NewTestBuilder().
		Title("Builder Level Title").
		Description("Builder Level Description").
		Epic("Builder Level Epic").
		Feature("Builder Level Feature").
		SetTestCase("BUILDER-TC-456").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), capturedT)

	// Verify test-level fields override builder-level fields
	capturedInfo := capturedT.GetCapturedAllureInfo("test_with_overrides")
	capturedLabels := capturedT.GetCapturedAllureLabels("test_with_overrides")
	capturedLinks := capturedT.GetCapturedAllureLinks("test_with_overrides")

	// Test-level values should win
	assert.Equal(t, "Test Level Title", capturedInfo.Title)
	assert.Equal(t, "Test Level Description", capturedInfo.Description)
	assert.Equal(t, "Test Level Epic", capturedLabels.Epic)
	assert.Equal(t, "Test Level Feature", capturedLabels.Feature)
	assert.Equal(t, "TEST-TC-123", capturedLinks.TestCase)
}

func TestTableBuilder_SliceFieldsCombined(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Test that slice fields are combined, not overridden
	// Builder sets Tags=["builder-tag1", "builder-tag2"]
	// Test sets Tags=["test-tag3"]
	// Result should contain all tags: ["builder-tag1", "builder-tag2", "test-tag3"]
	tests := []*Test{
		{
			Name: "test_combined_slices",
			AllureLabels: AllureLabels{
				Tags: []string{"test-tag3", "test-tag4"},
			},
			AllureLinks: AllureLinks{
				TmsLinks: []string{"TEST-TMS-789", "TEST-TMS-890"},
			},
			Request: &Request{
				Builders: []RequestBuilder{
					WithMethod(http.MethodGet),
					WithURI(testServerAddress),
				},
			},
			Expect: &Expect{
				Code: 200,
			},
		},
	}

	NewHTTPTestMaker().
		NewTestBuilder().
		Tags("builder-tag1", "builder-tag2").
		TmsLinks("BUILDER-TMS-123", "BUILDER-TMS-456").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), capturedT)

	// Verify slice fields are combined
	capturedLabels := capturedT.GetCapturedAllureLabels("test_combined_slices")
	capturedLinks := capturedT.GetCapturedAllureLinks("test_combined_slices")

	// Tags should contain both builder and test tags
	assert.Equal(t, []string{"builder-tag1", "builder-tag2", "test-tag3", "test-tag4"}, capturedLabels.Tags)
	// TmsLinks should contain both builder and test TMS links
	assert.Equal(t, []string{"BUILDER-TMS-123", "BUILDER-TMS-456", "TEST-TMS-789", "TEST-TMS-890"}, capturedLinks.TmsLinks)
}

func TestTableBuilder_TestExecuteAppliesStructValues(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Test that Test.Execute() applies AllureInfo/Links/Labels from struct
	testWithAllure := &Test{
		Name: "direct_execute_test",
		AllureInfo: AllureInformation{
			Title:       "Direct Execute Title",
			Description: "Direct Execute Description",
			Stage:       "Direct Execute Stage",
		},
		AllureLabels: AllureLabels{
			Epic:     "Direct Epic",
			Feature:  "Direct Feature",
			Story:    "Direct Story",
			Severity: "critical",
		},
		AllureLinks: AllureLinks{
			TestCase: "DIRECT-TC-999",
		},
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress),
			},
		},
		Expect: &Expect{
			Code: 200,
		},
	}

	// Execute the test directly
	testWithAllure.Execute(context.Background(), capturedT)

	// Verify allure fields from struct were applied
	capturedInfo := capturedT.GetCapturedAllureInfo("direct_execute_test")
	capturedLabels := capturedT.GetCapturedAllureLabels("direct_execute_test")
	capturedLinks := capturedT.GetCapturedAllureLinks("direct_execute_test")

	assert.Equal(t, "Direct Execute Title", capturedInfo.Title)
	assert.Equal(t, "Direct Execute Description", capturedInfo.Description)
	assert.Equal(t, "Direct Execute Stage", capturedInfo.Stage)
	assert.Equal(t, "Direct Epic", capturedLabels.Epic)
	assert.Equal(t, "Direct Feature", capturedLabels.Feature)
	assert.Equal(t, "Direct Story", capturedLabels.Story)
	assert.Equal(t, allure.SeverityType("critical"), capturedLabels.Severity)
	assert.Equal(t, "DIRECT-TC-999", capturedLinks.TestCase)
}
