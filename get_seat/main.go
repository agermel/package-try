package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"get_seat/tools" // 假设工具包放在 get_seat/tools 目录下

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 创建 Fyne 应用
	a := app.New()
	w := a.NewWindow("抢位系统 by 长期素食")

	// 定义 UI 组件
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	errorLabel := widget.NewLabel("") // 创建错误提示标签

	loginButton := widget.NewButton("登录", func() {
		// 获取用户名和密码
		username := usernameEntry.Text
		password := passwordEntry.Text

		// 调用工具包中的 Login 函数进行登录
		client, err := tools.Login(username, password)

		if err != nil {
			logOutput := fmt.Sprintf("登录失败: %v", err)
			log.Println(logOutput)
			// 显示错误信息并清空输入框
			errorLabel.SetText("登录失败: 您输入的用户名或密码有误")
			usernameEntry.SetText("") // 清空用户名
			passwordEntry.SetText("") // 清空密码
			return
		}

		// 登录成功后，显示预约界面
		showReservePage(w, client, a, errorLabel)
	})

	// 设置主界面布局
	w.SetContent(container.NewVBox(
		widget.NewLabel("用户名:"),
		usernameEntry,
		widget.NewLabel("密码:"),
		passwordEntry,
		loginButton,
		errorLabel, // 显示错误信息的标签
	))

	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

// 预定页面显示
func showReservePage(w fyne.Window, client *http.Client, a fyne.App, errorLabel *widget.Label) {
	// 创建一个新的页面来显示预定信息

	startTimeEntry := widget.NewEntry()
	startTimeEntry.SetPlaceHolder("请输入开始时间（例如：09:00）") //半透明提示

	endTimeEntry := widget.NewEntry()
	endTimeEntry.SetPlaceHolder("请输入结束时间（例如：12:00）")

	reserveButton := widget.NewButton("预定", func() {
		// 获取时间输入
		startTime := startTimeEntry.Text
		endTime := endTimeEntry.Text

		var reserveNeed tools.Reserve_need
		var err error
		n := 0
		// 获取可用座位
		for {
			n++
			reserveNeed, err = tools.All_rooms_reserve_find(client, startTime, endTime)
			if err != nil {
				errorLabel.SetText(err.Error())
				break
			}
			if reserveNeed.DevID != "" {
				break
			}
			text := fmt.Sprintf("第%d次抢位。。。", n)
			errorLabel.SetText(text)
			time.Sleep(time.Second * 10)
		}

		// 发起预约请求
		if reserveNeed.DevID != "" {
			tools.Reserve(reserveNeed)
			showReservationSuccessPage(w, reserveNeed, client, a, errorLabel)
		}
	})

	// 设置预约页面内容
	w.SetContent(container.NewVBox(
		widget.NewLabel("开始时间:"),
		startTimeEntry,
		widget.NewLabel("结束时间:"),
		endTimeEntry,
		reserveButton,
		errorLabel, // 显示错误信息的标签
	))
}

// 预约成功后的页面显示
func showReservationSuccessPage(w fyne.Window, reserveNeed tools.Reserve_need, client *http.Client, a fyne.App, errorLabel *widget.Label) {
	// 获取最近的预约历史信息
	rsvID, seat, err := tools.Search_reserve_history(client)
	if err != nil {
		errorLabel.SetText(err.Error())
		// 更新 UI，显示错误信息
		return
	}

	// 显示预约成功后的界面
	w.SetContent(container.NewVBox(
		widget.NewLabel("预约成功!"),
		widget.NewLabel(fmt.Sprintf("座位号: %s", seat)),
		widget.NewButton("确认", func() {
			// 点击确认后关闭程序
			log.Println("预约已确认，程序将于10秒钟后退出")
			time.Sleep(time.Second * 10)
			a.Quit() // 关闭应用
		}),
		widget.NewButton("取消预约", func() {
			// 取消预定功能
			cancelReservation(client, rsvID, w, a, errorLabel)
		}),
	))
}

// 取消预约
func cancelReservation(client *http.Client, rsvID string, w fyne.Window, a fyne.App, errorLabel *widget.Label) {
	// 调用取消预约的方法
	tools.Cancel(client, rsvID)

	// 更新UI，提示取消成功
	logOutput := "预约已取消"
	log.Println(logOutput)
	// 显示取消预约的提示
	errorLabel.SetText(logOutput)

	// 更新 UI
	w.SetContent(container.NewVBox(
		widget.NewLabel(logOutput),
		widget.NewButton("确认", func() {
			// 点击确认后关闭程序
			log.Println("预约已取消，程序将退出")
			a.Quit() // 关闭应用
		}),
	))
}
