package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("PDF signature maker")

	statusLabel := widget.NewLabel("")
	statusLabel.Hide() // 默认隐藏

	pdfEntry := widget.NewEntry()
	pdfEntry.SetPlaceHolder("Enter PDF file path")

	pngEntry := widget.NewEntry()
	pngEntry.SetPlaceHolder("Enter PNG file path")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Enter parameter")
	descEntry.SetText("scale:0.2, pos:br, rotation:0, off:0 0") // 设置默认值

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

		//selectedPages := []string{"1", "2"}
		err = api.AddWatermarksFile(pdfPath, "out.pdf", nil, wm, nil)
		if err != nil {
			statusLabel.SetText("Error applying watermark: " + err.Error())
			statusLabel.Show()
			return
		}
		statusLabel.SetText("Watermark applied successfully")
		statusLabel.Show()
	})

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.SetContent(container.NewVBox(
		pdfEntry,
		pngEntry,
		descEntry,
		submitButton,
		statusLabel,
	))

	myWindow.ShowAndRun()
}
