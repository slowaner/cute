package cute

func (qt *cute) setAllureInformation(t allureProvider) {
	// Log main vars to allureProvider
	qt.setLabelsAllure(t)
	qt.setInfoAllure(t)
	qt.setLinksAllure(t)
}

func (qt *cute) setLinksAllure(t linksAllureProvider) {
	if qt.allureLinks.Issue != "" {
		t.SetIssue(qt.allureLinks.Issue)
	}

	if qt.allureLinks.TestCase != "" {
		t.SetTestCase(qt.allureLinks.TestCase)
	}

	if qt.allureLinks.Link != nil {
		t.Link(qt.allureLinks.Link)
	}

	if qt.allureLinks.TmsLink != "" {
		t.TmsLink(qt.allureLinks.TmsLink)
	}

	if len(qt.allureLinks.TmsLinks) > 0 {
		t.TmsLinks(qt.allureLinks.TmsLinks...)
	}
}

func (qt *cute) setLabelsAllure(t labelsAllureProvider) {
	if qt.allureLabels.ID != "" {
		t.ID(qt.allureLabels.ID)
	}

	if qt.allureLabels.SuiteLabel != "" {
		t.AddSuiteLabel(qt.allureLabels.SuiteLabel)
	}

	if qt.allureLabels.SubSuite != "" {
		t.AddSubSuite(qt.allureLabels.SubSuite)
	}

	if qt.allureLabels.ParentSuite != "" {
		t.AddParentSuite(qt.allureLabels.ParentSuite)
	}

	if qt.allureLabels.Story != "" {
		t.Story(qt.allureLabels.Story)
	}

	if qt.allureLabels.Tag != "" {
		t.Tag(qt.allureLabels.Tag)
	}

	if qt.allureLabels.AllureID != "" {
		t.AllureID(qt.allureLabels.AllureID)
	}

	if qt.allureLabels.Severity != "" {
		t.Severity(qt.allureLabels.Severity)
	}

	if qt.allureLabels.Owner != "" {
		t.Owner(qt.allureLabels.Owner)
	}

	if qt.allureLabels.Lead != "" {
		t.Lead(qt.allureLabels.Lead)
	}

	if qt.allureLabels.Label != nil {
		t.Label(qt.allureLabels.Label)
	}

	if len(qt.allureLabels.Labels) != 0 {
		t.Labels(qt.allureLabels.Labels...)
	}

	if qt.allureLabels.Feature != "" {
		t.Feature(qt.allureLabels.Feature)
	}

	if qt.allureLabels.Epic != "" {
		t.Epic(qt.allureLabels.Epic)
	}

	if len(qt.allureLabels.Tags) != 0 {
		t.Tags(qt.allureLabels.Tags...)
	}

	if qt.allureLabels.Layer != "" {
		t.Layer(qt.allureLabels.Layer)
	}
}

func (qt *cute) setInfoAllure(t infoAllureProvider) {
	if qt.allureInfo.Title != "" {
		t.Title(qt.allureInfo.Title)
	}

	if qt.allureInfo.Description != "" {
		t.Description(qt.allureInfo.Description)
	}

	if qt.allureInfo.Stage != "" {
		t.Stage(qt.allureInfo.Stage)
	}
}

func (it *Test) setAllureInformation(tp allureProvider) {
	// Log main vars to allureProvider
	it.setLabelsAllure(tp)
	it.setInfoAllure(tp)
	it.setLinksAllure(tp)
}

func (it *Test) setInfoAllure(tp infoAllureProvider) {
	if it.AllureInfo.Title != "" {
		tp.Title(it.AllureInfo.Title)
	}

	if it.AllureInfo.Description != "" {
		tp.Description(it.AllureInfo.Description)
	}

	if it.AllureInfo.Stage != "" {
		tp.Stage(it.AllureInfo.Stage)
	}
}

func (it *Test) setLinksAllure(tp linksAllureProvider) {
	if it.AllureLinks.Issue != "" {
		tp.SetIssue(it.AllureLinks.Issue)
	}

	if it.AllureLinks.TestCase != "" {
		tp.SetTestCase(it.AllureLinks.TestCase)
	}

	if it.AllureLinks.Link != nil {
		tp.Link(it.AllureLinks.Link)
	}

	if it.AllureLinks.TmsLink != "" {
		tp.TmsLink(it.AllureLinks.TmsLink)
	}

	if len(it.AllureLinks.TmsLinks) > 0 {
		tp.TmsLinks(it.AllureLinks.TmsLinks...)
	}
}

func (it *Test) setLabelsAllure(tp labelsAllureProvider) {
	if it.AllureLabels.ID != "" {
		tp.ID(it.AllureLabels.ID)
	}

	if it.AllureLabels.SuiteLabel != "" {
		tp.AddSuiteLabel(it.AllureLabels.SuiteLabel)
	}

	if it.AllureLabels.SubSuite != "" {
		tp.AddSubSuite(it.AllureLabels.SubSuite)
	}

	if it.AllureLabels.ParentSuite != "" {
		tp.AddParentSuite(it.AllureLabels.ParentSuite)
	}

	if it.AllureLabels.Story != "" {
		tp.Story(it.AllureLabels.Story)
	}

	if it.AllureLabels.Tag != "" {
		tp.Tag(it.AllureLabels.Tag)
	}

	if it.AllureLabels.AllureID != "" {
		tp.AllureID(it.AllureLabels.AllureID)
	}

	if it.AllureLabels.Severity != "" {
		tp.Severity(it.AllureLabels.Severity)
	}

	if it.AllureLabels.Owner != "" {
		tp.Owner(it.AllureLabels.Owner)
	}

	if it.AllureLabels.Lead != "" {
		tp.Lead(it.AllureLabels.Lead)
	}

	if it.AllureLabels.Label != nil {
		tp.Label(it.AllureLabels.Label)
	}

	if len(it.AllureLabels.Labels) != 0 {
		tp.Labels(it.AllureLabels.Labels...)
	}

	if it.AllureLabels.Feature != "" {
		tp.Feature(it.AllureLabels.Feature)
	}

	if it.AllureLabels.Epic != "" {
		tp.Epic(it.AllureLabels.Epic)
	}

	if len(it.AllureLabels.Tags) != 0 {
		tp.Tags(it.AllureLabels.Tags...)
	}

	if it.AllureLabels.Layer != "" {
		tp.Layer(it.AllureLabels.Layer)
	}
}
