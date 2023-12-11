# Sneakfeed Service

Sneakfeed is a social media application like Twitter.
This project aims to play with unfamiliar techstacks such as Mongodb and Imagekit APIs.

## API List

### No auth zone

- GET /ping -> Health check
- POST /register -> Register as a user
- POST /login -> Login as a user
- POST /refresh -> Refresh the access token with refresh token

### Auth zone

- GET /posts -> Get all posts
- GET /posts/:postID -> Get a single post
- POST /posts -> Create a new post
- GET /posts/:postID/comments -> Get all comments in post
- POST /posts/:postID/comments -> Add a new comment to the post
- POST /posts/:postID/like -> Like a post

- GET /profiles -> Get current user profile
- PATCH /profiles -> Update profile image

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
go run .
```
