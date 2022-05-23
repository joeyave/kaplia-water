package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	MongoClient *mongo.Client
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Price       int                `bson:"price"`
	PhotoURL    string             `bson:"photoURL"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}

var products = []*Product{
	{
		ID:          primitive.NewObjectID(),
		Price:       69990,
		PhotoURL:    "./img/cafe/bottle.png",
		Title:       "Вода",
		Description: "Бутиль очищеної води 18,9 л.️",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       149990,
		PhotoURL:    "./img/cafe/pump.webp",
		Title:       "Помпа",
		Description: "Механічна помпа для бутилю.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       799990,
		PhotoURL:    "./img/cafe/filter.png",
		Title:       "Кулер",
		Description: "Кулер з підігрівом води.",
	},
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*Product, error) {
	return products, nil
}

func (r *ProductRepository) FindOneByID(ctx context.Context, ID primitive.ObjectID) (*Product, error) {
	for _, product := range products {
		if product.ID == ID {
			return product, nil
		}
	}

	return nil, mongo.ErrNoDocuments
}
