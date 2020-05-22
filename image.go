package opencv_helper

import (
	"bytes"
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"
)

// TemplateMatchMode is the type of the template matching operation.
type TemplateMatchMode int

const (
	// TmSqdiff maps to TM_SQDIFF
	TmSqdiff TemplateMatchMode = iota
	// TmSqdiffNormed maps to TM_SQDIFF_NORMED
	TmSqdiffNormed
	// TmCcorr maps to TM_CCORR
	TmCcorr
	// TmCcorrNormed maps to TM_CCORR_NORMED
	TmCcorrNormed
	// TmCcoeff maps to TM_CCOEFF
	TmCcoeff
	// TmCcoeffNormed maps to TM_CCOEFF_NORMED
	TmCcoeffNormed
)

const DefaultMatchMode = TmCcoeffNormed

var fillColor = color.RGBA{R: 255, G: 255, B: 255, A: 0}

func FindImageLocationFromRaw(source, search *bytes.Buffer, confidence float32, matchMode ...TemplateMatchMode) (loc image.Point, err error) {
	var pathnameImage, pathnameTpl string
	if pathnameImage, pathnameTpl, err = checkAndSave(source, search); err != nil {
		return image.Point{}, err
	}
	return FindImageLocationFromDisk(pathnameImage, pathnameTpl, confidence, matchMode...)
}

func FindImageLocationFromDisk(source, search string, confidence float32, matchMode ...TemplateMatchMode) (loc image.Point, err error) {
	if len(matchMode) == 0 {
		matchMode = []TemplateMatchMode{DefaultMatchMode}
	}
	var matImage, matTpl gocv.Mat
	if matImage, matTpl, err = getMats(source, search, gocv.IMReadGrayScale); err != nil {
		return image.Point{}, err
	}
	defer func() {
		_ = matImage.Close()
		_ = matTpl.Close()
	}()

	return getMatchingLocation(matImage, matTpl, confidence, matchMode[0])
}

func FindAllImageLocationsFromDisk(source, search string, confidence float32, matchMode ...TemplateMatchMode) (locs []image.Point, err error) {
	if len(matchMode) == 0 {
		matchMode = []TemplateMatchMode{DefaultMatchMode}
	}
	var matImage, matTpl gocv.Mat
	if matImage, matTpl, err = getMats(source, search, gocv.IMReadGrayScale); err != nil {
		return nil, err
	}
	defer func() {
		_ = matImage.Close()
		_ = matTpl.Close()
	}()

	var loc image.Point
	if loc, err = getMatchingLocation(matImage, matTpl, confidence, matchMode[0]); err != nil {
		return nil, err
	}
	widthTpl := matTpl.Cols()
	heightTpl := matTpl.Rows()

	locs = make([]image.Point, 0, 9)
	locs = append(locs, loc)

	gocv.FillPoly(&matImage, getPts(loc, widthTpl, heightTpl), fillColor)

	loc, err = getMatchingLocation(matImage, matTpl, confidence, matchMode[0])
	for ; err == nil; loc, err = getMatchingLocation(matImage, matTpl, confidence, matchMode[0]) {
		locs = append(locs, loc)
		gocv.FillPoly(&matImage, getPts(loc, widthTpl, heightTpl), fillColor)
	}

	return locs, nil
}

// getPts 根据图片坐标和宽高，获取填充区域
func getPts(loc image.Point, width, height int) [][]image.Point {
	return [][]image.Point{
		{
			image.Pt(loc.X, loc.Y),
			image.Pt(loc.X, loc.Y+height),
			image.Pt(loc.X+width, loc.Y+height),
			image.Pt(loc.X+width, loc.Y),
		},
	}
}

func FindImageRectFromDisk(source, search string, confidence float32, matchMode ...TemplateMatchMode) (rect image.Rectangle, err error) {
	var matTpl gocv.Mat
	if _, matTpl, err = getMats(source, search, gocv.IMReadGrayScale); err != nil {
		return image.Rectangle{}, err
	}
	defer func() {
		_ = matTpl.Close()
	}()

	var loc image.Point
	if loc, err = FindImageLocationFromDisk(source, search, confidence, matchMode...); err != nil {
		return image.Rectangle{}, err
	}
	rect = image.Rect(loc.X, loc.Y, loc.X+matTpl.Cols(), loc.Y+matTpl.Rows())
	return
}

func FindAllImageRectsFromDisk(source, search string, confidence float32, matchMode ...TemplateMatchMode) (rects []image.Rectangle, err error) {
	var matTpl gocv.Mat
	if _, matTpl, err = getMats(source, search, gocv.IMReadGrayScale); err != nil {
		return nil, err
	}
	defer func() {
		_ = matTpl.Close()
	}()

	var locs []image.Point
	if locs, err = FindAllImageLocationsFromDisk(source, search, confidence, matchMode...); err != nil {
		return nil, err
	}

	rects = make([]image.Rectangle, 0, len(locs))
	for i := range locs {
		r := image.Rect(locs[i].X, locs[i].Y, locs[i].X+matTpl.Cols(), locs[i].Y+matTpl.Rows())
		rects = append(rects, r)
	}
	return
}

// checkAndSave 检查保存路径并保存原图和目标图到该路径下
func checkAndSave(source, search *bytes.Buffer) (pathnameImage, pathnameTpl string, err error) {
	if err = checkStoreDirectory(); err != nil {
		return "", "", err
	}
	pathnameImage = filepath.Join(storeDirectory, GenFilename())
	if err = ioutil.WriteFile(pathnameImage, source.Bytes(), 0666); err != nil {
		return "", "", err
	}
	pathnameTpl = filepath.Join(storeDirectory, GenFilename())
	if err = ioutil.WriteFile(pathnameTpl, search.Bytes(), 0666); err != nil {
		return "", "", err
	}

	return
}

// getMats 从指定路径获取原图和目标图的 `gocv.Mat`
func getMats(nameImage, nameTpl string, flags gocv.IMReadFlag) (matImage, matTpl gocv.Mat, err error) {
	matImage = gocv.IMRead(nameImage, flags)
	if matImage.Empty() {
		return gocv.Mat{}, gocv.Mat{}, fmt.Errorf("invalid read %s", nameImage)
	}
	matTpl = gocv.IMRead(nameTpl, flags)
	if matTpl.Empty() {
		return gocv.Mat{}, gocv.Mat{}, fmt.Errorf("invalid read %s", nameTpl)
	}
	return
}

// getMatchingLocation 获取匹配的图片位置
func getMatchingLocation(matImage gocv.Mat, matTpl gocv.Mat, confidence float32, matchMode TemplateMatchMode) (loc image.Point, err error) {
	if confidence > 1 {
		confidence = 1.0
	}
	// TM_SQDIFF：该方法使用平方差进行匹配，最好匹配为 0。值越大匹配结果越差。
	// TM_SQDIFF_NORMED：该方法使用归一化的平方差进行匹配，最佳匹配也在结果为0处。
	// TmCcoeff 将模版对其均值的相对值与图像对其均值的相关值进行匹配,1表示完美匹配,-1表示糟糕的匹配,0表示没有任何相关性(随机序列)。
	minVal, maxVal, minLoc, maxLoc := getMatchingResult(matImage, matTpl, matchMode)

	// fmt.Println(matchMode[0], "\t", minVal, maxVal, "\t", minLoc, maxLoc)
	// fmt.Printf("%s\t %.10f \t %.10f \t %v \t %v \n", matchMode[0], minVal, maxVal, minLoc, maxLoc)

	var val float32
	val, loc = getValLoc(minVal, maxVal, minLoc, maxLoc, matchMode)

	if val >= confidence {
		return loc, nil
	} else {
		return image.Point{}, errors.New("no such target search image")
	}
}

// getMatchingResult 匹配图片并返回匹配值和位置
func getMatchingResult(matImage gocv.Mat, matTpl gocv.Mat, matchMode TemplateMatchMode) (minVal float32, maxVal float32, minLoc image.Point, maxLoc image.Point) {
	matResult, tmpMask := gocv.NewMat(), gocv.NewMat()
	defer func() {
		_ = matResult.Close()
		_ = tmpMask.Close()
	}()
	gocv.MatchTemplate(matImage, matTpl, &matResult, gocv.TemplateMatchMode(matchMode), tmpMask)
	minVal, maxVal, minLoc, maxLoc = gocv.MinMaxLoc(matResult)
	return
}

// getValLoc 根据不同的匹配模式返回匹配值和位置
func getValLoc(minVal float32, maxVal float32, minLoc image.Point, maxLoc image.Point, matchMode TemplateMatchMode) (val float32, loc image.Point) {
	val, loc = maxVal, maxLoc

	switch matchMode {
	case TmSqdiff, TmSqdiffNormed:
		// 平方差，最佳匹配为 0
		val = minVal
		// minVal = 8
		// for val >= 1 {
		// 	val -= 1
		// }
		if val >= 1 {
			val = float32(math.Mod(float64(val), 1))
		}
		val = 1 - val
		loc = minLoc
	case TmCcoeff:
		// TmCcoeff 将模版对其均值的相对值与图像对其均值的相关值进行匹配,1表示完美匹配,-1表示糟糕的匹配,0表示没有任何相关性(随机序列)。
		// maxVal = 5064792.5000000000
		_, frac := math.Modf(float64(val))
		val = float32(frac)
	case TmCcorr:
		// maxVal = 50553512.0000000000
		_, frac := math.Modf(float64(val))
		val = float32(frac)
	}
	// fmt.Println("匹配度", val)
	return
}
