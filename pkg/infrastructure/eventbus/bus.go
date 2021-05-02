package eventbus

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE

type Fields map[string]interface{}

type Bus interface {
	Publish(eventType string, body Fields)
}
