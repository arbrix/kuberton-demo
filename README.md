# kuberton-demo
Demo monolith system as initial task

## How to use

### Build

```
docker build -t <your-image-name> .
```

### Run

```
docker run --rm -it -p 3000:3000 <your-image-name>
```

Open: [shop](http://localhost:3000)

### Configuration

Key | Default Value | Description
---|---|---
PORT | `3000` |
BANNER_COLOR | "css property" |

## API

### Routes

Method | Route | Description
---|---|---
`GET` | `/` | home page (product list, link to the cart)
`GET`| `/product/{id}` | product page, select quantity, add to the cart. User `?json=true` for obtaining response at JSON format
`GET`| `/rate` | return list of supported rates at JSON format
`GET`| `/convert/{currency_id}/{price}` | return converted Money(price) from USD -> {currency_id}
`POST` | `/setCurrency` | change user currency preference
`GET` | `/static/*` | static files server
`GET` | `/_healthz` | container health check
