// Code generated by clubbygen.
// GENERATED FILE DO NOT EDIT
// +build clubby_strict

package swupdate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cesanta.com/clubby"
	"cesanta.com/clubby/endpoint"
	"cesanta.com/clubby/frame"
	"cesanta.com/common/go/ourjson"
	"cesanta.com/common/go/ourtrace"
	"github.com/cesanta/errors"
	"golang.org/x/net/trace"

	"github.com/cesanta/ucl"
	"github.com/cesanta/validate-json/schema"
	"github.com/golang/glog"
)

var _ = bytes.MinRead
var _ = fmt.Errorf
var emptyMessage = ourjson.RawMessage{}
var _ = ourtrace.New
var _ = trace.New

const ServiceID = "http://cesanta.com/clubby/service/v1/SWUpdate"

type ListSectionsResult struct {
	Section  *string `json:"section,omitempty"`
	Version  *string `json:"version,omitempty"`
	Writable *bool   `json:"writable,omitempty"`
}

type UpdateArgsSig struct {
	Alg *string `json:"alg,omitempty"`
	V   *string `json:"v,omitempty"`
}

type UpdateArgs struct {
	Blob      *string        `json:"blob,omitempty"`
	Blob_type *string        `json:"blob_type,omitempty"`
	Blob_url  *string        `json:"blob_url,omitempty"`
	Section   *string        `json:"section,omitempty"`
	Sig       *UpdateArgsSig `json:"sig,omitempty"`
	Version   *string        `json:"version,omitempty"`
}

type Service interface {
	ListSections(ctx context.Context) ([]ListSectionsResult, error)
	Update(ctx context.Context, args *UpdateArgs) error
}

type Instance interface {
	Call(context.Context, string, *frame.Command) (*frame.Response, error)
	TraceCall(context.Context, string, *frame.Command) (context.Context, trace.Trace, func(*error))
}

type _validators struct {
	// This comment prevents gofmt from aligning types in the struct.
	ListSectionsResult *schema.Validator
	// This comment prevents gofmt from aligning types in the struct.
	UpdateArgs *schema.Validator
}

var (
	validators     *_validators
	validatorsOnce sync.Once
)

func initValidators() {
	validators = &_validators{}

	loader := schema.NewLoader()

	service, err := ucl.Parse(bytes.NewBuffer(_ServiceDefinition))
	if err != nil {
		panic(err)
	}
	// Patch up shortcuts to be proper schemas.
	for _, v := range service.(*ucl.Object).Find("methods").(*ucl.Object).Value {
		if s, ok := v.(*ucl.Object).Find("result").(*ucl.String); ok {
			for kk := range v.(*ucl.Object).Value {
				if kk.Value == "result" {
					v.(*ucl.Object).Value[kk] = &ucl.Object{
						Value: map[ucl.Key]ucl.Value{
							ucl.Key{Value: "type"}: s,
						},
					}
				}
			}
		}
		if v.(*ucl.Object).Find("args") == nil {
			continue
		}
		args := v.(*ucl.Object).Find("args").(*ucl.Object)
		for kk, vv := range args.Value {
			if s, ok := vv.(*ucl.String); ok {
				args.Value[kk] = &ucl.Object{
					Value: map[ucl.Key]ucl.Value{
						ucl.Key{Value: "type"}: s,
					},
				}
			}
		}
	}
	var s *ucl.Object
	_ = s // avoid unused var error
	validators.ListSectionsResult, err = schema.NewValidator(service.(*ucl.Object).Find("methods").(*ucl.Object).Find("ListSections").(*ucl.Object).Find("result"), loader)
	if err != nil {
		panic(err)
	}
	s = &ucl.Object{
		Value: map[ucl.Key]ucl.Value{
			ucl.Key{Value: "properties"}: service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Update").(*ucl.Object).Find("args"),
			ucl.Key{Value: "type"}:       &ucl.String{Value: "object"},
		},
	}
	if req, found := service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Update").(*ucl.Object).Lookup("required_args"); found {
		s.Value[ucl.Key{Value: "required"}] = req
	}
	validators.UpdateArgs, err = schema.NewValidator(s, loader)
	if err != nil {
		panic(err)
	}
}

func NewClient(i Instance, addr string) Service {
	validatorsOnce.Do(initValidators)
	return &_Client{i: i, addr: addr}
}

type _Client struct {
	i    Instance
	addr string
}

func (c *_Client) ListSections(pctx context.Context) (res []ListSectionsResult, err error) {
	cmd := &frame.Command{
		Cmd: "/v1/SWUpdate.ListSections",
	}
	ctx, tr, finish := c.i.TraceCall(pctx, c.addr, cmd)
	defer finish(&err)
	_ = tr
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Status != 0 {
		return nil, errors.Trace(&endpoint.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}

	tr.LazyPrintf("res: %s", ourjson.LazyJSON(&resp))

	bb, err := resp.Response.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal result as JSON: %+v", err)
	} else {
		rv, err := ucl.Parse(bytes.NewReader(bb))
		if err == nil {
			if err := validators.ListSectionsResult.Validate(rv); err != nil {
				glog.Warningf("Got invalid result for ListSections: %+v", err)
				return nil, errors.Annotatef(err, "invalid response for ListSections")
			}
		}
	}
	var r []ListSectionsResult
	err = resp.Response.UnmarshalInto(&r)
	if err != nil {
		return nil, errors.Annotatef(err, "unmarshaling response")
	}
	return r, nil
}

func (c *_Client) Update(pctx context.Context, args *UpdateArgs) (err error) {
	cmd := &frame.Command{
		Cmd: "/v1/SWUpdate.Update",
	}
	ctx, tr, finish := c.i.TraceCall(pctx, c.addr, cmd)
	defer finish(&err)
	_ = tr

	tr.LazyPrintf("args: %s", ourjson.LazyJSON(&args))
	cmd.Args = ourjson.DelayMarshaling(args)
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		v, err := ucl.Parse(bytes.NewReader(b))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.UpdateArgs.Validate(v); err != nil {
				glog.Warningf("Sending invalid args for Update: %+v", err)
				return errors.Annotatef(err, "invalid args for Update")
			}
		}
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&endpoint.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

func RegisterService(i *clubby.Instance, impl Service) error {
	validatorsOnce.Do(initValidators)
	s := &_Server{impl}
	i.RegisterCommandHandler("/v1/SWUpdate.ListSections", s.ListSections)
	i.RegisterCommandHandler("/v1/SWUpdate.Update", s.Update)
	i.RegisterService(ServiceID, _ServiceDefinition)
	return nil
}

type _Server struct {
	impl Service
}

func (s *_Server) ListSections(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	r, err := s.impl.ListSections(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	bb, err := json.Marshal(r)
	if err == nil {
		v, err := ucl.Parse(bytes.NewBuffer(bb))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(bb), err)
		} else {
			if err := validators.ListSectionsResult.Validate(v); err != nil {
				glog.Warningf("Returned invalid response for ListSections: %+v", err)
				return nil, errors.Annotatef(err, "server generated invalid responce for ListSections")
			}
		}
	}
	return r, nil
}

func (s *_Server) Update(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		if v, err := ucl.Parse(bytes.NewReader(b)); err != nil {
			glog.Errorf("Failed to parse valid JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.UpdateArgs.Validate(v); err != nil {
				glog.Warningf("Got invalid args for Update: %+v", err)
				return nil, errors.Annotatef(err, "invalid args for Update")
			}
		}
	}
	var args UpdateArgs
	if len(cmd.Args) > 0 {
		if err := cmd.Args.UnmarshalInto(&args); err != nil {
			return nil, errors.Annotatef(err, "unmarshaling args")
		}
	}
	return nil, s.impl.Update(ctx, &args)
}

var _ServiceDefinition = json.RawMessage([]byte(`{
  "doc": "SWUpdate service provides a way to update device's software.",
  "methods": {
    "ListSections": {
      "doc": "Returns a list of components of the device's software. Each section is updated individually.",
      "result": {
        "items": {
          "properties": {
            "section": {
              "type": "string"
            },
            "version": {
              "type": "string"
            },
            "writable": {
              "type": "boolean"
            }
          },
          "type": "object"
        },
        "type": "array"
      }
    },
    "Update": {
      "args": {
        "blob": {
          "doc": "Image as a string, if appropriate.",
          "type": "string"
        },
        "blob_type": {
          "doc": "type of the blob. Valid values: manifest, zip",
          "type": "string"
        },
        "blob_url": {
          "doc": "URL pointing to the image if it's too big to fit in the ` + "`" + `blob` + "`" + `.",
          "type": "string"
        },
        "section": {
          "doc": "Name of the section to update.",
          "type": "string"
        },
        "sig": {
          "doc": "Hash or signature for the image that can be used to verify its integrity.",
          "properties": {
            "alg": {
              "type": "string"
            },
            "v": {
              "type": "string"
            }
          },
          "required": [
            "alg",
            "v"
          ],
          "type": "object"
        },
        "version": {
          "doc": "Optional version of the new image.",
          "type": "string"
        }
      },
      "doc": "Instructs the device to update a given section."
    }
  },
  "name": "/v1/SWUpdate",
  "namespace": "http://cesanta.com/clubby/service",
  "visibility": "private"
}`))
