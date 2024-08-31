package playlist

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func UpdatePlaylist(ctx context.Context, client *spotify.Client, userID string) error {
	playlistResult, err := client.GetPlaylistsForUser(ctx, userID, spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		return fmt.Errorf("failed to get playlists for user: %w", err)
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

	tracksResult, err := client.CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Market("GB"), spotify.Offset(0))
	if err != nil {
		return fmt.Errorf("failed to get current user's tracks: %w", err)
	}

	fmt.Println("Songs:")
	for _, track := range tracksResult.Tracks {
		fmt.Println("\t", track.Name)
	}

	trackIDs := make([]spotify.ID, len(tracksResult.Tracks))
	for i, track := range tracksResult.Tracks {
		trackIDs[i] = track.ID
	}

	if err := client.ReplacePlaylistTracks(ctx, appPlaylist.ID, trackIDs...); err != nil {
		return fmt.Errorf("failed to replace playlist tracks: %w", err)
	}

	return nil
}
