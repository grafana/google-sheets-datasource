//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by conversion-gen. DO NOT EDIT.

package v1

import (
	unsafe "unsafe"

	googlesheets "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Datasource)(nil), (*googlesheets.Datasource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Datasource_To_googlesheets_Datasource(a.(*Datasource), b.(*googlesheets.Datasource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*googlesheets.Datasource)(nil), (*Datasource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_googlesheets_Datasource_To_v1_Datasource(a.(*googlesheets.Datasource), b.(*Datasource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DatasourceList)(nil), (*googlesheets.DatasourceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DatasourceList_To_googlesheets_DatasourceList(a.(*DatasourceList), b.(*googlesheets.DatasourceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*googlesheets.DatasourceList)(nil), (*DatasourceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_googlesheets_DatasourceList_To_v1_DatasourceList(a.(*googlesheets.DatasourceList), b.(*DatasourceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DatasourceSpec)(nil), (*googlesheets.DatasourceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec(a.(*DatasourceSpec), b.(*googlesheets.DatasourceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*googlesheets.DatasourceSpec)(nil), (*DatasourceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec(a.(*googlesheets.DatasourceSpec), b.(*DatasourceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DatasourceStatus)(nil), (*googlesheets.DatasourceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus(a.(*DatasourceStatus), b.(*googlesheets.DatasourceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*googlesheets.DatasourceStatus)(nil), (*DatasourceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus(a.(*googlesheets.DatasourceStatus), b.(*DatasourceStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_Datasource_To_googlesheets_Datasource(in *Datasource, out *googlesheets.Datasource, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_Datasource_To_googlesheets_Datasource is an autogenerated conversion function.
func Convert_v1_Datasource_To_googlesheets_Datasource(in *Datasource, out *googlesheets.Datasource, s conversion.Scope) error {
	return autoConvert_v1_Datasource_To_googlesheets_Datasource(in, out, s)
}

func autoConvert_googlesheets_Datasource_To_v1_Datasource(in *googlesheets.Datasource, out *Datasource, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_googlesheets_Datasource_To_v1_Datasource is an autogenerated conversion function.
func Convert_googlesheets_Datasource_To_v1_Datasource(in *googlesheets.Datasource, out *Datasource, s conversion.Scope) error {
	return autoConvert_googlesheets_Datasource_To_v1_Datasource(in, out, s)
}

func autoConvert_v1_DatasourceList_To_googlesheets_DatasourceList(in *DatasourceList, out *googlesheets.DatasourceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]googlesheets.Datasource)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1_DatasourceList_To_googlesheets_DatasourceList is an autogenerated conversion function.
func Convert_v1_DatasourceList_To_googlesheets_DatasourceList(in *DatasourceList, out *googlesheets.DatasourceList, s conversion.Scope) error {
	return autoConvert_v1_DatasourceList_To_googlesheets_DatasourceList(in, out, s)
}

func autoConvert_googlesheets_DatasourceList_To_v1_DatasourceList(in *googlesheets.DatasourceList, out *DatasourceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Datasource)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_googlesheets_DatasourceList_To_v1_DatasourceList is an autogenerated conversion function.
func Convert_googlesheets_DatasourceList_To_v1_DatasourceList(in *googlesheets.DatasourceList, out *DatasourceList, s conversion.Scope) error {
	return autoConvert_googlesheets_DatasourceList_To_v1_DatasourceList(in, out, s)
}

func autoConvert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec(in *DatasourceSpec, out *googlesheets.DatasourceSpec, s conversion.Scope) error {
	out.AuthType = in.AuthType
	out.APIKey = in.APIKey
	out.DefaultProject = in.DefaultProject
	out.JWT = in.JWT
	out.ClientEmail = in.ClientEmail
	out.TokenURI = in.TokenURI
	out.AuthenticationType = in.AuthenticationType
	out.PrivateKeyPath = in.PrivateKeyPath
	out.PrivateKey = in.PrivateKey
	return nil
}

// Convert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec is an autogenerated conversion function.
func Convert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec(in *DatasourceSpec, out *googlesheets.DatasourceSpec, s conversion.Scope) error {
	return autoConvert_v1_DatasourceSpec_To_googlesheets_DatasourceSpec(in, out, s)
}

func autoConvert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec(in *googlesheets.DatasourceSpec, out *DatasourceSpec, s conversion.Scope) error {
	out.AuthType = in.AuthType
	out.APIKey = in.APIKey
	out.DefaultProject = in.DefaultProject
	out.JWT = in.JWT
	out.ClientEmail = in.ClientEmail
	out.TokenURI = in.TokenURI
	out.AuthenticationType = in.AuthenticationType
	out.PrivateKeyPath = in.PrivateKeyPath
	out.PrivateKey = in.PrivateKey
	return nil
}

// Convert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec is an autogenerated conversion function.
func Convert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec(in *googlesheets.DatasourceSpec, out *DatasourceSpec, s conversion.Scope) error {
	return autoConvert_googlesheets_DatasourceSpec_To_v1_DatasourceSpec(in, out, s)
}

func autoConvert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus(in *DatasourceStatus, out *googlesheets.DatasourceStatus, s conversion.Scope) error {
	return nil
}

// Convert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus is an autogenerated conversion function.
func Convert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus(in *DatasourceStatus, out *googlesheets.DatasourceStatus, s conversion.Scope) error {
	return autoConvert_v1_DatasourceStatus_To_googlesheets_DatasourceStatus(in, out, s)
}

func autoConvert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus(in *googlesheets.DatasourceStatus, out *DatasourceStatus, s conversion.Scope) error {
	return nil
}

// Convert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus is an autogenerated conversion function.
func Convert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus(in *googlesheets.DatasourceStatus, out *DatasourceStatus, s conversion.Scope) error {
	return autoConvert_googlesheets_DatasourceStatus_To_v1_DatasourceStatus(in, out, s)
}