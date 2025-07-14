package auction

import (
	"context"
	"go.uber.org/zap"
	"os"
	"time"

	"acution_dlcdev/configuration/logger"
	"acution_dlcdev/internal/entity/auction_entity"
	"acution_dlcdev/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

// getAuctionInterval retorna a duração da vida útil de um leilão
func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil || duration <= 0 {
		return 5 * time.Minute // padrão
	}
	return duration
}

// getAuctionCheckInterval retorna o intervalo para checar leilões vencidos
func getAuctionCheckInterval() time.Duration {
	checkInterval := os.Getenv("AUCTION_CHECK_INTERVAL")
	duration, err := time.ParseDuration(checkInterval)
	if err != nil || duration <= 0 {
		return 1 * time.Minute // padrão
	}
	return duration
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

// CloseExpiredAuctions fecha todos os leilões cujo timestamp + intervalo do leilão já passou, marcando status Completed
func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) error {
	auctionInterval := getAuctionInterval()
	nowUnix := time.Now().Unix()

	// timestamp limite para considerar expiração: agora - duração do leilão
	expirationLimit := nowUnix - int64(auctionInterval.Seconds())

	filter := bson.M{
		"timestamp": bson.M{"$lte": expirationLimit},
		"status":    auction_entity.Active,
	}

	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	res, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Failed to close expired auctions", err)
		return err
	}

	logger.Info("Closed %d expired auctions", zap.Int64("modified_count", res.ModifiedCount))
	return nil
}

// StartAuctionExpirationChecker inicia uma goroutine que verifica periodicamente se existem leilões vencidos para fechar
func (ar *AuctionRepository) StartAuctionExpirationChecker(ctx context.Context) {
	ticker := time.NewTicker(getAuctionCheckInterval())
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopping auction expiration checker")
				ticker.Stop()
				return
			case <-ticker.C:
				err := ar.CloseExpiredAuctions(context.Background())
				if err != nil {
					logger.Error("Error during CloseExpiredAuctions", err)
				}
			}
		}
	}()
}
