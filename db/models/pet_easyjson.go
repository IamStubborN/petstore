// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels(in *jlexer.Lexer, out *PetList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PetList, 0, 8)
			} else {
				*out = PetList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *Pet
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(Pet)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels(out *jwriter.Writer, in PetList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PetList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PetList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PetList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PetList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels(l, v)
}
func easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels1(in *jlexer.Lexer, out *Pet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "category":
			(out.Category).UnmarshalEasyJSON(in)
		case "name":
			out.Name = string(in.String())
		case "photo_urls":
			if in.IsNull() {
				in.Skip()
				out.PhotoURLs = nil
			} else {
				in.Delim('[')
				if out.PhotoURLs == nil {
					if !in.IsDelim(']') {
						out.PhotoURLs = make([]string, 0, 4)
					} else {
						out.PhotoURLs = []string{}
					}
				} else {
					out.PhotoURLs = (out.PhotoURLs)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.PhotoURLs = append(out.PhotoURLs, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]Tag, 0, 2)
					} else {
						out.Tags = []Tag{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Tag
					easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels2(in, &v5)
					out.Tags = append(out.Tags, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "status":
			out.Status = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels1(out *jwriter.Writer, in Pet) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"category\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Category).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	if len(in.PhotoURLs) != 0 {
		const prefix string = ",\"photo_urls\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v6, v7 := range in.PhotoURLs {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	if len(in.Tags) != 0 {
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v8, v9 := range in.Tags {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels2(out, v9)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels1(l, v)
}
func easyjson14a1085DecodeGithubComIamStubborNPetstoreDbModels2(in *jlexer.Lexer, out *Tag) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson14a1085EncodeGithubComIamStubborNPetstoreDbModels2(out *jwriter.Writer, in Tag) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}
