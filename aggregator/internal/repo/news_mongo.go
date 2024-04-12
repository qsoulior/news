package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/aggregator/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (any, error) {
		resultNews := new(entity.News)

		filter := bson.D{{Key: "link", Value: news.Link}}
		err := n.collection.FindOne(ctx, filter).Decode(resultNews)
		if err == mongo.ErrNoDocuments {
			return n.collection.InsertOne(ctx, news)
		}

		if err != nil {
			return nil, err
		}

		if news.PublishedAt.After(resultNews.PublishedAt) {
			filter := bson.D{{Key: "link", Value: resultNews.Link}}
			return n.collection.ReplaceOne(ctx, filter, news)
		}

		return nil, nil
	}, txnOptions)

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

	filter := bson.D{{Key: "_id", Value: objectID}}
	err = n.collection.FindOne(ctx, filter).Decode(news)
	if err != nil {
		return nil, fmt.Errorf("n.collection.FindOne.Decode: %w", err)
	}

	return news, nil
}

func (n *newsMongo) GetByQuery(ctx context.Context, query Query, opts Options) ([]entity.News, int, error) {
	match := make(bson.D, 0, 5)
	if query.Title {
		match = append(match, bson.E{
			Key:   "title",
			Value: primitive.Regex{Pattern: query.Text, Options: "i"},
		})
	} else {
		match = append(match, bson.E{
			Key:   "$text",
			Value: bson.D{{Key: "$search", Value: query.Text}},
		})
	}

	// TODO: $in
	match = append(match, bson.E{
		Key:   "source",
		Value: primitive.Regex{Pattern: query.Source, Options: "i"},
	})

	// TODO: authors $elem_match
	// match = append(match, bson.E{
	// 	Key:   "authors",
	// 	Value: "",
	// })

	tags := make(bson.A, len(query.Tags))
	for i, tag := range query.Tags {
		tags[i] = primitive.Regex{Pattern: fmt.Sprintf("^%s$", tag), Options: "i"}
	}

	if len(tags) > 0 {
		match = append(match, bson.E{
			Key:   "tags",
			Value: bson.D{{Key: "$all", Value: tags}},
		})
	}

	categories := make(bson.A, len(query.Categories))
	for i, category := range query.Categories {
		categories[i] = primitive.Regex{Pattern: fmt.Sprintf("^%s$", category), Options: "i"}
	}

	if len(categories) > 0 {
		match = append(match, bson.E{
			Key:   "categories",
			Value: bson.D{{Key: "$all", Value: categories}},
		})
	}

	matchStage := bson.D{{
		Key:   "$match",
		Value: match,
	}}

	paginationStage := bson.D{{
		Key: "$facet",
		Value: bson.D{
			{
				Key: "results",
				Value: bson.A{
					bson.D{{Key: "$skip", Value: opts.Skip}},
					bson.D{{Key: "$limit", Value: opts.Limit}},
				},
			},
			{
				Key: "total_results",
				Value: bson.A{
					bson.D{{Key: "$count", Value: "count"}},
				},
			},
		},
	}}

	unwindStage := bson.D{{
		Key:   "$unwind",
		Value: bson.D{{Key: "path", Value: "$total_results"}},
	}}

	projectStage := bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "results", Value: true},
			{Key: "total_count", Value: "$total_results.count"},
		},
	}}

	pipeline := mongo.Pipeline{matchStage, paginationStage, unwindStage, projectStage}
	cursor, err := n.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("n.collection.Aggregate: %w", err)
	}

	var res struct {
		Results    []entity.News `bson:"results"`
		TotalCount int           `bson:"total_count"`
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, 0, fmt.Errorf("cursor.Decode: %w", err)
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor.Err: %w", err)
	}

	return res.Results, res.TotalCount, nil
}
