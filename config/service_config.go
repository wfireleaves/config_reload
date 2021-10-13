package config

import "fmt"

type Service struct {
	Version string `xml:"version"`
	Name    string `xml:"name"`
}

func (s *Service) String() string {
	return fmt.Sprintf("version:%s name:%s\n", s.Version, s.Name)
}
