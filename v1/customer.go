package v1



type Customer struct {
	Owner string `json:"owner"`
	Nama string `json:"nama"`
	Kontak Kontak `json:"kontak"`
	AlamatPengiriman string `json:"alamatPengiriman"`
}