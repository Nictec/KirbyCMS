package youtube

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/user"
    "path/filepath"

    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "KirbyCMS/shortcuts"
    "google.golang.org/api/youtube/v3"
)

//Auth functions

func GetClient(scope string, w http.ResponseWriter, request *http.Request) *http.Client {
    ctx := context.Background()

    //read the config
    b, err := ioutil.ReadFile("client_secret.json")
    if err != nil {
        shortcuts.HttpError(w, 500, "could not read client_secret.json")
        return nil
    }

    // If modifying the scope, delete your previously saved credentials
    // at ~/.credentials/youtube-go.json
    config, err := google.ConfigFromJSON(b, scope)
    if err != nil {
        shortcuts.HttpError(w, 500, "Unable to parse client secret file to config")
        return nil
    }

    //set the redirect uri
    config.RedirectURL = "http://localhost:8000/oauth"

    cacheFile, err := TokenCacheFile()
    if err != nil {
        shortcuts.HttpError(w, 500, "Unable to get path to cached credential file.")
    } 
    tok, err := tokenFromFile(cacheFile)
    if err != nil {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        Authenticate(authURL, w, request)
        return nil
    } else {
        return config.Client(ctx, tok)
    }
} 

func Authenticate(url string, w http.ResponseWriter, r *http.Request){
    http.Redirect(w, r, url, 301)
} 

func ExchangeToken(config *oauth2.Config, code string) (*oauth2.Token) {
    tok, err := config.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Fatalf("Unable to retrieve token %v", err)
    }
    return tok
} 

func TokenCacheFile() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", err
    }
    tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
    os.MkdirAll(tokenCacheDir, 0700)
    return filepath.Join(tokenCacheDir,
        url.QueryEscape("youtube-go.json")), err
} 

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    t := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(t)
    defer f.Close()
    return t, err
}

func SaveToken(file string, token *oauth2.Token) {
    fmt.Println("trying to save token")
    fmt.Printf("Saving credential file to: %s\n", file)
    f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
} 

//action functions
func Upload(w http.ResponseWriter, client *http.Client){
    _, err := youtube.New(client)
    if err != nil {
        shortcuts.HttpError(w, 500, "API Error: "+err.Error())
    }

}