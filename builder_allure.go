package cute

import (
	"fmt"

	"github.com/ozontech/allure-go/pkg/allure"
)

func (qt *cute) Parallel() AllureBuilder {
	qt.parallel = true

	return qt
}

func (qt *cute) Title(title string) AllureBuilder {
	qt.allureInfo.Title = title

	return qt
}

func (qt *cute) Epic(epic string) AllureBuilder {
	qt.allureLabels.Epic = epic

	return qt
}

func (qt *cute) Titlef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.Title = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Descriptionf(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.Description = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Stage(stage string) AllureBuilder {
	qt.allureInfo.Stage = stage

	return qt
}

func (qt *cute) Stagef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.Stage = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Layer(value string) AllureBuilder {
	qt.allureLabels.Layer = value

	return qt
}

func (qt *cute) TmsLink(tmsLink string) AllureBuilder {
	qt.allureLinks.TmsLink = tmsLink

	return qt
}

func (qt *cute) TmsLinks(tmsLinks ...string) AllureBuilder {
	qt.allureLinks.TmsLinks = append(qt.allureLinks.TmsLinks, tmsLinks...)

	return qt
}

func (qt *cute) SetIssue(issue string) AllureBuilder {
	qt.allureLinks.Issue = issue

	return qt
}

func (qt *cute) SetTestCase(testCase string) AllureBuilder {
	qt.allureLinks.TestCase = testCase

	return qt
}

func (qt *cute) Link(link *allure.Link) AllureBuilder {
	qt.allureLinks.Link = link

	return qt
}

func (qt *cute) ID(value string) AllureBuilder {
	qt.allureLabels.ID = value

	return qt
}

func (qt *cute) AllureID(value string) AllureBuilder {
	qt.allureLabels.AllureID = value

	return qt
}

func (qt *cute) AddSuiteLabel(value string) AllureBuilder {
	qt.allureLabels.SuiteLabel = value

	return qt
}

func (qt *cute) AddSubSuite(value string) AllureBuilder {
	qt.allureLabels.SubSuite = value

	return qt
}

func (qt *cute) AddParentSuite(value string) AllureBuilder {
	qt.allureLabels.ParentSuite = value

	return qt
}

func (qt *cute) Story(value string) AllureBuilder {
	qt.allureLabels.Story = value

	return qt
}

func (qt *cute) Tag(value string) AllureBuilder {
	qt.allureLabels.Tag = value

	return qt
}

func (qt *cute) Severity(value allure.SeverityType) AllureBuilder {
	qt.allureLabels.Severity = value

	return qt
}

func (qt *cute) Owner(value string) AllureBuilder {
	qt.allureLabels.Owner = value

	return qt
}

func (qt *cute) Lead(value string) AllureBuilder {
	qt.allureLabels.Lead = value

	return qt
}

func (qt *cute) Label(label *allure.Label) AllureBuilder {
	qt.allureLabels.Label = label

	return qt
}

func (qt *cute) Labels(labels ...*allure.Label) AllureBuilder {
	qt.allureLabels.Labels = labels

	return qt
}

func (qt *cute) Description(description string) AllureBuilder {
	qt.allureInfo.Description = description

	return qt
}

func (qt *cute) Tags(tags ...string) AllureBuilder {
	qt.allureLabels.Tags = tags

	return qt
}

func (qt *cute) Feature(feature string) AllureBuilder {
	qt.allureLabels.Feature = feature

	return qt
}
