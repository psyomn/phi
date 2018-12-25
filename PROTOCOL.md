# PROTOCOL

Slight specifications for the protocol. HTTPS seems an okay enough
protocol for sharing bigger files accross the network, when
considering complexity to implement.

All the bellow actions assume that you're logged in. A login token is
passed through headers of the requests.

## POST

POST /register/
    username/password
    RETURN 200
    RETURN

Register user to serice.

POST /login/
    username/password
    RETURN <login-token>

Login user to service.


POST /upload/yyyy/mm/dd/hh/mm/ss
    <multipart-data>

Assuming that the device store on itself the date that we are
interested in. The date can't be preserved via other means, so we pass
it in the url. The server can then modify the modtime accordingly.

## HEAD

HEAD /status/
    RETURN 200

Status of service


HEAD /exist/yyyy/mm/filename
    Auth-Token: Basic <login-token>
    RETURN 200 for yes, 404 for no

Does the file exist?

## GET

GET /md5/yyyy/mm/dd/filename
    Auth-Token: Basic <login-token>
    RETURN md5 string of the file

May be used later for checking what pictures have been shared between
machines.
