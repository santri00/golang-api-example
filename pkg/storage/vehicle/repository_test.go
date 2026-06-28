package vehicle

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/santri00/golang-api-example/pkg/storage/cachekeys"
	"github.com/santri00/golang-api-example/pkg/storage/entities"
)

func newTestRepo(t *testing.T) (Repository, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return New(client), client
}

func sampleVehicleEntity() *entities.Vehicle {
	return &entities.Vehicle{
		BasicRecord: &entities.VehicleBasicRecord{
			Title:              "2022 Toyota Camry",
			Price:              "24995",
			Miles:              10000,
			CoverImageLink:     "https://example.com/camry.jpg",
			InventoryAddedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		DetailedRecord: &entities.VehicleDetailedRecord{
			Vin:   "1HGBH41JXMN109186",
			Brand: "Toyota",
			Model: "Camry",
			Year:  2022,
			Type:  "Sedan",
		},
	}
}

// --- Create ---

func TestCreate_NilVehicleRecord(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.Create(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil vehicle record, got nil")
	}
}

func TestCreate_NilBasicRecord(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.Create(context.Background(), &entities.Vehicle{
		BasicRecord:    nil,
		DetailedRecord: &entities.VehicleDetailedRecord{},
	})
	if err == nil {
		t.Fatal("expected error for nil basic record, got nil")
	}
}

func TestCreate_NilDetailedRecord(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.Create(context.Background(), &entities.Vehicle{
		BasicRecord:    &entities.VehicleBasicRecord{},
		DetailedRecord: nil,
	})
	if err == nil {
		t.Fatal("expected error for nil detailed record, got nil")
	}
}

func TestCreate_Success(t *testing.T) {
	repo, client := newTestRepo(t)
	ctx := context.Background()
	vehicle := sampleVehicleEntity()

	stockId, err := repo.Create(ctx, vehicle)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stockId == "" {
		t.Fatal("expected non-empty stock ID, got empty string")
	}

	// Verify the vehicle hash was written to Redis.
	hashKey := "vehicle:" + stockId
	hashLen, _ := client.HLen(ctx, hashKey).Result()
	if hashLen == 0 {
		t.Errorf("expected hash key %q in Redis, not found", hashKey)
	}

	// Verify the basic record was appended to the vehicle list.
	listLen, _ := client.LLen(ctx, cachekeys.VehicleRecordListCacheKey).Result()
	if listLen != 1 {
		t.Errorf("expected list length 1, got %d", listLen)
	}

	// Verify stock ID was stamped onto the basic record.
	if vehicle.BasicRecord.StockID != stockId {
		t.Errorf("basic record StockID not updated: got %q, want %q", vehicle.BasicRecord.StockID, stockId)
	}
}

// --- GetAll ---

func TestGetAll_ZeroLimit(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, _, err := repo.GetAll(context.Background(), 0, 0)
	if err == nil {
		t.Fatal("expected error for limit=0, got nil")
	}
}

func TestGetAll_EmptyStore(t *testing.T) {
	repo, _ := newTestRepo(t)
	records, total, err := repo.GetAll(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if records != nil {
		t.Errorf("expected nil records for empty store, got %+v", records)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestGetAll_ReturnsRecordsAndCount(t *testing.T) {
	repo, client := newTestRepo(t)
	ctx := context.Background()

	// Seed two records via the client.
	for i, title := range []string{"Car A", "Car B"} {
		br := entities.VehicleBasicRecord{StockID: "stock-" + string(rune('1'+i)), Title: title}
		b, _ := json.Marshal(br)
		client.RPush(ctx, cachekeys.VehicleRecordListCacheKey, string(b))
	}

	records, total, err := repo.GetAll(ctx, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("got total %d, want 2", total)
	}
	if records == nil || len(*records) != 2 {
		t.Errorf("got %v records, want 2", records)
	}
}

func TestGetAll_PaginationOffsetAndLimit(t *testing.T) {
	repo, client := newTestRepo(t)
	ctx := context.Background()

	for i, title := range []string{"Car A", "Car B", "Car C"} {
		br := entities.VehicleBasicRecord{StockID: "stock-" + string(rune('1'+i)), Title: title}
		b, _ := json.Marshal(br)
		client.RPush(ctx, cachekeys.VehicleRecordListCacheKey, string(b))
	}

	// LRANGE 0 1 returns indices 0 and 1 (2 records).
	records, total, err := repo.GetAll(ctx, 0, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("got total %d, want 3", total)
	}
	if records == nil || len(*records) != 2 {
		t.Errorf("got %d records for limit=1, want 2 (LRANGE stop is inclusive)", len(*records))
	}
}

// --- Get ---

func TestGet_EmptyStockId(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty stock ID, got nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	result, err := repo.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for missing record, got %+v", result)
	}
}

func TestGet_Success(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	vehicle := sampleVehicleEntity()
	stockId, err := repo.Create(ctx, vehicle)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	result, err := repo.Get(ctx, stockId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected vehicle, got nil")
	}
	if result.BasicRecord == nil {
		t.Fatal("expected basic record, got nil")
	}
	if result.BasicRecord.StockID != stockId {
		t.Errorf("got StockID %q, want %q", result.BasicRecord.StockID, stockId)
	}
	if result.DetailedRecord == nil {
		t.Fatal("expected detailed record, got nil")
	}
	if result.DetailedRecord.Brand != "Toyota" {
		t.Errorf("got Brand %q, want %q", result.DetailedRecord.Brand, "Toyota")
	}
}

func TestGet_RoundTrip(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	original := sampleVehicleEntity()
	stockId, err := repo.Create(ctx, original)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	result, err := repo.Get(ctx, stockId)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if result.BasicRecord.Title != original.BasicRecord.Title {
		t.Errorf("Title mismatch: got %q, want %q", result.BasicRecord.Title, original.BasicRecord.Title)
	}
	if result.DetailedRecord.Vin != original.DetailedRecord.Vin {
		t.Errorf("Vin mismatch: got %q, want %q", result.DetailedRecord.Vin, original.DetailedRecord.Vin)
	}
	if result.DetailedRecord.Year != original.DetailedRecord.Year {
		t.Errorf("Year mismatch: got %d, want %d", result.DetailedRecord.Year, original.DetailedRecord.Year)
	}
}
