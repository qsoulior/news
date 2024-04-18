package repo

import (
	"context"
	"errors"
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
		return nil, ErrInvalidID
	}

	news := new(entity.News)

	filter := bson.D{{Key: "_id", Value: objectID}}
	err = n.collection.FindOne(ctx, filter).Decode(news)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("n.collection.FindOne.Decode: %w", err)
	}

	return news, nil
}

var sortVariants = map[SortOption]bson.D{
	SortPublishedAtDesc: {{Key: "published_at", Value: -1}},
	SortPublishedAtAsc:  {{Key: "published_at", Value: 1}},
	SortRelevanceDesc:   {{Key: "score", Value: -1}},
	SortRelevanceAsc:    {{Key: "score", Value: 1}},
}

func (n *newsMongo) parseQuery(query Query) bson.D {
	doc := make(bson.D, 0, 6)
	if query.Text != "" {
		if query.Title {
			doc = append(doc, bson.E{
				Key:   "title",
				Value: primitive.Regex{Pattern: query.Text, Options: "i"},
			})
		} else {
			doc = append(doc, bson.E{
				Key:   "$text",
				Value: bson.D{{Key: "$search", Value: query.Text}},
			})
		}
	}

	if len(query.Sources) > 0 {
		doc = append(doc, bson.E{
			Key:   "source",
			Value: bson.D{{Key: "$in", Value: query.Sources}},
		})
	}

	authors := make(bson.A, len(query.Authors))
	for i, tag := range query.Authors {
		authors[i] = primitive.Regex{Pattern: fmt.Sprintf("^%s$", tag), Options: "i"}
	}

	if len(authors) > 0 {
		doc = append(doc, bson.E{
			Key:   "authors",
			Value: bson.D{{Key: "$all", Value: authors}},
		})
	}

	tags := make(bson.A, len(query.Tags))
	for i, tag := range query.Tags {
		tags[i] = primitive.Regex{Pattern: fmt.Sprintf("^%s$", tag), Options: "i"}
	}

	if len(tags) > 0 {
		doc = append(doc, bson.E{
			Key:   "tags",
			Value: bson.D{{Key: "$all", Value: tags}},
		})
	}

	dateCond := make(bson.D, 0, 2)
	if query.DateFrom != nil {
		dateCond = append(dateCond, bson.E{Key: "$gte", Value: *query.DateFrom})
	}

	if query.DateTo != nil {
		dateCond = append(dateCond, bson.E{Key: "$lt", Value: *query.DateTo})
	}

	if len(dateCond) > 0 {
		doc = append(doc, bson.E{
			Key:   "published_at",
			Value: dateCond,
		})
	}

	return doc
}

func (n *newsMongo) parseOptions(opts Options) mongo.Pipeline {
	pipeline := make(mongo.Pipeline, 0, 4)

	// sort stage
	if opts.Sort.IsRelevance() {
		pipeline = append(pipeline, bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}},
			}},
		}})
	}

	sort, ok := sortVariants[opts.Sort]
	if !ok {
		sort = sortVariants[SortDefault]
	}

	pipeline = append(pipeline, bson.D{{
		Key:   "$sort",
		Value: sort,
	}})

	// pagination stage
	pipeline = append(pipeline, bson.D{{
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
	}})

	pipeline = append(pipeline, bson.D{{
		Key:   "$unwind",
		Value: "$total_results",
	}})

	return pipeline
}

func (n *newsMongo) GetByQuery(ctx context.Context, query Query, opts Options) ([]entity.NewsHead, int, error) {
	pipeline := make(mongo.Pipeline, 0, 6)

	// match stage
	pipeline = append(pipeline, bson.D{{
		Key:   "$match",
		Value: n.parseQuery(query),
	}})

	pipeline = append(pipeline, n.parseOptions(opts)...)

	// project stage
	pipeline = append(pipeline, bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "results", Value: bson.D{
				{Key: "_id", Value: true},
				{Key: "title", Value: true},
				{Key: "description", Value: true},
				{Key: "source", Value: true},
				{Key: "published_at", Value: true},
			}},
			{Key: "total_count", Value: "$total_results.count"},
		},
	}})

	cursor, err := n.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("n.collection.Aggregate: %w", err)
	}

	var val struct {
		Results    []entity.NewsHead `bson:"results"`
		TotalCount int               `bson:"total_count"`
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&val)
		if err != nil {
			return nil, 0, fmt.Errorf("cursor.Decode: %w", err)
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor.Err: %w", err)
	}

	if val.Results == nil {
		val.Results = make([]entity.NewsHead, 0)
	}

	return val.Results, val.TotalCount, nil
}
