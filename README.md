# No Premium, No Problem

A short application to create a playlist of 50 of a user's most recently liked songs. A playlist of this length strikes a better balance between listening to recent favourites with older ones when the non-premium user's play order is set to shuffle, and particularly when a user's liked songs number in the hundreds.

## Execution

1. In Spotify, create a public or private playlist named "No Premium No Problem". If you want to use a different name for the playlist you can change the playlist name string in line 60 of [`main.go`](./main.go)

2. Create a file named `env.sh` which contains the following environment variables:

    ```bash
    #!/bin/bash
    export SPOTIFY_ID="<Spotify client ID"
    export SPOTIFY_SECRET="<Spotify client secret>"
    export USER="<username>"
    ```

3. Source the file:

    ```bash
    source ./env.sh
    ```

4. (Optional) Build the `main` program:

    ```bash
    go build main.go
    ```

5. Run with the main program with one of the following commands:

    ```bash
    go run main.go  # if you did **not** run step 3
    ./main          # if you did run step 3
    ```
