package cute

import (
	"net/http"
	"slices"
)

func (qt *cute) CreateTableTest() MiddlewareTable {
	qt.isTableTest = true

	return qt
}

func (qt *cute) PutNewTest(name string, r *http.Request, expect *Expect) TableTest {
	// Validate, that first step is empty
	if qt.tests[qt.countTests].Request.Base == nil &&
		len(qt.tests[qt.countTests].Request.Builders) == 0 {
		qt.tests[qt.countTests].Expect = expect
		qt.tests[qt.countTests].Name = name
		qt.tests[qt.countTests].Request.Base = r

		return qt
	}

	newTest := createDefaultTest(qt.baseProps, qt.reqInfo, qt.reqInfoT, qt.respInfo, qt.respInfoT)
	newTest.Expect = expect
	newTest.Name = name
	newTest.Request.Base = r
	qt.tests = append(qt.tests, newTest)
	qt.countTests++ // async?

	return qt
}

func (qt *cute) PutTests(tests ...*Test) TableTest {
	for _, test := range tests {
		// Fill common fields
		qt.fillProps(test)

		// Validate, that first step is empty
		if qt.tests[qt.countTests].Request.Base == nil &&
			len(qt.tests[qt.countTests].Request.Builders) == 0 {
			qt.tests[qt.countTests] = test

			continue
		}

		qt.tests = append(qt.tests, test)
		qt.countTests++
	}

	return qt
}

func (qt *cute) fillProps(t *Test) {
	var (
		baseRequestInformation   []RequestInformation
		baseRequestInformationT  []RequestInformationT
		baseResponseInformation  []ResponseInformation
		baseResponseInformationT []ResponseInformationT
	)
	if qt.baseProps != nil {
		if qt.baseProps.httpClient != nil {
			t.httpClient = qt.baseProps.httpClient
		}

		if qt.baseProps.jsonMarshaler != nil {
			t.jsonMarshaler = qt.baseProps.jsonMarshaler
		}

		if t.Middleware == nil {
			t.Middleware = createMiddlewareFromTemplate(qt.baseProps.middleware)
		} else {
			t.Middleware.After = append(t.Middleware.After, qt.baseProps.middleware.After...)
			t.Middleware.AfterT = append(t.Middleware.AfterT, qt.baseProps.middleware.AfterT...)
			t.Middleware.Before = append(t.Middleware.Before, qt.baseProps.middleware.Before...)
			t.Middleware.BeforeT = append(t.Middleware.BeforeT, qt.baseProps.middleware.BeforeT...)
		}
		if len(qt.baseProps.requestInformation) > 0 {
			baseRequestInformation = qt.baseProps.requestInformation
		}
		if len(qt.baseProps.requestInformationT) > 0 {
			baseRequestInformationT = qt.baseProps.requestInformationT
		}
		if len(qt.baseProps.responseInformation) > 0 {
			baseResponseInformation = qt.baseProps.responseInformation
		}
		if len(qt.baseProps.responseInformationT) > 0 {
			baseResponseInformationT = qt.baseProps.responseInformationT
		}
	}

	// Merge information handlers (APPEND pattern - combine builder + test handlers)
	t.RequestInformation = slices.Concat(baseRequestInformation, qt.reqInfo, t.RequestInformation)
	t.RequestInformationT = slices.Concat(baseRequestInformationT, qt.reqInfoT, t.RequestInformationT)
	t.ResponseInformation = slices.Concat(baseResponseInformation, qt.respInfo, t.ResponseInformation)
	t.ResponseInformationT = slices.Concat(baseResponseInformationT, qt.respInfoT, t.ResponseInformationT)
}

func (qt *cute) NextTest() NextTestBuilder {
	qt.countTests++ // async?

	qt.reqInfo = nil
	qt.reqInfoT = nil
	qt.respInfo = nil
	qt.respInfoT = nil
	qt.tests = append(qt.tests, createDefaultTest(qt.baseProps, qt.reqInfo, qt.reqInfoT, qt.respInfo, qt.respInfoT))

	return qt
}

func (qt *cute) RequestInformationTable(handlers ...RequestInformation) MiddlewareTable {
	qt.reqInfo = append(qt.reqInfo, handlers...)

	return qt
}

func (qt *cute) RequestInformationTTable(handlers ...RequestInformationT) MiddlewareTable {
	qt.reqInfoT = append(qt.reqInfoT, handlers...)

	return qt
}

func (qt *cute) ResponseInformationTable(handlers ...ResponseInformation) MiddlewareTable {
	qt.respInfo = append(qt.respInfo, handlers...)

	return qt
}

func (qt *cute) ResponseInformationTTable(handlers ...ResponseInformationT) MiddlewareTable {
	qt.respInfoT = append(qt.respInfoT, handlers...)

	return qt
}
