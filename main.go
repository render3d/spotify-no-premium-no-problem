package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/render3d/spotify-no-premium-no-problem/playlist"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	scopes = []string{
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserReadPrivate,
	}
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes...))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Printf("Please log in to Spotify by visiting the following page in your browser:\n%s\n", url)

	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	fmt.Printf("You are logged in as: %s\n", user.ID)

	if err := playlist.UpdatePlaylist(context.Background(), client, user.ID); err != nil {
		log.Fatalf("Failed to update playlist: %v", err)
	}

	fmt.Println("Done.")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatalf("Failed to get token: %v", err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintln(w, "Login Completed!")
	ch <- client
}
