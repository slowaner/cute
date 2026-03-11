package cute

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// paramCapture holds captured parameters and attachments - shared across parent and child captureT instances
type paramCapture struct {
	params      []*allure.Parameter
	attachments []*allure.Attachment
}

// captureT embeds provider.T and captures all parameters passed to WithParameters and attachments
type captureT struct {
	provider.T
	captured *paramCapture
}

// newCaptureT creates a new captureT wrapper
func newCaptureT(t provider.T) *captureT {
	return &captureT{
		T: t,
		captured: &paramCapture{
			params:      make([]*allure.Parameter, 0),
			attachments: make([]*allure.Attachment, 0),
		},
	}
}

// captureStepCtx wraps provider.StepCtx and captures parameters and attachments added within steps
type captureStepCtx struct {
	provider.StepCtx
	captured *paramCapture
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
	// Capture parameters passed to WithNewStep
	c.captured.params = append(c.captured.params, params...)

	// Wrap the step callback to capture parameters added within the step
	wrappedStep := func(stepCtx provider.StepCtx) {
		wrappedCtx := &captureStepCtx{
			StepCtx:  stepCtx,
			captured: c.captured,
		}
		step(wrappedCtx)
	}

	c.T.WithNewStep(name, wrappedStep, params...)
}

// Run wraps the test body to capture nested parameters and attachments
func (c *captureT) Run(testName string, testBody func(provider.T), tags ...string) *allure.Result {
	return c.T.Run(testName, func(innerT provider.T) {
		// Wrap the inner T to capture its parameters and attachments too
		// Share the same paramCapture pointer so nested params/attachments are propagated back
		wrappedInnerT := &captureT{
			T:        innerT,
			captured: c.captured,
		}
		testBody(wrappedInnerT)
	}, tags...)
}

// Helper function to find parameter by key
func getParameterValue(params []*allure.Parameter, key string) any {
	for _, p := range params {
		if p.Name == key {
			return p.Value
		}
	}
	return nil
}

// Helper function to find attachment by name
func getAttachment(attachments []*allure.Attachment, name string) *allure.Attachment {
	for _, a := range attachments {
		if a.Name == name {
			return a
		}
	}
	return nil
}
