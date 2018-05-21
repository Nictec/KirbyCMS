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
    "path"
    "path/filepath"

    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "KirbyCMS/shortcuts"
    "google.golang.org/api/youtube/v3"
    "gopkg.in/mgo.v2/bson"
    "KirbyCMS/helpers/db"
    "KirbyCMS/models"
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

//goroutines
func handleUpload(w http.ResponseWriter, obj interface{}, file interface{}, objId bson.ObjectId){
    response, err := obj.Media(file).Do()
    if err != nil {
        shortcuts.HttpError(w, 500, "fatal error a upload thread: " + err.Error() + "aborting upload")
        return
    }
    c := db.Session.DB("KirbyCMS").C("videos")
    c.Update(bson.M{"_id":objId}, bson.M{"$set":bson.M{"yt_id":response}})
    c.Update(bson.M{"_id":objId}, bson.M{"$set":bson.M{"status":"uploaded"}})
}
//callable stuff
/*
    The following function does the following things:
    -------------------------------------------------

    1. initializing a new youtube client
    2. initializing all youtube specific structs
    3. creating a new request
    4. opening the cached file
    5. creating a new db entry for the vid
    6. opening a new thread and starting the upload
*/
func Upload(w http.ResponseWriter, client *http.Client, title string, description string, catId string, filename string, privacy string, keywords string){
    const MEDIAPATH = "media/videos/"
    service, err := youtube.New(client)
    if err != nil {
        shortcuts.HttpError(w, 500, "API Error: "+err.Error())
        return
    }
    upload := &youtube.Video{
        Snippet: &youtube.VideoSnippet{
            Title: title,
            Description: description,
            CategoryId: catId,
        },
        Status: &youtube.VideoStatus{PrivacyStatus: privacy},
    }

    // The API returns a 400 Bad Request response if tags is an empty string.
    if strings.Trim(keywords, "") != "" {
            upload.Snippet.Tags = strings.Split(*keywords, ",")
    }

    call := service.Videos.Insert("snippet,status", upload)
    fullpath := path.Join(MEDIAPATH, filename)
    file, err := os.Open(fullpath)
    defer file.Close()
    if err != nil{
        shortcuts.HttpError(w, 500, "Error opening cached videofile: " + err.Error())
        return
    }
    c := db.Session.DB("KirbyCMS").C("videos")
    initData := models.Video{
        Title: title,
        Description: description,
        CatId: catId,
        Privacy: privacy,
        Status: "uploading",
    }
    initData.Id = bson.NewObjectId()
    err = c.Insert(initData)
    if err != nil {
        http.Error(w, 500, "Database error: "+ err.Error())
        return
    }
    go handleUpload(w, call, file, initData.Id)
}