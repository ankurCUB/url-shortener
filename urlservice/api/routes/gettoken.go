package routes

import (
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/adityasvat/url-shortener/database"
	"github.com/gofiber/fiber/v2"
)

type response2 struct {
	Start           int        `json:"start"`
	End             int        `json:"end"`
}

func GetTokenRange(c *fiber.Ctx) error {

	var endpoint = fmt.Sprintf("%s/api/create", os.Getenv("TOKEN_SERVICE"))
	resp, err := http.Post(endpoint, "application/json", nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send POST request",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	var createResp response2
	// if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "failed to decode response",
	// 	})
	// }

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "failed to read response body",
            "details": err.Error(),
        })
    }

    if err := json.Unmarshal(body, &createResp); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "failed to decode response",
            "details": err.Error(),
        })
    }

    fmt.Printf("Decoded response: %#v\n", createResp)

	r := database.CreateClient(2)
	defer r.Close()

	// Store start and end values in Redis
	if err := r.Set(database.Ctx, "start", createResp.Start, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to store start value",
			"details": err.Error(),
		})
	}
	if err := r.Set(database.Ctx, "end", createResp.End, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to store end value",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(createResp)
}