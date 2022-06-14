
# goqsm
Go http client to manage Qsan QSM models.

## Install
```
go get -u github.com/QsanJohnson/goqsm
```

## Usage
```
import "github.com/QsanJohnson/goqsm"

client := goqsm.NewClient("192.168.xxx.xxx")
systemAPI := goqsm.NewSystem(client)

authClient, err := client.GetAuthClient(ctx, "admin", "1234")
volumeAPI := goqsm.NewVolume(authClient)
```


## Testing

You have to create a test.conf file for go test. The following is an example,
```
QSM_IP = 192.168.xxx.xxx
QSM_USERNAME = admin
QSM_PASSWORD = 1234
```

Then execute 'go test'
