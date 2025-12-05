package mkv

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type LocaleInfo struct {
	language.Tag
}

func (li *LocaleInfo) UnmarshalJSON(b []byte) error {
	call := strings.Trim(string(b), `"`)
	t, err := language.Parse(call)
	if err != nil {
		return err
	}
	li.Tag = t
	return nil
}

func FromIETFName(ietf string) (LocaleInfo, error) {
	lTag, err := language.Parse(ietf)
	if err != nil {
		return LocaleInfo{}, err
	}
	return LocaleInfo{Tag: lTag}, nil
}

func (li *LocaleInfo) LangLocalizedName() string {
	base, _ := li.Base()
	return display.Self.Name(base)
}

func (li *LocaleInfo) LangEnglishName() string {
	base, _ := li.Base()
	return display.English.Languages().Name(base)
}

func (li *LocaleInfo) RegionLocalizedName() string {
	region, _ := li.Region()
	return display.Self.Name(region)
}

func (li *LocaleInfo) RegionEnglishName() string {
	region, _ := li.Region()
	return display.English.Regions().Name(region)
}

func (li *LocaleInfo) RegionCode() string {
	region, _ := li.Region()
	return region.String()
}
