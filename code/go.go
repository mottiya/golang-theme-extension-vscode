package main

import (
	"context"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"math/cmplx"
	"net/http"
	"time"
)

// comment

/*
multi-line comment
*/

const defaultTimeout = 5 * time.Second

const (
	_ = iota
	a
)

type UserName = string
type UserID string
type Role int

const (
	RoleGuest Role = iota
	RoleAdmin
	RoleAuditor
)

type User struct {
	ID      UserID
	Name    string
	Active  bool
	Balance float64
}

type AuditEvent struct {
	User
	Role    Role
	Tags    []string
	Aliases [3]string
	Done    chan struct{}
	Closer  NamedCloser
}

type NamedCloser interface {
	io.Closer
	Name() string
}

type Repository[T any] interface {
	Find(ctx context.Context, id UserID) (T, error)
	Save(ctx context.Context, value T) error
}

type CommandAction func(ctx context.Context) error

type Flag struct {
	Name  string
	Usage string
}

type Command struct {
	Name   string
	Usage  string
	Flags  []Flag
	Action CommandAction
}

type App struct {
	Commands []Command
}

type MemoryRepository struct {
	users map[UserID]User
}

func NewMemoryRepository(users map[UserID]User) *MemoryRepository {
	return &MemoryRepository{users: users}
}

func identity[T comparable](value T) T {
	return value
}

func joinLabels(prefix string, limit int, labels ...string) string {
	if limit > len(labels) {
		limit = len(labels)
	}

	result := prefix
	for index, label := range labels[:limit] {
		if index == 0 {
			result += ":"
		} else {
			result += ","
		}
		result += label
	}
	return result
}

// syntaxSamples touches [User], [fmt.Println], and [http.Server] for doc links.
//
//go:noinline
func syntaxSamples(input any) (summary string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			summary = fmt.Sprint(recovered)
		}
	}()

	raw := `raw\nstring`
	runeValue := 'A'
	imaginary := 1 + 2i
	_ = cmplx.Abs(complex(real(imaginary), imag(imaginary)))
	_ = runeValue

	names := make([]string, 0, 3)
	names = append(names, raw, "admin")
	copied := make([]string, len(names))
	copy(copied, names)
	fixedNames := [3]string{"root", "guest", "admin"}
	copied = append(copied, fixedNames[0])
	summary = joinLabels("names", 2, names...)
	summary = joinLabels(summary, len(fixedNames), fixedNames[:]...)

	cache := map[string]int{"first": 1, "second": 2}
	delete(cache, "missing")

	counter := new(int)
	*counter = len(names) + cap(copied)

	updates := make(chan int, 2)
	updates <- *counter
	updates <- cache["first"]
	close(updates)

scan:
	for value := range updates {
		switch {
		case value == 0:
			continue scan
		case value < 0:
			break scan
		default:
			summary = fmt.Sprintf("%s:%d", summary, value)
		}
	}

retry:
	if *counter < 0 {
		goto retry
	}

	switch value := input.(type) {
	case string:
		summary = value
	case User:
		summary = value.Name
	case nil:
		summary = "nil"
	default:
		summary = fmt.Sprintf("%T", value)
	}

	event := AuditEvent{
		User: User{
			ID:      identity[UserID]("42"),
			Name:    UserName("Ada"),
			Active:  true,
			Balance: 128.50,
		},
		Role:    RoleAdmin,
		Tags:    copied,
		Aliases: fixedNames,
		Done:    make(chan struct{}),
		Closer:  nil,
	}
	_ = event.User.Active
	_ = event.Aliases[0]
	close(event.Done)
	if event.Closer != nil {
		defer event.Closer.Close()
	}

	payload, _ := stdjson.Marshal(map[string]any{
		"role":    event.Role,
		"name":    event.Name,
		"tags":    event.Tags,
		"aliases": event.Aliases,
		"ready":   true,
	})

	switch event.Role {
	case RoleGuest:
		summary += ":guest"
	case RoleAdmin:
		summary += ":admin"
		fallthrough
	default:
		summary += ":" + string(payload)
	}

	return summary
}

func (r *MemoryRepository) Find(ctx context.Context, id UserID) (user User, err error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}
	var ok bool
	user, ok = r.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func (r *MemoryRepository) Save(ctx context.Context, user User) error {
	if user.ID == "" {
		return fmt.Errorf("empty user id")
	}

	r.users[user.ID] = user
	return nil
}

func handler(repo Repository[User]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
		defer cancel()

		user, err := repo.Find(ctx, UserID(r.URL.Query().Get("id")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		_, _ = fmt.Fprintf(w, "%s: %.2f", user.Name, user.Balance)
	}
}

func cmdServer(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fmt.Println("starting server command")
	return nil
}

func main() {
	app := App{
		Commands: []Command{
			{
				Name:  "server",
				Usage: "Run the service as a server",
				Flags: []Flag{
					{
						Name:  "automigrate",
						Usage: "Run database migrations before starting.",
					},
				},
				Action: cmdServer,
			},
		},
	}

	if err := app.Commands[0].Action(context.Background()); err != nil {
		panic(err)
	}

	repo := NewMemoryRepository(map[UserID]User{
		"42": {ID: "42", Name: "Ada", Active: true, Balance: 128.50},
	})
	saveMethod := repo.Save
	findExpression := (*MemoryRepository).Find
	_, _ = saveMethod, findExpression
	_ = syntaxSamples(User{Name: "Grace"})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler(repo),
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}
