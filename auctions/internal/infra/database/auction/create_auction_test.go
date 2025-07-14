package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoTestDB(t *testing.T) *mongo.Database {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		uri = "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	require.NoError(t, err)
	return client.Database("auctions")
}

func TestCloseExpiredAuctions(t *testing.T) {
	db := setupMongoTestDB(t)
	repo := NewAuctionRepository(db)
	ctx := context.Background()

	// Limpa a coleção antes do teste
	_, _ = repo.Collection.DeleteMany(ctx, bson.M{})

	// Cria leilão vencido
	auction := &AuctionEntityMongo{
		Id:          "testid",
		ProductName: "Produto Teste",
		Category:    "Teste",
		Description: "Leilão para teste automático",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-10 * time.Minute).Unix(), // vencido
	}
	_, err := repo.Collection.InsertOne(ctx, auction)
	require.NoError(t, err)

	// Executa fechamento
	err = repo.CloseExpiredAuctions(ctx)
	require.NoError(t, err)

	// Busca e valida se foi fechado
	var updated AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": "testid"}).Decode(&updated)
	require.NoError(t, err)
	require.Equal(t, auction_entity.Completed, updated.Status)
}
