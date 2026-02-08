package protocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
)

// varExtractRegex 匹配变量引用格式：{{inputs.xxx}} 或 {{steps.node_id.output}}
var varExtractRegex = regexp.MustCompile(`\{\{\s*([\w\.]+)\s*\}\}`)

// Validate 执行全量静态语义校验
func Validate(proto *RunlyProtocol) error {
	// 1. 构建节点快速索引，用于 O(1) 查找
	nodeMap := make(map[string]Node)
	for _, node := range proto.Topology.Nodes {
		nodeMap[node.ID] = node
	}

	// 2. 检查拓扑连通性（起始节点、逻辑分支、末端节点）
	if err := validateTopology(proto, nodeMap); err != nil {
		return err
	}

	// 3. 检查变量引用一致性
	if err := validateVariables(proto, nodeMap); err != nil {
		return err
	}

	// 4. 检查外部资源（Skills/Knowledge）引用有效性
	if err := validateResourceLinks(proto, nodeMap); err != nil {
		return err
	}

	return nil
}

// validateTopology 验证拓扑结构的完整性
func validateTopology(proto *RunlyProtocol, nodeMap map[string]Node) error {
	// 验证 StartAt 节点是否存在
	if _, ok := nodeMap[proto.Topology.StartAt]; !ok {
		// 🚩 拓扑起始节点 [%s] 未定义
		return fmt.Errorf(i18n.T("errors.start_node_missing"), proto.Topology.StartAt)
	}

	// 遍历所有节点，验证其下游跳转 ID
	for _, node := range proto.Topology.Nodes {
		// 收集所有潜在跳转路径
		targets := []string{node.OnSuccess, node.OnFailure}
		if node.Type == "LOGIC_GATE" {
			for _, rule := range node.Rules {
				targets = append(targets, rule.Next)
			}
		}

		for _, t := range targets {
			// 跳过终点标记
			if t == "" || t == "terminate" || t == "terminate_error" {
				continue
			}
			// 检查下游节点是否存在
			if _, exists := nodeMap[t]; !exists {
				// 📍 节点 [%s] 引用了不存在的下游目标: %s
				return fmt.Errorf(i18n.T("errors.node_not_found"), node.ID, t)
			}
		}
	}
	return nil
}

// validateVariables 验证所有变量引用的源头是否合法
func validateVariables(proto *RunlyProtocol, nodeMap map[string]Node) error {
	for _, node := range proto.Topology.Nodes {
		// 序列化 Config 进行静态扫描，查找 {{...}} 占位符
		rawConfig := fmt.Sprintf("%v", node.Config)
		matches := varExtractRegex.FindAllStringSubmatch(rawConfig, -1)

		for _, match := range matches {
			path := match[1]
			parts := strings.Split(path, ".")

			switch parts[0] {
			case "inputs":
				// 检查 Dictionary.Inputs 域
				if !hasInputParam(proto.Dictionary.Inputs, parts[1]) {
					// ⌨️ 节点 [%s] 引用了 Dictionary 中未定义的输入参数: %s
					return fmt.Errorf(i18n.T("errors.input_ref_missing"), node.ID, parts[1])
				}
			case "steps":
				// 检查 Steps 引用格式及引用的节点是否存在
				if len(parts) < 3 {
					// 🔗 节点 [%s] 的变量引用格式错误: %s
					return fmt.Errorf(i18n.T("errors.var_format_err"), node.ID, path)
				}
				refNodeID := parts[1]
				if _, exists := nodeMap[refNodeID]; !exists {
					// 📍 节点 [%s] 引用了不存在的对象: %s
					return fmt.Errorf(i18n.T("errors.node_not_found"), node.ID, refNodeID)
				}
			}
		}
	}
	return nil
}

// validateResourceLinks 验证节点对 Skill 和 Knowledge 的引用
func validateResourceLinks(proto *RunlyProtocol, nodeMap map[string]Node) error {
	for _, node := range proto.Topology.Nodes {
		// 技能引用检查
		if node.Type == "SKILL_CALL" {
			ref, _ := node.Config["skill_ref"].(string)
			if !hasSkillID(proto.Skills, ref) {
				// 🛠️ 节点 [%s] 引用的技能 [%s] 未在 skills 域定义
				return fmt.Errorf(i18n.T("errors.skill_ref_missing"), node.ID, ref)
			}
		}

		// 知识库引用检查
		if node.Type == "AI_TASK" {
			ref, ok := node.Config["knowledge_ref"].(string)
			if ok && ref != "" {
				if !hasKnowledgeID(proto.Knowledge, ref) {
					// 📚 节点 [%s] 引用的知识库 [%s] 未在 knowledge 域定义
					return fmt.Errorf(i18n.T("errors.kb_ref_missing"), node.ID, ref)
				}
			}
		}
	}
	return nil
}

// 辅助查询逻辑
func hasInputParam(params []Parameter, name string) bool {
	for _, p := range params {
		if p.Name == name {
			return true
		}
	}
	return false
}

func hasSkillID(skills []SkillResource, id string) bool {
	for _, s := range skills {
		if s.ID == id {
			return true
		}
	}
	return false
}

func hasKnowledgeID(kb []KnowledgeResource, id string) bool {
	for _, k := range kb {
		if k.ID == id {
			return true
		}
	}
	return false
}
