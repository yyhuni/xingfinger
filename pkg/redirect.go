package pkg

import "fmt"

type RedirectPolicy string

const (
	RedirectPolicyNever RedirectPolicy = "never"
	RedirectPolicyHTTP  RedirectPolicy = "http"
	RedirectPolicyAll   RedirectPolicy = "all"
)

func ParseRedirectPolicy(raw string) (RedirectPolicy, error) {
	policy := RedirectPolicy(raw)
	switch policy {
	case "", RedirectPolicyNever:
		return RedirectPolicyNever, nil
	case RedirectPolicyHTTP:
		return RedirectPolicyHTTP, nil
	case RedirectPolicyAll:
		return RedirectPolicyAll, nil
	default:
		return "", fmt.Errorf("invalid redirect policy: %s", raw)
	}
}

func (p RedirectPolicy) normalized() RedirectPolicy {
	if p == "" {
		return RedirectPolicyAll
	}
	return p
}

func (p RedirectPolicy) FollowHTTP() bool {
	switch p.normalized() {
	case RedirectPolicyHTTP, RedirectPolicyAll:
		return true
	default:
		return false
	}
}

func (p RedirectPolicy) FollowContent() bool {
	return p.normalized() == RedirectPolicyAll
}
