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
