package aduan

// Aduan
type Aduan struct {
	ID            string          `json:"id" bson:"ticketId"`
	MongoID       any             `json:"-" bson:"_id,omitempty"`
	NISN          string          `json:"nisn" bson:"nisn"`
	NamaSiswa     string          `json:"namaSiswa" bson:"namaSiswa"`
	Kelas         string          `json:"kelas" bson:"kelas"`
	Messages      []Message       `json:"messages" bson:"messages"`
	Status        string          `json:"status" bson:"status"`
	AdminNama     string          `json:"adminNama" bson:"adminNama"`
	CreatedAt     string          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     string          `json:"updatedAt" bson:"updatedAt"`
	StatusHistory []StatusHistory `json:"statusHistory,omitempty" bson:"statusHistory,omitempty"`
	Walas         string          `json:"walas,omitempty" bson:"walas,omitempty"`
}

// Message
type Message struct {
	ID        string `json:"id,omitempty" bson:"id,omitempty"`
	From      string `json:"from,omitempty" bson:"from,omitempty"`
	Role      string `json:"role" bson:"role"`
	Isi       string `json:"isi" bson:"message"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
}

type StatusHistory struct {
	Status    string `json:"status" bson:"status"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
	Role      string `json:"role" bson:"role"`
	UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
	Note      string `json:"note" bson:"note"`
}

// AduanCreateRequest
type AduanCreateRequest struct {
	Isi string `json:"isi" binding:"required"`
}

// AduanRespondRequest
type AduanRespondRequest struct {
	AduanID string `json:"aduanId" binding:"required"`
	Isi     string `json:"isi" binding:"required"`
}

// AduanStatusRequest
type AduanStatusRequest struct {
	AduanID string `json:"aduanId" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=open in_progress closed pending"`
}

// AduanFilter
type AduanFilter struct {
	NISN      string
	Status    string
	AdminRole string
	AdminNama string
	Page      int
	Limit     int
}

