package main

import (
	_ "github.com/AliFnieer/morph/internal/converters/archive"
	archiveconv "github.com/AliFnieer/morph/internal/converters/archive"
	_ "github.com/AliFnieer/morph/internal/converters/audio"
	audioconv "github.com/AliFnieer/morph/internal/converters/audio"
	_ "github.com/AliFnieer/morph/internal/converters/data"
	dataconv "github.com/AliFnieer/morph/internal/converters/data"
	_ "github.com/AliFnieer/morph/internal/converters/epub"
	epubconv "github.com/AliFnieer/morph/internal/converters/epub"
	_ "github.com/AliFnieer/morph/internal/converters/html"
	htmlconv "github.com/AliFnieer/morph/internal/converters/html"
	_ "github.com/AliFnieer/morph/internal/converters/image"
	_ "github.com/AliFnieer/morph/internal/converters/odf"
	odfconv "github.com/AliFnieer/morph/internal/converters/odf"
	_ "github.com/AliFnieer/morph/internal/converters/office"
	officeconv "github.com/AliFnieer/morph/internal/converters/office"
	_ "github.com/AliFnieer/morph/internal/converters/ooxml"
	ooxmlconv "github.com/AliFnieer/morph/internal/converters/ooxml"
	pdfconv "github.com/AliFnieer/morph/internal/converters/pdf"
	_ "github.com/AliFnieer/morph/internal/converters/raw"
	rawconv "github.com/AliFnieer/morph/internal/converters/raw"
	_ "github.com/AliFnieer/morph/internal/converters/rtf"
	rtfconv "github.com/AliFnieer/morph/internal/converters/rtf"
	textconv "github.com/AliFnieer/morph/internal/converters/text"
	_ "github.com/AliFnieer/morph/internal/converters/vector"
	vectorconv "github.com/AliFnieer/morph/internal/converters/vector"
	_ "github.com/AliFnieer/morph/internal/converters/video"
	videoconv "github.com/AliFnieer/morph/internal/converters/video"
)

func init() {
	_ = textconv.Conv{}
	_ = pdfconv.Conv{}
	_ = htmlconv.Conv{}
	_ = archiveconv.Conv{}
	_ = officeconv.Conv{}
	_ = rawconv.Conv{}
	_ = dataconv.Conv{}
	_ = ooxmlconv.Conv{}
	_ = odfconv.Conv{}
	_ = epubconv.Conv{}
	_ = rtfconv.Conv{}
	_ = vectorconv.Conv{}
	_ = audioconv.Conv{}
	_ = videoconv.Conv{}
}
