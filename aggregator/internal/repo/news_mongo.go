package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/aggregator/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type newsMongo struct {
	collection *mongo.Collection
}

func NewNewsMongo(database *mongo.Database) News {
	return &newsMongo{
		collection: database.Collection("news"),
	}
}

func (n *newsMongo) Create(ctx context.Context, news entity.News) error {
	_, err := n.collection.InsertOne(ctx, news)
	if err != nil {
		return fmt.Errorf("n.collection.InsertOne: %w", err)
	}
	return nil
}

func (n *newsMongo) ReplaceOrCreate(ctx context.Context, news entity.News) error {
	session, err := n.collection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("client.StartSession: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (any, error) {
		resultNews := new(entity.News)
		err := n.collection.FindOne(ctx, bson.M{"link": news.Link}).Decode(resultNews)
		if err == mongo.ErrNoDocuments {
			return n.collection.InsertOne(ctx, resultNews)
		}

		if err != nil {
			return nil, err
		}

		if news.PublishedAt.After(resultNews.PublishedAt) {
			return n.collection.ReplaceOne(ctx, bson.M{"link": resultNews.Link}, news)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("session.WithTransaction: %w", err)
	}

	return nil
}

func (n *newsMongo) CreateMany(ctx context.Context, news []entity.News) error {
	documents := make([]any, len(news))
	for i, v := range news {
		documents[i] = v
	}

	_, err := n.collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("n.collection.InsertMany: %w", err)
	}
	return nil
}

func (n *newsMongo) GetByID(ctx context.Context, id string) (*entity.News, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex: %w", err)
	}

	news := new(entity.News)
	err = n.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(news)
	if err != nil {
		return nil, fmt.Errorf("n.collection.FindOne.Decode: %w", err)
	}

	return news, nil
}

func (n *newsMongo) GetByQuery(ctx context.Context, query Query, opts Options) ([]entity.News, error) {
	filter := bson.M{
		"title":  fmt.Sprintf("/%s/i", query.Title),
		"source": fmt.Sprintf("/%s/i", query.Source),
	}

	findOpts := options.Find()
	findOpts.SetSkip(int64(opts.Skip))
	findOpts.SetLimit(int64(opts.Limit))

	cursor, err := n.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("n.collection.Find: %w", err)
	}

	var news []entity.News
	err = cursor.All(ctx, &news)
	if err != nil {
		return nil, fmt.Errorf("cursor.All: %w", err)
	}

	return news, nil
}
