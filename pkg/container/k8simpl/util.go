package k8simpl

import (
	"fmt"
	networkingv1 "k8s.io/api/networking/v1"
)

func int32Ptr(i int32) *int32                                           { return &i }
func pathTypePtr(pathType networkingv1.PathType) *networkingv1.PathType { return &pathType }

func deploymentName(appID, envID string) string {
	return fmt.Sprintf("nsapp-%s-%s", appID, envID)
}
