(function () {
  if (typeof window === 'undefined') return;

  if (!window.runtime) {
    window.runtime = {
      EventsOnMultiple: function () {},
      EventsOn: function () {},
      EventsOff: function () {},
      EventsOffAll: function () {},
      EventsEmit: function () {},
      LogPrint: console.log,
      LogTrace: console.trace,
      LogDebug: console.debug,
      LogInfo: console.info,
      LogWarning: console.warn,
      LogError: console.error,
    };
  }

  if (!window.go) {
    var params = new URLSearchParams(window.location.search);
    var isEmpty = params.get('empty') === '1';

    var mockDir = isEmpty ? '' : 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2';
    var mockLocal = isEmpty
      ? []
      : [
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
        ];

    var mockWorkshop = isEmpty
      ? []
      : [
          {
            name: '2847592011.vpk',
            path: 'D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2\\left4dead2\\addons\\workshop\\2847592011.vpk',
            size: 84200000,
            modified_time: 1720000000,
            is_valid: true,
            is_disabled: false,
            addontitle: '黑色狂欢节重制版：Extended Remix',
            display_name: '黑色狂欢节重制版：Extended Remix',
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
        ];

    var mockRemote = [
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
    ];

    window.go = {
      main: {
        App: {
          GetLastDirectory: function () {
            return Promise.resolve(mockDir);
          },
          BrowseDirectory: function () {
            return Promise.resolve(mockDir);
          },
          GetDownloadsDirectory: function () {
            return Promise.resolve('C:\\Users\\User\\Downloads');
          },
          ScanVpkFiles: function () {
            return Promise.resolve({
              success: true,
              local_vpks: mockLocal,
              workshop_vpks: mockWorkshop,
            });
          },
          GetSftpSettings: function () {
            return Promise.resolve({
              address: '192.168.1.100:22',
              username: 'steam',
              password: '••••••••',
              addons_path: '/opt/l4d2/left4dead2/addons',
            });
          },
          SaveSftpSettings: function () {
            return Promise.resolve({ success: true });
          },
          ListRemoteVpks: function () {
            return Promise.resolve({
              success: true,
              files: mockRemote,
            });
          },
          DeleteVpkFile: function () {
            return Promise.resolve({ success: true });
          },
          DeleteMultipleVpkFiles: function () {
            return Promise.resolve({ success: true, count: 1 });
          },
          EnableVpkFile: function () {
            return Promise.resolve({ success: true });
          },
          ExportVpks: function () {
            return Promise.resolve({ success: true, exported_count: 1 });
          },
          ImportCompressedVpks: function () {
            return Promise.resolve({ success: true, imported_count: 1 });
          },
          LocalizeWorkshopVpks: function () {
            return Promise.resolve({ success: true, count: 1 });
          },
          OpenDirectory: function () {
            return Promise.resolve();
          },
          OpenExternalUrl: function () {
            return Promise.resolve();
          },
          LogFromJs: function () {
            return Promise.resolve();
          },
          RestartL4d2Server: function () {
            return Promise.resolve({ success: true, duration_s: 1.2 });
          },
        },
      },
    };

    function waitForEl(selector, cb) {
      var el = document.querySelector(selector);
      if (el) return cb(el);
      var tries = 0;
      var timer = setInterval(function () {
        tries++;
        var found = document.querySelector(selector);
        if (found || tries > 30) {
          clearInterval(timer);
          if (found) cb(found);
        }
      }, 100);
    }

    if (params.get('compact') === '1' || params.get('view') === 'compact') {
      try { localStorage.setItem('vpk-view-mode', 'compact'); } catch(e){}
    } else {
      try { localStorage.setItem('vpk-view-mode', 'card'); } catch(e){}
    }

    // Selection is handled directly in store on load when ?selected=1
    if (params.get('settings') === '1') {
      waitForEl('.settings-btn', function (el) {
        el.click();
      });
    }
    if (params.get('sftp') === '1') {
      waitForEl('.sftp-toggle-btn', function (el) {
        el.click();
      });
    }
  }
})();
