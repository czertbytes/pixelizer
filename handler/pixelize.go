package handler

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"mime"
	"net/http"
	"strconv"

	"github.com/czertbytes/pixelizer/pixelizer"
)

const maxFileSize = 10 << (10 * 2) // 10 MiB

type pixErr struct {
	Value string `json:"error"`
}

type Pixelize struct {
}

func (h *Pixelize) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, pixErr{"method not allowed"})
		return
	}

	if r.ContentLength > maxFileSize {
		errorResponse(w, http.StatusBadRequest, pixErr{"image is too large"})
		return
	}

	inFile, inHeader, err := r.FormFile("file")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, pixErr{"reading form file failed"})
		return
	}

	var inImage image.Image
	switch inHeader.Header.Get("Content-Type") {
	case "image/png", "image/jpeg":
		inImage, _, err = image.Decode(inFile)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, pixErr{"decoding image format failed"})
			return
		}
	default:
		errorResponse(w, http.StatusBadRequest, pixErr{"unsupported content type"})
		return
	}

	accept := r.Header.Get("Accept")
	outMediaType, _, err := mime.ParseMediaType(accept)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, pixErr{"bad output media type"})
		return
	}

	blockSize := 8
	blockSizeParam := r.URL.Query().Get("block-size")
	if len(blockSizeParam) > 0 {
		v, err := strconv.ParseInt(blockSizeParam, 10, 64)
		if v < 0 || v > 128 || err != nil {
			errorResponse(w, http.StatusBadRequest, pixErr{"block size parameter is not valid"})
			return
		}

		blockSize = int(v)
	}

	outImage := pixelizer.NewSimplePixelizer(blockSize).Pixelize(inImage)

	switch outMediaType {
	case "image/jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, outImage, &jpeg.Options{Quality: 80}); err != nil {
			errorResponse(w, http.StatusInternalServerError, pixErr{"encoding final image failed"})
			return
		}
		return
	case "image/png":
		fallthrough
	default:
		w.Header().Set("Content-Type", "image/png")
		if err := png.Encode(w, outImage); err != nil {
			errorResponse(w, http.StatusInternalServerError, pixErr{"encoding final image failed"})
		}
		return
	}
}

func errorResponse(w http.ResponseWriter, code int, e pixErr) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	b, err := json.Marshal(e)
	if err != nil {
		w.Write([]byte(`"error":"internal server error"`))
		return
	}
	w.Write(b)
}
