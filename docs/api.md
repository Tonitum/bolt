# API Design

## Main service

### /new

Accepted Method:(s) POST

Accepts a long URL and an optional user object that is the main write 
path for the service. The long url has a hash assigned to it and if the user
was provided, associates it with the user. 

If the user object is provided, the user token must be provided as a header
keyed TOKEN.

Returns the short url.

Example payload
```json
{
  "url": "string",
  "alias": "string", // optional
  "user": { // optional
    "id": "string" 
    },
}
```


Example return
```json
{
  "short_url": "string"
}
```

### /<short url>

Accepted Method(s): GET, POST (future, for user updates), DELETE (future, for user deletes)

If the short url exists, redirects the client to the long url

## Auth (future)

### /login

POST

### /logout

POST

### /auth

POST

