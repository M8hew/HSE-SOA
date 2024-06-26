package handlers

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	pb "user_service/api/proto"
)

type rsaKeys struct {
	jwtPrivate *rsa.PrivateKey
	jwtPublic  *rsa.PublicKey
}

func parseYAMLConfig(yamlConfigPath string) (map[string]any, error) {
	yamlFile, err := os.ReadFile(yamlConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config.yaml: %v", err)
	}

	var configMap map[string]any
	if err := yaml.Unmarshal(yamlFile, &configMap); err != nil {
		return nil, fmt.Errorf("error parsing config.yaml: %v", err)
	}

	return configMap, nil
}

func parseJWTKeys(jwtPublicPath, jwtPrivatePath string) (rsaKeys, error) {
	private, err := os.ReadFile(jwtPrivatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading jwtPrivateFile: %v", err)
		os.Exit(1)
	}

	public, err := os.ReadFile(jwtPublicPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading jwtPublicFile: %v", err)
		os.Exit(1)
	}

	jwtPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jwtPublic, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return rsaKeys{
		jwtPrivate: jwtPrivate,
		jwtPublic:  jwtPublic,
	}, nil
}

func connectContentService() (*pb.ContentServiceClient, error) {
	contentServicePort := os.Getenv("CONTENT_SERVICE_PORT")
	if contentServicePort == "" {
		return nil, errors.New("content service grpc port unknown")
	}

	conn, err := grpc.NewClient("content_service:"+contentServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewContentServiceClient(conn)
	return &client, nil
}

func connectStatService() (*pb.StatServiceClient, error) {
	statServicePort := os.Getenv("STAT_SERVICE_RPC_PORT")
	if statServicePort == "" {
		return nil, errors.New("stat service grpc port unknown")
	}

	serverAddr := "stat_service:" + statServicePort
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewStatServiceClient(conn)
	return &client, nil
}

func NewServerHandler(configPath string) (*ServerHandler, error) {
	configMap, err := parseYAMLConfig(configPath)
	if err != nil {
		return nil, err
	}

	var (
		ok                 bool
		jwtPublicFile      string
		jwtPrivateFile     string
		secSessionLifetime int
	)
	if secSessionLifetime, ok = configMap["user_sessions_lifetime"].(int); !ok {
		return nil, fmt.Errorf("error parsing user_session_lifetime field in config")
	}
	if jwtPublicFile, ok = configMap["jwt_public_path"].(string); !ok {
		return nil, fmt.Errorf("error parsing jwt_public_path field in config")
	}
	if jwtPrivateFile, ok = configMap["jwt_private_path"].(string); !ok {
		return nil, fmt.Errorf("error parsing jwt_private_path field in config")
	}

	keys, err := parseJWTKeys(jwtPublicFile, jwtPrivateFile)
	if err != nil {
		return nil, err
	}

	dbWrapper_, err := NewDBWrapper()
	if err != nil {
		return nil, err
	}

	contentClient, err := connectContentService()
	if err != nil {
		return nil, err
	}

	statClient, err := connectStatService()
	if err != nil {
		return nil, err
	}

	kafkaProd, err := setKafkaUp()
	if err != nil {
		return nil, err
	}

	serverHandler := ServerHandler{
		userSessionLifeTime: time.Duration(secSessionLifetime) * time.Second,
		keys:                keys,
		db:                  dbWrapper_,
		contentService:      *contentClient,
		statService:         *statClient,
		kafkaProducer:       kafkaProd,
	}
	return &serverHandler, nil
}
