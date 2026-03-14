package cute

import (
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"strings"
	"testing"
)

type paramCaptureGetter interface {
	newParamCapturing(name string) *paramCapture
}

// paramCapture holds captured parameters, attachments, and allure metadata - shared across parent and child captureT instances
type paramCapture struct {
	params      []*allure.Parameter
	attachments []*allure.Attachment

	// Allure metadata captures
	allureInfo   AllureInformation
	allureLinks  AllureLinks
	allureLabels AllureLabels
}

// captureT embeds provider.T and captures all parameters passed to WithParameters and attachments
type captureT struct {
	provider.T
	parent   paramCaptureGetter
	captured *paramCapture
	children map[string]*paramCapture
}

// newCaptureT creates a new captureT wrapper
func newCaptureT(t provider.T) *captureT {
	captT := captureT{
		T:        t,
		children: make(map[string]*paramCapture),
	}
	captT.captured = captT.newParamCapturing(t.Name())
	return &captT
}

func (c *captureT) CapturedParamsCount() int {
	var cnt int
	for _, capture := range c.children {
		cnt += len(capture.params)
	}
	return cnt
}

func (c *captureT) CapturedAttachmentsCount() int {
	var cnt int
	for _, capture := range c.children {
		cnt += len(capture.attachments)
	}
	return cnt
}

func (c *captureT) getParams(testName, paramName string) []any {
	child, ok := c.children[testName]
	if !ok {
		return nil
	}
	var result []any
	for _, param := range child.params {
		if param.Name == paramName {
			result = append(result, param)
		}
	}
	return result
}

func (c *captureT) findAllParamsByTestName(testName string) []*allure.Parameter {
	var result []*allure.Parameter
	for childName, capture := range c.children {
		if strings.Contains(childName, testName) {
			result = append(result, capture.params...)
		}
	}
	return result
}

func (c *captureT) findParamByName(paramName string) any {
	for _, capture := range c.children {
		param := findParameterValue(capture.params, paramName)
		if param != nil {
			return param
		}
	}
	return nil
}

func (c *captureT) findAttachmentByName(attachmentName string) *allure.Attachment {
	for _, capture := range c.children {
		param := findAttachment(capture.attachments, attachmentName)
		if param != nil {
			return param
		}
	}
	return nil
}

func (c *captureT) newParamCapturing(name string) *paramCapture {
	if c.parent != nil {
		return c.parent.newParamCapturing(name)
	}
	if _, ok := c.children[name]; ok {
		panic("param capture already exists")
	}
	c.children[name] = new(paramCapture)
	return c.children[name]
}

// WithParameters captures parameters and forwards to embedded provider.T
func (c *captureT) WithParameters(parameters ...*allure.Parameter) {
	c.captured.params = append(c.captured.params, parameters...)
	c.T.WithParameters(parameters...)
}

// WithNewParameters forwards to embedded provider.T
func (c *captureT) WithNewParameters(kv ...interface{}) {
	c.T.WithNewParameters(kv...)
}

// WithAttachments captures attachments and forwards to embedded provider.T
func (c *captureT) WithAttachments(attachments ...*allure.Attachment) {
	c.captured.attachments = append(c.captured.attachments, attachments...)
	c.T.WithAttachments(attachments...)
}

// WithNewAttachment captures attachment and forwards to embedded provider.T
func (c *captureT) WithNewAttachment(name string, mimeType allure.MimeType, content []byte) {
	c.captured.attachments = append(c.captured.attachments, allure.NewAttachment(name, mimeType, content))
	c.T.WithNewAttachment(name, mimeType, content)
}

// WithNewStep captures parameters and forwards to embedded provider.T
func (c *captureT) WithNewStep(name string, step func(sCtx provider.StepCtx), params ...*allure.Parameter) {
	// Wrap the step callback to capture parameters added within the step
	wrappedStep := func(stepCtx provider.StepCtx) {
		wrappedCtx := &captureStepCtx{
			StepCtx: stepCtx,
			parent:  c,
		}
		wrappedCtx.captured = wrappedCtx.newParamCapturing(wrappedCtx.Name() + "/" + name)
		// Capture parameters passed to WithNewStep
		wrappedCtx.captured.params = append(wrappedCtx.captured.params, params...)
		step(wrappedCtx)
	}

	c.T.WithNewStep(name, wrappedStep, params...)
}

// Run wraps the test body to capture nested parameters and attachments
func (c *captureT) Run(testName string, testBody func(provider.T), tags ...string) *allure.Result {
	return c.T.Run(testName, func(innerT provider.T) {
		// Wrap the inner T to capture its parameters and attachments too.
		// Share the same paramCapture pointer so nested params/attachments are propagated back
		wrappedInnerT := &captureT{
			T:      innerT,
			parent: c,
		}
		wrappedInnerT.captured = wrappedInnerT.newParamCapturing(wrappedInnerT.Name())
		testBody(wrappedInnerT)
	}, tags...)
}

// captureStepCtx wraps provider.StepCtx and captures parameters and attachments added within steps
type captureStepCtx struct {
	provider.StepCtx
	parent   paramCaptureGetter
	captured *paramCapture
}

func (cs *captureStepCtx) newParamCapturing(name string) *paramCapture {
	return cs.parent.newParamCapturing(name)
}

// WithParameters captures parameters and forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithParameters(parameters ...*allure.Parameter) {
	cs.captured.params = append(cs.captured.params, parameters...)
	cs.StepCtx.WithParameters(parameters...)
}

// WithAttachments captures attachments and forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithAttachments(attachments ...*allure.Attachment) {
	cs.captured.attachments = append(cs.captured.attachments, attachments...)
	cs.StepCtx.WithAttachments(attachments...)
}

// WithNewAttachment captures attachment and forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithNewAttachment(name string, mimeType allure.MimeType, content []byte) {
	cs.captured.attachments = append(cs.captured.attachments, allure.NewAttachment(name, mimeType, content))
	cs.StepCtx.WithNewAttachment(name, mimeType, content)
}

// WithNewParameters forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithNewParameters(kv ...interface{}) {
	cs.StepCtx.WithNewParameters(kv...)
}

func (cs *captureStepCtx) WithNewStep(name string, step func(sCtx provider.StepCtx), params ...*allure.Parameter) {
	// Wrap the step callback to capture parameters added within the step
	wrappedStep := func(stepCtx provider.StepCtx) {
		wrappedCtx := &captureStepCtx{
			StepCtx: stepCtx,
			parent:  cs,
		}
		wrappedCtx.captured = wrappedCtx.newParamCapturing(wrappedCtx.Name() + "/" + name)
		// Capture parameters passed to WithNewStep
		wrappedCtx.captured.params = append(wrappedCtx.captured.params, params...)
		step(wrappedCtx)
	}

	cs.StepCtx.WithNewStep(name, wrappedStep, params...)
}

// Allure Info methods - infoAllureProvider
func (c *captureT) Title(args ...interface{}) {
	title := fmt.Sprint(args...)
	c.captured.allureInfo.Title = title
	c.T.Title(args...)
}

func (c *captureT) Titlef(format string, args ...interface{}) {
	title := fmt.Sprintf(format, args...)
	c.captured.allureInfo.Title = title
	c.T.Titlef(format, args...)
}

func (c *captureT) Description(args ...interface{}) {
	desc := fmt.Sprint(args...)
	c.captured.allureInfo.Description = desc
	c.T.Description(args...)
}

func (c *captureT) Descriptionf(format string, args ...interface{}) {
	desc := fmt.Sprintf(format, args...)
	c.captured.allureInfo.Description = desc
	c.T.Descriptionf(format, args...)
}

func (c *captureT) Stage(args ...interface{}) {
	stage := fmt.Sprint(args...)
	c.captured.allureInfo.Stage = stage
	c.T.Stage(args...)
}

func (c *captureT) Stagef(format string, args ...interface{}) {
	stage := fmt.Sprintf(format, args...)
	c.captured.allureInfo.Stage = stage
	c.T.Stagef(format, args...)
}

// Allure Labels methods - labelsAllureProvider
func (c *captureT) ID(value string) {
	c.captured.allureLabels.ID = value
	c.T.ID(value)
}

func (c *captureT) AllureID(value string) {
	c.captured.allureLabels.AllureID = value
	c.T.AllureID(value)
}

func (c *captureT) Epic(value string) {
	c.captured.allureLabels.Epic = value
	c.T.Epic(value)
}

func (c *captureT) Layer(value string) {
	c.captured.allureLabels.Layer = value
	c.T.Layer(value)
}

func (c *captureT) AddSuiteLabel(value string) {
	c.captured.allureLabels.SuiteLabel = value
	c.T.AddSuiteLabel(value)
}

func (c *captureT) AddSubSuite(value string) {
	c.captured.allureLabels.SubSuite = value
	c.T.AddSubSuite(value)
}

func (c *captureT) AddParentSuite(value string) {
	c.captured.allureLabels.ParentSuite = value
	c.T.AddParentSuite(value)
}

func (c *captureT) Feature(value string) {
	c.captured.allureLabels.Feature = value
	c.T.Feature(value)
}

func (c *captureT) Story(value string) {
	c.captured.allureLabels.Story = value
	c.T.Story(value)
}

func (c *captureT) Tag(value string) {
	c.captured.allureLabels.Tag = value
	c.T.Tag(value)
}

func (c *captureT) Tags(values ...string) {
	c.captured.allureLabels.Tags = values
	c.T.Tags(values...)
}

func (c *captureT) Severity(value allure.SeverityType) {
	c.captured.allureLabels.Severity = value
	c.T.Severity(value)
}

func (c *captureT) Owner(value string) {
	c.captured.allureLabels.Owner = value
	c.T.Owner(value)
}

func (c *captureT) Lead(value string) {
	c.captured.allureLabels.Lead = value
	c.T.Lead(value)
}

func (c *captureT) Label(label *allure.Label) {
	c.captured.allureLabels.Label = label
	c.T.Label(label)
}

func (c *captureT) Labels(labels ...*allure.Label) {
	c.captured.allureLabels.Labels = labels
	c.T.Labels(labels...)
}

// Allure Links methods - linksAllureProvider
func (c *captureT) SetIssue(issue string) {
	c.captured.allureLinks.Issue = issue
	c.T.SetIssue(issue)
}

func (c *captureT) SetTestCase(testCase string) {
	c.captured.allureLinks.TestCase = testCase
	c.T.SetTestCase(testCase)
}

func (c *captureT) Link(link *allure.Link) {
	c.captured.allureLinks.Link = link
	c.T.Link(link)
}

func (c *captureT) TmsLink(tmsCase string) {
	c.captured.allureLinks.TmsLink = tmsCase
	c.T.TmsLink(tmsCase)
}

func (c *captureT) TmsLinks(tmsCases ...string) {
	c.captured.allureLinks.TmsLinks = tmsCases
	c.T.TmsLinks(tmsCases...)
}

// GetCapturedAllureInfo retrieves captured allure info for a test
func (c *captureT) GetCapturedAllureInfo(testName string) AllureInformation {
	if capture, ok := c.children[testName]; ok {
		return capture.allureInfo
	}
	return AllureInformation{}
}

// GetCapturedAllureLabels retrieves captured allure labels for a test
func (c *captureT) GetCapturedAllureLabels(testName string) AllureLabels {
	if capture, ok := c.children[testName]; ok {
		return capture.allureLabels
	}
	return AllureLabels{}
}

// GetCapturedAllureLinks retrieves captured allure links for a test
func (c *captureT) GetCapturedAllureLinks(testName string) AllureLinks {
	if capture, ok := c.children[testName]; ok {
		return capture.allureLinks
	}
	return AllureLinks{}
}

// Helper function to find parameter by key
func findParameterValue(params []*allure.Parameter, key string) any {
	for _, p := range params {
		if p.Name == key {
			return p.Value
		}
	}
	return nil
}

// Helper function to find attachment by name
func findAttachment(attachments []*allure.Attachment, name string) *allure.Attachment {
	for _, a := range attachments {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func buildCapturedFinalName(t testing.TB, testName string) string {
	return t.Name() + "/#00/" + testName
}
