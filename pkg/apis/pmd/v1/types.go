package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +gencrd=config
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +genclient:noStatus

// Definition of our CRD BGPAsNumber class
type PMDAsNumber struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               PMDAsNumberSpec   `json:"spec"`
	Status             PMDAsNumberStatus `json:"status,omitempty"`
}

type PMDAsNumberSpec struct {
	AsNumber string `json:"asnumber"`
	Enable   bool   `json:"enable"`
}

type PMDAsNumberStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// k8s List Type
type PMDAsNumberList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []PMDAsNumber `json:"items"`
}

// +gencrd=config
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +genclient:noStatus

// Definition of our CRD BGPRoute class
type PMDRoute struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               PMDRouteSpec   `json:"spec"`
	Status             PMDRouteStatus `json:"status,omitempty"`
}

// +gencrd=state

type PMDRouteSpec struct {
	Prefix  string `json:"prefix"`
	Length  uint32 `json:"length"`
	Counter uint32 `json:"counter"`
}

// +gencrd=state

type PMDRouteStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// k8s List Type
type PMDRouteList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []PMDRoute `json:"items"`
}
