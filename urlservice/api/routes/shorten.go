package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/adityasvat/url-shortener/database"
	"github.com/adityasvat/url-shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// implement rate limiting
	// for a user query, check if the IP in DB
	// 10 req/30 min
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err() 
	} else {
		// valInt, _ := strconv.Atoi(val)
		// Uncomment for rate
		// if valInt <= 0 {
		// 	limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
		// 	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		// 		"error":            "Rate limit exceeded",
		// 		"rate_limit_reset": limit / time.Nanosecond / time.Minute,
		// 	})
		// }
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "haha... nice try",
		})
	}

	// enforce SSL https
	// all url will be converted to https before storing in database
	body.URL = helpers.EnforceHTTP(body.URL)

	// implement collision checks
	var id string

	r3 := database.CreateClient(2)
	defer r3.Close()
	start, err := r3.Get(database.Ctx, "start").Result()
	if err == redis.Nil {
		// If start value doesn't exist, fetch token range
		if err := GetTokenRange(c); err != nil {
			return err
		}
		// Retry
		start, err = r3.Get(database.Ctx, "start").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "failed to retrieve token range",
				"details": err.Error(),
			})
		}
	}
	end, err := r3.Get(database.Ctx, "end").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to retrieve token range",
			"details": err.Error(),
		})
	}

	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)
	if startInt > endInt {
		// Fetch new token range if start exceeds end
		if err := GetTokenRange(c); err != nil {
			return err
		}
		// Retry
		start, err = r3.Get(database.Ctx, "start").Result()
		if err != nil {
			return err
		}
		end, err = r3.Get(database.Ctx, "end").Result()
		if err != nil {
			return err
		}
		startInt, _ = strconv.Atoi(start)
		endInt, _ = strconv.Atoi(end)
	}
	// Generate UUID based on fetched start value
	id = generateShortUUID(startInt)

	// Increment start value in Redis after successfully sending the response
	startInt++
	if err := r3.Set(database.Ctx, "start", strconv.Itoa(startInt), 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to increment start value",
			"details": err.Error(),
		})
	}

	r := database.CreateClient(0)
	defer r.Close()

	// not needed, unique urls using token
	// val, _ = r.Get(database.Ctx, id).Result()
	// // check if the user provided short is already in use
	// if val != "" {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": "URL short already in use",
	// 	})
	// }
	
	if body.Expiry == 0 {
		body.Expiry = 24 // default expiry of 24 hours
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}
	// respond with the url, short, expiry in hours, calls remaining and time to reset
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	r2.Decr(database.Ctx, c.IP())
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}

func generateShortUUID(id int) string {
    charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    base := len(charset)
    res := make([]byte, 6)
    for i := range res {
        res[i] = 'a'
    }
    id -= 1
    for i := 5; i >= 0; i-- {
        power := 1
        for j := 0; j < i; j++ {
            power *= base
        }
        index := id / power
        res[5-i] = charset[index]
        id %= power
    }
    return string(res)
}
