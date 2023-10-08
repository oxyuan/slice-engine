package _const

import mapset "github.com/deckarep/golang-set"

var (
	CHAT_NOISE_SET = mapset.NewSet()

	LEFT_BRACKET_SET = mapset.NewSet("(", "（", "【", "[", "{", "《", "<", "〈", "「", "『", "〖", "〔", "〘", "〚", "﹤")

	RIGHT_BRACKET_SET = mapset.NewSet("）", ")", "】", "]", "}", "》", ">", "〉", "」", "』", "〗", "〕", "〙", "〛", "﹥", "﹚", "﹜", "﹞")

	CHAT_NOISE_COMMON_SET = mapset.NewSet(
		",", "。", "*", "%", "、", ":", "：", "\n", "#", "`", "~", "!", "$", "^", "&", "+", "=", "|", "-", ".", "·", "'", ",", "\\", "\"", "/", "?", "？", "”", "“", "￥", "¥",
		"﹃", "﹄", "﹁", "﹂", "﹏", "﹍", "﹎", "﹋", "﹌", "﹊", "﹉", "﹐", "﹑", "﹒", "﹔", "﹕", "", "；", ";", "﹖", "﹗", "﹢", "﹣", "﹦", "﹨", "﹩", "﹪", "﹫", "﹟", "﹠", "﹡",
	)

	CONNECTOR_SET = mapset.NewSet("-")

	CHINESE_NUMBER_SET = mapset.NewSet("一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "廿", "百", "千")

	STAGE_SET = mapset.NewSet("前", "中", "后", "晚")

	GRADE_SUFFIX_SET = mapset.NewSet("型", "度", "期", "级", "组", "阴性", "阳性")

	DESCRIPTIVE_WORD_SET = mapset.NewSet("阴性", "阳性")

	XINA_SET = mapset.NewSet("Γ", "Δ", "Θ", "Λ", "Ξ", "Π", "Σ", "Υ", "Φ", "Ψ", "Ω", "λ", "μ")

	DESC_WORD_SET = mapset.NewSet("其它", "其他", "部分", "不伴", "多发性", "相关性", "功能")

	UNIT_SUFFIX_SET = mapset.NewSet(
		"％", "%", "umol/l", "/g/l", "/L", "/l", "%l", "cm", "w", "mmhg", "mmol/l", "mmol/L", "mmol", "ml", "mg", "mg/l", "mg/L", "mg/dl", "mg/dL", "mg/d",
		"次", "个", "型", "年", "月", "日", "点", "周", "瓶", "管", "级", "期", "区", "酯酶",
		"分", "片", "支", "盒", "袋", "包", "套", "组", "张", "根", "条", "块", "颗", "粒", "枚", "只", "株", "头", "斤", "升", "毫升", "克", "克",
	)

	STATE_UNIT_SET = mapset.NewSet(
		"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
		"Ⅰ", "Ⅱ", "Ⅲ", "Ⅳ", "Ⅴ", "Ⅵ", "Ⅶ", "Ⅷ", "Ⅸ", "Ⅹ", "Ⅺ", "Ⅻ",
	)

	SIMILAR_KIND_SET = mapset.NewSet(
		mapset.NewSet("口", "咽", "唇", "喉", "眼", "耳", "鼻", "额", "眉"),
		mapset.NewSet("腰", "头", "肘", "肩", "背", "胳", "脖", "面", "颅", "颈", "腹", "颌", "臂", "腋", "胸", "胸壁", "腿", "躯干", "脑室", "纵隔"),
		mapset.NewSet("趾", "趾骨", "踝", "踝关节"),
		mapset.NewSet("头", "脑"),
		mapset.NewSet("眼", "视网膜", "睑"),
		mapset.NewSet("脐", "腰", "腹", "肚", "脂", "脊", "胸", "臀", "下肢", "足"),
		mapset.NewSet("肺", "肾", "脾", "脏", "肛", "肝", "胃", "胆", "胞", "胎", "胰", "腺", "肠"),
		mapset.NewSet("气体", "烟雾", "蒸汽", "气味"),
		mapset.NewSet("切割", "针刺", "穿孔", "出血", "切口", "创伤", "创口", "创面", "清创", "加压"),
		mapset.NewSet("器械", "材料", "装置", "设备", "工具"),
	)

	ADAPTIVE_WORD_SET = mapset.NewSet("部", "区", "侧", "上", "下", "左", "右", "前", "后", "大", "小")

	OPT_WORD_SET = mapset.NewSet("部分")

	ESCAPE_SYMBOL_SET = mapset.NewSet("＼", "\\")

	ROMAN_CHARACTER_SET = mapset.NewSet("Ⅰ", "Ⅱ", "Ⅲ", "Ⅳ", "Ⅴ", "Ⅵ", "Ⅶ", "Ⅷ", "Ⅸ", "Ⅹ", "Ⅺ", "Ⅻ")
)

func init() {
	for v := range LEFT_BRACKET_SET.Iter() {
		CHAT_NOISE_SET.Add(v)
	}
	for v := range RIGHT_BRACKET_SET.Iter() {
		CHAT_NOISE_SET.Add(v)
	}
	for v := range CHAT_NOISE_COMMON_SET.Iter() {
		CHAT_NOISE_SET.Add(v)
	}
	for v := range CHINESE_NUMBER_SET.Iter() {
		STAGE_SET.Add(v)
	}
}
