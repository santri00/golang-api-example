package vehicle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/santri00/golang-api-example/pkg/models"
	"github.com/santri00/golang-api-example/pkg/storage/entities"
	storagevehicle "github.com/santri00/golang-api-example/pkg/storage/vehicle"
)

// mockRepository is a test double for storagevehicle.Repository.
type mockRepository struct {
	createFn func(ctx context.Context, v *entities.Vehicle) (string, error)
	getAllFn  func(ctx context.Context, offset, limit int) (*[]entities.VehicleBasicRecord, int, error)
	getFn    func(ctx context.Context, stockId string) (*entities.Vehicle, error)
}

var _ storagevehicle.Repository = (*mockRepository)(nil)

func (m *mockRepository) Create(ctx context.Context, v *entities.Vehicle) (string, error) {
	return m.createFn(ctx, v)
}

func (m *mockRepository) GetAll(ctx context.Context, offset, limit int) (*[]entities.VehicleBasicRecord, int, error) {
	return m.getAllFn(ctx, offset, limit)
}

func (m *mockRepository) Get(ctx context.Context, stockId string) (*entities.Vehicle, error) {
	return m.getFn(ctx, stockId)
}

func sampleVehicle() *models.VehicleDetailedRecord {
	return &models.VehicleDetailedRecord{
		Vehicle: models.Vehicle{
			Title:              "2022 Toyota Camry",
			Price:              "24995",
			Miles:              10000,
			CoverImageLink:     "https://example.com/camry.jpg",
			InventoryAddedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Vin:   "1HGBH41JXMN109186",
		Brand: "Toyota",
		Model: "Camry",
		Year:  2022,
		Type:  "Sedan",
	}
}

// --- AddVehicle ---

func TestAddVehicle_NilVehicle(t *testing.T) {
	svc := NewService(&mockRepository{})
	_, err := svc.AddVehicle(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil vehicle, got nil")
	}
}

func TestAddVehicle_RepoError(t *testing.T) {
	repo := &mockRepository{
		createFn: func(_ context.Context, _ *entities.Vehicle) (string, error) {
			return "", errors.New("redis failure")
		},
	}
	svc := NewService(repo)
	_, err := svc.AddVehicle(context.Background(), sampleVehicle())
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

func TestAddVehicle_Success(t *testing.T) {
	wantStockId := "test-stock-uuid"
	repo := &mockRepository{
		createFn: func(_ context.Context, _ *entities.Vehicle) (string, error) {
			return wantStockId, nil
		},
	}
	svc := NewService(repo)
	resp, err := svc.AddVehicle(context.Background(), sampleVehicle())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.StockId != wantStockId {
		t.Errorf("got StockId %q, want %q", resp.StockId, wantStockId)
	}
}

// --- ListVehicles ---

func TestListVehicles_RepoError(t *testing.T) {
	repo := &mockRepository{
		getAllFn: func(_ context.Context, _, _ int) (*[]entities.VehicleBasicRecord, int, error) {
			return nil, 0, errors.New("redis failure")
		},
	}
	svc := NewService(repo)
	_, err := svc.ListVehicles(context.Background(), 0, -1)
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

func TestListVehicles_EmptyRecords(t *testing.T) {
	empty := &[]entities.VehicleBasicRecord{}
	repo := &mockRepository{
		getAllFn: func(_ context.Context, _, _ int) (*[]entities.VehicleBasicRecord, int, error) {
			return empty, 0, nil
		},
	}
	svc := NewService(repo)
	resp, err := svc.ListVehicles(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Records != nil {
		t.Errorf("expected nil Records for empty result, got %v", resp.Records)
	}
}

func TestListVehicles_WithRecords(t *testing.T) {
	records := &[]entities.VehicleBasicRecord{
		{StockID: "s1", Title: "Car A"},
		{StockID: "s2", Title: "Car B"},
	}
	repo := &mockRepository{
		getAllFn: func(_ context.Context, offset, limit int) (*[]entities.VehicleBasicRecord, int, error) {
			return records, 2, nil
		},
	}
	svc := NewService(repo)
	resp, err := svc.ListVehicles(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Metadata.TotalCount != 2 {
		t.Errorf("got TotalCount %d, want 2", resp.Metadata.TotalCount)
	}
	if resp.Records == nil || len(*resp.Records) != 2 {
		t.Errorf("got %v records, want 2", resp.Records)
	}
	if resp.Metadata.Offset != 0 || resp.Metadata.Limit != 10 {
		t.Errorf("unexpected metadata: %+v", resp.Metadata)
	}
}

// --- GetVehicleByStockId ---

func TestGetVehicleByStockId_RepoError(t *testing.T) {
	repo := &mockRepository{
		getFn: func(_ context.Context, _ string) (*entities.Vehicle, error) {
			return nil, errors.New("redis failure")
		},
	}
	svc := NewService(repo)
	_, err := svc.GetVehicleByStockId(context.Background(), "abc-123")
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

func TestGetVehicleByStockId_NotFound(t *testing.T) {
	repo := &mockRepository{
		getFn: func(_ context.Context, _ string) (*entities.Vehicle, error) {
			return nil, nil
		},
	}
	svc := NewService(repo)
	result, err := svc.GetVehicleByStockId(context.Background(), "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for missing record, got %+v", result)
	}
}

func TestGetVehicleByStockId_Success(t *testing.T) {
	entity := &entities.Vehicle{
		BasicRecord: &entities.VehicleBasicRecord{
			StockID: "abc-123",
			Title:   "2022 Toyota Camry",
			Miles:   10000,
		},
		DetailedRecord: &entities.VehicleDetailedRecord{
			Vin:   "1HGBH41JXMN109186",
			Brand: "Toyota",
			Model: "Camry",
			Year:  2022,
		},
	}
	repo := &mockRepository{
		getFn: func(_ context.Context, stockId string) (*entities.Vehicle, error) {
			return entity, nil
		},
	}
	svc := NewService(repo)
	result, err := svc.GetVehicleByStockId(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected vehicle record, got nil")
	}
	if result.Vehicle.StockID != "abc-123" {
		t.Errorf("got StockID %q, want %q", result.Vehicle.StockID, "abc-123")
	}
}
