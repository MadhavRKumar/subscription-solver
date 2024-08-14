package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/MadhavRKumar/subscription-solver/internal/subscriptions"
    "github.com/MadhavRKumar/subscription-solver/internal/middleware/cors"
)

type SubscriptionsHandler struct {
	store subscriptionStore
}

type subscriptionStore interface {
	Add(uuid string, subscription subscriptions.Subscription) (subscriptions.Subscription, error)
	Get(uuid string) (subscriptions.Subscription, error)
	List() ([]subscriptions.Subscription, error)
	Update(uuid string, subscription subscriptions.Subscription) error
	Remove(uuid string) error
}

func main() {
	store := subscriptions.NewMemStore()
	handler := NewSubscriptionsHandler(store)

	router := gin.New()
    router.Use(cors.Cors())
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

	router.GET("/subscriptions", handler.ListSubscription)
	router.GET("/subscriptions/:id", handler.GetSubscription)
	router.POST("/subscriptions", handler.CreateSubscription)
	router.DELETE("/subscriptions/:id", handler.DeleteSubscription)
	router.PATCH("/subscriptions/:id", handler.UpdateSubscription)

	err := router.Run()
	if err != nil {
		panic(err)
	}
}

func NewSubscriptionsHandler(s subscriptionStore) *SubscriptionsHandler {
	return &SubscriptionsHandler{
		store: s,
	}
}

func (h *SubscriptionsHandler) CreateSubscription(c *gin.Context) {
	var newSubscription subscriptions.Subscription

	if err := c.BindJSON(&newSubscription); err != nil {

		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			return
		}
		return
	}

	uuid := uuid.New().String()

	sub, err := h.store.Add(uuid, newSubscription)
	if err != nil {
		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			return
		}

		return
	}

	c.IndentedJSON(http.StatusCreated, sub)
}

func (h *SubscriptionsHandler) ListSubscription(c *gin.Context) {
	subscriptions, err := h.store.List()
	if err != nil {
		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			return
		}

		return
	}

	c.IndentedJSON(http.StatusOK, subscriptions)
}

func (h *SubscriptionsHandler) GetSubscription(c *gin.Context) {
	uuid := c.Param("id")

	sub, err := h.store.Get(uuid)
	if err != nil {

		if _, ok := err.(*subscriptions.NotFoundError); ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			return
		}
		return
	}

	c.IndentedJSON(http.StatusOK, sub)
}

func (h *SubscriptionsHandler) DeleteSubscription(c *gin.Context) {
	uuid := c.Param("id")

	err := h.store.Remove(uuid)
	if err != nil {
		if _, ok := err.(*subscriptions.NotFoundError); ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			panic(err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionsHandler) UpdateSubscription(c *gin.Context) {
    uuid := c.Param("id")

    var updatedSubscription subscriptions.Subscription

    if err := c.BindJSON(&updatedSubscription); err != nil {
        if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
            return
        }
        return
    }

    err := h.store.Update(uuid, updatedSubscription)

    if err != nil {
        if _, ok := err.(*subscriptions.NotFoundError); ok {
            c.AbortWithStatus(http.StatusNotFound)
            return
        }

        if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
            return
        }
        return
    }

    c.Status(http.StatusNoContent)
}
