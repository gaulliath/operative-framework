package pictures

import (
	"encoding/json"
	"os"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

type PictureExifModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

type Exif struct {
	BitsPerSample                    []int    `json:"BitsPerSample"`
	ColorSpace                       []int    `json:"ColorSpace"`
	DateTime                         string   `json:"DateTime"`
	ExifIFDPointer                   []int    `json:"ExifIFDPointer"`
	ExifVersion                      string   `json:"ExifVersion"`
	ImageLength                      []int    `json:"ImageLength"`
	ImageWidth                       []int    `json:"ImageWidth"`
	Orientation                      []int    `json:"Orientation"`
	PhotometricInterpretation        []int    `json:"PhotometricInterpretation"`
	PixelXDimension                  []int    `json:"PixelXDimension"`
	PixelYDimension                  []int    `json:"PixelYDimension"`
	ResolutionUnit                   []int    `json:"ResolutionUnit"`
	SamplesPerPixel                  []int    `json:"SamplesPerPixel"`
	Software                         string   `json:"Software"`
	ThumbJPEGInterchangeFormat       []int    `json:"ThumbJPEGInterchangeFormat"`
	ThumbJPEGInterchangeFormatLength []int    `json:"ThumbJPEGInterchangeFormatLength"`
	XResolution                      []string `json:"XResolution"`
	YResolution                      []string `json:"YResolution"`
}

func PushPictureExifModule(s *session.Session) *PictureExifModule {
	mod := PictureExifModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Target file", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *PictureExifModule) Name() string {
	return "picture.exif"
}

func (module *PictureExifModule) Description() string {
	return "View exif data on selected picture"
}

func (module *PictureExifModule) Author() string {
	return "Tristan Granier"
}

func (module *PictureExifModule) GetType() string {
	return "file"
}

func (module *PictureExifModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *PictureExifModule) Start() {
	paramEnterprise, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramEnterprise.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	fname := target.GetName()
	separator := target.GetSeparator()

	f, err := os.Open(fname)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	s, _ := x.MarshalJSON()
	var e Exif
	err = json.Unmarshal(s, &e)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)

	t.AppendRow(table.Row{
		"ThumbJPEGInterchangeFormatLength",
		e.ThumbJPEGInterchangeFormatLength,
	})
	t.AppendRow(table.Row{
		"XResolution",
		e.XResolution,
	})
	t.AppendRow(table.Row{
		"YResolution",
		e.YResolution,
	})
	t.AppendRow(table.Row{
		"ResolutionUnit",
		e.ResolutionUnit,
	})
	t.AppendRow(table.Row{
		"ExifVersion",
		e.ExifVersion,
	})
	t.AppendRow(table.Row{
		"ColorSpace",
		e.ColorSpace,
	})
	t.AppendRow(table.Row{
		"PixelXDimension",
		e.PixelXDimension,
	})
	t.AppendRow(table.Row{
		"PixelYDimension",
		e.PixelYDimension,
	})
	t.AppendRow(table.Row{
		"ImageWidth",
		e.ImageWidth,
	})
	t.AppendRow(table.Row{
		"ImageLength:",
		e.ImageLength,
	})
	t.AppendRow(table.Row{
		"PhotometricInterpretation",
		e.PhotometricInterpretation,
	})
	t.AppendRow(table.Row{
		"Software",
		e.Software,
	})
	t.AppendRow(table.Row{
		"DateTime",
		e.DateTime,
	})
	t.AppendRow(table.Row{
		"SamplesPerPixel",
		e.SamplesPerPixel,
	})
	t.AppendRow(table.Row{
		"ExifIFDPointer:",
		e.ExifIFDPointer,
	})
	t.AppendRow(table.Row{
		"ThumbJPEGInterchangeFormat",
		e.ThumbJPEGInterchangeFormat,
	})
	t.AppendRow(table.Row{
		"BitsPerSample:",
		e.BitsPerSample,
	})

	result := session.TargetResults{
		Header: "Date Time" + separator + "Exif Version" + separator + "Software",
		Value:  e.DateTime + separator + e.ExifVersion + separator + e.Software,
	}
	target.Save(module, result)

	module.sess.Stream.Render(t)
}
