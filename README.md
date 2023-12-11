# Sneakfeed Service

Sneakfeed is a social media application like Twitter.
This project aims to play with unfamiliar techstacks such as Mongodb and Imagekit APIs.

## API List

### No auth zone

- `GET /ping` -> Health check
- `POST /register` -> Register as a user
- `POST /login` -> Login as a user
- `POST /refresh` -> Refresh the access token with refresh token

### Auth zone

- `GET /posts` -> Get all posts
- `GET /posts/:postID` -> Get a single post
- `POST /posts` -> Create a new post
- `GET /posts/:postID/comments` -> Get all comments in post
- `POST /posts/:postID/comments` -> Add a new comment to the post
- `POST /posts/:postID/like` -> Like a post

- `GET /profiles` -> Get current user profile
- `PATCH /profiles` -> Update profile image

## Environment Variables

- ACCESS_TOKEN_SECRET -> { ACCESS_TOKEN_SECRET - can be any }
- REFRESH_TOKEN_SECRET -> { REFRESH_TOKEN_SECRET - can be any}
- IMAGEKIT_PUBLIC_KEY -> { IMAGEKIT_PUBLIC_KEY }
- IMAGEKIT_PRIVATE_KEY -> { IMAGEKIT_PRIVATE_KEY }
- IMAGEKIT_ENDPOINT_URL -> https://ik.imagekit.io/ { YOUR_USERNAME }
- MONGODB_USERNAME -> { MONGODB_USERNAME }
- MONGODB_PASSWORD -> { MONGODB_PASSWORD }
- DATABASE_NAME -> { DATABASE_NAME }

## How to run the project locally?

Create .env file on based path and add environment variable on above
OR export variable as environment variables

Then run the following command on your based path

```
go mod tidy
go run .
```

If you see below text, meaning it works.

```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

Pinged your deployment. You successfully connected to MongoDB!
[GIN-debug] GET    /ping                     --> github.com/gin-contrib/cors.New.func1 (3 handlers)
[GIN-debug] POST   /register                 --> github.com/tipbk/sneakfeed-service/handler.UserHandler.Register-fm (4 handlers)
[GIN-debug] POST   /login                    --> github.com/tipbk/sneakfeed-service/handler.UserHandler.Login-fm (4 handlers)
[GIN-debug] POST   /refresh                  --> github.com/tipbk/sneakfeed-service/handler.UserHandler.RefreshToken-fm (4 handlers)
[GIN-debug] GET    /posts                    --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.GetPosts-fm (5 handlers)
[GIN-debug] GET    /posts/:postID            --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.GetPostByID-fm (5 handlers)
[GIN-debug] GET    /posts/:postID/comments   --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.GetCommentByPostID-fm (5 handlers)
[GIN-debug] POST   /posts                    --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.CreatePost-fm (5 handlers)
[GIN-debug] POST   /posts/:postID/comments   --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.AddComment-fm (5 handlers)
[GIN-debug] GET    /profiles                 --> github.com/tipbk/sneakfeed-service/handler.UserHandler.GetProfile-fm (5 handlers)
[GIN-debug] PATCH  /profiles                 --> github.com/tipbk/sneakfeed-service/handler.UserHandler.PartiallyUpdateProfile-fm (5 handlers)
[GIN-debug] POST   /posts/:postID/like       --> github.com/tipbk/sneakfeed-service/handler.ContentHandler.ToggleLikePostByID-fm (5 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
```
