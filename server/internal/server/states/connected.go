package states

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/internal/server"
	"server/internal/server/db"
	"server/pkg/packets"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Connected struct {
	client  server.ClientInterfacer
	logger  *log.Logger
	queries *db.Queries
	dbCtx   context.Context
}

func (c *Connected) Name() string {
	return "Connected"
}

func (c *Connected) SetClient(client server.ClientInterfacer) {
	c.client = client
	loggingPrefix := fmt.Sprintf("Client %d [%s]: ", client.Id(), c.Name())
	c.logger = log.New(log.Writer(), loggingPrefix, log.LstdFlags)
	c.queries = client.DbTx().Queries
	c.dbCtx = client.DbTx().Ctx
}

func (c *Connected) OnEnter() {
	c.client.SocketSend(packets.NewId(c.client.Id()))
	//create a new user in the database
	user, err := c.client.DbTx().Queries.CreateUser(c.client.DbTx().Ctx, db.CreateUserParams{
		Username:     "test",
		PasswordHash: "test123",
	})

	if err != nil {
		c.logger.Fatalf("Failed to create user: %s", err)
	} else {
		c.logger.Printf("Created user: %v", user)
	}
}

func (c *Connected) HandleMessage(senderId uint64, message packets.Msg) {
	switch message := message.(type) {
	case *packets.Packet_LoginRequest:
		c.handleLoginRequest(senderId, message)
	case *packets.Packet_RegisterRequest:
		c.handleRegisterRequest(senderId, message)
	}
}

func (c *Connected) OnExit() {

}

func (c *Connected) handleLoginRequest(senderId uint64, message *packets.Packet_LoginRequest) {
	if senderId != c.client.Id() {
		c.logger.Printf("Received login request from another client(Id %d)", senderId)
		return
	}
	username := message.LoginRequest.Username
	genericFailMessage := packets.NewDenyResponse("Incorrect username or password")

	user, err := c.queries.GetUserByUsername(c.dbCtx, strings.ToLower(username))
	if err != nil {
		c.logger.Printf("Error getting user by username: %v", err)
		c.client.SocketSend(genericFailMessage)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(message.LoginRequest.Password))
	if err != nil {
		c.logger.Printf("Incorrect password for user: %s", username)
		c.client.SocketSend(genericFailMessage)
		return
	}

	c.logger.Printf("User %s logged in successfully!", username)
	c.client.SocketSend(packets.NewOkResponse())
}

func (c *Connected) handleRegisterRequest(senderId uint64, message *packets.Packet_RegisterRequest) {
	if senderId != c.client.Id() {
		c.logger.Printf("Received login request from another client(Id %d)", senderId)
		return
	}
	username := message.RegisterRequest.Username
	err := validateUsername(username)
	if err != nil {
		reason := fmt.Sprintf("Invalid username: %v", err)
		c.logger.Println(reason)
		c.client.SocketSend(packets.NewDenyResponse(reason))
		return
	}

	if _, err := c.queries.GetUserByUsername(c.dbCtx, strings.ToLower(username)); err != nil {
		c.logger.Printf("User already exists: %v", err)
		c.client.SocketSend(packets.NewDenyResponse("User already exists"))
		return
	}

	genericFailMessage := packets.NewDenyResponse("Failed to register user(internal server error) - pls try again later")

	// add new user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(message.RegisterRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.logger.Printf("Failed to hash password: %v", err)
		c.client.SocketSend(genericFailMessage)
		return
	}
	_, err = c.queries.CreateUser(c.dbCtx, db.CreateUserParams{
		Username:     strings.ToLower(username),
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		c.logger.Printf("Failed to create user: %v", err)
		c.client.SocketSend(genericFailMessage)
		return
	}
	c.logger.Printf("User %s registered successfully!", username)
	c.client.SocketSend(packets.NewOkResponse())
}

func validateUsername(username string) error {
	if len(username) <= 0 {
		return errors.New("empty")
	}
	if len(username) > 20 {
		return errors.New("Too long")
	}
	if username != strings.TrimSpace(username) {
		return errors.New("Leading or Trailing whitespace")
	}
	return nil
}
