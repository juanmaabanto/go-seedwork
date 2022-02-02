package seedwork

import (
	"context"
	"strconv"

	"github.com/juanmaabanto/go-seedwork/seedwork/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBaseRepository interface {
	Count(ctx context.Context, filter interface{}) (int64, error)
	DeleteById(ctx context.Context, id string) (int64, error)
	FilterBy(ctx context.Context, filter interface{}, receiver []interface{}) error
	FindById(ctx context.Context, id string, receiver interface{}) error
	FindOne(ctx context.Context, filter interface{}, receiver interface{}) error
	InsertMany(ctx context.Context, documents []interface{}) ([]string, error)
	InsertOne(ctx context.Context, document interface{}) (string, error)
	Paginated(ctx context.Context, filter interface{}, sort interface{}, pageSize int64, start int64, receiver interface{}) error
	UpdateOne(ctx context.Context, document interface{}) error
}

type BaseRepository struct {
	collection *mongo.Collection
}

func NewBaseRepository(connection database.MongoConnection, document IDocument) *BaseRepository {
	repository := &BaseRepository{
		collection: connection.Database.Collection(document.GetCollectionName()),
	}

	return repository
}

func (repo BaseRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	result, err := repo.collection.CountDocuments(ctx, filter)

	return result, err
}

func (repo BaseRepository) DeleteById(ctx context.Context, id string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	var result *mongo.DeleteResult

	if err != nil {
		result, err = repo.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	} else {
		result, err = repo.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objID}})
	}

	return result.DeletedCount, err
}

func (repo BaseRepository) FilterBy(ctx context.Context, filter interface{}, receiver []interface{}) error {
	cursor, err := repo.collection.Find(ctx, filter)

	if err != nil {
		return err
	}

	cursor.Decode(receiver)

	return nil
}

func (repo BaseRepository) FindById(ctx context.Context, id string, receiver interface{}) error {
	var result *mongo.SingleResult
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		val, errVal := strconv.Atoi(id)

		if errVal == nil {
			result = repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: val}})
		} else {
			result = repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
		}
	} else {
		result = repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}})
	}

	return result.Decode(receiver)
}

func (repo BaseRepository) FindOne(ctx context.Context, filter interface{}, receiver interface{}) error {
	result := repo.collection.FindOne(ctx, filter)

	if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
		return result.Err()
	}

	err := result.Decode(receiver)

	if err != nil {
		return err
	}

	return nil
}

func (repo BaseRepository) InsertMany(ctx context.Context, documents []interface{}) ([]string, error) {
	result, err := repo.collection.InsertMany(ctx, documents)

	if err != nil {
		panic(err)
	}

	array := []string{}

	for i := range result.InsertedIDs {
		array = append(array, result.InsertedIDs[i].(primitive.ObjectID).String())
	}

	return array, err
}

func (repo BaseRepository) InsertOne(ctx context.Context, document interface{}) (string, error) {
	result, err := repo.collection.InsertOne(ctx, document)

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (repo BaseRepository) Paginated(ctx context.Context, filter interface{}, sort interface{}, pageSize int64, start int64, receiver interface{}) error {
	options := options.Find()

	options.SetSort(sort)
	options.SetSkip(start)
	options.SetLimit(pageSize)

	cursor, err := repo.collection.Find(ctx, filter, options)

	if err != nil {
		return err
	}

	cursor.Decode(receiver)

	return nil
}

func (repo BaseRepository) UpdateOne(ctx context.Context, document interface{}) error {
	_, err := repo.collection.UpdateOne(ctx, bson.D{}, document)

	return err
}

func (repository BaseRepository) GetCollection() *mongo.Collection {
	return repository.collection
}
