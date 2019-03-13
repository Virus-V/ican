System:
  TaskPath: ./conf/tasks.d/
  Subjects:
    - "叮!来走访一下你的任务完成的怎么样了呢~"
    - "偷偷瞄一眼,有木有努力呀?"
    - "你努力奋斗的样子,真的让人心动呢~"
    - "ONE DAY, I WILL SAY \"I DID IT\""
  Bodys: 
    - "嗨! {{.Name}}~, 距离目标达成还有{{days .NowT .DeadlineT}}天! 准备的怎么样啦?(doge\n\n你的目标:{{.Content}}\n你的金额:{{.Money}}\n截止期限:{{.DeadlineT.Format \"2006年01月02日\"}}"
    - "{{.Name}}~ 好久不见. 回望过去的{{days .CreateT .NowT}}天, 你创造了奇迹!\n\n你的目标:{{.Content}}\n你的金额:{{.Money}}\n截止期限:{{.DeadlineT.Format \"2006年01月02日\"}}"
Mail:
  From: virusv@qq.com
  FromName: ICan!
  SMTPAddr: smtp.qq.com
  SMTPPort: 587
  SMTPUsername: virusv@qq.com
  SMTPPassword: SMTP登陆密码
Log: