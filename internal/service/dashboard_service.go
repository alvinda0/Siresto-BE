package service

import (
	"project-name/internal/entity"
	"project-name/internal/repository"
	"time"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetHomeStats(companyID, branchID uuid.UUID) (*entity.DashboardHomeResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *dashboardService) GetHomeStats(companyID, branchID uuid.UUID) (*entity.DashboardHomeResponse, error) {
	now := time.Now()
	
	// Hitung range waktu
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	
	// Total items per tanggal (7 hari terakhir)
	totalItemsByDate, err := s.dashboardRepo.GetTotalItemsByDate(companyID, branchID, 7)
	if err != nil {
		return nil, err
	}
	
	// Revenue per tanggal (7 hari terakhir)
	revenueByDate, err := s.dashboardRepo.GetRevenueByDate(companyID, branchID, 7)
	if err != nil {
		return nil, err
	}
	
	// Best selling daily
	bestSellingDaily, err := s.dashboardRepo.GetBestSellingItems(companyID, branchID, startOfDay, endOfDay, 10)
	if err != nil {
		return nil, err
	}
	
	// Best selling weekly
	bestSellingWeekly, err := s.dashboardRepo.GetBestSellingItems(companyID, branchID, startOfWeek, endOfWeek, 10)
	if err != nil {
		return nil, err
	}
	
	// Best selling monthly
	bestSellingMonthly, err := s.dashboardRepo.GetBestSellingItems(companyID, branchID, startOfMonth, endOfMonth, 10)
	if err != nil {
		return nil, err
	}
	
	// Complimentary items (bulan ini)
	complimentaryItems, err := s.dashboardRepo.GetComplimentaryItems(companyID, branchID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}
	
	return &entity.DashboardHomeResponse{
		TotalItemsByDate:   totalItemsByDate,
		RevenuByDate:       revenueByDate,
		BestSellingDaily:   bestSellingDaily,
		BestSellingWeekly:  bestSellingWeekly,
		BestSellingMonthly: bestSellingMonthly,
		ComplimentaryItems: complimentaryItems,
	}, nil
}
