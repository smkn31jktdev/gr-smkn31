package mongodb

import (
	"context"
	"errors"
	"strconv"

	aduanmodel "be-gr31/internal/model/aduan"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AduanStore struct {
	coll *mongo.Collection
}

func NewAduanStore(db *mongo.Database) *AduanStore {
	return &AduanStore{
		coll: db.Collection("aduan_siswa"),
	}
}

func (s *AduanStore) Create(ctx context.Context, data *aduanmodel.Aduan) error {
	_, err := s.coll.InsertOne(ctx, data)
	return err
}

func (s *AduanStore) Update(ctx context.Context, data *aduanmodel.Aduan) error {
	filter := bson.M{"ticketId": data.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := s.coll.ReplaceOne(ctx, filter, data, opts)
	return err
}

func (s *AduanStore) FindByID(ctx context.Context, id string) (*aduanmodel.Aduan, error) {
	filter := bson.M{"ticketId": id}
	var result aduanmodel.Aduan
	err := s.coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (s *AduanStore) ListPaged(ctx context.Context, filter aduanmodel.AduanFilter, pageSize int, pageState string) ([]aduanmodel.Aduan, string, error) {
	where := bson.M{}
	if filter.NISN != "" {
		where["nisn"] = filter.NISN
	}
	if filter.Status != "" {
		if filter.Status == "open" {
			where["status"] = bson.M{"$in": []string{"open", "pending"}}
		} else {
			where["status"] = filter.Status
		}
	}
	if filter.AdminNama != "" {
		where["adminNama"] = filter.AdminNama
	}
	if filter.Wali != "" {
		where["wali"] = filter.Wali
	}

	offset := 0
	if pageState != "" {
		if val, err := strconv.Atoi(pageState); err == nil {
			offset = val
		}
	}

	opts := options.Find().
		SetLimit(int64(pageSize)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"updatedAt": -1})

	cursor, err := s.coll.Find(ctx, where, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var items []aduanmodel.Aduan
	if err := cursor.All(ctx, &items); err != nil {
		return nil, "", err
	}

	// If the items returned is equal to pageSize, it might have more items, so provide offset state
	nextState := ""
	if len(items) == pageSize {
		nextState = strconv.Itoa(offset + len(items))
	}

	return items, nextState, nil
}
