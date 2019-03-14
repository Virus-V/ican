package main

import (
	"bytes"
	"fmt"
	"html/template"
	"ican/types"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
	"vgof/core"

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

// Application 服务接口
type ican struct {
	logger   *zap.Logger
	sendMail types.SendMailService
	// 目标文件目录
	taskPath string
	subject  []interface{} // Subject列表
	body     []interface{} // body模板列表
}

// 用户任务
type userTask struct {
	SendEnable bool      `yaml:"send_enable"` // 发送使能
	Create     string    `yaml:"create"`      // 开始时间
	Deadline   string    `yaml:"deadline"`    // 结束时间
	PrevSend   string    `yaml:"prev_send"`   // 上一次发送提醒邮件的时间
	Period     string    `yaml:"period"`      // 发送邮件的间隔
	Money      string    `yaml:"money"`       // 金额
	Name       string    `yaml:"name"`        // 任务创建者的名字
	Address    string    `yaml:"address"`     // 任务创建者的邮箱
	Content    string    `yaml:"content"`     // 目标内容
	NowT       time.Time `yaml:"-"`
	DeadlineT  time.Time `yaml:"-"`
	CreateT    time.Time `yaml:"-"`
}

var _ core.ApplicationService = (*ican)(nil)

func (a *ican) Main(s core.Service) {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("Some Error: ", zap.Any("panic", r))
		}
	}()
	a.logger.Info("I Can!")
	a.logger.Info("Config:", zap.String("TaskPath", a.taskPath))
	// ==== 以下代码以1分钟为周期执行
	dura, err := time.ParseDuration("1m")
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(dura)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 获得所有task
	for {
		select {
		case <-ticker.C:
			taskFiles := a.getTaskFiles()
			for _, v := range taskFiles {
				task := a.getUserTask(v)
				if task.SendEnable == false { // 检查邮件发送开关是否使能
					continue
				}
				if task.PrevSend == "" { // 上次发送时间为空
					task.PrevSend = time.Now().Format(time.RFC3339)
					// 更新任务文件
					a.updateUserTask(v, task)
				}
				if a.checkTaskPeriod(task) { // 发送了邮件,更新
					a.logger.Info("Update task file", zap.String("File", v))
					a.updateUserTask(v, task)
				}
			}
		case <-sigs:
			a.logger.Info("Receive exit signal")
			return
		}
	}
}

// 检查任务提醒周期是否到达
// 返回为真则表示发送了邮件,更新了下次发送时间,需要写回到文件
func (a *ican) checkTaskPeriod(task *userTask) bool {
	var err error
	// 截止时间
	task.DeadlineT, err = time.Parse(time.RFC3339, task.Deadline)
	if err != nil {
		panic(err)
	}
	// 任务创建时间
	task.CreateT, err = time.Parse(time.RFC3339, task.Create)
	if err != nil {
		panic(err)
	}
	// 当前时间
	task.NowT = time.Now()
	var prev time.Time
	// 上次发送时间
	prev, err = time.Parse(time.RFC3339, task.PrevSend)
	if err != nil {
		panic(err)
	}

	var dur time.Duration
	// 计算时间间隔
	if strings.ToUpper(task.Period) == "AUTO" {
		timeSpan := task.DeadlineT.Sub(task.CreateT)
		if int(timeSpan.Hours()) < 8760 { // 任务跨度小时数小于1年
			dur, _ = time.ParseDuration("360h") // 15天
		} else {
			dur, _ = time.ParseDuration("720h") // 30天
		}
	} else {
		// 解析用户指定的时间跨度
		dur, err = time.ParseDuration(task.Period)
		if err != nil {
			panic(err)
		}
	}

	// 下次发送时间
	nextSend := prev.Add(dur)
	// 是否到达截止日期
	if nextSend.Unix() > task.DeadlineT.Unix() {
		a.logger.Info("Task arrival deadline.", zap.String("Name", task.Name), zap.String("Content", task.Content))
		// 给我自己发送提醒邮件
		addr := types.MailAddr{
			Name: "任务到期提醒",
			Addr: "virusv@live.com",
		}
		err := a.sendMail.SendTo(addr, fmt.Sprintf("%s任务到期", task.Name), fmt.Sprintf("金额:%s\n任务内容:%s\n", task.Money, task.Content))
		if err != nil {
			a.logger.Error("Send email to administrator failed", zap.Error(err), zap.String("Name", task.Name), zap.String("Content", task.Content))
			return false
		}
		task.SendEnable = false
		return true
	}
	// 如果下次发送时间大于当前时间,则不发送
	if nextSend.Unix() > task.NowT.Unix() {
		return false
	}
	// 需要发送邮件
	// 获得邮件内容和subject
	subject, body := a.getEmailContent(task)
	a.logger.Info("Send Email", zap.String("Subject", subject), zap.String("Body", body))
	// 发送出去
	addr := types.MailAddr{
		Name: task.Name,
		Addr: task.Address,
	}
	err = a.sendMail.SendTo(addr, subject, body)
	if err != nil {
		a.logger.Error("Send email failed", zap.Error(err), zap.String("Name", task.Name), zap.String("Content", task.Content))
		return false
	}
	a.logger.Info("Send email success", zap.String("Name", task.Name), zap.String("Content", task.Content))
	// 更新信息
	task.PrevSend = task.NowT.Format(time.RFC3339)
	return true
}

// 生成subject和body
func (a *ican) getEmailContent(task *userTask) (string, string) {
	var subject, body bytes.Buffer
	// 模板中使用的函数
	funcs := template.FuncMap{
		"days": func(s, e time.Time) int {
			dur := e.Sub(s)
			return int(dur.Hours() / 24)
		},
	}
	rand.Seed(time.Now().Unix())
	// 随机选择一个主题模板
	subjectTmp := a.subject[rand.Intn(len(a.subject))]
	subT := template.New("email subject").Funcs(funcs)
	subT, err := subT.Parse(subjectTmp.(string))
	if err != nil {
		panic(err)
	}
	subT.Execute(&subject, task)
	// 解析body
	bodyTmp := a.body[rand.Intn(len(a.body))]
	bodyT := template.New("email body").Funcs(funcs)
	bodyT, err = bodyT.Parse(bodyTmp.(string))
	if err != nil {
		panic(err)
	}
	bodyT.Execute(&body, task)
	return subject.String(), body.String()
}

// getUserTask 获得用户任务对象
func (a *ican) getUserTask(file string) *userTask {
	taskObj := &userTask{}
	taskContent, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	// 解析
	err = yaml.Unmarshal(taskContent, taskObj)
	if err != nil {
		panic(err)
	}
	return taskObj
}

// 更新用户任务
func (a *ican) updateUserTask(file string, task *userTask) {
	content, err := yaml.Marshal(task)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file, content, 0644)
	if err != nil {
		panic(err)
	}
}

// getTaskFiles 获得目标文件的列表
func (a *ican) getTaskFiles() []string {
	taskFiles := make([]string, 0)
	// 获取所有目标文件
	taskPath, err := ioutil.ReadDir(a.taskPath)
	if err != nil {
		panic(err)
	}
	for _, file := range taskPath {
		if file.IsDir() {
			continue
		} else {
			ext := path.Ext(file.Name())
			if ext != ".yaml" {
				continue
			}
			taskFiles = append(taskFiles, strings.TrimRight(a.taskPath, "/")+"/"+file.Name())
		}
	}
	return taskFiles
}
