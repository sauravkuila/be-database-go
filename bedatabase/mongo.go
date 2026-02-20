package bedatabase

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func ConnectMongo(cfg DbConfig) (*mongo.Client, error) {
	var mongoURI string
	var clientOpts *options.ClientOptions

	timeout := cfg.ConnectTimeout
	if timeout <= 0 {
		timeout = DEFAULT_CONNECT_TIMEOUT
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if cfg.QueryParams == nil {
		cfg.QueryParams = make(map[string]string)
	}

	switch cfg.Mode {

	// TUNNEL MODE (SSM / SSH → localhost)
	case ModeTunnel:
		var (
			directConnection bool
			retryWrites      bool
			tlsValue         bool
		)

		// ignore the query param in uri as library directly handles it in SetDirect option
		if v, found := cfg.QueryParams["directConnection"]; !found {
			directConnection = true
		} else {
			directConnection = v == "true"
			// Remove "directConnection" from QueryParams map if present
			delete(cfg.QueryParams, "directConnection")
		}

		// ignore the query param in uri as library directly handles it in SetRetryWrites option
		if v, found := cfg.QueryParams["retryWrites"]; !found {
			retryWrites = false
		} else {
			retryWrites = v == "true"
			// Remove "retryWrites" from QueryParams map if present
			delete(cfg.QueryParams, "retryWrites")
		}

		// ignore the query param in uri as library directly handles it in tlsConfig.InsecureSkipVerify option
		if v, found := cfg.QueryParams["tls"]; !found {
			tlsValue = true
		} else {
			tlsValue = v == "true"
			// Remove "tls" from QueryParams map if present
			delete(cfg.QueryParams, "tls")
		}

		if val, found := cfg.QueryParams["authSource"]; !found {
			cfg.QueryParams["authSource"] = "admin"
		} else { // not needed since we are modifying the map in place, but added for clarity
			cfg.QueryParams["authSource"] = val
		}

		if val, found := cfg.QueryParams["tlsAllowInvalidHostnames"]; !found {
			cfg.QueryParams["tlsAllowInvalidHostnames"] = "true"
		} else { // not needed since we are modifying the map in place, but added for clarity
			cfg.QueryParams["tlsAllowInvalidHostnames"] = val
		}

		mongoURI = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)

		mongoURI = prepareUriWithParams(mongoURI, cfg.QueryParams)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: tlsValue,
		}

		clientOpts = options.Client().
			ApplyURI(mongoURI).
			SetTLSConfig(tlsConfig).
			SetDirect(directConnection).
			SetRetryWrites(retryWrites).
			SetServerSelectionTimeout(5 * time.Second)

	// DIRECT ATLAS MODE (mongodb+srv)
	case ModeAtlas:
		var (
			retryWrites  bool
			writeConcern *writeconcern.WriteConcern
		)

		// ignore the query param in uri as library directly handles it in SetRetryWrites option
		if v, found := cfg.QueryParams["retryWrites"]; !found {
			retryWrites = false
		} else {
			retryWrites = v == "true"
			// Remove "retryWrites" from QueryParams map if present
			delete(cfg.QueryParams, "retryWrites")
		}

		// ignore the query param in uri as library directly handles it in SetWriteConcern option
		if v, found := cfg.QueryParams["w"]; !found {
			writeConcern = writeconcern.Majority()
		} else {
			if v == "majority" {
				writeConcern = writeconcern.Majority()
			}
			// Remove "w" from QueryParams map if present
			delete(cfg.QueryParams, "w")
		}

		// cfg.QueryParams["retryWrites"] = "true"
		// cfg.QueryParams["w"] = "majority"

		mongoURI = fmt.Sprintf(
			"mongodb+srv://%s:%s@%s/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Database,
		)

		mongoURI = prepareUriWithParams(mongoURI, cfg.QueryParams)

		clientOpts = options.Client().
			ApplyURI(mongoURI).
			SetRetryWrites(retryWrites).
			SetWriteConcern(writeConcern)

	// PRIVATE VPC / PEERING MODE
	case ModePrivate:
		var (
			retryWrites  bool
			writeConcern *writeconcern.WriteConcern
		)

		// ignore the query param in uri as library directly handles it in SetRetryWrites option
		if v, found := cfg.QueryParams["retryWrites"]; !found {
			retryWrites = false
		} else {
			retryWrites = v == "true"
			// Remove "retryWrites" from QueryParams map if present
			delete(cfg.QueryParams, "retryWrites")
		}

		// ignore the query param in uri as library directly handles it in SetWriteConcern option
		if v, found := cfg.QueryParams["w"]; !found {
			writeConcern = writeconcern.Majority()
		} else {
			if v == "majority" {
				writeConcern = writeconcern.Majority()
			}
			// Remove "w" from QueryParams map if present
			delete(cfg.QueryParams, "w")
		}

		if val, found := cfg.QueryParams["authSource"]; !found {
			cfg.QueryParams["authSource"] = "admin"
		} else { // not needed since we are modifying the map in place, but added for clarity
			cfg.QueryParams["authSource"] = val
		}

		mongoURI = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s",
			cfg.User,
			cfg.Password,
			cfg.Host, // private IP or DNS
			cfg.Port,
			cfg.Database,
		)

		mongoURI = prepareUriWithParams(mongoURI, cfg.QueryParams)

		clientOpts = options.Client().
			ApplyURI(mongoURI).
			SetRetryWrites(retryWrites).
			SetWriteConcern(writeConcern)

	default:
		return nil, fmt.Errorf("unknown mongo mode: %s", cfg.Mode)
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	if cfg.ValidatePing {
		if err := client.Ping(ctx, nil); err != nil {
			client.Disconnect(ctx)
			return nil, fmt.Errorf("mongo ping error: %w", err)
		}
		log.Println("ping successful")
	}

	log.Printf("Mongo connected [%s]", cfg.Mode)
	return client, nil
}

func prepareUriWithParams(baseUri string, params map[string]string) string {
	if len(params) > 0 {
		queryString := "?"
		paramList := make([]string, 0)
		for key, value := range params {
			paramList = append(paramList, fmt.Sprintf("%s=%s", key, value))
		}
		queryString += strings.Join(paramList, "&")
		baseUri += queryString
	}
	fmt.Println("Connecting to MongoDB with URI: ", baseUri)
	return baseUri
}
