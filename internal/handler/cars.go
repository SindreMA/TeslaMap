package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sindrema/teslamap/internal/db"
	"github.com/sindrema/teslamap/internal/model"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) ListCars(w http.ResponseWriter, r *http.Request) {
	cars, err := db.ListCars(r.Context(), h.DB)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if cars == nil {
		cars = []model.Car{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *Handler) GetCarPosition(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	car, err := db.GetCar(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if car == nil {
		http.Error(w, "car not found", http.StatusNotFound)
		return
	}

	pos, err := db.GetLatestPosition(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	resp := model.CarPosition{
		Car:      *car,
		Position: pos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
