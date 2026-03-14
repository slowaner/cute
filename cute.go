package cute

import (
	"context"
	"strings"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/allure_manager/manager"
	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type cute struct {
	baseProps *HTTPTestMaker

	parallel bool

	allureInfo   *AllureInformation
	allureLinks  *AllureLinks
	allureLabels *AllureLabels

	countTests int // Общее количество тестов.

	isTableTest bool
	tests       []*Test

	// stored information callbacks for next added tests
	reqInfo   []RequestInformation
	reqInfoT  []RequestInformationT
	respInfo  []ResponseInformation
	respInfoT []ResponseInformationT
}

// AllureInformation stores test description metadata for test reports
// Fields are applied during test execution and override builder-level info in table tests
type AllureInformation struct {
	Title       string
	Description string
	Stage       string
}

// AllureLabels stores label metadata for test reports
// Scalar fields (Epic, Feature, Story, etc.) override builder-level values in table tests
// Slice fields (Tags, Labels) are combined with builder-level values for table tests
type AllureLabels struct {
	ID          string
	Feature     string
	Epic        string
	Tag         string
	Tags        []string
	SuiteLabel  string
	SubSuite    string
	ParentSuite string
	Story       string
	Severity    allure.SeverityType
	Owner       string
	Lead        string
	Label       *allure.Label
	Labels      []*allure.Label
	AllureID    string
	Layer       string
}

// AllureLinks stores link metadata for test reports
// Fields are applied during test execution and combined with builder-level links for table tests
type AllureLinks struct {
	Issue    string
	TestCase string
	Link     *allure.Link
	TmsLink  string
	TmsLinks []string
}

func (qt *cute) ExecuteTest(ctx context.Context, t tProvider) []ResultsHTTPBuilder {
	var internalT allureProvider

	if t == nil {
		panic("could not start test without testing.T")
	}

	stepCtx, isStepCtx := t.(provider.StepCtx)
	if isStepCtx {
		return qt.executeTestsInsideStep(ctx, stepCtx)
	}

	switch v := t.(type) {
	case provider.T:
		internalT = v
	case *testing.T:
		newT := createAllureT(v)
		if !qt.isTableTest {
			defer newT.FinishTest() //nolint
		}
		internalT = newT
	default:
		panic("could not start test without testing.T or provider.T")
	}

	if qt.parallel {
		internalT.Parallel()
	}

	return qt.executeTests(ctx, internalT)
}

func createAllureT(t *testing.T) *common.Common {
	var (
		newT        = common.NewT(t)
		callers     = strings.Split(t.Name(), "/")
		providerCfg = manager.NewProviderConfig().
				WithFullName(t.Name()).
				WithPackageName("package").
				WithSuiteName(t.Name()).
				WithRunner(callers[0])
		newProvider = manager.NewProvider(providerCfg)
	)

	newProvider.NewTest(t.Name(), "package")

	newT.SetProvider(newProvider)
	newT.Provider.TestContext()

	return newT
}

// executeTests is method for run tests
// It's could be table tests or usual tests
func (qt *cute) executeTests(ctx context.Context, allureProvider allureProvider) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= qt.countTests; i++ {
		currentTest := qt.tests[i]

		// Execute by new T for table tests
		if qt.isTableTest {
			tableTestName := currentTest.Name

			allureProvider.Run(tableTestName, func(inT provider.T) {
				// Set current test info
				qt.setAllureInformation(inT)
				inT.Title(tableTestName)
				currentTest.setAllureInformation(inT)

				res = append(res, qt.executeInsideAllure(ctx, inT, currentTest))
			})
		} else {
			currentTest.Name = allureProvider.Name()

			// set labels
			qt.setAllureInformation(allureProvider)

			res = append(res, qt.executeInsideAllure(ctx, allureProvider, currentTest))
		}
	}

	return res
}

// executeInsideAllure is method for run test inside allure
// It's could be table tests or usual tests
func (qt *cute) executeInsideAllure(ctx context.Context, allureProvider allureProvider, currentTest *Test) ResultsHTTPBuilder {
	resT := currentTest.executeInsideAllure(ctx, allureProvider)

	// Remove from base struct all asserts
	currentTest.clearFields()

	return resT
}

// executeTestsInsideStep is method for run group of tests inside provider.StepCtx
func (qt *cute) executeTestsInsideStep(ctx context.Context, stepCtx provider.StepCtx) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= qt.countTests; i++ {
		currentTest := qt.tests[i]

		result := currentTest.executeInsideStep(ctx, stepCtx)

		// Remove from base struct all asserts
		currentTest.clearFields()

		res = append(res, result)
	}

	return res
}
