package service

import (
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
)

type SearchService struct {
	repo postgres.Search
}

func NewSearchService(repo postgres.Search) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) SearchJob(searchQuery string) ([]models.SearchResult, error) {
	return s.repo.SearchJob(searchQuery)
}

func (s *SearchService) SearchOrg(searchQuery string) ([]models.SearchResult, error) {
	return s.repo.SearchOrg(searchQuery)
}
