package robot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/opentdp/wechat-rest/args"
	"github.com/opentdp/wechat-rest/wcferry"
	"github.com/opentdp/wechat-rest/wclient/model"
)

func helpHandler() {

	handlers["/help"] = &Handler{
		Level:    0,
		ChatAble: true,
		RoomAble: true,
		Describe: "查看帮助信息",
		Callback: func(msg *wcferry.WxMsg) string {
			user := args.Usr.Member[msg.Sender]
			// 生成指令菜单
			helper := []string{}
			for k, v := range handlers {
				if v.Level > 0 {
					if user == nil || v.Level > user.Level {
						continue // 没有权限
					}
				}
				if msg.IsGroup {
					if v.RoomAble { // 群聊指令
						helper = append(helper, k+" "+v.Describe)
					}
				} else {
					if v.ChatAble { // 私聊指令
						helper = append(helper, k+" "+v.Describe)
					}
				}
			}
			sort.Strings(helper)
			text := strings.Join(helper, "\n") + "\n"
			// 模型运行时信息
			if model.GetUser(msg.Sender).AiArgot != "" {
				text += "唤醒词 " + model.GetUser(msg.Sender).AiArgot + "，"
			}
			if len(args.LLM.Models) > 0 {
				text += "对话模型 " + model.GetUserModel(msg.Sender).Family + "，"
				text += fmt.Sprintf("上下文长度 %d/%d", model.CountHistory(msg.Sender), args.LLM.HistoryNum)
			}
			return text
		},
	}

}
