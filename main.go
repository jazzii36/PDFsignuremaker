package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/font"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"os"
	"path/filepath"
)

func main() {

	files := font.FindFontFile("ZHS")
	if len(files) > 0 {
		os.Setenv("FYNE_FONT", files[0])
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("PDF Signature Maker")

	statusLabel := widget.NewLabel("")
	statusLabel.Hide() // 默认隐藏

	pdfEntry := widget.NewEntry()
	pdfEntry.SetPlaceHolder("Enter or choose PDF path...")

	pngEntry := widget.NewEntry()
	pngEntry.SetPlaceHolder("Enter or choose PNG/jpg/jpeg path...")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Enter parameter")
	descEntry.SetText("scale:0.15, pos:br, rotation:0, off:0 0") // 设置默认值

	submitButton := widget.NewButton("Submit", func() {
		statusLabel.Hide() // 每次提交时隐藏状态标签
		pdfPath := pdfEntry.Text
		pngPath := pngEntry.Text
		desc := descEntry.Text

		// 添加水印逻辑
		wm, err := pdfcpu.ParseImageWatermarkDetails(pngPath, desc, false, 100)
		if err != nil {
			statusLabel.SetText("Error creating watermark: " + err.Error())
			statusLabel.Show()
			return
		}

		_, filename := filepath.Split(pdfPath)
		//outPath := filepath.Join(dir, file)

		// 获取当前可执行文件的路径
		//executablePath, _ := os.Executable()
		//executableDir := filepath.Dir(executablePath)
		//outPath = filepath.Join(executableDir, outPath)

		err = api.AddWatermarksFile(pdfPath, filename, nil, wm, nil)
		if err != nil {
			statusLabel.SetText("Error applying watermark: " + err.Error())
			statusLabel.Show()
			return
		}
		statusLabel.SetText("Watermark applied successfully")
		statusLabel.Show()
	})

	// 文件选择器按钮
	pdfOpenButton := widget.NewButton("选择pdf", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err == nil && r != nil {
				pdfEntry.SetText(r.URI().Path())
			}
		}, myWindow)
	})

	pngOpenButton := widget.NewButton("选择图片", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err == nil && r != nil {
				pngEntry.SetText(r.URI().Path())
			}
		}, myWindow)
	})

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetContent(container.NewVBox(
		container.NewGridWithColumns(2,
			pdfEntry,
			pdfOpenButton,
			pngEntry,
			pngOpenButton,
		),
		descEntry,
		submitButton,
		statusLabel,
	))

	myWindow.ShowAndRun()
}
