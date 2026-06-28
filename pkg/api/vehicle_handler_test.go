package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/santri00/golang-api-example/pkg/models"
	vehicle "github.com/santri00/golang-api-example/pkg/services/vehicle"
)

// mockVehicleService is a test double for vehicle.Service.
type mockVehicleService struct {
	addVehicleFn          func(ctx context.Context, v *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error)
	listVehiclesFn        func(ctx context.Context, offset, limit int) (*models.ListVehicleResponse, error)
	getVehicleByStockIdFn func(ctx context.Context, stockId string) (*models.VehicleDetailedRecord, error)
}

var _ vehicle.Service = (*mockVehicleService)(nil)

func (m *mockVehicleService) AddVehicle(ctx context.Context, v *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error) {
	return m.addVehicleFn(ctx, v)
}

func (m *mockVehicleService) ListVehicles(ctx context.Context, offset, limit int) (*models.ListVehicleResponse, error) {
	return m.listVehiclesFn(ctx, offset, limit)
}

func (m *mockVehicleService) GetVehicleByStockId(ctx context.Context, stockId string) (*models.VehicleDetailedRecord, error) {
	return m.getVehicleByStockIdFn(ctx, stockId)
}

func newTestServer(svc vehicle.Service) *Server {
	return &Server{ctx: context.Background(), vehicleService: svc}
}

// routedRequest routes GET requests through mux so {stock_id} vars are populated.
func routedRequest(handler http.HandlerFunc, pattern, url string) (*httptest.ResponseRecorder, *http.Request) {
	r := mux.NewRouter()
	r.HandleFunc(pattern, handler).Methods(http.MethodGet)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr, req
}

func TestAddVehicles(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svc        vehicle.Service
		wantStatus int
	}{
		{
			name: "valid request returns 201",
			body: `{"title":"2022 Toyota Camry","price":"24995","miles":10000,"vin":"1HGBH41JXMN109186","brand":"Toyota","model":"Camry","year":2022}`,
			svc: &mockVehicleService{
				addVehicleFn: func(_ context.Context, _ *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error) {
					return &models.AddVehicleResponse{StockId: "test-uuid"}, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid JSON returns 400",
			body:       `not-valid-json`,
			svc:        &mockVehicleService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body returns 400",
			body:       ``,
			svc:        &mockVehicleService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error returns 500",
			body: `{}`,
			svc: &mockVehicleService{
				addVehicleFn: func(_ context.Context, _ *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error) {
					return nil, errors.New("storage failure")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/inventory/v1/vehicles", bytes.NewBufferString(tc.body))
			rr := httptest.NewRecorder()
			newTestServer(tc.svc).addVehicles(rr, req)
			if rr.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d; body: %s", rr.Code, tc.wantStatus, rr.Body.String())
			}
		})
	}
}

func TestListVehicles(t *testing.T) {
	successSvc := &mockVehicleService{
		listVehiclesFn: func(_ context.Context, offset, limit int) (*models.ListVehicleResponse, error) {
			return &models.ListVehicleResponse{
				Metadata: models.Metadata{Offset: offset, Limit: limit, TotalCount: 1},
				Records:  &models.Vehicles{{StockID: "abc", Title: "Car"}},
			}, nil
		},
	}

	tests := []struct {
		name       string
		query      string
		svc        vehicle.Service
		wantStatus int
	}{
		{
			name:       "no params uses defaults and returns 200",
			query:      "",
			svc:        successSvc,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid offset and limit returns 200",
			query:      "?offset=0&limit=10",
			svc:        successSvc,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid offset returns 400",
			query:      "?offset=abc",
			svc:        &mockVehicleService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid limit returns 400",
			query:      "?limit=xyz",
			svc:        &mockVehicleService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "both invalid returns 400 with multiple errors",
			query:      "?offset=abc&limit=xyz",
			svc:        &mockVehicleService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service error returns 500",
			query: "",
			svc: &mockVehicleService{
				listVehiclesFn: func(_ context.Context, _, _ int) (*models.ListVehicleResponse, error) {
					return nil, errors.New("redis down")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/inventory/v1/vehicles"+tc.query, nil)
			rr := httptest.NewRecorder()
			newTestServer(tc.svc).listVehicles(rr, req)
			if rr.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d; body: %s", rr.Code, tc.wantStatus, rr.Body.String())
			}
		})
	}
}

func TestListVehicles_ResponseBody(t *testing.T) {
	svc := &mockVehicleService{
		listVehiclesFn: func(_ context.Context, offset, limit int) (*models.ListVehicleResponse, error) {
			return &models.ListVehicleResponse{
				Metadata: models.Metadata{Offset: offset, Limit: limit, TotalCount: 2},
				Records:  &models.Vehicles{{StockID: "s1"}, {StockID: "s2"}},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/inventory/v1/vehicles?offset=0&limit=5", nil)
	rr := httptest.NewRecorder()
	newTestServer(svc).listVehicles(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rr.Code)
	}

	var resp models.ListVehicleResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Metadata.TotalCount != 2 {
		t.Errorf("got total_count %d, want 2", resp.Metadata.TotalCount)
	}
}

func TestGetVehicleByStockId(t *testing.T) {
	tests := []struct {
		name       string
		stockId    string
		svc        vehicle.Service
		wantStatus int
	}{
		{
			name:    "vehicle found returns 200",
			stockId: "abc-123",
			svc: &mockVehicleService{
				getVehicleByStockIdFn: func(_ context.Context, _ string) (*models.VehicleDetailedRecord, error) {
					return &models.VehicleDetailedRecord{}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "vehicle not found returns 404",
			stockId: "missing-id",
			svc: &mockVehicleService{
				getVehicleByStockIdFn: func(_ context.Context, _ string) (*models.VehicleDetailedRecord, error) {
					return nil, nil
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:    "service error returns 500",
			stockId: "abc-123",
			svc: &mockVehicleService{
				getVehicleByStockIdFn: func(_ context.Context, _ string) (*models.VehicleDetailedRecord, error) {
					return nil, errors.New("redis error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := newTestServer(tc.svc)
			rr, _ := routedRequest(
				srv.getVehicleByStockId,
				"/inventory/v1/vehicles/{stock_id}",
				"/inventory/v1/vehicles/"+tc.stockId,
			)
			if rr.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d; body: %s", rr.Code, tc.wantStatus, rr.Body.String())
			}
		})
	}
}

func TestGetVehicleByStockId_MissingStockId(t *testing.T) {
	// Call handler directly (no mux) so mux.Vars returns empty map → empty stockId.
	req := httptest.NewRequest(http.MethodGet, "/inventory/v1/vehicles/", nil)
	rr := httptest.NewRecorder()
	newTestServer(&mockVehicleService{}).getVehicleByStockId(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want 400", rr.Code)
	}
}
