package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"


    "github.com/MadhavRKumar/subscription-solver/internal/subscriptions"
)

var subscriptions_list = []subscriptions.Subscription{
    { Name: "Netflix", ProfileLimit: 1, Cost: 2500 },
}


type SubscriptionsHandler struct {
    store subscriptionStore
}

type subscriptionStore interface {
    Add(uuid string, subscription subscriptions.Subscription) (subscriptions.Subscription, error)
    Get(uuid string) (subscriptions.Subscription, error)
    List() (map[string]subscriptions.Subscription, error)
    Update(uuid string, subscription subscriptions.Subscription) error
    Remove(uuid string) error
}

func NewSubscriptionsHandler(s subscriptionStore) *SubscriptionsHandler {
    return &SubscriptionsHandler{
        store: s,
    }
}

func (h *SubscriptionsHandler) CreateSubscription(c *gin.Context) {
    var newSubscription subscriptions.Subscription

    if err := c.BindJSON(&newSubscription); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    uuid := uuid.New().String()

    sub, err := h.store.Add(uuid, newSubscription)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    } 

    c.IndentedJSON(http.StatusCreated, sub)
}

func (h *SubscriptionsHandler) ListSubscription(c *gin.Context) { 
    subscriptions, err := h.store.List()

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.IndentedJSON(http.StatusOK, subscriptions)
}



func main() {

    store := subscriptions.NewMemStore()
    handler := NewSubscriptionsHandler(store)

    router := gin.Default()
    router.GET("/subscriptions", handler.ListSubscription)
    router.GET("/subscriptions/:id", getSubscription)
    router.POST("/subscriptions", handler.CreateSubscription)
    router.DELETE("/subscriptions/:id", deleteSubscription)

    router.Run()
}

func getSubscriptions(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, subscriptions_list)
}

func getSubscription(c *gin.Context) {
}


func deleteSubscription(c *gin.Context) {
}
