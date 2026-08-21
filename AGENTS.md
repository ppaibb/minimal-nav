# Minimal Nav - 项目核心设计规范与约束 (Design Rules)

> 本文件记录本项目的审美底线与架构铁律，后续任何 AI Agent、开发者均须严格遵守，严禁违反！

---

## 🚫 绝对红线与禁区 (Strict Prohibitions)

1. **坚决不要滥用“圆角胶囊（Capsules / Pill Buttons）”与“悬浮泡泡卡片”**：
   - 严禁添加类似药丸胶囊、圆角悬浮浮岛（Floating Islands）、厚重描边胶囊等元素。
   - 这种设计显得过度包装与油腻，严重破坏页面的干练与高级质感。

2. **坚决不要滥用大面积渐变色块与浮夸阴影**：
   - 严禁大面积花哨渐变色（如蓝紫渐变、彩虹渐变底色）。
   - 严禁多层浓重投影（heavy drop-shadows）。

---

## 🏛️ 提倡的设计语言 (Recommended Aesthetic)

1. **冷峻直线的建筑与出版物排版 (Swiss & Editorial Line Layout)**：
   - 使用**极细直线分割线（`border-b border-border/50` 或 `border-zinc-200`）**划分空间，干净利落。
   - 保持几何直角或极小圆角（`rounded-sm` / 直线冷峻感）。

2. **纯粹清晰的字阶与无衬线排版 (Crisp Sans-Serif Typography)**：
   - 标题与正文字迹清晰挺拔，层级分明。
   - 巧妙利用中英文字号差、等宽代码字体（`font-mono`）点缀，而不是依赖色块堆砌。

---

## 📝 提交规范 (Git Rules)

- 提交 Git 时，必须编写中文 Log，包含：**问题或需求描述**、**修复或实现思路**。
