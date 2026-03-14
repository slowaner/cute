package cute

import (
	"context"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestPlainBuilder_AllureInfoPropagation(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	NewHTTPTestMaker().
		NewTestBuilder().
		Title("Test Title").
		Description("Test Description").
		Stage("Test Stage").
		Create().
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress),
		).
		ExpectStatus(200).
		ExecuteTest(context.Background(), capturedT)

	// Verify captured allure info
	testName := capturedT.Name()
	capturedInfo := capturedT.GetCapturedAllureInfo(testName)

	res := allureT.GetResult()
	t.Log(res.Name)
	assert.Equal(t, "Test Title", capturedInfo.Title)
	assert.Equal(t, "Test Description", capturedInfo.Description)
	assert.Equal(t, "Test Stage", capturedInfo.Stage)
}

func TestPlainBuilder_AllureLabelsPropagation(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	NewHTTPTestMaker().
		NewTestBuilder().
		Epic("Test Epic").
		Feature("Test Feature").
		Story("Test Story").
		Tags("tag1", "tag2").
		Severity(allure.CRITICAL).
		Owner("test-owner").
		Create().
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress),
		).
		ExpectStatus(200).
		ExecuteTest(context.Background(), capturedT)

	// Verify captured allure labels
	testName := capturedT.Name()
	capturedLabels := capturedT.GetCapturedAllureLabels(testName)

	assert.Equal(t, "Test Epic", capturedLabels.Epic)
	assert.Equal(t, "Test Feature", capturedLabels.Feature)
	assert.Equal(t, "Test Story", capturedLabels.Story)
	assert.Equal(t, []string{"tag1", "tag2"}, capturedLabels.Tags)
	assert.Equal(t, allure.CRITICAL, capturedLabels.Severity)
	assert.Equal(t, "test-owner", capturedLabels.Owner)
}

func TestPlainBuilder_AllureLinksPropagation(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	testLink := &allure.Link{
		Name: "TestLink",
		URL:  "http://example.com",
		Type: "issue",
	}

	NewHTTPTestMaker().
		NewTestBuilder().
		SetIssue("ISSUE-123").
		SetTestCase("TC-456").
		Link(testLink).
		TmsLink("TMS-789").
		Create().
		RequestBuilder(
			WithMethod(http.MethodGet),
			WithURI(testServerAddress),
		).
		ExpectStatus(200).
		ExecuteTest(context.Background(), capturedT)

	// Verify captured allure links
	testName := capturedT.Name()
	capturedLinks := capturedT.GetCapturedAllureLinks(testName)

	assert.Equal(t, "ISSUE-123", capturedLinks.Issue)
	assert.Equal(t, "TC-456", capturedLinks.TestCase)
	assert.Equal(t, testLink, capturedLinks.Link)
	assert.Equal(t, "TMS-789", capturedLinks.TmsLink)
}

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
				Issue: "ISSUE-1",
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
				Issue: "ISSUE-2",
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
	test1Name := buildCapturedFinalName(t, "test1")
	captured1Info := capturedT.GetCapturedAllureInfo(test1Name)
	captured1Labels := capturedT.GetCapturedAllureLabels(test1Name)
	captured1Links := capturedT.GetCapturedAllureLinks(test1Name)

	assert.Equal(t, "Test 1 Title", captured1Info.Title)
	assert.Equal(t, "Epic1", captured1Labels.Epic)
	assert.Equal(t, "ISSUE-1", captured1Links.Issue)

	// Verify test2 captured data
	test2Name := buildCapturedFinalName(t, "test2")
	captured2Info := capturedT.GetCapturedAllureInfo(test2Name)
	captured2Labels := capturedT.GetCapturedAllureLabels(test2Name)
	captured2Links := capturedT.GetCapturedAllureLinks(test2Name)

	assert.Equal(t, "Test 2 Title", captured2Info.Title)
	assert.Equal(t, "Epic2", captured2Labels.Epic)
	assert.Equal(t, "ISSUE-2", captured2Links.Issue)
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
	testName := buildCapturedFinalName(t, "test_with_own_fields")
	capturedLabels := capturedT.GetCapturedAllureLabels(testName)

	assert.Equal(t, "Builder Epic", capturedLabels.Epic)
	assert.Equal(t, "Builder Feature", capturedLabels.Feature)
	assert.Equal(t, "Test Story Override", capturedLabels.Story)
}
