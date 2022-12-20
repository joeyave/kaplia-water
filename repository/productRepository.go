package repository

import (
	"context"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	MongoClient *mongo.Client
}

type Product struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Price              int                `bson:"price"`
	WholesalePrice     int                `bson:"wholesalePrice"`
	WholesaleThreshold int                `bson:"wholesaleThreshold"`
	PhotoURL           string             `bson:"photoURL"`
	Title              string             `bson:"title"`
	Description        string             `bson:"description"`
}

func (p *Product) GetPriceBasedOnCount(count int) int {
	productPrice := p.Price

	if p.WholesalePrice != 0 && p.WholesaleThreshold != 0 {
		if count >= p.WholesaleThreshold {
			productPrice = p.WholesalePrice
		}
	}
	return productPrice
}

var products = []*Product{
	{
		ID:                 objectIDFromHex("639789d7cbeb25111d000000"),
		Price:              75000,
		WholesalePrice:     60000,
		WholesaleThreshold: 2,
		PhotoURL:           "./img/cafe/water.png",
		Title:              "Вода 18,9л",
		Description:        "Бутиль очищеної води 18,9 л.️",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       160000,
		PhotoURL:    "./img/cafe/pump.png",
		Title:       "Помпа",
		Description: "Механічна помпа для бутилю з краником.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       350000,
		PhotoURL:    "./img/cafe/epump.png",
		Title:       "Помпа електрична",
		Description: "Електронна помпа для бутилю.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       370000,
		PhotoURL:    "./img/cafe/dispenser.png",
		Title:       "Диспенсер",
		Description: "Диспенсер для води.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       350000,
		PhotoURL:    "./img/cafe/bottle.png",
		Title:       "Бутиль",
		Description: "Полікарбонатний бутиль для води 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       370000,
		PhotoURL:    "./img/cafe/bottle-with-holder.png",
		Title:       "Бутиль з ручкою",
		Description: "Полікарбонатний бутиль для води з ручкою 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       80000,
		PhotoURL:    "./img/cafe/funnel.png",
		Title:       "Воронка",
		Description: "Воронка для бутилю 18.9 л.",
	},
	{
		ID:          primitive.NewObjectID(),
		Price:       80000,
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

func objectIDFromHex(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return id
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
