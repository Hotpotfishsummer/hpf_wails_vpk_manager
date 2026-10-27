// 开发环境与无后端浏览器预览的 Mock 数据
// 仅在 window.go 不存在时激活，不影响实际 Wails 运行

export function initDevMock() {
  if (typeof window === 'undefined') return

  if (!(window as any).runtime) {
    ;(window as any).runtime = {
      EventsOnMultiple: () => {},
      EventsOn: () => {},
      EventsOff: () => {},
      EventsOffAll: () => {},
      EventsEmit: () => {},
      LogPrint: console.log,
      LogTrace: console.trace,
      LogDebug: console.debug,
      LogInfo: console.info,
      LogWarning: console.warn,
      LogError: console.error,
    }
  }

  if (!(window as any).go) {
    const mockState = {
      dir: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2',
      localVpks: [
        {
          name: 'miku_zoey.vpk',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\miku_zoey.vpk',
          size: 24500000,
          modified_time: 1718000000,
          is_valid: true,
          is_disabled: false,
          addontitle: '初音未来 (替换佐伊 Survivor)',
          display_name: '初音未来 (替换佐伊 Survivor)',
        },
        {
          name: 'urban_flight.vpk',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\urban_flight.vpk',
          size: 345000000,
          modified_time: 1717000000,
          is_valid: true,
          is_disabled: false,
          addontitle: '都市逃亡 (Urban Flight)',
          display_name: '都市逃亡 (Urban Flight)',
          maps: [
            { name: 'uf1_boulevard', display_name: '1: 林荫大道' },
            { name: 'uf2_rooftops', display_name: '2: 屋顶穿梭' },
            { name: 'uf3_harbor', display_name: '3: 港口突围' },
            { name: 'uf4_airfield', display_name: '4: 机场逃离' },
          ],
        },
        {
          name: 'hd_weapons_pack.vpk',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\hd_weapons_pack.vpk',
          size: 156000000,
          modified_time: 1715000000,
          is_valid: true,
          is_disabled: false,
          addontitle: '全武器高清贴图重置包 v2.0',
          display_name: '全武器高清贴图重置包 v2.0',
        },
        {
          name: 'tank_metal_bgm.vpk.disabled',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\tank_metal_bgm.vpk.disabled',
          size: 8900000,
          modified_time: 1710000000,
          is_valid: true,
          is_disabled: true,
          addontitle: 'Tank 登场重金属摇滚 BGM',
          display_name: 'Tank 登场重金属摇滚 BGM',
        },
      ],
      workshopVpks: [
        {
          name: '2847592011.vpk',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\workshop\\2847592011.vpk',
          size: 84200000,
          modified_time: 1720000000,
          is_valid: true,
          is_disabled: false,
          addontitle: '黑色狂欢节重制版：Extended Remix',
          maps: [
            { name: 'c2m1_highway', display_name: '高速公路' },
            { name: 'c2m2_fairgrounds', display_name: '游乐场' },
            { name: 'c2m3_coaster', display_name: '过山车' },
            { name: 'c2m4_barns', display_name: '谷仓' },
            { name: 'c2m5_concert', display_name: '音乐会' },
          ],
        },
        {
          name: '1948201944.vpk',
          path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\workshop\\1948201944.vpk',
          size: 12400000,
          modified_time: 1719000000,
          is_valid: true,
          is_disabled: false,
          addontitle: '智能生还者 Bot AI 优化模组',
          display_name: '智能生还者 Bot AI 优化模组',
        },
      ],
      remoteVpks: [
        {
          name: 'competitive_configs.vpk',
          path: '/opt/l4d2/left4dead2/addons/competitive_configs.vpk',
          size: 4200000,
          modified_time: 1721000000,
          is_disabled: false,
        },
        {
          name: 'sourcemod_base.vpk',
          path: '/opt/l4d2/left4dead2/addons/sourcemod_base.vpk',
          size: 15800000,
          modified_time: 1720500000,
          is_disabled: false,
        },
      ],
    }

    // 默认空目录还是模拟目录可通过 url 参数切换: ?mock=1 开启填充，默认开启方便排查
    const params = new URLSearchParams(window.location.search)
    const isMockLoaded = params.get('empty') !== '1'

    ;(window as any).go = {
      main: {
        App: {
          GetLastDirectory: async () => (isMockLoaded ? mockState.dir : ''),
          BrowseDirectory: async () => mockState.dir,
          GetDownloadsDirectory: async () => 'C:\\Users\\User\\Downloads',
          ScanVpkFiles: async () => ({
            success: true,
            local_vpks: isMockLoaded ? mockState.localVpks : [],
            workshop_vpks: isMockLoaded ? mockState.workshopVpks : [],
          }),
          GetSftpSettings: async () => ({
            address: '192.168.1.100:22',
            username: 'steam',
            password: '••••••••',
            addons_path: '/opt/l4d2/left4dead2/addons',
          }),
          SaveSftpSettings: async () => ({ success: true }),
          ListRemoteVpks: async () => ({
            success: true,
            files: mockState.remoteVpks,
          }),
          DeleteVpkFile: async () => ({ success: true }),
          DeleteMultipleVpkFiles: async () => ({ success: true, count: 1 }),
          EnableVpkFile: async () => ({ success: true }),
          ExportVpks: async () => ({ success: true, exported_count: 1 }),
          ImportCompressedVpks: async () => ({ success: true, imported_count: 1 }),
          ImportFiles: async (paths: string[]) => ({ success: true, imported_count: paths?.length || 1 }),
          LocalizeWorkshopVpks: async () => ({ success: true, count: 1 }),
          OpenDirectory: async () => {},
          OpenExternalUrl: async () => {},
          LogFromJs: async () => {},
          RestartL4d2Server: async () => ({ success: true, duration_s: 1.2 }),
          DetectL4D2Directories: async (_deepScan = false) => ({
            success: true,
            clients: [
              {
                path: mockState.dir,
                exe_path: `${mockState.dir}\\left4dead2.exe`,
                source: 'Steam 游戏库',
                source_name: 'steam',
                mod_count: 24,
                has_workshop: true,
                last_modified: 1718000000,
              },
              {
                path: 'D:\\Games\\Left4Dead2_NonSteam',
                exe_path: 'D:\\Games\\Left4Dead2_NonSteam\\left4dead2.exe',
                source: '系统注册表/记录',
                source_name: 'registry',
                mod_count: 5,
                has_workshop: false,
                last_modified: 1715000000,
              },
            ],
          }),
          ValidateL4D2Directory: async (_dir: string) => true,
          LaunchGame: async () => ({ success: true, message: 'Mock: 游戏已成功启动' }),
          IsSteamClient: async (dir: string = '') => dir.toLowerCase().includes('steam'),
          GetLaunchSettings: async () => ({
            novid: true,
            chinese_lang: true,
            direct_launch: false,
            custom_args: '',
            server_connect_cmd: 'connect 192.168.1.200',
            network_rates_cmd: "cl_updaterate 100; cl_cmdrate 100; cl_interp 0.01; cl_interp_ratio 1",
          }),
          SaveLaunchSettings: async (_settings: any) => ({ success: true }),
          GetVpkGroups: async () => {
            try {
              const saved = localStorage.getItem('vpk-groups-dev')
              return saved ? JSON.parse(saved) : []
            } catch {
              return []
            }
          },
          SaveVpkGroups: async (groups: any) => {
            try {
              localStorage.setItem('vpk-groups-dev', JSON.stringify(groups))
            } catch {}
            return { success: true }
          },
          GetAiSettings: async () => {
            try {
              const saved = localStorage.getItem('ai-settings-dev')
              if (saved) return JSON.parse(saved)
            } catch {}
            return {
              api_key: '',
              base_url: 'https://api.openai.com/v1',
              model: 'gpt-4o-mini',
              temperature: 0.7,
              max_tokens: 500,
            }
          },
          SaveAiSettings: async (settings: any) => {
            try {
              localStorage.setItem('ai-settings-dev', JSON.stringify(settings))
            } catch {}
            return { success: true }
          },
          TestAiConnection: async (settings: any) => {
            if (!settings?.api_key) {
              return { success: false, error: 'Mock: API Key 不能为空' }
            }
            return { success: true, message: `Mock: 连接成功！模型 [${settings.model || 'gpt-4o-mini'}] 响应: OK` }
          },
          PredictGroupName: async (_files: any[]) => {
            return {
              success: true,
              predicted_name: '生还者模组套装',
              reason: 'Mock 智能预测',
              source: 'ai',
            }
          },
          AiChat: async () => {
            return {
              success: true,
              content: 'Mock: AI 响应成功',
            }
          },
        },
      },
    }
  }
}
