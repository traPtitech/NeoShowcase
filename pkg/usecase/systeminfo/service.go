package systeminfo

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

type ServiceConfig struct {
	AdditionalLinks []*domain.AdditionalLink
}

type Service interface {
	GetSystemInfo() (*domain.SystemInfo, error)
}

type service struct {
	c *ServiceConfig

	backend domain.Backend
	appRepo domain.ApplicationRepository
	sshConf domain.SSHConfig
	pubKey  *ssh.PublicKeys
}

func NewService(
	c *ServiceConfig,
	backend domain.Backend,
	appRepo domain.ApplicationRepository,
	sshConf domain.SSHConfig,
	pubKey *ssh.PublicKeys,
) Service {
	return &service{
		c:       c,
		backend: backend,
		appRepo: appRepo,
		sshConf: sshConf,
		pubKey:  pubKey,
	}
}

func (s *service) GetSystemInfo() (*domain.SystemInfo, error) {
	domains := s.backend.AvailableDomains()
	existingApps, err := s.appRepo.GetApplications(context.Background(), domain.GetApplicationCondition{})
	if err != nil {
		return nil, err
	}
	for _, ad := range domains {
		ad.AlreadyBound = ad.IsAlreadyBound(existingApps)
	}

	ports := s.backend.AvailablePorts()
	ver, rev := cli.GetVersion()

	return &domain.SystemInfo{
		PublicKey: domain.Base64EncodedPublicKey(s.pubKey.Signer.PublicKey()) + " neoshowcase",
		SSHInfo: struct {
			Host string
			Port int
		}{
			Host: s.sshConf.Host,
			Port: s.sshConf.Port,
		},
		AvailableDomains: domains,
		AvailablePorts:   ports,
		AdditionalLinks:  s.c.AdditionalLinks,
		Version:          ver,
		Revision:         rev,
	}, nil
}
