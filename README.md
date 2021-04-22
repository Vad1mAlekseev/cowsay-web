# Cowsay-web

Based on the [cowsay-cli program](https://en.wikipedia.org/wiki/Cowsay).

## Features

- Ability to enter custom text (for example: /bunny?text=HelloWorld!)
- Implementation of all original figures
- Displays a random message using [fortune](https://en.wikipedia.org/wiki/Fortune_(Unix))

## Run

To run this project locally, you need Docker:

```bash
docker-compose up
```

Then visit localhost:8080 to see cowsay-web. If you need to see grafana, visit localhost:3000.

## CLI usage

```
$ curl "cowsay-web.xyz/eyes?mode=plain"

_________________________________________
/ You will pay for your sins. If you have \
| already paid, please disregard this     |
\ message.                                /
 -----------------------------------------
    \
     \
                                   .::!!!!!!!:.
  .!!!!!:.                        .:!!!!!!!!!!!!
  ~~~~!!!!!!.                 .:!!!!!!!!!UWWW$$$ 
      :$$NWX!!:           .:!!!!!!XUWW$$$$$$$$$P 
      $$$$$##WX!:      .<!!!!UW$$$$"  $$$$$$$$# 
      $$$$$  $$$UX   :!!UW$$$$$$$$$   4$$$$$* 
      ^$$$B  $$$$\     $$$$$$$$$$$$   d$$R" 
        "*$bd$$$$      '*$$$$$$$$$$$o+#" 
             """"          """"""" 

```
