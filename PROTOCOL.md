# PROTOCOL

Slight specifications for the protocol. HTTPS seems an okay enough
protocol for sharing bigger files accross the network, when
considering complexity to implement.

All the bellow actions assume that you're logged in. A login token is
passed through headers of the requests.

## POST

```nocode
POST /register/
    Register user to serice.
    {"username": username,
     "password": password}
    RETURN 400, if errors (eg, user exists, etc)
    RETURN 200, with empty body

POST /login/
    {"username": username,
     "password": password}
    RETURN 200
        {"token": token}

Login user to service.

POST /upload/<filename>/<timestamp>
    Header/Authorization: token <token>
    BODY: IMAGE DATA
    RETURN 400, on bad credentials
    RETURN 200, on success
    RETURN 500, on server go boom

POST /manifest/
    {"token": token,
     "files": [filelist]}
    RETURN 400, on bad credentials
    RETURN 500, on server go boom
    RETURN 200, on success
        {"filename_1": <md5sum>,
         "filename_2": <md5sum>,
         ...
         "filename_n": <md5sum>}
```

Assuming that the device store on itself the date that we are
interested in. The date can't be preserved via other means, so we pass
it in the url. The server can then modify the modtime accordingly.

## HEAD

HEAD /status/
    RETURN 200

Status of service

May be used later for checking what pictures have been shared between
machines.
