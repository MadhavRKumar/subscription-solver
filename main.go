package main

import (
    "context"
    "os"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscription struct {
    ID int `json:"id" uri:"id"`
    Name string `json:"name"`
    ProfileLimit int32 `json:"profileLimit"`
    Cost int32 `json:"cost"`
}

var subscriptions = []Subscription{
    { ID: 1, Name: "Netflix", ProfileLimit: 1, Cost: 2500 },
}

var conn *pgxpool.Pool


func main() {

    var err error
    conn, err = pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))

    if err != nil {
        os.Exit(1)
    }
    defer conn.Close()

    router := gin.Default()
    router.GET("/subscriptions", getSubscriptions)
    router.GET("/subscriptions/:id", getSubscription)
    router.POST("/subscriptions", postSubscriptions)
    router.DELETE("/subscriptions/:id", deleteSubscription)

    router.Run()
}

func getSubscriptions(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, subscriptions)
}

func getSubscription(c *gin.Context) {
    var subscription Subscription

    if err := c.ShouldBindUri(&subscription); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }

    for _, s := range subscriptions {
        if s.ID == subscription.ID {
            c.IndentedJSON(http.StatusOK, s)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Subscription not found"})

}

func postSubscriptions(c *gin.Context) {
    var newSubscription Subscription

    if err := c.BindJSON(&newSubscription); err != nil {
        return
    }

    err := conn.QueryRow(context.Background(), "INSERT into subscriptions ('name', 'profileLimit', 'cost') VALUES ($1, $2, $3)", newSubscription.Name, newSubscription.ProfileLimit, newSubscription.Cost).Scan()

    if err != nil {
        fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
        c.AbortWithError(http.StatusInternalServerError, err)
    }

    subscriptions = append(subscriptions, newSubscription)
    c.IndentedJSON(http.StatusCreated, newSubscription)
}

func deleteSubscription(c *gin.Context) {
    var subscription Subscription

    if err := c.ShouldBindUri(&subscription); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }

    for index, s := range subscriptions {
        if s.ID == subscription.ID {
            newSubscriptions := make([]Subscription, 0)
            newSubscriptions = append(newSubscriptions, subscriptions[:index]...)
            subscriptions = append(newSubscriptions, subscriptions[index+1:]...)

            c.Status(http.StatusNoContent)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Subscription not found"})

}
