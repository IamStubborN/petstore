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

func easyjson120d1ca2DecodeGithubComIamStubborNPetstoreDbModels(in *jlexer.Lexer, out *Order) {
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
		case "pet_id":
			out.PetID = int64(in.Int64())
		case "user_id":
			out.UserID = int64(in.Int64())
		case "quantity":
			out.Quantity = int32(in.Int32())
		case "ship_date":
			out.ShipDate = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "complete":
			out.Complete = bool(in.Bool())
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
func easyjson120d1ca2EncodeGithubComIamStubborNPetstoreDbModels(out *jwriter.Writer, in Order) {
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
		const prefix string = ",\"pet_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.PetID))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.UserID))
	}
	{
		const prefix string = ",\"quantity\":"
		out.RawString(prefix)
		out.Int32(int32(in.Quantity))
	}
	{
		const prefix string = ",\"ship_date\":"
		out.RawString(prefix)
		out.String(string(in.ShipDate))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"complete\":"
		out.RawString(prefix)
		out.Bool(bool(in.Complete))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Order) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson120d1ca2EncodeGithubComIamStubborNPetstoreDbModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Order) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson120d1ca2EncodeGithubComIamStubborNPetstoreDbModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Order) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson120d1ca2DecodeGithubComIamStubborNPetstoreDbModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Order) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson120d1ca2DecodeGithubComIamStubborNPetstoreDbModels(l, v)
}
