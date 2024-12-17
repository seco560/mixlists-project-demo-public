# Mixlists Project Demo

This is a very simple proof of concept for a playlist analysis project I have in mind. This currently does one simple thing.
It takes a JSON of playlists obtained from [Exported Spotify Data](https://www.spotify.com/us/account/privacy/) (at the bottom of this page, you can find a button for making a data export request) and prints out all songs that appear with the playlists they appear in. 
You can adjust if you want all songs printed or just songs that appear more than n times, where n is passed as a parameter. 

It's very easy to pass your own exported data from Spotify, but I've included a snippet of mine in order to test the program. It currently does not save output to a file, just prints it to console, the test-output.txt file was obtained via passing regular output: 
`go run . > test-output.txt`

This experiment was borne out of frustration with having too many playlists and forgetting which songs are which. With a structure like this, it's very easy to implement a lookup function that shows what playlists a certain song is in. 

For now it's barebones, and I'm publishing it like this to get some feedback on my go code. The first aspect that begs improvement
is code organisation, as well as clarity when handling indeces. A potential integration with Spotify API is the end goal, but for now it does its very basic job. 

eg. I forgot which playlists I added Threshold to!
`user@prompt % grep "Threshold with" -A 15 test-output.txt`

Song neighbours I found useful for transitions into and out of the song.