package cute

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"strings"
	"testing"
)

type paramCaptureGetter interface {
	newParamCapturing(name string) *paramCapture
}

// paramCapture holds captured parameters and attachments - shared across parent and child captureT instances
type paramCapture struct {
	params      []*allure.Parameter
	attachments []*allure.Attachment
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
