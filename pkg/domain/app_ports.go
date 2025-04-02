package domain

import (
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
)

type PortPublicationProtocol string

const (
	PortPublicationProtocolTCP PortPublicationProtocol = "tcp"
	PortPublicationProtocolUDP PortPublicationProtocol = "udp"
)

type AvailablePort struct {
	StartPort int
	EndPort   int
	Protocol  PortPublicationProtocol
}

type AvailablePortSlice []*AvailablePort

func isValidPort(port int) error {
	if port < 0 || 65535 < port {
		return errors.New("invalid port (needs to be within 0 to 65535)")
	}
	return nil
}

func (ap *AvailablePort) Validate() error {
	if err := isValidPort(ap.StartPort); err != nil {
		return errors.Wrap(err, "invalid start port")
	}
	if err := isValidPort(ap.EndPort); err != nil {
		return errors.Wrap(err, "invalid end port")
	}
	if ap.EndPort < ap.StartPort {
		return errors.New("end port comes before start port")
	}
	return nil
}

func (s AvailablePortSlice) IsAvailable(port int, protocol PortPublicationProtocol) bool {
	return lo.ContainsBy(s, func(ap *AvailablePort) bool {
		return ap.Protocol == protocol && ap.StartPort <= port && port <= ap.EndPort
	})
}

type PortPublication struct {
	InternetPort    int
	ApplicationPort int
	Protocol        PortPublicationProtocol
}

func (p *PortPublication) Validate() error {
	if err := isValidPort(p.InternetPort); err != nil {
		return errors.Wrap(err, "invalid internet port")
	}
	if err := isValidPort(p.ApplicationPort); err != nil {
		return errors.Wrap(err, "invalid application port")
	}
	return nil
}

func (p *PortPublication) ConflictsWith(existing []*Application) bool {
	return lo.ContainsBy(existing, func(app *Application) bool {
		return lo.ContainsBy(app.PortPublications, func(used *PortPublication) bool {
			return used.InternetPort == p.InternetPort && used.Protocol == p.Protocol
		})
	})
}

func (p *PortPublication) Compare(other *PortPublication) bool {
	if p.Protocol != other.Protocol {
		return strings.Compare(string(p.Protocol), string(other.Protocol)) < 0
	}
	return p.InternetPort < other.InternetPort
}
