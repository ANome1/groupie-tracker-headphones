package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sync"
)

// In-memory store for categories (thread-safe)
var (
	categories      = []string{"Pays", "Ville", "Animal", "Métier", "Objet", "Prénom"}
	categoriesMutex sync.RWMutex
)

// PetitBacHandler renders the game page
func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	// Parse templates (adjust paths as necessary for your project structure)
	tmpl, err := template.ParseFiles(
		"templates/games/petitbac.html",
		"templates/components/header.html",
		"templates/components/footer.html",
		"templates/components/scoreboard.html",
	)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass initial data if needed
	data := map[string]interface{}{
		"Categories": categories,
	}

	tmpl.Execute(w, data)
}

// GetCategoriesHandler returns the list of categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categoriesMutex.RLock()
	defer categoriesMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// AddCategoryHandler adds a new category
func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Category string
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Category == "" {
		http.Error(w, "Category cannot be empty", http.StatusBadRequest)
		return
	}

	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()

	// Check for duplicates
	for _, c := range categories {
		if c == req.Category {
			http.Error(w, "Category already exists", http.StatusConflict)
			return
		}
	}

	categories = append(categories, req.Category)
	w.WriteHeader(http.StatusOK)
}

// DeleteCategoryHandler removes a category
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Category string
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()

	newCategories := make([]string, 0, len(categories))
	found := false
	for _, c := range categories {
		if c == req.Category {
			found = true
			continue
		}
		newCategories = append(newCategories, c)
	}

	if !found {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	categories = newCategories
	w.WriteHeader(http.StatusOK)
}
