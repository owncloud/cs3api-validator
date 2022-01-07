package featurecontext

import "fmt"

func (f *FeatureContext) GetPublicShareToken(publicShare string) (string, error) {
	token, ok := f.PublicSharesToken[publicShare]
	if !ok {
		return "", fmt.Errorf("publicshare \"%s\" is not known", publicShare)
	}
	return token, nil
}
