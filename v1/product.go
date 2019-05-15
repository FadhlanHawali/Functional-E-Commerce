package v1

type Product struct {
	NamaBarang string `json:"namaBarang"`
	Deskripsi string `json:"deskripsi"`
	Quantity string `json:"quantity"`
	Harga int `json:"harga"`
	UrlGambar string `json:"urlGambar"`
}
