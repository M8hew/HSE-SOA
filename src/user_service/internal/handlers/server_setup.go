package handlers

import (
	"crypto/rsa"
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

func connectGRPCServer() (*pb.ContentServiceClient, error) {
	contentServicePort := os.Getenv("CONTENT_SERVICE_PORT")
	conn, err := grpc.Dial("content_service:"+contentServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewContentServiceClient(conn)
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

	client, err := connectGRPCServer()
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
		contentService:      *client,
		kafkaProducer:       kafkaProd,
	}
	return &serverHandler, nil
}
