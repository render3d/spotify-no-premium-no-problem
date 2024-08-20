package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/zmb3/spotify/v2/auth"
	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	scopes = "playlist-read-private user-library-read playlist-modify-private user-read-email user-read-private"
	auth   = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes))
	ch     = make(chan *spotify.Client)
	state  = "abc123"
)

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:\n", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	playlistResult, err := client.GetPlaylistsForUser(context.Background(), user.ID, spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		log.Fatal(err)
	}

	var appPlaylist spotify.SimplePlaylist
	fmt.Println("Playlists:")
	for _, playlist := range playlistResult.Playlists {
		fmt.Println("\t", playlist.Name)
		if playlist.Name == "No Premium No Problem" {
			appPlaylist = playlist
			break
		}
	}

	tracksResult, err := client.CurrentUsersTracks(context.Background(), spotify.Limit(50), spotify.Market("GB"), spotify.Offset(0))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Songs:")
	for _, track := range tracksResult.Tracks {
		fmt.Println("\t", track.Name)
	}

	trackIds := make([]spotify.ID, len(tracksResult.Tracks))

    // Extract names from student objects
    for i, track := range tracksResult.Tracks {
        trackIds[i] = track.ID
    }

	err = client.ReplacePlaylistTracks(context.Background(), appPlaylist.ID, trackIds...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
	os.Exit(0)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
