package trip

import (
	"bytes"
	"context"

	"github.com/jung-kurt/gofpdf"
	"github.com/korroziea/taxi/user-service/internal/domain"
	triphandler "github.com/korroziea/taxi/user-service/internal/handler/trip"
)

type Repo interface {
	CheckWalletBalance(ctx context.Context, cost int64) error
}

type HTTPAdapter interface {
	Trips(ctx context.Context, userID string) ([]domain.Trip, error)
}

type Adapter interface {
	StartTrip(ctx context.Context, trip domain.StartTrip) error
	CancelTrip(ctx context.Context) error
}

type Service struct {
	// repo    Repo
	httpAdapter HTTPAdapter
	adapter     Adapter
}

func New(adapter Adapter, httpAdapter HTTPAdapter) *Service {
	service := &Service{
		// repo:    repo,
		httpAdapter: httpAdapter,
		adapter:     adapter,
	}

	return service
}

func (s *Service) StartTrip(ctx context.Context, trip domain.StartTrip) error {
	// if err := s.repo.CheckWalletBalance(ctx, 0); err != nil { // todo: add cost
	// 	return fmt.Errorf("repo.CheckWalletBalance: %w", err)
	// }

	return s.adapter.StartTrip(ctx, trip)
}

func (s *Service) CancelTrip(ctx context.Context) error {
	return s.adapter.CancelTrip(ctx)
}

func (s *Service) Trips(ctx context.Context, userID string) ([]domain.Trip, error) {
	return s.httpAdapter.Trips(ctx, userID)
}

func (s *Service) Cost(ctx context.Context) (int64, error) {
	return 0, nil
}

func (s *Service) Report(ctx context.Context) ([]byte, error) {
	trips, err := s.httpAdapter.Trips(ctx, triphandler.FromContext(ctx))
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Заголовок
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Car's Report")
	pdf.Ln(12)

	// Таблица с машинами
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(80, 10, "Trip ID", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 10, "Car Number", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 10, "Car Color", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 10, "Car ID", "1", 0, "", false, 0, "")
	pdf.CellFormat(60, 10, "Created At", "1", 0, "", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, trip := range trips {
		pdf.CellFormat(80, 10, trip.ID, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, trip.CarNumber, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 10, trip.CarColor, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, trip.CarID, "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 10, trip.CreatedAt.String(), "1", 0, "", false, 0, "")
		pdf.Ln(10)
	}

	// Сохраняем PDF в буфер
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
