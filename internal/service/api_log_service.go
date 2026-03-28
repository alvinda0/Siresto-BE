package service

import (
	"project-name/internal/entity"
	"project-name/internal/repository"
	"project-name/pkg"
	"errors"
)

type APILogService interface {
	CreateLog(log *entity.APILog) error
	GetAllLogs(page, limit int, method, companyID, branchID string) ([]entity.APILogListDTO, *pkg.PaginationMeta, error)
	GetLogByID(id string, companyID, branchID string) (*entity.APILogDetailDTO, error)
}

type apiLogService struct {
	repo repository.APILogRepository
}

func NewAPILogService(repo repository.APILogRepository) APILogService {
	return &apiLogService{repo: repo}
}

func (s *apiLogService) CreateLog(log *entity.APILog) error {
	return s.repo.Create(log)
}

func (s *apiLogService) GetAllLogs(page, limit int, method, companyID, branchID string) ([]entity.APILogListDTO, *pkg.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	logs, total, err := s.repo.FindAll(page, limit, method, companyID, branchID)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTO (without response_body)
	logDTOs := make([]entity.APILogListDTO, len(logs))
	for i, log := range logs {
		logDTOs[i] = log.ToListDTO()
	}

	// Calculate pagination meta
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &pkg.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return logDTOs, meta, nil
}

func (s *apiLogService) GetLogByID(id string, companyID, branchID string) (*entity.APILogDetailDTO, error) {
	log, err := s.repo.FindByID(id, companyID, branchID)
	if err != nil {
		return nil, errors.New("log not found")
	}
	
	// Convert to detail DTO (with response_body)
	detailDTO := log.ToDetailDTO()
	return &detailDTO, nil
}
