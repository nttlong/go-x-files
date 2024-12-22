package Configreader

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Server configuration
	CacheServer       string `yaml:"cache_server"` // Optional commented out line
	CachePrefix       string `yaml:"cache_prefix"`
	HostURL           string `yaml:"host_url"`
	Db                string `yaml:"db"`
	DbCodx            string `yaml:"db_codx_"`
	AdminDbName       string `yaml:"admin_db_name"`
	JinjaTemplatesDir string `yaml:"jinja_templates_dir"`
	StaticDir         string `yaml:"static_dir"`
	Bind              string `yaml:"bind"`
	HostDir           string `yaml:"host_dir"`
	Workers           int    `yaml:"workers"`
	OnPremiseTenant   string `yaml:"on_premise_tenant"`

	// Basic Auth configuration
	BasicAuth struct {
		UI struct {
			Realm    string `yaml:"realm"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"ui"`
	} `yaml:"basic_auth"`

	// Server settings
	SecurityMode            string `yaml:"security_mode"`
	TimeoutKeepAlive        int    `yaml:"timeout_keep_alive"`
	TimeoutGracefulShutdown int    `yaml:"timeout_graceful_shutdown"`
	ServerType              string `yaml:"server_type"`
	WorkerClass             string `yaml:"worker_class"`
	H2MaxConcurrentStreams  int    `yaml:"h2_max_concurrent_streams"`

	// JWT configuration
	Jwt struct {
		SecretKey                string `yaml:"secret_key"`
		Algorithm                string `yaml:"algorithm"`
		AccessTokenExpireMinutes int    `yaml:"access_token_expire_minutes"`
		SecretUser               string `yaml:"secret_user"`
		SecretPassword           string `yaml:"secret_password"`
	} `yaml:"jwt"`

	// Elastic Search configuration
	ElasticSearch struct {
		Server                 []string `yaml:"server"`
		PrefixIndexLvCodx_     string   `yaml:"prefix_index_: lv-codx"`
		PrefixIndex            string   `yaml:"prefix_index"`
		IndexMaxAnalyzedOffset int      `yaml:"index_max_analyzed_offset"`
		MaxAnalyzedOffset      int      `yaml:"max_analyzed_offset"`
		FieldContent           string   `yaml:"field_content"`
	} `yaml:"elastic_search"`

	// Allowed office file extensions
	ExtOfficeFile []string `yaml:"ext_office_file"`

	// Allowed video file extensions
	ExtVideoFile []string `yaml:"ext_video_file"`

	// Allowed image file extensions
	ExtImageFile []string `yaml:"ext_image_file"`

	// RabbitMQ configuration
	Rabbitmq struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		Msg      string `yaml:"msg"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`

	// Service mappings
	Services []string `yaml:"services"`

	// File storage configuration
	FileStoragePath    string `yaml:"file_storage_path"`
	FileStorageEncrypt bool   `yaml:"file_storage_encrypt"`
	SharedStorage      string `yaml:"shared_storage"`

	// Minio storage configuration
	MinioStorage struct {
		EndPoint  string `yaml:"end_point"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Secure    bool   `yaml:"secure"`
	} `yaml:"minio_storage"`

	// Dataset path
	DatasetPath string `yaml:"dataset_path"`

	// Server for generating thumbnails of office documents
	ServerThumbOffice string `yaml:"server_thumb_office"`

	// Host for process services
	ProcessServicesHost string `yaml:"process_services_host"`

	// Supported cloud providers
	CloudsSupport []string `yaml:"clouds_support"`

	// Remote service URLs
	RemoteOffice  string `yaml:"remote_office"`
	RemoteOffice_ string `yaml:"remote_office_"`
	RemotePdf     string `yaml:"remote_pdf"`
	RemotePdf_    string `yaml:"remote_pdf_"`
	RemoteVideo   string `yaml:"remote_video"`
	RemoteOcr     string `yaml:"remote_ocr"`
	TikaServer    string `yaml:"tika_server"`
	RemoteThumb   string `yaml:"remote_thumb"`
	Log           struct {
		Path string `yaml:"path"`
		Rot  int    `yaml:"rote"`
		Size int    `yaml:"size"`
		Fmt  string `yaml:"format"`
	} `yaml:"log"`
}

var (
	config     *Config
	configOnce sync.Once
	configErr  error
)

func GetAppPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func LoadConfig(filename string) *Config {
	configOnce.Do(func() {
		// Read the YAML file
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		// Unmarshal the YAML data into the Config struct
		config = &Config{}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			panic(err)
		}

	})

	return config
}

// print the config as jon format with indentation
func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}
	jsonData, err := json.MarshalIndent(c, "", "    ") // Use MarshalIndent for pretty printing
	if err != nil {
		panic(fmt.Sprintf("Error marshaling config to JSON: %v", err))
	}
	return string(jsonData)
}
