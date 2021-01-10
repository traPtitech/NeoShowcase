package k8simpl

import "fmt"

func int32Ptr(i int32) *int32 { return &i }

func deploymentName(appID string) string {
	return fmt.Sprintf("nsapp-%s", appID)
}
