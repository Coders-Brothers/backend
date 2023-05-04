package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicListenAddr                    string
	PublicURL                           string
	FrontendURL                         string
	PostgresURL                         string
	JwtSigningKey                       string
	DisableUiProxy                      bool
	AdminUiBackend                      string
	BookingUiBackend                    string
	SMTPHost                            string
	SMTPPort                            int
	SMTPSenderAddress                   string
	SMTPStartTLS                        bool
	SMTPInsecureSkipVerify              bool
	SMTPAuth                            bool
	SMTPAuthUser                        string
	SMTPAuthPass                        string
	MockSendmail                        bool
	PrintConfig                         bool
	Development                         bool
	InitOrgName                         string
	InitOrgDomain                       string
	InitOrgUser                         string
	InitOrgPass                         string
	InitOrgCountry                      string
	InitOrgLanguage                     string
	OrgSignupEnabled                    bool
	OrgSignupDomain                     string
	OrgSignupAdmin                      string
	OrgSignupMaxUsers                   int
	OrgSignupDelete                     bool
	LoginProtectionMaxFails             int
	LoginProtectionSlidingWindowSeconds int
	LoginProtectionBanMinutes           int
}

var _configInstance *Config
var _configOnce sync.Once

func GetConfig() *Config {
	_configOnce.Do(func() {
		_configInstance = &Config{}
		_configInstance.ReadConfig()
	})
	return _configInstance
}

func (c *Config) ReadConfig() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
		return
	}
	POSTGRES_URL := os.Getenv("POSTGRES_URL")
	JWT_SIGNING_KEY := os.Getenv("JWT_SIGNING_KEY")
	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_SENDER_ADDRESS := os.Getenv("SMTP_SENDER_ADDRESS")
	// SMTP_PORT := os.Getenv("SMTP_PORT")
	SMTP_AUTH_USER := os.Getenv("SMTP_AUTH_USER")
	INIT_ORG_NAME := os.Getenv("INIT_ORG_NAME")
	INIT_ORG_DOMAIN := os.Getenv("INIT_ORG_DOMAIN")
	INIT_ORG_COUNTRY := os.Getenv("INIT_ORG_COUNTRY")
	INIT_ORG_LANGUAGE := os.Getenv("INIT_ORG_LANGUAGE")
	ORG_SIGNUP_DOMAIN := os.Getenv("ORG_SIGNUP_DOMAIN")
	ORG_SIGNUP_ADMIN := os.Getenv("ORG_SIGNUP_ADMIN")
	INIT_ORG_USER := os.Getenv("INIT_ORG_USER")
	INIT_ORG_PASS := os.Getenv("INIT_ORG_PASS")
	// ORG_SIGNUP_MAX_USERS := os.Getenv("ORG_SIGNUP_MAX_USERS")
	log.Println(POSTGRES_URL)

	log.Println("Reading config...")
	c.Development = (c.getEnv("DEV", "0") == "1")
	c.PublicListenAddr = c.getEnv("PUBLIC_LISTEN_ADDR", "0.0.0.0:8080")
	c.PublicURL = strings.TrimSuffix(c.getEnv("PUBLIC_URL", "http://localhost:8080"), "/") + "/"
	if c.Development {
		c.FrontendURL = c.getEnv("FRONTEND_URL", "http://localhost:3000")
	} else {
		c.FrontendURL = c.getEnv("FRONTEND_URL", "http://localhost:8080")
	}
	c.FrontendURL = strings.TrimSuffix(c.FrontendURL, "/") + "/"
	c.DisableUiProxy = (c.getEnv("DISABLE_UI_PROXY", "0") == "1")
	c.AdminUiBackend = c.getEnv("ADMIN_UI_BACKEND", "localhost:3000")
	c.BookingUiBackend = c.getEnv("BOOKING_UI_BACKEND", "localhost:3001")
	c.PostgresURL = c.getEnv("POSTGRES_URL", POSTGRES_URL)
	c.JwtSigningKey = c.getEnv("JWT_SIGNING_KEY", JWT_SIGNING_KEY)
	c.SMTPHost = c.getEnv("SMTP_HOST", SMTP_HOST)
	c.SMTPPort = c.getEnvInt("SMTP_PORT", 12)
	c.SMTPStartTLS = (c.getEnv("SMTP_START_TLS", "0") == "1")
	c.SMTPInsecureSkipVerify = (c.getEnv("SMTP_INSECURE_SKIP_VERIFY", "0") == "1")
	c.SMTPAuth = (c.getEnv("SMTP_AUTH", "0") == "1")
	c.SMTPAuthUser = c.getEnv("SMTP_AUTH_USER", SMTP_AUTH_USER)
	c.SMTPAuthPass = c.getEnv("SMTP_AUTH_PASS", "")
	c.SMTPSenderAddress = c.getEnv("SMTP_SENDER_ADDRESS", SMTP_SENDER_ADDRESS)
	c.MockSendmail = (c.getEnv("MOCK_SENDMAIL", "0") == "1")
	c.PrintConfig = (c.getEnv("PRINT_CONFIG", "0") == "1")
	c.InitOrgName = c.getEnv("INIT_ORG_NAME", INIT_ORG_NAME)
	c.InitOrgDomain = c.getEnv("INIT_ORG_DOMAIN", INIT_ORG_DOMAIN)
	c.InitOrgUser = c.getEnv("INIT_ORG_USER", INIT_ORG_USER)
	c.InitOrgPass = c.getEnv("INIT_ORG_PASS", INIT_ORG_PASS)
	c.InitOrgCountry = c.getEnv("INIT_ORG_COUNTRY", INIT_ORG_COUNTRY)
	c.InitOrgLanguage = c.getEnv("INIT_ORG_LANGUAGE", INIT_ORG_LANGUAGE)
	c.OrgSignupEnabled = (c.getEnv("ORG_SIGNUP_ENABLED", "1") == "1")
	c.OrgSignupDomain = c.getEnv("ORG_SIGNUP_DOMAIN", ORG_SIGNUP_DOMAIN)
	c.OrgSignupAdmin = c.getEnv("ORG_SIGNUP_ADMIN", ORG_SIGNUP_ADMIN)
	c.OrgSignupMaxUsers = c.getEnvInt("ORG_SIGNUP_MAX_USERS", 20)
	c.OrgSignupDelete = (c.getEnv("ORG_SIGNUP_DELETE", "0") == "1")
	c.LoginProtectionMaxFails = c.getEnvInt("LOGIN_PROTECTION_MAX_FAILS", 10)
	c.LoginProtectionSlidingWindowSeconds = c.getEnvInt("LOGIN_PROTECTION_SLIDING_WINDOW_SECONDS", 600)
	c.LoginProtectionBanMinutes = c.getEnvInt("LOGIN_PROTECTION_BAN_MINUTES", 5)
}

func (c *Config) Print() {
	s, _ := json.MarshalIndent(c, "", "\t")
	log.Println("Using config:\n" + string(s))
}

func (c *Config) getEnv(key, defaultValue string) string {
	res := os.Getenv(key)
	if res == "" {
		return defaultValue
	}
	return res
}

func (c *Config) getEnvInt(key string, defaultValue int) int {
	val, err := strconv.Atoi(c.getEnv(key, strconv.Itoa(defaultValue)))
	if err != nil {
		log.Fatal("Could not parse " + key + " to int")
	}
	return val
}
