package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID                int     `json:"id"`
	Email             string  `json:"email"`
	Password          string  `json:"password,omitempty"`
	Name              string  `json:"name"`
	Phone             string  `json:"phone"`
	UserType          string  `json:"user_type"`
	Address           string  `json:"address"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	DonationDays      string  `json:"donation_days"`
	ResetToken        string  `json:"reset_token,omitempty"`
	ResetTokenExpires string  `json:"reset_token_expires,omitempty"`
}

type Donation struct {
	ID                   int     `json:"id"`
	DonorID              int     `json:"donor_id"`
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	Category             string  `json:"category"`
	Quantity             int     `json:"quantity"`
	ExpiryDate           string  `json:"expiry_date"`
	PickupAddress        string  `json:"pickup_address"`
	PickupLatitude       float64 `json:"pickup_latitude"`
	PickupLongitude      float64 `json:"pickup_longitude"`
	Status               string  `json:"status"`
	ReservedBy           int     `json:"reserved_by"`
	ReservedAt           string  `json:"reserved_at"`
	PickupTime           string  `json:"pickup_time"`
	PickupPersonName     string  `json:"pickup_person_name"`
	PickupPersonID       string  `json:"pickup_person_id"`
	VerificationCode     string  `json:"verification_code"`
	CompletedAt          string  `json:"completed_at"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	DonorConfirmed       bool    `json:"donor_confirmed"`
	RecipientConfirmed   bool    `json:"recipient_confirmed"`
	BusinessConfirmed    bool    `json:"business_confirmed"`
	BusinessConfirmedAt  string  `json:"business_confirmed_at"`
	DonorConfirmedAt     string  `json:"donor_confirmed_at"`
	RecipientConfirmedAt string  `json:"recipient_confirmed_at"`
	Weight               float64 `json:"weight"`
	DonationReason       string  `json:"donation_reason"`
	ContactInfo          string  `json:"contact_info"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(mysql:3306)/food_donation_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/email/{email}", getUserByEmail).Methods("GET")
	r.HandleFunc("/users/name/{name}", getUserByName).Methods("GET")
	r.HandleFunc("/users/type/{type}", getUsersByType).Methods("GET")
	r.HandleFunc("/users/phone/{phone}", getUsersByPhone).Methods("GET")
	r.HandleFunc("/users/address/{address}", getUsersByAddress).Methods("GET")
	r.HandleFunc("/users/days/{days}", getUsersByDonationDays).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/donations", getDonations).Methods("GET")
	r.HandleFunc("/donations/{id}", getDonation).Methods("GET")
	r.HandleFunc("/donations/donor/{donor_id}", getDonationsByDonor).Methods("GET")
	r.HandleFunc("/donations/category/{category}", getDonationsByCategory).Methods("GET")
	r.HandleFunc("/donations/status/{status}", getDonationsByStatus).Methods("GET")
	r.HandleFunc("/donations/reserved/{user_id}", getDonationsReservedBy).Methods("GET")
	r.HandleFunc("/donations/title/{title}", getDonationsByTitle).Methods("GET")
	r.HandleFunc("/donations/address/{address}", getDonationsByAddress).Methods("GET")
	r.HandleFunc("/donations/confirmed/donor", getDonationsConfirmedByDonor).Methods("GET")
	r.HandleFunc("/donations/confirmed/recipient", getDonationsConfirmedByRecipient).Methods("GET")
	r.HandleFunc("/donations/confirmed/business", getDonationsConfirmedByBusiness).Methods("GET")
	r.HandleFunc("/donations", createDonation).Methods("POST")
	r.HandleFunc("/donations/{id}", updateDonation).Methods("PUT")
	r.HandleFunc("/donations/{id}", deleteDonation).Methods("DELETE")

	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var u User
	err := db.QueryRow(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
	if err != nil {
		http.Error(w, "Usuario no encontrado", 404)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	var u User
	err := db.QueryRow(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE email=?`, email).
		Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
	if err != nil {
		http.Error(w, "Usuario no encontrado", 404)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func getUserByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE name LIKE ?`, "%"+name+"%")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	_, err := db.Exec(`INSERT INTO users (email, password, name, phone, user_type, address, 
		latitude, longitude, donation_days) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email, u.Password, u.Name, u.Phone, u.UserType, u.Address,
		u.Latitude, u.Longitude, u.DonationDays)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario creado"})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	_, err := db.Exec(`UPDATE users SET email=?, name=?, phone=?, user_type=?, address=?, 
		latitude=?, longitude=?, donation_days=? WHERE id=?`,
		u.Email, u.Name, u.Phone, u.UserType, u.Address,
		u.Latitude, u.Longitude, u.DonationDays, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func getDonations(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var d Donation
	err := db.QueryRow(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE id=?`, id).
		Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
	if err != nil {
		http.Error(w, "Donación no encontrada", 404)
		return
	}
	json.NewEncoder(w).Encode(d)
}

func createDonation(w http.ResponseWriter, r *http.Request) {
	var d Donation
	json.NewDecoder(r.Body).Decode(&d)
	_, err := db.Exec(`INSERT INTO donations (donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, pickup_time, weight,
		donation_reason, contact_info) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.DonorID, d.Title, d.Description, d.Category, d.Quantity, d.ExpiryDate,
		d.PickupAddress, d.PickupLatitude, d.PickupLongitude, d.Status, d.PickupTime,
		d.Weight, d.DonationReason, d.ContactInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "Donación creada"})
}

func updateDonation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var d Donation
	json.NewDecoder(r.Body).Decode(&d)
	_, err := db.Exec(`UPDATE donations SET title=?, description=?, category=?, quantity=?, expiry_date=?,
		pickup_address=?, pickup_latitude=?, pickup_longitude=?, status=?, reserved_by=?,
		pickup_time=?, pickup_person_name=?, pickup_person_id=?, verification_code=?,
		donor_confirmed=?, recipient_confirmed=?, business_confirmed=?, weight=?,
		donation_reason=?, contact_info=? WHERE id=?`,
		d.Title, d.Description, d.Category, d.Quantity, d.ExpiryDate,
		d.PickupAddress, d.PickupLatitude, d.PickupLongitude, d.Status, d.ReservedBy,
		d.PickupTime, d.PickupPersonName, d.PickupPersonID, d.VerificationCode,
		d.DonorConfirmed, d.RecipientConfirmed, d.BusinessConfirmed, d.Weight,
		d.DonationReason, d.ContactInfo, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func deleteDonation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.Exec("DELETE FROM donations WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func getDonationsByDonor(w http.ResponseWriter, r *http.Request) {
	donorID := mux.Vars(r)["donor_id"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE donor_id=?`, donorID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE category=?`, category)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsByStatus(w http.ResponseWriter, r *http.Request) {
	status := mux.Vars(r)["status"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE status=?`, status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getUsersByType(w http.ResponseWriter, r *http.Request) {
	userType := mux.Vars(r)["type"]
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE user_type=?`, userType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getDonationsReservedBy(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE reserved_by=?`, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getUsersByPhone(w http.ResponseWriter, r *http.Request) {
	phone := mux.Vars(r)["phone"]
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE phone LIKE ?`, "%"+phone+"%")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getUsersByAddress(w http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE address LIKE ?`, "%"+address+"%")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func getDonationsByTitle(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE title LIKE ?`, "%"+title+"%")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsByAddress(w http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE pickup_address LIKE ?`, "%"+address+"%")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsConfirmedByDonor(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE donor_confirmed = 1`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsConfirmedByRecipient(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE recipient_confirmed = 1`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationsConfirmedByBusiness(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, donor_id, title, description, category, quantity, expiry_date,
		pickup_address, pickup_latitude, pickup_longitude, status, reserved_by, reserved_at,
		pickup_time, pickup_person_name, pickup_person_id, verification_code, completed_at,
		created_at, updated_at, donor_confirmed, recipient_confirmed, business_confirmed,
		business_confirmed_at, donor_confirmed_at, recipient_confirmed_at, weight,
		donation_reason, contact_info FROM donations WHERE business_confirmed = 1`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var donations []Donation
	for rows.Next() {
		var d Donation
		rows.Scan(&d.ID, &d.DonorID, &d.Title, &d.Description, &d.Category, &d.Quantity, &d.ExpiryDate,
			&d.PickupAddress, &d.PickupLatitude, &d.PickupLongitude, &d.Status, &d.ReservedBy, &d.ReservedAt,
			&d.PickupTime, &d.PickupPersonName, &d.PickupPersonID, &d.VerificationCode, &d.CompletedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DonorConfirmed, &d.RecipientConfirmed, &d.BusinessConfirmed,
			&d.BusinessConfirmedAt, &d.DonorConfirmedAt, &d.RecipientConfirmedAt, &d.Weight,
			&d.DonationReason, &d.ContactInfo)
		donations = append(donations, d)
	}
	json.NewEncoder(w).Encode(donations)
}

func getUsersByDonationDays(w http.ResponseWriter, r *http.Request) {
	days := mux.Vars(r)["days"]
	rows, err := db.Query(`SELECT id, email, name, phone, user_type, address, latitude, longitude, 
		created_at, updated_at, donation_days FROM users WHERE JSON_CONTAINS(donation_days, ?)`, `"`+days+`"`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.UserType, &u.Address,
			&u.Latitude, &u.Longitude, &u.CreatedAt, &u.UpdatedAt, &u.DonationDays)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}
