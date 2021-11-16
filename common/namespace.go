package common

import (
	v1 "k8s.io/api/core/v1"
)

type NamespaceQuery struct {
	namespaces []string
}

func NewNamespaceQuery(namespace []string) *NamespaceQuery {
	return &NamespaceQuery{
		namespaces: namespace,
	}
}

func NewSingleNamespaceQuery(namespace string) *NamespaceQuery {
	return &NamespaceQuery{[]string{namespace}}
}

func (n NamespaceQuery) ToRequestParam() string {
	if len(n.namespaces) > 0 {
		return n.namespaces[0]
	}
	return v1.NamespaceAll
}

func (n NamespaceQuery) Matches(namespace string) bool {
	for _, queryNamespace := range n.namespaces {
		if queryNamespace == namespace {
			return true
		}
	}
	return false
}
