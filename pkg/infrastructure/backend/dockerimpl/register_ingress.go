package dockerimpl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/volatiletech/null/v8"
)

var ingressTmpl = template.Must(template.New("ingress").Parse(ingressTmplString))

const ingressTmplString = `
http:
  routers:
    nsapp-{{.AppID}}-{{.branchID}}:
      service: "nsapp-{{.AppID}}-{{.branchID}}"
      rule: "Host('{{.Host}}')"
  services:
    nsapp-{{.AppID}}-{{.branchID}}:
      loadBalancer:
        servers:
        - url: http://{{.Destination}}:{{.Port}}
`

func (b *dockerBackend) RegisterIngress(ctx context.Context, appID string, branchID string, host string, destination null.String, port null.Int) error {
	conf := filepath.Join(b.ingressConfDir, containerName(appID, branchID)+".yaml")

	data := map[string]interface{}{
		"AppID":       appID,
		"BranchID":    branchID,
		"Host":        host,
		"Destination": containerName(appID, branchID),
		"Port":        80,
	}
	if destination.Valid {
		data["Destination"] = destination.String
	}
	if port.Valid {
		data["Port"] = port.Int
	}

	f, err := os.OpenFile(conf, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	return ingressTmpl.Execute(f, data)
}
