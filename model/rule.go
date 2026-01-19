package model

type Rule struct {
  Id          string
  Severity    string
  Tags        []Tag
  Description string
  Name        string
}
