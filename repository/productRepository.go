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
		Price:       55000,
		PhotoURL:    "./img/cafe/water.png",
		Title:       "Вода 18,9л",
		Description: "Бутиль очищеної води 18,9 л.️",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       155000,
		PhotoURL:    "./img/cafe/pump.png",
		Title:       "Помпа",
		Description: "Механічна помпа для бутилю з краником.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       250000,
		PhotoURL:    "./img/cafe/epump.png",
		Title:       "Помпа електрична",
		Description: "Електронна помпа для бутилю.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       300000,
		PhotoURL:    "./img/cafe/dispenser.png",
		Title:       "Диспенсер",
		Description: "Диспенсер для води.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       290000,
		PhotoURL:    "./img/cafe/bottle.png",
		Title:       "Бутиль",
		Description: "Полікарбонатний бутиль для води 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       290000,
		PhotoURL:    "./img/cafe/bottle-with-holder.png",
		Title:       "Бутиль з ручкою",
		Description: "Полікарбонатний бутиль для води з ручкою 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       120000,
		PhotoURL:    "./img/cafe/funnel.png",
		Title:       "Воронка",
		Description: "Воронка для бутилю 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       155000,
		PhotoURL:    "./img/cafe/holder.png",
		Title:       "Ручка",
		Description: "Ручка для переносу бутилю.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       3200000,
		PhotoURL:    "./img/cafe/cooler.png",
		Title:       "Кулер підлоговий",
		Description: "Електронний підлоговий кулер.",
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
