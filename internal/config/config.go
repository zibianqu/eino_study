package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	VectorDB VectorDBConfig `mapstructure:"vectordb"`
	Eino     EinoConfig     `mapstructure:"eino"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type VectorDBConfig struct {
	Type      string `mapstructure:"type"`
	Dimension int    `mapstructure:"dimension"`
}

type EinoConfig struct {
	LLM       LLMConfig       `mapstructure:"llm"`
	Embedding EmbeddingConfig `mapstructure:"embedding"`
	Splitter  SplitterConfig  `mapstructure:"splitter"`
	Retriever RetrieverConfig `mapstructure:"retriever"`
}

type LLMConfig struct {
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

type EmbeddingConfig struct {
	Provider  string `mapstructure:"provider"`
	APIKey    string `mapstructure:"api_key"`
	Model     string `mapstructure:"model"`
	Dimension int    `mapstructure:"dimension"`
}

type SplitterConfig struct {
	ChunkSize    int `mapstructure:"chunk_size"`
	ChunkOverlap int `mapstructure:"chunk_overlap"`
}

type RetrieverConfig struct {
	TopK               int     `mapstructure:"top_k"`
	SimilarityThreshold float64 `mapstructure:"similarity_threshold"`
}

type LogConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// GetConfig returns the global config
func GetConfig() *Config {
	return GlobalConfig
}