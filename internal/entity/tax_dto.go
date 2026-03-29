package entity

import "github.com/google/uuid"

type CreateTaxRequest struct {
	BranchID   *uuid.UUID `json:"branch_id"` // optional, null = company level
	NamaPajak  string     `json:"nama_pajak" binding:"required"`
	TipePajak  string     `json:"tipe_pajak" binding:"required,oneof=sc pb1"`
	Presentase float64    `json:"presentase" binding:"required,min=0,max=100"`
	Deskripsi  string     `json:"deskripsi"`
	Status     string     `json:"status" binding:"omitempty,oneof=active inactive"`
	Prioritas  int        `json:"prioritas"`
}

type UpdateTaxRequest struct {
	NamaPajak  string  `json:"nama_pajak"`
	TipePajak  string  `json:"tipe_pajak" binding:"omitempty,oneof=sc pb1"`
	Presentase float64 `json:"presentase" binding:"omitempty,min=0,max=100"`
	Deskripsi  string  `json:"deskripsi"`
	Status     string  `json:"status" binding:"omitempty,oneof=active inactive"`
	Prioritas  int     `json:"prioritas"`
}

type TaxResponse struct {
	ID          uuid.UUID  `json:"id"`
	CompanyID   uuid.UUID  `json:"company_id"`
	CompanyName string     `json:"company_name"`
	BranchID    *uuid.UUID `json:"branch_id"`
	BranchName  *string    `json:"branch_name"`
	NamaPajak   string     `json:"nama_pajak"`
	TipePajak   string     `json:"tipe_pajak"`
	Presentase  float64    `json:"presentase"`
	Deskripsi   string     `json:"deskripsi"`
	Status      string     `json:"status"`
	Prioritas   int        `json:"prioritas"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}
