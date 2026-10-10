<template>
  <transition name="fade">
    <div v-if="show" class="modal-overlay" @click.self="emit('update:show', false)">
      <div class="modal-content guide-modal">
        <!-- 弹窗顶部 Header -->
        <div class="modal-header">
          <div class="title-with-icon">
            <div class="header-icon-box">
              <SvgIcon name="book-open" :size="20" class="header-icon" />
            </div>
            <div>
              <h2 class="modal-title">{{ $t('userGuideTitle') }}</h2>
              <p class="modal-subtitle">{{ $t('userGuideSubtitle') }}</p>
            </div>
          </div>
          <button class="modal-close-btn" @click="emit('update:show', false)" :title="$t('close')">
            <SvgIcon name="x" :size="16" />
          </button>
        </div>

        <!-- 弹窗主体 Body -->
        <div class="modal-body">
          <!-- 左侧/顶部 Tab 导航 -->
          <div class="guide-nav-bar">
            <button
              v-for="(tab, index) in tabs"
              :key="tab.id"
              class="guide-nav-btn"
              :class="{ active: currentTab === tab.id }"
              @click="currentTab = tab.id"
            >
              <SvgIcon :name="tab.icon" :size="14" />
              <span>{{ tab.label }}</span>
              <span v-if="currentTab === tab.id" class="nav-active-pill"></span>
            </button>
          </div>

          <!-- 内容展示区 -->
          <div class="guide-content-area">
            <!-- Tab 1: 快速上手 -->
            <div v-if="currentTab === 'quick-start'" class="tab-panel">
              <div class="panel-hero">
                <div class="hero-badge">STEP BY STEP</div>
                <h3 class="panel-title">三步搞定 L4D2 模组日常管理</h3>
                <p class="panel-desc">从定位游戏到轻松整理数十个模组，无需繁琐操作。</p>
              </div>

              <div class="steps-flow">
                <div class="step-card">
                  <div class="step-num">01</div>
                  <div class="step-info">
                    <div class="step-header">
                      <SvgIcon name="compass" :size="16" class="step-icon" />
                      <h4>定位游戏目录</h4>
                    </div>
                    <p>
                      点击目录栏的<strong>“自动搜索”</strong>按钮，程序会在 0.1 秒内自动扫描电脑中的运行进程、系统记录、桌面快捷方式与磁盘目录，无论 Steam 版还是任意非 Steam 免安装/对战平台版均可精准检测，弹窗一键选择！
                    </p>
                    <div class="tip-box">
                      <SvgIcon name="info" :size="13" />
                      <span>若客户端位于深层冷门路径，可在弹窗内点击“全盘深度搜索”。</span>
                    </div>
                  </div>
                </div>

                <div class="step-card">
                  <div class="step-num">02</div>
                  <div class="step-info">
                    <div class="step-header">
                      <SvgIcon name="grid" :size="16" class="step-icon" />
                      <h4>浏览与智能筛选</h4>
                    </div>
                    <p>
                      主界面清晰划分<strong>本地 VPK</strong>与<strong>创意工坊</strong>。支持关键词即时过滤（支持中文显示名、文件名、备注以及三方地图关卡代号），支持按大小、修改时间或名称排序。
                    </p>
                  </div>
                </div>

                <div class="step-card">
                  <div class="step-num">03</div>
                  <div class="step-info">
                    <div class="step-header">
                      <SvgIcon name="zap" :size="16" class="step-icon" />
                      <h4>毫秒级启停与批量操作</h4>
                    </div>
                    <p>
                      点击模组开关即可<strong>启用 / 禁用</strong>，底层通过重命名为 <code>.disabled</code> 实现，游戏引擎即时跳过，<strong>无需重新下载或复制大文件</strong>；勾选多个模组还可使用底部操作栏进行批量管理。
                    </p>
                  </div>
                </div>

                <div class="step-card">
                  <div class="step-num">04</div>
                  <div class="step-info">
                    <div class="step-header">
                      <SvgIcon name="play" :size="16" class="step-icon" />
                      <h4>智能启动与参数优化</h4>
                    </div>
                    <p>
                      点击目录栏右侧<strong>“启动游戏”</strong>即可开玩！若为 Steam 客户端将唤起 Steam 保证联机与云存档，若为独立客户端将直接运行游戏主程序。
                    </p>
                    <div class="tip-box">
                      <SvgIcon name="info" :size="13" />
                      <span>在<strong>“设置”</strong>中可自由开启<strong>跳过开头动画 (-novid)</strong>、<strong>强制中文启动 (-language schinese)</strong>及自定义额外参数。</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tab 2: 模组管理与本地化 -->
            <div v-if="currentTab === 'mods'" class="tab-panel">
              <div class="panel-hero">
                <div class="hero-badge">WORKFLOW</div>
                <h3 class="panel-title">本地化、导入与备份导出</h3>
                <p class="panel-desc">告别 Steam 订阅失效困扰，随时备份与导入模组包。</p>
              </div>

              <div class="feature-cards-grid">
                <div class="feature-item-card">
                  <div class="feature-card-top">
                    <div class="feature-icon-box icon-purple">
                      <SvgIcon name="box" :size="18" />
                    </div>
                    <h4>创意工坊一键本地化</h4>
                  </div>
                  <p>
                    Steam 创意工坊下载的 VPK 均为一串无意义的纯数字编号，且若作者下架模组将导致游戏内失效。
                  </p>
                  <div class="feature-highlight">
                    <strong>解决方案：</strong>在工坊列表中勾选模组，点击底部<strong>“本地化”</strong>，程序会自动解析其真实标题、复制重命名到本地 addons 目录，彻底脱离工坊订阅依赖！
                  </div>
                </div>

                <div class="feature-item-card">
                  <div class="feature-card-top">
                    <div class="feature-icon-box icon-green">
                      <SvgIcon name="import" :size="18" />
                    </div>
                    <h4>压缩包与模组导入 (支持拖拽)</h4>
                  </div>
                  <p>
                    从 Gamemaps、群文件下载了 <code>.zip</code>、<code>.7z</code>、<code>.tar.zst</code> 模组压缩包或 <code>.vpk</code> 文件？
                  </p>
                  <div class="feature-highlight">
                    <strong>超便捷支持：</strong>直接把一个或多个压缩包、<code>.vpk</code> 文件<strong>从桌面或文件管理器拖入窗口任意位置即可自动解压并安装</strong>！亦可点击目录栏<strong>“导入模组”</strong>手动选取。
                  </div>
                </div>

                <div class="feature-item-card">
                  <div class="feature-card-top">
                    <div class="feature-icon-box icon-blue">
                      <SvgIcon name="export" :size="18" />
                    </div>
                    <h4>高效压缩打包导出 (Export)</h4>
                  </div>
                  <p>
                    想要打包一套自己精心搭配的武器与人物模型分享给开黑队友？
                  </p>
                  <div class="feature-highlight">
                    多选目标模组后点击<strong>“导出所选”</strong>，采用业界顶尖的 <code>tar.zst</code> 高压缩算法，压缩比极高、打包飞速，自动生成时间戳归档包。
                  </div>
                </div>

                <div class="feature-item-card">
                  <div class="feature-card-top">
                    <div class="feature-icon-box icon-purple">
                      <SvgIcon name="layers" :size="18" />
                    </div>
                    <h4>模组自定义分组与 AI 命名</h4>
                  </div>
                  <p>
                    将多个人物、武器或地图模组组合为一个主题分组，支持 12 种主题色彩标记、折叠聚合与一键全选。
                  </p>
                  <div class="feature-highlight">
                    <strong>智能省心：</strong>创建分组时支持<strong>AI 大模型智能推断</strong>与<strong>本地规则启发式保底</strong>，自动分析模组特征并预填贴切分组名称！
                  </div>
                </div>
              </div>
            </div>

            <!-- Tab 3: 地图指令与备注 -->
            <div v-if="currentTab === 'maps'" class="tab-panel">
              <div class="panel-hero">
                <div class="hero-badge">CAMPAIGN & NOTES</div>
                <h3 class="panel-title">地图控制台指令与自定义备注</h3>
                <p class="panel-desc">游戏内找不到三方图入口？模组太多记不住用途？一招搞定！</p>
              </div>

              <div class="split-guide-row">
                <div class="guide-box">
                  <div class="guide-box-header">
                    <SvgIcon name="map" :size="18" class="text-accent" />
                    <h4>第三方战役地图一键建图</h4>
                  </div>
                  <p>
                    很多优秀的第三方战役地图没有制作游戏主菜单战役海报，导致新手安装后在主菜单找不到。
                  </p>
                  <ul class="guide-list">
                    <li>本程序会自动穿透解析 VPK 内部的 mission 与 bsp 关卡文件。</li>
                    <li>在地图模组卡片上会清晰列出所有<strong>章节名称与内部代号</strong>。</li>
                    <li>点击<strong>“复制指令”</strong>，进游戏按 <code>~</code> 调出控制台直接粘贴即可开图！</li>
                  </ul>
                  <div class="code-preview-box">
                    <code>map c1m1_hotel</code>
                    <span class="code-tip">← 控制台回车直接秒开</span>
                  </div>
                </div>

                <div class="guide-box">
                  <div class="guide-box-header">
                    <SvgIcon name="edit" :size="18" class="text-accent" />
                    <h4>自定义中文备注 (Notes)</h4>
                  </div>
                  <p>
                    当安装了上百个模组时，英文原名很难区分是哪个角色的模型或哪把枪的声效。
                  </p>
                  <ul class="guide-list">
                    <li>点击卡片上的<strong>“附加备注”</strong>按钮，输入你熟悉的中文。</li>
                    <li>例如：输入 <code>[生化危机吉尔]</code> 或 <code>[连喷打击声]</code>。</li>
                    <li>备注会自动融合至卡片标题，且<strong>支持在顶部搜索框即时匹配</strong>！</li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Tab 4: 远程服务器同步 -->
            <div v-if="currentTab === 'remote'" class="tab-panel">
              <div class="panel-hero">
                <div class="hero-badge">SERVER MANAGER</div>
                <h3 class="panel-title">SFTP 远程服务器同步与运维</h3>
                <p class="panel-desc">为开黑腐竹量身定制：双栏分屏、双向模组同步、一键重启服务器。</p>
              </div>

              <div class="server-steps-grid">
                <div class="server-step-card">
                  <div class="card-step-badge">1. 配置连接</div>
                  <p>
                    在右上角<strong>“设置”</strong>中填入远程 Linux 游戏服务器的 SFTP 地址、账号、密码及 addons 绝对路径。
                  </p>
                </div>
                <div class="server-step-card">
                  <div class="card-step-badge">2. 开启分屏</div>
                  <p>
                    点击顶部<strong>“远程同步”</strong>按钮，主界面立刻化身左右分屏工作区，支持鼠标自由拖拽分割线调整面板宽度。
                  </p>
                </div>
                <div class="server-step-card">
                  <div class="card-step-badge">3. 双向同步</div>
                  <p>
                    本地选中的模组一键上传至服务器；右侧远程面板中的模组也可一键下载回本地，两端无缝流转。
                  </p>
                </div>
                <div class="server-step-card">
                  <div class="card-step-badge">4. 远程重启</div>
                  <p>
                    添加新图或插件后，在远程面板点击<strong>“重启服务”</strong>，自动执行远程容器热重启并反馈耗时。
                  </p>
                </div>
              </div>
            </div>

            <!-- Tab 5: 实用技巧与常见问题 -->
            <div v-if="currentTab === 'faq'" class="tab-panel">
              <div class="panel-hero">
                <div class="hero-badge">FAQ & TIPS</div>
                <h3 class="panel-title">实用技巧与常见疑问解答</h3>
                <p class="panel-desc">收集了日常使用中最常遇到的问题与解决方案。</p>
              </div>

              <div class="faq-accordion">
                <div class="faq-item">
                  <h4 class="faq-q">
                    <SvgIcon name="help-circle" :size="15" class="text-accent" />
                    <span>为什么启停模组速度几乎是瞬时的？</span>
                  </h4>
                  <p class="faq-a">
                    本管理器充分利用 Valve Source 引擎机制，禁用时仅将文件名后缀修改为 <code>.vpk.disabled</code>，启用时移除该后缀。这属于底层文件系统元数据重命名，不需要真正搬移或复制几百兆数据，不仅瞬间完成，还极大保护固态硬盘寿命。
                  </p>
                </div>

                <div class="faq-item">
                  <h4 class="faq-q">
                    <SvgIcon name="help-circle" :size="15" class="text-accent" />
                    <span>游戏内出现人物紫黑格子或闪退怎么办？</span>
                  </h4>
                  <p class="faq-a">
                    这是典型的“模组模型/材质冲突”。可在列表中勾选最近安装的模组逐一“禁用”排查；若是三方地图模组闪退，通常是未下载完整的地图依赖包（例如部分三方地图分为 Part 1 与 Part 2 两个 VPK，需全部安装）。
                  </p>
                </div>

                <div class="faq-item">
                  <h4 class="faq-q">
                    <SvgIcon name="help-circle" :size="15" class="text-accent" />
                    <span>有多个 L4D2 客户端时如何自由切换？</span>
                  </h4>
                  <p class="faq-a">
                    随时点击顶部目录栏旁边的<strong>“自动搜索”</strong>图标，程序会再次搜索全部客户端；或者直接点击“浏览文件夹”挑选任意目标路径。
                  </p>
                </div>

                <div class="faq-item">
                  <h4 class="faq-q">
                    <SvgIcon name="help-circle" :size="15" class="text-accent" />
                    <span>如何配置 OpenAI 或本地大模型 (Ollama / LM Studio)？</span>
                  </h4>
                  <p class="faq-a">
                    点击右上角<strong>“设置”</strong>进入<strong>“AI 大模型接入”</strong>，填写 Base URL（例如 <code>https://api.openai.com/v1</code> 或本地 Ollama <code>http://localhost:11434/v1</code>）、API Key 及模型名，点击“测试 AI 连接”验证。未配置 API Key 时系统会自动采用内置的领域规则推断分组名称，保证 0 门槛与离线可用。
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 弹窗底部 Footer -->
        <div class="modal-footer">
          <div class="footer-hint">
            <SvgIcon name="info" :size="14" />
            <span>按 ESC 键或点击蒙层任意位置可随时关闭此窗口</span>
          </div>

          <div class="footer-actions">
            <button
              v-if="currentTabIndex > 0"
              type="button"
              class="btn-footer-secondary"
              @click="handlePrevTab"
            >
              <SvgIcon name="chevron-up" :size="12" style="transform: rotate(-90deg)" />
              <span>{{ $t('guidePrev') }}</span>
            </button>

            <button
              v-if="currentTabIndex < tabs.length - 1"
              type="button"
              class="btn-footer-primary"
              @click="handleNextTab"
            >
              <span>{{ $t('guideNext') }}</span>
              <SvgIcon name="chevron-right" :size="13" />
            </button>

            <button
              v-else
              type="button"
              class="btn-footer-primary finish-btn"
              @click="emit('update:show', false)"
            >
              <SvgIcon name="check" :size="14" />
              <span>{{ $t('guideGotIt') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'update:show', value: boolean): void }>()

const { t: $t } = useI18n()

const currentTab = ref<string>('quick-start')

const tabs = computed(() => [
  { id: 'quick-start', label: $t('guideTabQuickStart'), icon: 'zap' },
  { id: 'mods', label: $t('guideTabMods'), icon: 'box' },
  { id: 'maps', label: $t('guideTabMaps'), icon: 'map' },
  { id: 'remote', label: $t('guideTabRemote'), icon: 'globe' },
  { id: 'faq', label: $t('guideTabFaq'), icon: 'help-circle' },
])

const currentTabIndex = computed(() => {
  return tabs.value.findIndex((t) => t.id === currentTab.value)
})

function handlePrevTab() {
  if (currentTabIndex.value > 0) {
    currentTab.value = tabs.value[currentTabIndex.value - 1].id
  }
}

function handleNextTab() {
  if (currentTabIndex.value < tabs.value.length - 1) {
    currentTab.value = tabs.value[currentTabIndex.value + 1].id
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.78);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 24px;
}

.guide-modal {
  width: 100%;
  max-width: 820px;
  max-height: 88vh;
  background: var(--bg-card, #131d2e);
  border: 1px solid var(--border-medium, #2a3a52);
  border-radius: var(--radius-lg, 12px);
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7);
  animation: modal-pop 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

@keyframes modal-pop {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(8px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-header {
  padding: 16px 22px;
  border-bottom: 1px solid var(--border-subtle, #1e2c40);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(0, 0, 0, 0.25);
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon-box {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md, 8px);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%);
  border: 1px solid rgba(56, 189, 248, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary, #38bdf8);
}

.modal-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.3px;
}

.modal-subtitle {
  margin: 3px 0 0;
  font-size: 12px;
  color: var(--text-dim, #94a3b8);
}

.modal-close-btn {
  background: transparent;
  border: none;
  color: var(--text-dim, #94a3b8);
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm, 6px);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast, 0.15s ease);
}

.modal-close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.modal-body {
  padding: 0;
  overflow: hidden;
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 顶部 Tab 导航条 */
.guide-nav-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid var(--border-subtle, #1e2c40);
  overflow-x: auto;
}

.guide-nav-btn {
  appearance: none;
  background: transparent;
  border: 1px solid transparent;
  color: var(--text-dim, #94a3b8);
  padding: 6px 14px;
  border-radius: var(--radius-md, 8px);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  position: relative;
  transition: all var(--transition-fast, 0.15s ease);
  white-space: nowrap;
}

.guide-nav-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.guide-nav-btn.active {
  background: rgba(56, 189, 248, 0.12);
  color: var(--accent-primary, #38bdf8);
  border-color: rgba(56, 189, 248, 0.3);
}

/* 内容区域 */
.guide-content-area {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1;
  max-height: calc(88vh - 160px);
}

.panel-hero {
  margin-bottom: 20px;
}

.hero-badge {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(56, 189, 248, 0.15);
  color: var(--accent-primary, #38bdf8);
  border-radius: 4px;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.8px;
  margin-bottom: 6px;
}

.panel-title {
  margin: 0;
  font-size: 17px;
  font-weight: 700;
  color: #fff;
}

.panel-desc {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--text-dim, #94a3b8);
}

/* 步骤流卡片 */
.steps-flow {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.step-card {
  display: flex;
  gap: 16px;
  padding: 14px 18px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
  transition: border-color var(--transition-fast, 0.15s ease);
}

.step-card:hover {
  border-color: rgba(56, 189, 248, 0.35);
}

.step-num {
  font-size: 22px;
  font-weight: 800;
  color: rgba(56, 189, 248, 0.4);
  font-family: monospace;
  line-height: 1;
}

.step-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.step-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
}

.step-icon {
  color: var(--accent-primary, #38bdf8);
}

.step-info p {
  margin: 0;
  font-size: 13px;
  color: #cbd5e1;
  line-height: 1.6;
}

.step-info strong {
  color: #fff;
}

.step-info code {
  background: rgba(0, 0, 0, 0.4);
  padding: 2px 5px;
  border-radius: 4px;
  font-family: Consolas, monospace;
  color: #38bdf8;
  font-size: 12px;
}

.tip-box {
  margin-top: 6px;
  padding: 6px 10px;
  background: rgba(56, 189, 248, 0.08);
  border-left: 3px solid var(--accent-primary, #38bdf8);
  border-radius: 0 4px 4px 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #94a3b8;
}

/* 特性卡片网格 */
.feature-cards-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.feature-item-card {
  padding: 16px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
}

.feature-card-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.feature-card-top h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
}

.feature-icon-box {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-purple {
  background: rgba(168, 85, 247, 0.15);
  color: #c084fc;
}

.icon-green {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.icon-blue {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.feature-item-card p {
  margin: 0 0 10px;
  font-size: 13px;
  color: var(--text-dim, #94a3b8);
  line-height: 1.5;
}

.feature-highlight {
  background: rgba(0, 0, 0, 0.35);
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  color: #cbd5e1;
  line-height: 1.5;
}

.feature-highlight strong {
  color: var(--accent-primary, #38bdf8);
}

/* 分栏并排指南 (地图与备注) */
.split-guide-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.guide-box {
  padding: 16px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
  display: flex;
  flex-direction: column;
}

.guide-box-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.guide-box-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
}

.text-accent {
  color: var(--accent-primary, #38bdf8);
}

.guide-box p {
  margin: 0 0 10px;
  font-size: 12.5px;
  color: var(--text-dim, #94a3b8);
  line-height: 1.5;
}

.guide-list {
  margin: 0 0 12px;
  padding-left: 18px;
  font-size: 12px;
  color: #cbd5e1;
  line-height: 1.6;
}

.code-preview-box {
  margin-top: auto;
  background: rgba(0, 0, 0, 0.4);
  padding: 6px 10px;
  border-radius: 4px;
  border: 1px dashed rgba(56, 189, 248, 0.3);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.code-preview-box code {
  color: #38bdf8;
  font-family: Consolas, monospace;
  font-weight: 600;
  font-size: 13px;
}

.code-tip {
  font-size: 11px;
  color: #64748b;
}

/* 远程服务器 4 格步骤 */
.server-steps-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.server-step-card {
  padding: 14px 16px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
}

.card-step-badge {
  display: inline-block;
  padding: 2px 7px;
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  font-size: 11px;
  font-weight: 700;
  border-radius: 4px;
  margin-bottom: 8px;
}

.server-step-card p {
  margin: 0;
  font-size: 12.5px;
  color: #cbd5e1;
  line-height: 1.5;
}

/* FAQ 列表 */
.faq-accordion {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.faq-item {
  padding: 14px 18px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
}

.faq-q {
  margin: 0 0 8px;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.faq-a {
  margin: 0;
  font-size: 13px;
  color: var(--text-dim, #94a3b8);
  line-height: 1.6;
}

.faq-a code {
  color: #38bdf8;
  background: rgba(0, 0, 0, 0.3);
  padding: 1px 4px;
  border-radius: 3px;
}

/* 弹窗底部操作栏 */
.modal-footer {
  padding: 14px 22px;
  border-top: 1px solid var(--border-subtle, #1e2c40);
  background: rgba(0, 0, 0, 0.25);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-dim, #94a3b8);
}

.footer-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-footer-secondary {
  height: 34px;
  padding: 0 14px;
  border-radius: var(--radius-sm, 6px);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-dim, #94a3b8);
  border: 1px solid var(--border-medium, #2a3a52);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-footer-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-footer-primary {
  height: 34px;
  padding: 0 16px;
  border-radius: var(--radius-sm, 6px);
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: #fff;
  border: 1px solid rgba(56, 189, 248, 0.4);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(2, 132, 199, 0.3);
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-footer-primary:hover {
  filter: brightness(1.15);
  box-shadow: 0 0 14px rgba(56, 189, 248, 0.5);
}

.finish-btn {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
  border-color: rgba(52, 211, 153, 0.4);
  box-shadow: 0 2px 8px rgba(5, 150, 105, 0.3);
}

.finish-btn:hover {
  box-shadow: 0 0 14px rgba(52, 211, 153, 0.5);
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
