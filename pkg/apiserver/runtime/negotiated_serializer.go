// SPDX-License-Identifier: AGPL-3.0-only

package runtime

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
)

var _ runtime.NegotiatedSerializer = (*NegotiatedSerializer)(nil)

type NegotiatedSerializer struct {
	typer     runtime.ObjectTyper
	creator   runtime.ObjectCreater
	convertor runtime.ObjectConvertor
	scheme    *runtime.Scheme

	unstructuredScheme runtime.Codec
}

func NewNegotiatedSerializer(scheme *runtime.Scheme) *NegotiatedSerializer {
	return &NegotiatedSerializer{
		typer:     &ObjectTyper{},
		creator:   &objectCreator{},
		convertor: &objectConvertor{},

		scheme:             scheme,
		unstructuredScheme: unstructured.UnstructuredJSONScheme,
	}
}

func (n NegotiatedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{
		{
			MediaType:        "application/json",
			MediaTypeType:    "application",
			MediaTypeSubType: "json",
			EncodesAsText:    true,
			Serializer:       json.NewSerializer(json.DefaultMetaFactory, n.creator, n.typer, false),
			PrettySerializer: json.NewSerializer(json.DefaultMetaFactory, n.creator, n.typer, true),
			StrictSerializer: json.NewSerializerWithOptions(json.DefaultMetaFactory, n.creator, n.typer, json.SerializerOptions{
				Strict: true,
			}),
			StreamSerializer: &runtime.StreamSerializerInfo{
				EncodesAsText: true,
				Serializer:    json.NewSerializer(json.DefaultMetaFactory, n.creator, n.typer, false),
				Framer:        json.Framer,
			},
		},
		{
			MediaType:        "application/yaml",
			MediaTypeType:    "application",
			MediaTypeSubType: "yaml",
			EncodesAsText:    true,
			Serializer:       json.NewYAMLSerializer(json.DefaultMetaFactory, n.creator, n.typer),
			StrictSerializer: json.NewSerializerWithOptions(json.DefaultMetaFactory, n.creator, n.typer, json.SerializerOptions{
				Yaml:   true,
				Strict: true,
			}),
		},
	}
}

func (n NegotiatedSerializer) EncoderForVersion(serializer runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return versioning.NewCodec(n.unstructuredScheme, nil, n.convertor, n.scheme, n.scheme, n.scheme, gv, nil, "grdNegotiatedSerializer")
}

func (n NegotiatedSerializer) DecoderToVersion(serializer runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return versioning.NewCodec(nil, n.unstructuredScheme, runtime.UnsafeObjectConvertor(n.scheme), n.scheme, n.scheme, n.scheme, gv, nil, "grdNegotiatedSerializer")
}
