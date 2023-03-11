package service

import (
	"brutalITSM-BE-Users/pkg/repository"
)

type SearchService struct {
	repo repository.Search
}

func NewSearchService(repo repository.Search) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) SearchJob(searchQuery string) ([]string, error) {
	return s.repo.SearchJob(searchQuery)
}

func (s *SearchService) SearchOrg(searchQuery string) ([]string, error) {
	return s.repo.SearchOrg(searchQuery)
}
