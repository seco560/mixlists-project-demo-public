package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"errors"
	"bytes"
	"fmt"
	"os"
)

/** @author Seco Vlad 17.12.2024
* 	Mixlists Project - A way to gain insights into Spotify playlist data
*	This takes a JSON of playlists formatted per Spotify export and processes some useful data.
*	At its current point it prints all songs
**/

// structs defined precisely as json is structured
// TODO Interface and getter methods to minimize accessor nesting

type DecodedPlaylists struct {
	Playlists []SpotifyJsonPlaylist `json:"playlists"`
}

type SpotifyJsonPlaylist struct {
	Name string `json:"name"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Items []SpotifyJsonTrack `json:"items"`
}

type SpotifyJsonTrack struct {
	TrackDetails SpotifyJsonTrackData `json:"track"`
	Episode string `json:"episode"`
	Audiobook string `json:"audiobook"`
	LocalTrack string `json:"localTrack"`
	AddedDate string `json:"addedDate"`
}

type SpotifyJsonTrackData struct {
	TrackName string `json:"trackName"`
	ArtistName string `json:"artistName"`
	AlbumName string `json:"albumName"`
	TrackUri string `json:"trackUri"`
}

func main() {
	// change filename for input here
	playlists, err := parseAndLoadMixlists("demo-input/demo-mixlists.json")
	if err != nil {
		fmt.Printf("Error encountered: %s", err)
	}

	allPlaylists := playlists.Playlists

	// Print songs that have over n appearances,
	// Includes name of playlist, index in playlist and song's neighbours
	songCountMap := makeSongCounts(allPlaylists)
	printAppearancesAndNeighbours(songCountMap, 1, allPlaylists)

	// Initial simpler function - print all playlists a given song is in
	// Suppress output using silence boolean - true means no printing
	songTest := "the quiet things that no one ever knows"
	checkForSong(songTest, allPlaylists, true)
}

func parseAndLoadMixlists(filename string) (DecodedPlaylists, error) {
	mixlistsAndCo, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return DecodedPlaylists{}, errors.New("problem loading json file")
	}

	playlists := DecodedPlaylists{}
	decoder := json.NewDecoder(bytes.NewReader(mixlistsAndCo))

	if err := decoder.Decode(&playlists); err != nil {
		return DecodedPlaylists{}, errors.New("error decoding response body")
	}

	return playlists, nil
}

// Creates a map for each song's playlist occurrences.
// K - songName ; v - []playlistTitles
func makeSongCounts(playlists []SpotifyJsonPlaylist) map[string][]string {
	songCounts := make(map[string][]string)
	for _, playlist := range playlists {
		for i, track := range playlist.Items {
			songWithArtist := track.TrackDetails.ArtistName + " - " + track.TrackDetails.TrackName
			if _, ok := songCounts[songWithArtist]; !ok {
				songCounts[songWithArtist] = []string{}
			}
			songCounts[songWithArtist] = 
				append(songCounts[songWithArtist], playlist.Name + 
					" [#" + strconv.FormatInt(int64(i+1), 10) + "/" + strconv.FormatInt(int64(len(playlist.Items)), 10) + "]" +
					" last modified: " + playlist.LastModifiedDate)
		}
	}
	return songCounts
}

// To print appearances and neighbours for all songs simply pass n == 1
func printAppearancesAndNeighbours(songCounts map[string][]string, n int, playlists []SpotifyJsonPlaylist) {
	counter := 1 // an id of sorts
	for song, mixlists := range songCounts {
		if len(mixlists) >= n {
			fmt.Printf("%d) %s with ", counter, song)
			if len(mixlists) == 1{
				fmt.Printf("%d appearance:\n", len(mixlists))
			} else {
				fmt.Printf("%d appearances:\n", len(mixlists))
			}
			counter++
			// This part here is just to extract the index from playlist string
			// A tuple would have sufficed... or a two member slice, I suppose
			// It's probably the dodgiest part of this program, and uses FieldsFunc with ascii code checks to encode and decode index
			index := 0
			// Print all playlists first, then occurrences
			for _, pls := range songCounts[song] {
				fmt.Println(pls)
			}
			fmt.Println()
			for _, pls := range songCounts[song] {
				firstIndex := strings.FieldsFunc(pls, func (a rune) bool { return a == 35 }) // hashtag
				secondIndex := strings.FieldsFunc(firstIndex[1], func (a rune) bool { return a == 47 }) // forward slash
				intval, err := strconv.Atoi(secondIndex[0])
				if err != nil {
					fmt.Println("That was not a number mate")
				}
				index = intval - 1
				// Print occurrences in current playlist
				for _, playlist := range playlists {
					// Chop off the square bracket to get playlist name without index
					shortName := strings.FieldsFunc(pls, func (a rune) bool { return a == 91 })[0] // open square bracket
					shortName = shortName[:len(shortName) - 1] // chop off final character - leftover whitespace
					if playlist.Name == shortName {
						printFollowedPrecededBy(playlist, index)
					}
				}
			}
			fmt.Printf("---------------------------------------------------\n\n")
		}
	}
}

func printFollowedPrecededBy(playlist SpotifyJsonPlaylist, i int){
	if (len(playlist.Items) == 0 || len(playlist.Items) == 1) {
		fmt.Println("This playlist is either empty or it has only one song!")
		return
	}
	fmt.Printf("Song's neighbours in playlist: %s\n", playlist.Name)
	if i > 0 {
		// is not first track
		fmt.Println(playlist.Items[i - 1].TrackDetails.ArtistName + " - " + playlist.Items[i - 1].TrackDetails.TrackName)
	}
	fmt.Println(">>" + playlist.Items[i].TrackDetails.ArtistName + " - " + playlist.Items[i].TrackDetails.TrackName + "<<")
	if i < len(playlist.Items) - 1 {
		// is not last track
		fmt.Println(playlist.Items[i + 1].TrackDetails.ArtistName + " - " + playlist.Items[i + 1].TrackDetails.TrackName)
	}
	fmt.Println()
}

// First experiment - no longer used, kept around for posterity
func checkForSong(songName string, playlists []SpotifyJsonPlaylist, silence bool) {
	if !silence { fmt.Println("Test for song check...") }
	for _, playlist := range playlists {
		for i, track := range playlist.Items {
			if strings.Contains(strings.ToLower(track.TrackDetails.TrackName), strings.ToLower(songName)) {
				if !silence { 
					fmt.Printf("%s from %s is song #%d in playlist %s\n", 
					track.TrackDetails.TrackName, track.TrackDetails.ArtistName, i+1, playlist.Name)
				}
			}
		}
	}
}