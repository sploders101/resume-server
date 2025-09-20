package main

import "html/template"

type ResumeContext struct {
	IncludeDescriptions bool
	Resume              Resume
}

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
	AccentBg           template.CSS `yaml:"accentBg"`
	AccentText         template.CSS `yaml:"accentText"`
	TextOnAccentColor  template.CSS `yaml:"textOnAccentColor"`
	HighlightColor     template.CSS `yaml:"highlightColor"`
	HighlightTextColor template.CSS `yaml:"highlightTextColor"`
}

type LinkableKV struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
	Href  string `yaml:"href"`
}

type ProfExperience struct {
	CompanyName string        `yaml:"companyName"`
	Roles       []ProfExpRole `yaml:"roles"`
	Highlights  *string       `yaml:"highlights,omitempty"`
	Description *string       `yaml:"description,omitempty"`
}

type ProfExpRole struct {
	Title     string `yaml:"title"`
	StartedAt string `yaml:"startedAt"`
	EndedAt   string `yaml:"endedAt"`
}

type PersonalProject struct {
	Title       string  `yaml:"title"`
	Highlights  *string `yaml:"highlights"`
	Description *string `yaml:"description"`
}

type Certification struct {
	Title       string `yaml:"title"`
	ValidStart  string `yaml:"validStart"`
	ValidEnd    string `yaml:"validEnd"`
	Description string `yaml:"description"`
}
