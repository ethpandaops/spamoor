
(function() {
  // Authentication state mirrored from window.ethpandaops.authenticatoor.
  // In open mode (no auth provider configured), authDisabled stays true
  // and isAuthenticated is unconditionally true so all UI is unlocked.
  var authState = {
    token: null,
    user: null,
    expiresAt: null,
    isAuthenticated: false,
    authDisabled: false
  };

  window.addEventListener('DOMContentLoaded', function() {
    initControls();
    window.setInterval(updateTimers, 1000);
    initAuth();
  });

  var tooltipDict = {};
  var tooltipIdx = 1;
  var spamoor = window.spamoor = {
    initControls: initControls,
    renderRecentTime: renderRecentTime,
    tooltipDict: tooltipDict,
    getAuthToken: getAuthToken,
    isAuthenticated: isAuthenticated,
    authFetch: authFetch,
    authState: authState,
  };

  // Returns the ethpandaops.authenticatoor client when an auth provider is
  // configured AND its client.js has loaded. Returns null otherwise (open
  // mode).
  function authClient() {
    if (!window.spamoorConfig || !window.spamoorConfig.authProviderURL) return null;
    if (!window.ethpandaops || !window.ethpandaops.authenticatoor) return null;
    return window.ethpandaops.authenticatoor;
  }

  // Initialize authentication. In open mode, mark everything as authed
  // and bail. In remote mode, wire the login button and run checkLogin
  // (fragment → cache → silent iframe) in the background.
  function initAuth() {
    var client = authClient();

    if (!client) {
      authState.isAuthenticated = true;
      authState.authDisabled = true;
      updateAuthUI();
      return;
    }

    authState.authDisabled = false;
    updateAuthUI();

    var loginBtn = document.getElementById('loginBtn');
    if (loginBtn) {
      loginBtn.addEventListener('click', function(e) {
        e.preventDefault();
        client.login();
      });
    }

    client.checkLogin().then(function(info) {
      if (info && info.authenticated) {
        applyClientState(info);
      }
    }).catch(function() {
      // Stay unauthenticated; user can click Login.
    });

    // Pick up subsequent state changes (token refresh, logout, etc.)
    // without requiring a page reload.
    if (typeof client.onStateChange === 'function') {
      client.onStateChange(function(info) { applyClientState(info); });
    }
  }

  // Mirror the auth client's state into our local authState.
  function applyClientState(info) {
    authState.token = info.token || null;
    authState.user = info.user || info.email || null;
    authState.expiresAt = info.exp ? info.exp * 1000 : null;
    authState.isAuthenticated = !!info.authenticated;
    authState.authDisabled = false;
    updateAuthUI();
  }

  // Update UI based on auth state
  function updateAuthUI() {
    var loginBtn = document.getElementById('loginBtn');
    var userInfo = document.getElementById('userInfo');
    var userName = document.getElementById('userName');

    if (loginBtn && userInfo) {
      if (authState.isAuthenticated) {
        loginBtn.classList.add('d-none');
        userInfo.classList.remove('d-none');
        if (userName) {
          userName.textContent = authState.user || 'User';
        }
      } else {
        loginBtn.classList.remove('d-none');
        userInfo.classList.add('d-none');
      }
    }

    // Update page elements based on auth state
    document.querySelectorAll('[data-auth-required]').forEach(function(el) {
      if (authState.isAuthenticated) {
        el.classList.remove('d-none');
      } else {
        el.classList.add('d-none');
      }
    });

    document.querySelectorAll('[data-auth-hide]').forEach(function(el) {
      if (authState.isAuthenticated) {
        el.classList.add('d-none');
      } else {
        el.classList.remove('d-none');
      }
    });

    // Dispatch event for other scripts to listen
    window.dispatchEvent(new CustomEvent('authStateChanged', { detail: authState }));
  }

  // Get current auth token. In open mode there's no token to send;
  // returns null and authFetch leaves the request unauthenticated.
  function getAuthToken() {
    var client = authClient();
    if (client) {
      var t = client.getToken();
      if (t) return t;
    }
    return null;
  }

  // Check if user is authenticated. In open mode this is always true
  // (the backend treats every request as authorized).
  function isAuthenticated() {
    if (authState.authDisabled) return true;
    var client = authClient();
    if (client) return !!client.isLoggedIn();
    return false;
  }

  // Fetch with the current bearer token attached. On 401 the auth client
  // is asked to re-check (the client itself decides whether to silent-
  // refresh or surface a logged-out state).
  function authFetch(url, options) {
    options = options || {};
    options.headers = options.headers || {};

    var token = getAuthToken();
    if (token) {
      options.headers['Authorization'] = 'Bearer ' + token;
    }

    return fetch(url, options).then(function(response) {
      if (response.status === 401) {
        var client = authClient();
        if (client) {
          client.checkLogin().then(function(info) {
            if (info && info.authenticated) applyClientState(info);
          });
        }
      }
      return response;
    });
  }

  function initControls() {
    // init tooltips
    document.querySelectorAll('[data-bs-toggle="tooltip"]').forEach(initTooltip);
    cleanupTooltips();

    // init clipboard buttons
    document.querySelectorAll("[data-clipboard-text]").forEach(initCopyBtn);
    document.querySelectorAll("[data-clipboard-target]").forEach(initCopyBtn);
  }

  function initTooltip(el) {
    if($(el).data("tooltip-init"))
      return;
    //console.log("init tooltip", el);
    var idx = tooltipIdx++;
    $(el).data("tooltip-init", idx).attr("data-tooltip-idx", idx.toString());
    $(el).tooltip();
    var tooltip = bootstrap.Tooltip.getInstance(el);
    tooltipDict[idx] = {
      element: el,
      tooltip: tooltip,
    };
  }

  function cleanupTooltips() {
    Object.keys(spamoor.tooltipDict).forEach(function(idx) {
      var ref = spamoor.tooltipDict[idx];
      if(document.body.contains(ref.element)) return;
      ref.tooltip.dispose();
      delete spamoor.tooltipDict[idx];
    });
  }

  function initCopyBtn(el) {
    if($(el).data("clipboard-init"))
      return;
    $(el).data("clipboard-init", true);
    var clipboard = new ClipboardJS(el);
    clipboard.on("success", onClipboardSuccess);
    clipboard.on("error", onClipboardError);
  }

  function onClipboardSuccess(e) {
    var title = e.trigger.getAttribute("data-bs-original-title");
    var tooltip = bootstrap.Tooltip.getInstance(e.trigger);
    tooltip.setContent({ '.tooltip-inner': 'Copied!' });
    tooltip.show();
    setTimeout(function () {
      tooltip.setContent({ '.tooltip-inner': title });
    }, 1000);
  }

  function onClipboardError(e) {
    var title = e.trigger.getAttribute("data-bs-original-title");
    var tooltip = bootstrap.Tooltip.getInstance(e.trigger);
    tooltip.setContent({ '.tooltip-inner': 'Failed to Copy!' });
    tooltip.show();
    setTimeout(function () {
      tooltip.setContent({ '.tooltip-inner': title });
    }, 1000);
  }

  function updateTimers() {
    var timerEls = document.querySelectorAll("[data-timer]");
    timerEls.forEach(function(timerEl) {
      var time = timerEl.getAttribute("data-timer");
      var textEls = Array.prototype.filter.call(timerEl.querySelectorAll("*"), function(el) { return el.firstChild && el.firstChild.nodeType === 3 });
      var textEl = textEls.length ? textEls[0] : timerEl;
      
      textEl.innerText = renderRecentTime(time);
    });
  }

  function renderRecentTime(time) {
    var duration = time - Math.floor(new Date().getTime() / 1000);
    var timeStr= "";
    var absDuration = Math.abs(duration);

    if (absDuration < 1) {
      return "now";
    } else if (absDuration < 60) {
      timeStr = absDuration + " sec."
    } else if (absDuration < 60*60) {
      timeStr = (Math.floor(absDuration / 60)) + " min."
    } else if (absDuration < 24*60*60) {
      timeStr = (Math.floor(absDuration / (60 * 60))) + " hr."
    } else {
      timeStr = (Math.floor(absDuration / (60 * 60 * 24))) + " day."
    }
    if (duration < 0) {
      return timeStr + " ago";
    } else {
      return "in " + timeStr;
    }
  }

  
})()
