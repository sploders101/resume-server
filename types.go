package main

import "html/template"

type Resume struct {
	Theme            ResumeTheme       `yaml:"theme"`
	Name             string            `yaml:"name"`
	Role             string            `yaml:"role"`
	ContactInfo      []LinkableKV      `yaml:"contactInfo"`
	Links            []LinkableKV      `yaml:"links"`
	Objective        string            `yaml:"objective"`
	Skills           []string          `yaml:"skills"`
	ProfExperience   []ProfExperience  `yaml:"profExperience"`
	PersonalProjects []PersonalProject `yaml:"personalProjects"`
	Certifications   []Certification   `yaml:"certifications"`
}

type ResumeTheme struct {
	AccentBg          template.CSS `yaml:"accentBg"`
	AccentText        template.CSS `yaml:"accentText"`
	TextOnAccentColor template.CSS `yaml:"textOnAccentColor"`
}

type LinkableKV struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
	Href  string `yaml:"href"`
}

type ProfExperience struct {
	CompanyName string        `yaml:"companyName"`
	Roles       []ProfExpRole `yaml:"roles"`
	Description string        `yaml:"description"`
}

type ProfExpRole struct {
	Title     string `yaml:"title"`
	StartedAt string `yaml:"startedAt"`
	EndedAt   string `yaml:"endedAt"`
}

type PersonalProject struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Certification struct {
	Title       string `yaml:"title"`
	ValidStart  string `yaml:"validStart"`
	ValidEnd    string `yaml:"validEnd"`
	Description string `yaml:"description"`
}
