package protocol

import "time"

// RunlyProtocol 代表整个 .runly 文件的根结构，是 Runly 资产的终极载体
type RunlyProtocol struct {
	Manifest   Manifest            `yaml:"manifest" json:"manifest"`
	Knowledge  []KnowledgeResource `yaml:"knowledge,omitempty" json:"knowledge,omitempty"`
	Skills     []SkillResource     `yaml:"skills,omitempty" json:"skills,omitempty"`
	Dictionary Dictionary          `yaml:"dictionary" json:"dictionary"`
	Topology   Topology            `yaml:"topology" json:"topology"`
	Commerce   Commerce            `yaml:"commerce" json:"commerce"`
	Security   Security            `yaml:"security" json:"security"`
}

// 1. MANIFEST - 协议元数据，定义资产的身份与版本
type Manifest struct {
	URN        string    `yaml:"urn" json:"urn"`                 // 资源统一名称
	Title      string    `yaml:"title" json:"title"`             // 资产标题
	Version    string    `yaml:"version" json:"version"`         // 语义化版本
	Status     string    `yaml:"status" json:"status"`           // 状态: draft | published | deprecated
	Creator    Creator   `yaml:"creator" json:"creator"`         // 创作者信息
	CreatedAt  time.Time `yaml:"created_at" json:"created_at"`   // 创建时间
	UpdatedAt  time.Time `yaml:"updated_at" json:"updated_at"`   // 更新时间
	MinRuntime string    `yaml:"min_runtime" json:"min_runtime"` // 最低 CLI 版本要求
}

type Creator struct {
	MeID   string `yaml:"me_id" json:"me_id"`     // 来自 Runly Me 的唯一身份标识
	Name   string `yaml:"name" json:"name"`       // 专家显示名称
	PubKey string `yaml:"pub_key" json:"pub_key"` // 专家的 Ed25519 公钥 (Hex)
}

// 2. KNOWLEDGE - 知识资源契约
type KnowledgeResource struct {
	ID           string             `yaml:"id" json:"id"`
	ProviderType string             `yaml:"provider_type" json:"provider_type"` // SEMANTIC_API | VDB_DIRECT
	Description  string             `yaml:"description" json:"description"`
	Config       KnowledgeConfig    `yaml:"config" json:"config"`
	Injection    KnowledgeInjection `yaml:"injection" json:"injection"`
}

type KnowledgeConfig struct {
	Endpoint   string            `yaml:"endpoint" json:"endpoint"`
	Method     string            `yaml:"method" json:"method"`
	Timeout    int               `yaml:"timeout" json:"timeout"`
	MaxRetries int               `yaml:"max_retries" json:"max_retries"`
	Headers    map[string]string `yaml:"headers" json:"headers"`
	VDBParams  *VDBParams        `yaml:"vdb_params,omitempty" json:"vdb_params,omitempty"`
}

type VDBParams struct {
	IndexName      string  `yaml:"index_name" json:"index_name"`
	TopK           int     `yaml:"top_k" json:"top_k"`
	Threshold      float64 `yaml:"threshold" json:"threshold"`
	EmbeddingModel string  `yaml:"embedding_model" json:"embedding_model"`
}

type KnowledgeInjection struct {
	TargetVariable string `yaml:"target_variable" json:"target_variable"`
	MaxTokens      int    `yaml:"max_tokens" json:"max_tokens"`
	Format         string `yaml:"format" json:"format"`
	CacheTTL       int    `yaml:"cache_ttl" json:"cache_ttl"`
}

// 3. SKILLS - 技能/工具能力契约
type SkillResource struct {
	ID          string        `yaml:"id" json:"id"`
	Type        string        `yaml:"type" json:"type"`
	Description string        `yaml:"description" json:"description"`
	Config      SkillConfig   `yaml:"config" json:"config"`
	Contract    SkillContract `yaml:"contract" json:"contract"`
}

type SkillConfig struct {
	Endpoint   string            `yaml:"endpoint" json:"endpoint"`
	Method     string            `yaml:"method" json:"method"`
	Timeout    int               `yaml:"timeout" json:"timeout"`
	MaxRetries int               `yaml:"max_retries" json:"max_retries"`
	Headers    map[string]string `yaml:"headers" json:"headers"`
}

type SkillContract struct {
	Request  map[string]interface{} `yaml:"request" json:"request"`
	Response ResponseContract       `yaml:"response" json:"response"`
}

type ResponseContract struct {
	StrictMode bool                   `yaml:"strict_mode" json:"strict_mode"`
	Schema     map[string]interface{} `yaml:"schema" json:"schema"`
}

// 4. DICTIONARY - 运行时数据字典
type Dictionary struct {
	Inputs    []Parameter `yaml:"inputs" json:"inputs"`
	Artifacts []Artifact  `yaml:"artifacts" json:"artifacts"`
}

type Parameter struct {
	Name     string      `yaml:"name" json:"name"`
	Type     string      `yaml:"type" json:"type"`
	Pattern  string      `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Values   []string    `yaml:"values,omitempty" json:"values,omitempty"`
	Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Required bool        `yaml:"required" json:"required"`
}

type Artifact struct {
	ID          string                 `yaml:"id" json:"id"`
	Type        string                 `yaml:"type" json:"type"`
	Description string                 `yaml:"description" json:"description"`
	Schema      map[string]interface{} `yaml:"schema" json:"schema"`
}

// 5. TOPOLOGY - 任务调度拓扑图
type Topology struct {
	StartAt string `yaml:"start_at" json:"start_at"`
	Nodes   []Node `yaml:"nodes" json:"nodes"`
}

type Node struct {
	ID        string                 `yaml:"id" json:"id"`
	Type      string                 `yaml:"type" json:"type"` // SKILL_CALL | AI_TASK | HITL | LOGIC_GATE | TERMINUS
	Config    map[string]interface{} `yaml:"config" json:"config"`
	OnSuccess string                 `yaml:"on_success,omitempty" json:"on_success,omitempty"`
	OnFailure string                 `yaml:"on_failure,omitempty" json:"on_failure,omitempty"`
	Rules     []LogicRule            `yaml:"rules,omitempty" json:"rules,omitempty"`
}

type LogicRule struct {
	Condition string `yaml:"condition" json:"condition"`
	Next      string `yaml:"next" json:"next"`
}

// 6. COMMERCE - 商业授权与清算模块
type Commerce struct {
	Pricing    Pricing    `yaml:"pricing" json:"pricing"`
	Royalty    Royalty    `yaml:"royalty" json:"royalty"`
	Settlement Settlement `yaml:"settlement" json:"settlement"`
}

type Pricing struct {
	Mode     string  `yaml:"mode" json:"mode"` // FREE | PAY_PER_USE | SUBSCRIPTION
	Amount   float64 `yaml:"amount" json:"amount"`
	Currency string  `yaml:"currency" json:"currency"`
}

type Royalty struct {
	CreatorShare  float64 `yaml:"creator_share" json:"creator_share"`
	PlatformShare float64 `yaml:"platform_share" json:"platform_share"`
}

type Settlement struct {
	Trigger string `yaml:"trigger" json:"trigger"` // INSTANT | BATCH_MONTHLY
}

// 7. SECURITY - 资产指纹与数字签名
type Security struct {
	HashAlgo  string `yaml:"hash_algo" json:"hash_algo"` // SHA-256
	Signature string `yaml:"signature" json:"signature"` // Ed25519 Signature (Hex)
}
