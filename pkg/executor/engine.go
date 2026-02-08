package executor

import (
	"fmt"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
)

// Context 维护运行时数据域
type Context struct {
	Vars      map[string]interface{} // 存储 inputs 和各步骤的 outputs
	Artifacts map[string]interface{} // 存储最终交付物
}

// Engine 拓扑执行引擎
type Engine struct {
	Protocol *protocol.RunlyProtocol
	Context  *Context
}

// NewEngine 初始化引擎并注入初始输入
func NewEngine(p *protocol.RunlyProtocol, inputs map[string]interface{}) *Engine {
	return &Engine{
		Protocol: p,
		Context: &Context{
			Vars: map[string]interface{}{
				"inputs": inputs,
				"steps":  make(map[string]interface{}),
			},
			Artifacts: make(map[string]interface{}),
		},
	}
}

// Run 启动多语言感知的仿真运行
func (e *Engine) Run() error {
	ui.PrintHeader("executor.engine_header")

	currentNodeID := e.Protocol.Topology.StartAt
	for {
		if currentNodeID == "terminate" || currentNodeID == "" {
			break
		}

		node := e.findNode(currentNodeID)
		if node == nil {
			// 使用 i18n 报告节点未找到错误
			return fmt.Errorf(i18n.T("errors.node_not_found"), "SYSTEM", currentNodeID)
		}

		// 输出当前步骤：正在执行节点 [%s] (%s)
		ui.PrintStep("executor.step_executing", node.ID, node.Type)

		nextID, err := e.executeNode(node)
		if err != nil {
			if node.OnFailure != "" {
				// 打印跳转提示：条件不匹配或执行失败，正在跳转至错误处理分支
				ui.PrintStep("executor.node_jump", node.OnFailure)
				currentNodeID = node.OnFailure
				continue
			}
			return err
		}
		currentNodeID = nextID
	}

	ui.PrintSuccess("executor.execution_complete")
	return nil
}

func (e *Engine) findNode(id string) *protocol.Node {
	for _, n := range e.Protocol.Topology.Nodes {
		if n.ID == id {
			return &n
		}
	}
	return nil
}

func (e *Engine) executeNode(n *protocol.Node) (string, error) {
	steps := e.Context.Vars["steps"].(map[string]interface{})

	switch n.Type {
	case "SKILL_CALL":
		skillRef, _ := n.Config["skill_ref"].(string)
		// 输出：📡 正在连接服务端: %s
		ui.PrintStep("executor.skill_calling", skillRef)

		steps[n.ID] = map[string]interface{}{"output": "MOCK_SKILL_DATA"}
		return n.OnSuccess, nil

	case "AI_TASK":
		// 输出：🤖 正在执行 AI 推理任务...
		ui.PrintStep("executor.ai_processing")

		prompt, _ := n.Config["prompt"].(string)
		rendered := RenderTemplate(prompt, e.Context)

		steps[n.ID] = map[string]interface{}{"output": "AI_RESULT_FOR_" + rendered}
		return n.OnSuccess, nil

	case "HITL":
		instruction, _ := n.Config["instruction"].(string)
		// 输出：🧑‍💻 等待专家审核: %s
		ui.PrintStep("executor.hitl_waiting", instruction)
		// 输出：⌨️  按回车键 [Enter] 模拟专家授权...
		ui.PrintStep("executor.hitl_continue")
		fmt.Scanln()
		return n.OnSuccess, nil

	case "TERMINUS":
		artifactRef, _ := n.Config["artifact_ref"].(string)
		dataSource, _ := n.Config["data_source"].(string)

		// 模拟渲染最终数据
		finalData := RenderTemplate("{{"+dataSource+"}}", e.Context)
		e.Context.Artifacts[artifactRef] = finalData
		return "terminate", nil

	default:
		return n.OnSuccess, nil
	}
}
