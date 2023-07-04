package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/bxb454/instant-interpret/package/database"
	"github.com/bxb454/instant-interpret/package/translationservice"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"golang.org/x/text/language"
)

var ErrCacheMiss = errors.New("cache miss")

const (
	MsgTypeUserMessage = 1
	MsgTypeLangUpdate  = 2
)

var supportedLanguages = map[string]language.Tag{
	"en":    language.English,
	"es":    language.Spanish,
	"fr":    language.French,
	"de":    language.German,
	"it":    language.Italian,
	"ja":    language.Japanese,
	"ko":    language.Korean,
	"pt":    language.Portuguese,
	"ru":    language.Russian,
	"zh":    language.SimplifiedChinese,
	"zh-TW": language.TraditionalChinese,
	"ar":    language.Arabic,
	"cs":    language.Czech,
	"da":    language.Danish,
	"nl":    language.Dutch,
	"fi":    language.Finnish,
	"el":    language.Greek,
	"hi":    language.Hindi,
	"hu":    language.Hungarian,
	"id":    language.Indonesian,
	"no":    language.Norwegian,
	"pl":    language.Polish,
	"sv":    language.Swedish,
	"th":    language.Thai,
	"tr":    language.Turkish,
	"uk":    language.Ukrainian,
	"vi":    language.Vietnamese,
	"af":    language.Afrikaans,
	"am":    language.Amharic,
	"ro":    language.Romanian,
	"bg":    language.Bulgarian,
	"hr":    language.Croatian,
	"et":    language.Estonian,
	"fa":    language.Persian,
	"fil":   language.Filipino,
	"he":    language.Hebrew,
	"is":    language.Icelandic,
	"sw":    language.Swahili,
	"lv":    language.Latvian,
	"lt":    language.Lithuanian,
	"ms":    language.Malay,
	"sr":    language.Serbian,
	"sk":    language.Slovak,
	"sl":    language.Slovenian,
	"ur":    language.Urdu,
	"pu":    language.Punjabi,
	"ta":    language.Tamil,
	"te":    language.Telugu,
	"ne":    language.Nepali,
	"bn":    language.Bengali,
	"gu":    language.Gujarati,
}

type Client struct {
	ID                 string
	Conn               *websocket.Conn
	Pool               *Pool
	TranslationService *translationservice.TranslationService
	DB                 *database.Database
	Lang               string
	PreferredLang      string
	RedisClient        *redis.Client
}

type Message struct {
	Type             int    `json:"type"`
	Body             string `json:"body"`
	OriginalLanguage string `json:"originalLanguage"`
	Username         string `json:"username"`
}

func NewClient(conn *websocket.Conn, pool *Pool, ts *translationservice.TranslationService, db *database.Database, id string, lang string, preferredLang string) *Client {
	// Initialize Redis Client here
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &Client{
		ID:                 id,
		Conn:               conn,
		Pool:               pool,
		TranslationService: ts,
		DB:                 db,
		Lang:               lang,
		PreferredLang:      preferredLang,
		RedisClient:        rdb,
	}
}

func (c *Client) getTranslationFromCache(ctx context.Context, originalText string, targetLang string) (string, error) {
	key := fmt.Sprintf("%s:%s", originalText, targetLang)
	val, err := c.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func (c *Client) cacheTranslation(ctx context.Context, originalText string, targetLang string, translatedText string) error {
	key := fmt.Sprintf("%s:%s", originalText, targetLang)
	err := c.RedisClient.Set(ctx, key, translatedText, 0).Err() // 0 for no expiration
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{}
		log.Printf("Raw message: %s", string(p))
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println("Failed to parse message: ", err)
			continue
		}

		if message.Type == MsgTypeLangUpdate {
			_, ok := supportedLanguages[message.Body]
			if !ok {
				log.Println("Unsupported language: ", message.Body)
				continue
			}
			// Add log message here
			log.Printf("Updated language to: %s", message.Body)
			c.Lang = message.Body
			log.Println("Updated client language to: ", c.Lang) // Add this
			continue
		}

		if message.Type != MsgTypeUserMessage {
			log.Println("Unknown message type: ", message.Type)
			continue
		}

		// Detect language of the received message
		detectedLang, err := c.TranslationService.DetectLanguage(context.Background(), message.Body)
		if err != nil {
			log.Println("Failed to detect language: ", err)
			continue
		}

		// Iterate over all clients and send them translated messages
		for client := range c.Pool.Clients {
			langTag, ok := supportedLanguages[client.Lang]
			if !ok {
				log.Println("Unsupported language for translation: ", client.Lang)
				continue
			}

			// Check cache first
			ctx := context.Background()
			translatedText, err := c.getTranslationFromCache(ctx, message.Body, client.Lang)
			if err == ErrCacheMiss {
				// Translate message to the chosen language
				translatedText, err = c.TranslationService.TranslateText(ctx, langTag, message.Body)
				if err != nil {
					log.Println("Failed to translate text: ", err)
					continue
				}

				// Cache the translation
				if err = c.cacheTranslation(ctx, message.Body, client.Lang, translatedText); err != nil {
					log.Printf("Failed to cache translation: %v\n", err)
				}
			} else if err != nil {
				log.Printf("Failed to get translation from cache: %v\n", err)
			}

			// Send the translated message to the client
			translatedMessage := Message{
				Type:             message.Type,
				Body:             fmt.Sprintf("%s [translated from %s]", translatedText, detectedLang),
				OriginalLanguage: detectedLang,
				Username:         c.ID,
			}
			client.Conn.WriteJSON(translatedMessage)

			// Insert the message into the database
			_, err = c.DB.Exec("INSERT INTO messages (original_message, original_language, translated_message, translated_language) VALUES ($1, $2, $3, $4)",
				message.Body, detectedLang, translatedText, client.Lang)
			if err != nil {
				log.Printf("Failed to insert message into database: %v\n", err)
			}
			fmt.Printf("Message Sent to client %s: %+v\n", client.ID, translatedMessage)
		}
	}
}
